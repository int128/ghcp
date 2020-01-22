# ghcp [![CircleCI](https://circleci.com/gh/int128/ghcp.svg?style=shield)](https://circleci.com/gh/int128/ghcp) [![codecov](https://codecov.io/gh/int128/ghcp/branch/master/graph/badge.svg)](https://codecov.io/gh/int128/ghcp) [![GoDoc](https://godoc.org/github.com/int128/ghcp?status.svg)](https://godoc.org/github.com/int128/ghcp)

This is a release engineering tool for GitHub.
It depends on GitHub APIs and works without git installation.

It provides the following features:

- Commit files to a repository
- Fork a repository and commit files to the forked repository
- Create a pull request
- Upload files to GitHub Releases


## Getting Started

You can install the latest release from [GitHub Releases](https://github.com/int128/ghcp/releases) or Homebrew.

```sh
# GitHub Releases
curl -fL -o /tmp/ghcp.zip https://github.com/int128/ghcp/releases/download/v1.7.0/ghcp_linux_amd64.zip
unzip /tmp/ghcp.zip -d ~/bin

# Homebrew
brew install int128/ghcp/ghcp
```

You need to get a personal access token from the [settings](https://github.com/settings/tokens) and set it to `GITHUB_TOKEN` environment variable or `--token` option.


### Commit files to a branch

To commit files to the default branch:

```sh
ghcp commit -u OWNER -r REPO -m MESSAGE file1 file2
```

To commit files to `feature` branch:

```sh
ghcp commit -u OWNER -r REPO -b feature -m MESSAGE file1 file2
```

If `feature` branch does not exist, ghcp will create it from the default branch.

To create `feature` branch from `develop` branch:

```sh
ghcp commit -u OWNER -r REPO -b feature --parent=develop -m MESSAGE file1 file2
```

If `feature` branch already exists, ghcp will fail.
Currently only fast-forward is supported.

ghcp performs a commit operation as follows:

- An author and committer of a commit are set to the login user (depending on the token).
- If the branch has same files, do not create a new commit. It prevents an empty commit.
- It excludes `.git` directories.
- It does not support `.gitconfig`.

You can set the following options.

```
Flags:
  -b, --branch string    Name of the branch to create or update (default: the default branch of repository)
      --dry-run          Upload files but do not update the branch actually
  -h, --help             help for commit
  -m, --message string   Commit message (mandatory)
      --no-file-mode     Ignore executable bit of file and treat as 0644
      --no-parent        Create a commit without a parent
  -u, --owner string     GitHub repository owner (mandatory)
      --parent string    Create a commit from the parent branch/tag (default: fast-forward)
  -r, --repo string      GitHub repository name (mandatory)
```


### Fork the repository and commit files to a new branch

To fork repository `UPSTREAM/REPO` and create `feature` branch from the default branch:

```sh
ghcp fork-commit -u UPSTREAM -r REPO -b feature -m MESSAGE file1 file2
```

To fork repository `UPSTREAM/REPO` and create `feature` branch from `develop` branch of the upstream:

```sh
ghcp fork-commit -u UPSTREAM -r REPO -b feature --parent develop -m MESSAGE file1 file2
```

If the branch already exists, ghcp will fail.
Currently only fast-forward is supported.

You can set the following options.

```
Flags:
  -b, --branch string    Name of the branch to create (mandatory)
      --dry-run          Upload files but do not update the branch actually
  -h, --help             help for fork-commit
  -m, --message string   Commit message (mandatory)
      --no-file-mode     Ignore executable bit of file and treat as 0644
  -u, --owner string     Upstream repository owner (mandatory)
      --parent string    Upstream branch name (default: the default branch of the upstream repository)
  -r, --repo string      Upstream repository name (mandatory)
```


### Create a pull request

To create a pull request from `feature` branch to the default branch:

```sh
ghcp pull-request -u OWNER -r REPO -b feature --title TITLE --body BODY
```

To create a pull request from `feature` branch to the `develop` branch:

```sh
ghcp pull-request -u OWNER -r REPO -b feature --base develop --title TITLE --body BODY
```

To create a pull request from `feature` branch of `OWNER/REPO` repository to the default branch of `UPSTREAM/REPO` repository:

```sh
ghcp pull-request -u OWNER -r REPO -b feature --base-owner UPSTREAM --base-repo REPO --title TITLE --body BODY
```

To create a pull request from `feature` branch of `OWNER/REPO` repository to the default branch of `UPSTREAM/REPO` repository:

```sh
ghcp pull-request -u OWNER -r REPO -b feature --base-owner UPSTREAM --base-repo REPO --base feature --title TITLE --body BODY
```

If a pull request already exists, ghcp do nothing.

You can set the following options.

```
Flags:
      --base string         Base branch name (default: default branch of base repository)
      --base-owner string   Base repository owner (default: head)
      --base-repo string    Base repository name (default: head)
      --body string         Body of a pull request
  -b, --head string         Head branch name (mandatory)
  -u, --head-owner string   Head repository owner (mandatory)
  -r, --head-repo string    Head repository name (mandatory)
  -h, --help                help for pull-request
      --title string        Title of a pull request (mandatory)
```


### Release assets

To upload files to the release associated to tag `v1.0.0`:

```sh
ghcp release -u OWNER -r REPO -t v1.0.0 dist/
```

If the release does not exist, it will create a release.
If the tag does not exist, it will create a tag from the default branch and create a release.

You can set the following options.

```
Flags:
      --dry-run        Do not create a release and assets actually
  -h, --help           help for release
  -u, --owner string   GitHub repository owner (mandatory)
  -r, --repo string    GitHub repository name (mandatory)
  -t, --tag string     Tag name (mandatory)
```


## Usage

### Global options

You can set the following options.

```
Global Flags:
      --api string         GitHub API v3 URL (v4 will be inferred) [$GITHUB_API]
      --debug              Show debug logs
  -C, --directory string   Change to directory before operation
      --token string       GitHub API token [$GITHUB_TOKEN]
```

### GitHub Enterprise

You can set a GitHub API v3 URL by `GITHUB_API` environment variable or `--api` option.

```sh
export GITHUB_API=https://github.example.com/api/v3/
```

GitHub API v4 URL will be automatically inferred from the v3 URL by resolving the relative path `../graphql`.


## Contributions

This is an open source software.
Feel free to open issues and pull requests.

Author: [Hidetake Iwata](https://github.com/int128)
