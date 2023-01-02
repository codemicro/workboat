package core

import (
	"code.gitea.io/sdk/gitea"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"net/url"
	"strings"
)

type workflowError error

func newWorkflowError(err error) error {
	return workflowError(err)
}

func dispatch(event string, body map[string]any) {
	if err := runDispatch(event, body); err != nil {
		log.Error().Err(err).Stack().Msg("failed to run dispatch")
	}
}

func runDispatch(event string, body map[string]any) error {

	if !strings.EqualFold(event, "push") {
		return nil
	}

	rr, found := body["repository"]
	if !found {
		return errors.New("unable to decode body")
	}

	r, ok := rr.(map[string]any)
	if !ok {
		return errors.New("unable to decode body")
	}

	fullName := getStringFromBody("full_name", r)
	splitFullName := strings.Split(fullName, "/")

	if len(splitFullName) != 2 {
		return errors.New("malformed repository.full_name")
	}

	repoOwner, repoName := splitFullName[0], splitFullName[1]

	cloneURL := getStringFromBody("clone_url", r)
	if cloneURL == "" {
		return errors.New("could not get repository.clone_url from body")
	}

	commitSha := getStringFromBody("after", body)
	if commitSha == "" {
		return errors.New("missing commit sha in after key")
	}

	ref := getStringFromBody("ref", body)

	manifest, err := getWorkflowManifest(repoOwner, repoName, ref)
	if err != nil {
		if _, ok := err.(workflowError); ok {
			// TODO: Report this as a workflow error.
		}
		return errors.WithStack(err)
	}

	var selectedManifestEntry *workflowManifestEntry
	for _, e := range manifest.Entries {
		if e.Ref == ref {
			selectedManifestEntry = e
			break
		}
	}

	if selectedManifestEntry == nil {
		log.Debug().Msg("no manifest entry matching target")
		return nil
	}

	parsedCloneURL, err := url.Parse(cloneURL)
	if err != nil {
		return errors.WithStack(err)
	}

	enqueueDockerJob(&dockerJob{
		RepoOwner:     repoOwner,
		RepoName:      repoName,
		CommitSha:     commitSha,
		CloneURL:      parsedCloneURL,
		ManifestEntry: selectedManifestEntry,
	})

	if err := ReportRepoStatus(repoOwner, repoName, commitSha, &gitea.CreateStatusOption{
		State:       gitea.StatusPending,
		Description: "queued, waiting for runner to become available",
	}); err != nil {
		log.Warn().Err(err).Stack().Msg("unable to report status back to Gitea")
	}

	return nil
}

func getStringFromBody(key string, body map[string]any) string {
	rr, found := body[key]
	if found {
		r, ok := rr.(string)
		if ok {
			return r
		}
	}
	return ""
}

func getRepositoryNameFromBody(body map[string]any) (owner, name string) {
	rr, found := body["repository"]
	if !found {
		return
	}

	r, ok := rr.(map[string]any)
	if !ok {
		return
	}

	fullName := getStringFromBody("full_name", r)
	splitFullName := strings.Split(fullName, "/")

	if len(splitFullName) != 2 {
		// shouldn't happen, but just in case
		return
	}

	return splitFullName[0], splitFullName[1]
}
