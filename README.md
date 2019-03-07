# ghcp [![CircleCI](https://circleci.com/gh/int128/ghcp.svg?style=shield)](https://circleci.com/gh/int128/ghcp) [![GoDoc](https://godoc.org/github.com/int128/ghcp?status.svg)](https://godoc.org/github.com/int128/ghcp)

ghcp is a command to copy (commit and push) files to a repository on GitHub.
It depends on GitHub API and works without Git commands.


## Getting Started

Download [the latest release](https://github.com/int128/ghcp/releases) or install from brew tap:

```sh
brew tap int128/ghcp
brew install ghcp
```

Run `ghcp -h` to see help:

```
Usage: ghcp [options] [file or directory...]

Options:
      --debug            Show debug logs
      --dry-run          Upload files but do not update the branch actually
  -m, --message string   Commit message (mandatory)
  -u, --owner string     GitHub repository owner (mandatory)
  -r, --repo string      GitHub repository name (mandatory)
      --token string     GitHub API token [$GITHUB_TOKEN]
```

To upload the files and create a commit on the default branch (typically master):

```
% export GITHUB_TOKEN=YOUR_TOKEN
% ghcp -u int128 -r sandbox -m 'Example Commit' dist/
```

You need to set your personal access token by `-token` option or `GITHUB_TOKEN` environment variable.

ghcp does not create a new commit if the default branch has same files.
Therefore it prevents an empty commit.

ghcp does not respect the current Git state or Git config.
You need to always set owner and name of a repository.


## Recipes

### Working with CI

TODO

Here is an example for CircleCI:

```sh
version: 2
jobs:
  build:
    steps:
      - run: |
          curl -L -o /tmp/ghcp https://github.com/int128/ghcp/releases/download/$GHCP_VERSION/ghcp_linux_amd64
          /tmp/ghcp -u owner -r repo -m "release" index.html
```

### GitHub Pages

TODO

### Homebrew tap

You can release your formula to the Homebrew tap.

1. Create a repository with prefix `homebrew-`, e.g. [`homebrew-ghcp`](https://github.com/int128/homebrew-ghcp).
1. Create a formula as like:
```sh
dist_sha256=$(shasum -a 256 -b your_app | cut -f1 -d' ')
cat <<EOF
class YourApp < Formula
  desc "Your awesome application"
  homepage "https://github.com/YOUR/REPO"
  url "https://github.com/YOUR/REPO/releases/download/v1.0.0/your_app_darwin_amd64"
  version "v1.0.0"
  sha256 "${dist_sha256}"

  def install
    bin.install "your_app_darwin_amd64" => "your_app"
  end

  test do
    system "#{bin}/your_app -h"
  end
end
EOF
```
1. Release the formula to the tap repository.
```sh
ghcp -u int128 -r homebrew-ghcp -m v1.0.0 ghcp.rb
```
1. Test the formula by `brew tap` and `brew install`.

ghcp is released to the tap by using ghcp self.
See also the scripts in [.circleci](.circleci).


## Contributions

This is an open source software.
Feel free to open issues and pull requests.
