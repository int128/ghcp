# ghcp [![CircleCI](https://circleci.com/gh/int128/ghcp.svg?style=shield)](https://circleci.com/gh/int128/ghcp) [![codecov](https://codecov.io/gh/int128/ghcp/branch/master/graph/badge.svg)](https://codecov.io/gh/int128/ghcp) [![GoDoc](https://godoc.org/github.com/int128/ghcp?status.svg)](https://godoc.org/github.com/int128/ghcp)

This is a command to commit files to a GitHub repository.
It depends on GitHub APIs and works without git installation.


## Getting Started

Install the latest release from [here](https://github.com/int128/ghcp/releases) or Homebrew.

```sh
# GitHub Releases
curl -L -o ~/bin/ghcp https://github.com/int128/ghcp/releases/download/v1.4.0/ghcp_linux_amd64
chmod +x ~/bin/ghcp

# Homebrew
brew tap int128/ghcp
brew install ghcp
```

You need to get a personal access token from the [settings](https://github.com/settings/tokens) and set it to the `GITHUB_TOKEN` environment variable or `--token` option.

Let's see the following examples.


### Example: Release assets to GitHub Pages

To commit the files to the `gh-pages` branch:

```sh
ghcp commit -u OWNER -r REPO -b gh-pages -m MESSAGE index.html index.css
```


### Example: Release your Homebrew formula

You can release a Homebrew formula to a tap repository.

You need to create a repository with the prefix `homebrew-`, e.g. `homebrew-hello`.

Then create a formula file like:

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
ghcp commit -u OWNER -r homebrew-hello -m v1.0.0 hello.rb
```

Now you can install the formula from the repository.

```sh
brew tap OWNER/hello
brew install hello
```

See also [Makefile](Makefile).
ghcp is released to [the tap repository](https://github.com/int128/homebrew-ghcp) by using ghcp.


### Example: Bump version string

You can change version string in files such as README or build script.
For example,

```sh
# substitute version string in files
sed -i -e "s/version '[0-9.]*'/version '$TAG'/g" README.md build.gradle

# commit the changes to a branch
ghcp commit -u OWNER -r REPO -b bump-v1.1.0 -m v1.1.0 README.md build.gradle
```


## Usage

```
Usage:
  ghcp [command]

Available Commands:
  commit      Commit files to the branch
  help        Help about any command

Flags:
      --api string         GitHub API v3 URL (v4 will be inferred) [$GITHUB_API]
      --debug              Show debug logs
  -C, --directory string   Change to directory before operation
  -h, --help               help for main
      --token string       GitHub API token [$GITHUB_TOKEN]
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


### Behaviors

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


### Working with CI

Here is an example for CircleCI:

```yaml
version: 2
jobs:
  release:
    steps:
      - run: |
          mkdir -p ~/bin
          echo 'export PATH="$HOME/bin:$PATH"' >> $BASH_ENV
      - run: |
          curl -L -o ~/bin/ghcp https://github.com/int128/ghcp/releases/download/v1.3.0/ghcp_linux_amd64
          chmod +x ~/bin/ghcp
      - checkout
      # release the Homebrew formula
      - run: |
          ghcp -u "$CIRCLE_PROJECT_USERNAME" -r "homebrew-$CIRCLE_PROJECT_REPONAME" -m "$CIRCLE_TAG" hello.rb
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
