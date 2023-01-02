package core

import (
	"code.gitea.io/sdk/gitea"
	"context"
	"fmt"
	"github.com/codemicro/workboat/workboat/config"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	dockerClient "github.com/docker/docker/client"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"net/url"
	"time"
)

type dockerJob struct {
	RepoOwner     string
	RepoName      string
	CommitSha     string
	CloneURL      *url.URL
	ManifestEntry *workflowManifestEntry
}

var jobQueue = make(chan *dockerJob, 512)

func enqueueDockerJob(job *dockerJob) {
	jobQueue <- job
}

func StartJobConsumer() {
	go func() {
		for dj := range jobQueue {
			err := runDockerJob(dj)
			if err != nil {
				log.Error().Err(err).Stack().Msg("failed to run Docker job")
			}
		}
	}()
}

func runDockerJob(job *dockerJob) error {
	log.Info().Msgf("starting job for %s/%s@%s", job.RepoOwner, job.RepoName, job.ManifestEntry.Ref)

	if err := ReportRepoStatus(job.RepoOwner, job.RepoName, job.CommitSha, &gitea.CreateStatusOption{
		State:       gitea.StatusPending,
		Description: "running",
	}); err != nil {
		log.Warn().Err(err).Stack().Msg("unable to report status back to Gitea")
	}

	client, err := dockerClient.NewClientWithOpts(
		dockerClient.WithHost(config.Docker.Socket),
		dockerClient.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return errors.WithStack(err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
	defer cancel()

	// Arguments to the container should be: clone URL, auth string for that domain, user command
	authURL := &url.URL{
		Scheme: job.CloneURL.Scheme,
		User:   url.UserPassword(job.RepoOwner, config.Gitea.AccessToken),
		Host:   job.CloneURL.Host,
	}

	log.Info().Msg("creating container")

	// Spin up Docker container
	container, err := client.ContainerCreate(ctx, &container.Config{
		AttachStdout: true,
		AttachStderr: true,
		Env:          nil,
		Cmd:          []string{job.CloneURL.String(), authURL.String(), job.ManifestEntry.Command},
		Image:        config.Docker.Image,
	}, nil, nil, nil, "")
	if err != nil {
		return errors.WithStack(err)
	}

	shortContainerID := container.ID
	if len(shortContainerID) > 12 {
		shortContainerID = shortContainerID[:12]
	}

	log.Debug().Str("container_id", container.ID).Send()
	log.Info().Msg("starting container")

	if err := client.ContainerStart(ctx, container.ID, types.ContainerStartOptions{}); err != nil {
		_ = client.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{})
		return errors.WithStack(err)
	}

	if err := ReportRepoStatus(job.RepoOwner, job.RepoName, job.CommitSha, &gitea.CreateStatusOption{
		State:       gitea.StatusPending,
		Description: fmt.Sprintf("running (container ID: %s)", shortContainerID),
	}); err != nil {
		log.Warn().Err(err).Stack().Msg("unable to report status back to Gitea")
	}

	log.Info().Msg("waiting for container to end")

	okc, errc := client.ContainerWait(ctx, container.ID, "not-running")

	var exitStatusCode int64

	select {
	case err := <-errc:
		return errors.WithStack(err)
	case kc := <-okc:
		exitStatusCode = kc.StatusCode
	}

	log.Info().Msgf("finished %d", exitStatusCode)

	var statusOpt *gitea.CreateStatusOption

	if exitStatusCode == 0 {
		log.Info().Msgf("removing container %s", container.ID)
		if err := client.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{}); err != nil {
			return errors.WithStack(err)
		}
		log.Info().Msg("done")

		statusOpt = &gitea.CreateStatusOption{
			State:       gitea.StatusSuccess,
			Description: "finished",
		}
	} else {
		statusOpt = &gitea.CreateStatusOption{
			State:       gitea.StatusFailure,
			Description: fmt.Sprintf("failed with status code %d (container ID: %s)", exitStatusCode, shortContainerID),
		}
	}

	if err := ReportRepoStatus(job.RepoOwner, job.RepoName, job.CommitSha, statusOpt); err != nil {
		log.Warn().Err(err).Stack().Msg("unable to report status back to Gitea")
	}

	return nil
}
