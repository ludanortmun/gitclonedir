# Git Clone Directory

Simple CLI utility to download individual directories from GitHub repositories. 
Written in Go and leveraging the GitHub API and Google's [go-github library](https://github.com/google/go-github).

# Requirements

- Go 1.25 or higher

# Installation

To install the tool, run:

```bash
go install github.com/ludanortmun/gitclonedir@latest
```

# Usage

To download a directory from a GitHub repository, use the following command:

```bash
gitclonedir <github-url> [destination-path]
```

where:
- `<github-url>` is the URL of the GitHub repository and the path to the directory you want to download.
- `[destination-path]` is an optional argument specifying where to save the downloaded directory. If not provided, it defaults to the current directory.