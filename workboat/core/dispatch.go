package core

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"strings"
)

type workflowError error

func newWorkflowError(err error) error {
	return workflowError(err)
}

func dispatch(event string, body map[string]any) {
	if err := runDispatch(event, body); err != nil {
		log.Error().Stack().Err(err).Msg("failed to run dispatch")
	}
}

func runDispatch(event string, body map[string]any) error {

	if !strings.EqualFold(event, "push") {
		return nil
	}

	repoOwner, repoName := getRepositoryNameFromBody(body)

	if repoOwner == "" || repoName == "" {
		return nil
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
		return nil
	}

	enqueueDockerJob(&dockerJob{
		RepoOwner:     repoOwner,
		RepoName:      repoName,
		ManifestEntry: selectedManifestEntry,
	})

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
