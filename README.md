# ghcp [![CircleCI](https://circleci.com/gh/int128/ghcp.svg?style=shield)](https://circleci.com/gh/int128/ghcp) [![GoDoc](https://godoc.org/github.com/int128/ghcp?status.svg)](https://godoc.org/github.com/int128/ghcp)

ghcp is a command to copy files to a repository on GitHub, like `git commit` and `git push`.
It depends on GitHub API and works without Git commands.


## Getting Started

Download [the latest release](https://github.com/int128/ghcp/releases) or install by brew tap:

```sh
brew tap int128/ghcp
brew install ghcp
```

To copy the files in directory `dist/` to the default branch of the repository `int128/sandbox`:

```sh
ghcp -u int128 -r sandbox -m message dist/
```

To copy the files in directory `dist/` to the branch `gh-pages` of the repository `int128/sandbox`:

```sh
ghcp -u int128 -r sandbox -b gh-pages -m message dist/
```

You need to get your personal access token from [GitHub settings](https://github.com/settings/tokens) and set it by `GITHUB_TOKEN` environment variable or `--token` option.

### Usage

```
Usage: ghcp [options] [file or directory...]

Options:
  -b, --branch string      Branch name (default: default branch of repository)
      --debug              Show debug logs
  -C, --directory string   Change to directory before copy
      --dry-run            Upload files but do not update the branch actually
  -m, --message string     Commit message (mandatory)
  -u, --owner string       GitHub repository owner (mandatory)
  -r, --repo string        GitHub repository name (mandatory)
      --token string       GitHub API token [$GITHUB_TOKEN]
```

It does not create a new commit if the branch has same files.
Therefore it prevents an empty commit.

It does not respect the current Git config and Git state.
You need to always set owner and name of a repository.


## Recipes

### Working with CI

TODO

Here is an example for CircleCI:

```yaml
version: 2
jobs:
  build:
    steps:
      - run: |
          curl -L -o /tmp/ghcp https://github.com/int128/ghcp/releases/download/$GHCP_VERSION/ghcp_linux_amd64
          /tmp/ghcp -u owner -r repo -m "message" index.html
```

### GitHub Pages

TODO

### Homebrew tap

You can release your formula to a tap repository.

You need to create a repository with prefix `homebrew-`, for example `homebrew-hello`.

Generate a formula as like:

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

You can release the formula to the tap repository.

```sh
ghcp -u YOUR -r homebrew-hello -m v1.0.0 hello.rb
```

Now you can install the formula:

```sh
brew tap YOUR/hello
brew install hello
```

ghcp is released to [the tap repository](https://github.com/int128/homebrew-ghcp) by using ghcp self and CircleCI.
See also [Makefile](Makefile) and [.circleci/config.yaml](.circleci/config.yaml).


## Related works

You can upload files to GitHub Releases by [`ghr`](https://github.com/tcnksm/ghr).

You can generate change log from Git history by [`ghch`](https://github.com/Songmu/ghch).


## Contributions

This is an open source software.
Feel free to open issues and pull requests.
