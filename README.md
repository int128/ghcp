# ghcp [![CircleCI](https://circleci.com/gh/int128/ghcp.svg?style=shield)](https://circleci.com/gh/int128/ghcp) [![GoDoc](https://godoc.org/github.com/int128/ghcp?status.svg)](https://godoc.org/github.com/int128/ghcp)

ghcp is a command to commit and push files to a repository on GitHub.
It depends on GitHub API and works without Git commands.


## Getting Started

```
Usage: ghcp [options] [file or directory...]

Options:
  -debug
    	Show debug logs
  -dry-run
    	Upload files but do not update the branch actually
  -m string
    	Commit message (mandatory)
  -r string
    	GitHub repository name (mandatory)
  -token string
    	GitHub API token [$GITHUB_TOKEN]
  -u string
    	GitHub repository owner (mandatory)
```

You need to get a personal access token from GitHub settings.
You can set the token by `-token` option or `GITHUB_TOKEN` environment variable.

You need to specify repository owner and name.
ghcp does not read `.git` and `.gitconfig`.

ghcp will upload the files and create a commit on the default branch (typically master).

```
% export GITHUB_TOKEN=YOUR_TOKEN

% ghcp -u int128 -r sandbox -m 'Example Commit' dist/
```

ghcp does not create a new commit if the default branch has same files.
Therefore it prevents an empty commit.


## Contributions

This is an open source software.
Feel free to open issues and pull requests.
