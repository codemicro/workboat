package core

import (
	"github.com/pkg/errors"
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

	var entries []*workflowManifestEntry
	if err := yaml.Unmarshal([]byte(*fcont.Content), &entries); err != nil {
		return nil, newWorkflowError(err)
	}

	return &workflowManifest{
		Entries: entries,
	}, nil
}
