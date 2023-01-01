package core

import (
	"encoding/base64"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

type workflowManifest struct {
	Entries []*workflowManifestEntry
}

type workflowManifestEntry struct {
	Trigger string `yaml:"trigger"`
	Ref     string `yaml:"ref"`
	Command string `yaml:"command"`
}

func getWorkflowManifest(repoOwner, repoName, ref string) (*workflowManifest, error) {
	fcont, err := GetFileFromRepository(repoOwner, repoName, ".workboat/workflows.yml", ref)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if fcont.Content == nil {
		return nil, nil
	}

	sDec, err := base64.StdEncoding.DecodeString(*fcont.Content)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	log.Debug().Str("manifest_content", *fcont.Content).Send()

	var entries []*workflowManifestEntry
	if err := yaml.Unmarshal(sDec, &entries); err != nil {
		return nil, newWorkflowError(err)
	}

	return &workflowManifest{
		Entries: entries,
	}, nil
}
