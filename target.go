package main

import (
	"errors"
	"strings"
)

type GitHubTarget struct {
	owner      string
	repository string
	directory  string
	ref        string
}

// InferTargetFromUrl will take a valid GitHub URL and create a GitHubTarget object from it.
// A GitHub URL will take the form "https://github.com/{owner}/{repo}/(tree/<ref>)?/(<path>/<to>/<root>)?", where:
// - "tree/<ref>" is optional, if not present it defaults to the default branch of the repo
// - "<ref>" can either be a commit hash or branch
// - "<path>/<to>/<root>" is optional
// - If "<path>/<to>/<root>" is present, then "tree/<ref>" MUST be present
func InferTargetFromUrl(url string) (GitHubTarget, error) {
	target := GitHubTarget{}

	_url, ok := strings.CutPrefix(url, "https://github.com/")
	if !ok {
		return GitHubTarget{}, errors.New(`invalid GitHub URL`)
	}

	parts := strings.Split(_url, "/")

	if len(parts) < 2 {
		return GitHubTarget{}, errors.New(`invalid GitHub URL`)
	}
	target.owner = parts[0]
	target.repository = parts[1]

	// Nothing more to process, we are at the root of the repo in the default branch
	if len(parts) == 2 {
		return target, nil
	}

	// Otherwise, the URL will include at least the "/tree/<ref>" part
	if len(parts) < 4 || parts[2] != "tree" {
		return GitHubTarget{}, errors.New(`invalid GitHub URL`)
	}

	target.ref = parts[3]

	// If the target includes the <path>/<to>/<root> part.
	// It will be the rest of the string.
	if len(parts) > 4 {
		target.directory = strings.Join(parts[4:], "/")
	}

	return target, nil
}
