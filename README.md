# ghcp [![CircleCI](https://circleci.com/gh/int128/ghcp.svg?style=shield)](https://circleci.com/gh/int128/ghcp) [![codecov](https://codecov.io/gh/int128/ghcp/branch/master/graph/badge.svg)](https://codecov.io/gh/int128/ghcp) [![GoDoc](https://godoc.org/github.com/int128/ghcp?status.svg)](https://godoc.org/github.com/int128/ghcp)

ghcp is a command to copy files to a repository on GitHub, like `git commit` and `git push`.
It depends on GitHub API and works without Git commands.


## Getting Started

Install [the latest release](https://github.com/int128/ghcp/releases) as follows:

```sh
# GitHub Releases
curl -L -o /usr/local/bin/ghcp https://github.com/int128/ghcp/releases/download/${ghcp_version}/ghcp_linux_amd64

# Homebrew
brew tap int128/ghcp
brew install ghcp

# Go
go get github.com/int128/ghcp
```

To copy files in the directory `dist/` to the default branch of the repository `YOUR/REPO`:

```sh
ghcp -u YOUR -r REPO -m MESSAGE dist/
```

To copy files in the directory `dist/` to the branch `gh-pages` of the repository `YOUR/REPO`:

```sh
ghcp -u YOUR -r REPO -b gh-pages -m MESSAGE dist/
```

You need to get a personal access token from [GitHub settings](https://github.com/settings/tokens) and set it by `GITHUB_TOKEN` environment variable or `--token` option.

### Usage

```
Usage: ghcp [options] [file or directory...]

Options:
  -b, --branch string      Branch name (default: default branch of repository)
      --debug              Show debug logs
  -C, --directory string   Change to directory before copy
      --dry-run            Upload files but do not update the branch actually
  -m, --message string     Commit message (mandatory)
      --no-file-mode       Ignore executable bit of file and treat as 0644
  -u, --owner string       GitHub repository owner (mandatory)
  -r, --repo string        GitHub repository name (mandatory)
      --token string       GitHub API token [$GITHUB_TOKEN]
```

Author and comitter of a commit are set to the login user, that depends on the token.

It does not create a new commit if the branch has same files.
Therefore it prevents an empty commit.

It does not read the current Git config and Git state.
You need to always set owner and name of a repository.


## Use cases

### Release to GitHub Pages

You can release your site to [GitHub Pages](https://pages.github.com/).

You need to create a branch `gh-pages` on the repository before running ghcp.

To copy files to the `gh-pages` branch:

```sh
ghcp -u YOUR -r REPO -b gh-pages -m MESSAGE index.html
```

### Release to Homebrew tap

You can release your Homebrew formula to a tap repository.

You need to create a repository with prefix `homebrew-`, for example `homebrew-hello`.

You need to generate a formula.
For example, the following script will generate `hello.rb`:

```sh
cat > hello.rb <<EOF
class Hello < Formula
  desc "Your awesome application"
  homepage "https://github.com/YOUR/hello"
  url "https://github.com/YOUR/hello/releases/download/v1.0.0/hello_darwin_amd64"
  version "v1.0.0"
  sha256 "$(shasum -a 256 -b hello | cut -f1 -d' ')"

  def install
    bin.install "hello_darwin_amd64" => "hello"
  end

  test do
    system "#{bin}/hello -h"
  end
end
EOF
```

To copy the formula to the tap repository:

```sh
ghcp -u YOUR -r homebrew-hello -m v1.0.0 hello.rb
```

Now we can install the formula by the following commands:

```sh
brew tap YOUR/hello
brew install hello
```

See also [Makefile](Makefile) because ghcp is released to [the tap repository](https://github.com/int128/homebrew-ghcp) by using ghcp self.


## Contributions

This is an open source software.
Feel free to open issues and pull requests.
