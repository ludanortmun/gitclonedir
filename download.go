package main

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/google/go-github/v74/github"
)

type fileTree struct {
	name     string
	content  []byte
	children []*fileTree
}

func (f *fileTree) isDirectory() bool {
	return f.children != nil
}

type DownloadCommand struct {
	target     GitHubTarget
	outputPath string
	client     *github.Client
}

func NewDownloadCommand(target GitHubTarget, outputPath string) *DownloadCommand {
	return &DownloadCommand{
		target:     target,
		outputPath: outputPath,
		client:     github.NewClient(nil),
	}
}

func (cmd *DownloadCommand) Execute() error {
	files, err := cmd.startDownload(cmd.target.directory)
	if err != nil {
		return err
	}

	return saveToDisk(files, cmd.outputPath)
}

func saveToDisk(node *fileTree, path string) error {
	if node == nil {
		return errors.New("no files to save")
	}

	if !node.isDirectory() {
		return os.WriteFile(path+"/"+node.name, node.content, os.ModePerm)
	}

	newBasePath := path + "/" + node.name

	// Ensure the directory exists
	// We call this in the directory path to ensure it's only called once for each directory
	if err := os.MkdirAll(newBasePath, os.ModePerm); err != nil {
		return err
	}

	for _, child := range node.children {
		if err := saveToDisk(child, newBasePath); err != nil {
			return err
		}
	}

	return nil

}

// startDownload will start the download process for the given path.
// It will check if the path is a file or a directory and handle it accordingly.
func (cmd *DownloadCommand) startDownload(path string) (*fileTree, error) {
	opts := &github.RepositoryContentGetOptions{}
	if cmd.target.ref != "" {
		opts.Ref = cmd.target.ref
	}

	file, dir, _, err := cmd.client.Repositories.GetContents(
		context.Background(),
		cmd.target.owner,
		cmd.target.repository,
		path,
		opts,
	)

	if err != nil {
		return nil, err
	}

	// If the target is a file, download it directly
	if file != nil {
		return downloadFile(file)
	}

	// If the target is not a file, it will be a directory.
	pathParts := strings.Split(path, "/")
	dirname := pathParts[len(pathParts)-1]
	return cmd.downloadDirectory(dirname, dir)
}

// downloadFile will download the file from the GitHub repository using its download URL.
func downloadFile(item *github.RepositoryContent) (*fileTree, error) {
	res, err := http.Get(item.GetDownloadURL())
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, errors.New(res.Status)
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return &fileTree{
		name:     item.GetName(),
		content:  bytes,
		children: nil,
	}, nil
}

// downloadDirectory will download the contents of a directory from the GitHub repository.
// It will recursively download all files and directories within the given directory.
func (cmd *DownloadCommand) downloadDirectory(dirname string, items []*github.RepositoryContent) (*fileTree, error) {

	result := &fileTree{
		name:    dirname,
		content: nil,
	}

	children := make([]*fileTree, len(items))

	for i, item := range items {
		if item.GetType() == "file" {
			f, err := downloadFile(item)
			if err != nil {
				return nil, err
			}
			children[i] = f
		} else if item.GetType() == "dir" {
			dir, err := cmd.startDownload(item.GetPath())
			if err != nil {
				return nil, err
			}
			children[i] = dir
		}
	}

	result.children = children
	return result, nil
}
