# ghcp [![CircleCI](https://circleci.com/gh/int128/ghcp.svg?style=shield)](https://circleci.com/gh/int128/ghcp) [![codecov](https://codecov.io/gh/int128/ghcp/branch/master/graph/badge.svg)](https://codecov.io/gh/int128/ghcp) [![GoDoc](https://godoc.org/github.com/int128/ghcp?status.svg)](https://godoc.org/github.com/int128/ghcp)

This is a simple command to release files to a GitHub repository.
It depends on GitHub APIs and works without git installation.

Features:

- Commit files to a repository
- Fork a repository and commit files to the forked repository (for a pull request)
- Upload files to the GitHub Releases


## Getting Started

Install the latest release from [here](https://github.com/int128/ghcp/releases) or Homebrew.

```sh
# GitHub Releases
curl -fL -o /tmp/ghcp.zip https://github.com/int128/ghcp/releases/download/v1.5.1/ghcp_linux_amd64.zip
unzip /tmp/ghcp.zip -d ~/bin
rm /tmp/ghcp.zip

# Homebrew
brew tap int128/ghcp
brew install ghcp
```

You need to get a personal access token from the [settings](https://github.com/settings/tokens) and set it to the `GITHUB_TOKEN` environment variable or `--token` option.

Let's see the following examples.


### Example: Commit files to GitHub Pages

To commit the files to the `gh-pages` branch:

```sh
ghcp commit -u OWNER -r REPO -b gh-pages -m MESSAGE index.html index.css
```


### Example: Commit a Homebrew formula

Create a formula file like:

```rb
# hello.rb
class Hello < Formula
  desc "Your awesome application"
  homepage "https://github.com/OWNER/hello"
  url "https://github.com/OWNER/hello/releases/download/v1.0.0/hello_darwin_amd64"
  version "v1.0.0"
  sha256 "SHA256_SUM"

  def install
    bin.install "hello_darwin_amd64" => "hello"
  end

  test do
    system "#{bin}/hello -h"
  end
end
```

To commit the formula to the repository:

```sh
ghcp commit -u OWNER -r homebrew-hello -m "Release v1.0.0" hello.rb
```

Now you can install the formula from the repository.

```sh
brew install OWNER/hello/hello
```

See also [Makefile](Makefile).
ghcp is released to [the tap repository](https://github.com/int128/homebrew-ghcp) by using ghcp.


### Example: Bump version string

To change a version string of files in a repository:

```sh
# substitute version string in files
sed -i -e "s/version '[0-9.]*'/version '$TAG'/g" README.md build.gradle

# commit the changes to a branch
ghcp commit -u OWNER -r REPO -b bump-v1.1.0 -m "Bump the version to v1.1.0" README.md build.gradle
```


### Example: Fork a repository and create a branch for a pull request

To fork the upstream repository `UPSTREAM/REPO` and commit files to the branch `topic` of your repository `YOUR/REPO`:

```sh
ghcp fork-commit -u UPSTREAM -r REPO -b topic -m 'Add foo' foo.txt
```

You can manually create a pull request of the created branch.


### Example: Release assets

To upload files to the release associated to the tag `v1.0.0`:

```sh
ghcp release -u OWNER -r REPO -v v1.0.0 dist/
```

If the release does not exist, it will create a release.
If the tag does not exist, it will create a tag from the master branch and a release.


## Usage

```
Usage:
  ghcp [command]

Available Commands:
  commit      Commit files to the branch
  fork-commit Fork the repository and commit files to a branch
  help        Help about any command
  release     Release files to the repository

Flags:
      --api string         GitHub API v3 URL (v4 will be inferred) [$GITHUB_API]
      --debug              Show debug logs
  -C, --directory string   Change to directory before operation
  -h, --help               help for ghcp
      --token string       GitHub API token [$GITHUB_TOKEN]
      --version            version for ghcp
```

```
Usage:
  ghcp commit [flags] FILES...

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

```
Usage:
  ghcp fork-commit [flags] FILES...

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

```
Usage:
  ghcp release [flags] FILES...

Examples:
  To release files to the tag:
    ghcp release -u OWNER -r REPO -t TAG FILES...


Flags:
      --dry-run        Do not create a release and assets actually
  -h, --help           help for release
  -u, --owner string   GitHub repository owner (mandatory)
  -r, --repo string    GitHub repository name (mandatory)
  -t, --tag string     Tag name (mandatory)
```


### Behavior of `commit`

To commit files to the default branch:

```sh
ghcp commit -u OWNER -r REPO -m MESSAGE FILES...
```

To commit files to the branch:

```sh
ghcp commit -u OWNER -r REPO -b BRANCH -m MESSAGE FILES...
```

If the branch does not exist, ghcp creates a branch from the default branch.
It the branch exists, ghcp updates the branch by fast-forward.

To commit files to a new branch from the parent branch:

```sh
ghcp commit -u OWNER -r REPO -b BRANCH --parent PARENT -m MESSAGE FILES...
```

If the branch exists, ghcp cannot update the branch by fast-forward and will fail.

To commit files to a new branch without any parent:

```sh
ghcp commit -u OWNER -r REPO -b BRANCH --no-parent -m MESSAGE FILES...
```

If the branch exists, ghcp cannot update the branch by fast-forward and will fail.

ghcp performs a commit operation as follows:

- An author and committer of a commit are set to the login user depending on the token.
- It does not create a new commit if the branch has same files, that prevents an empty commit.
- It does not read the current Git config (`.gitconfig`) and Git state (`.git`) and you need to always set owner and name of a repository.
- It excludes `.git` directories.


### Behavior of `fork-commit`

To fork the repository and commit files to a branch:

```sh
ghcp fork-commit -u UPSTREAM -r REPO -b BRANCH -m MESSAGE FILES...
```

If the branch does not exist, ghcp creates a branch from the default branch of the upstream repository.
It the branch exists, ghcp updates the branch by fast-forward.

To fork the repository and commit files to a branch from the branch of the upstream:

```sh
ghcp fork-commit -u UPSTREAM -r REPO -b BRANCH --parent PARENT -m MESSAGE FILES...
```

If the branch does not exist, ghcp creates a branch from the branch of the upstream repository.
If the branch exists, ghcp cannot update the branch by fast-forward and will fail.


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
