# ghcp [![CircleCI](https://circleci.com/gh/int128/ghcp.svg?style=shield)](https://circleci.com/gh/int128/ghcp) [![codecov](https://codecov.io/gh/int128/ghcp/branch/master/graph/badge.svg)](https://codecov.io/gh/int128/ghcp) [![GoDoc](https://godoc.org/github.com/int128/ghcp?status.svg)](https://godoc.org/github.com/int128/ghcp)

This is a command to commit files to a GitHub repository.
It depends on GitHub APIs and works without git installation.


## Getting Started

Install the latest release from [here](https://github.com/int128/ghcp/releases) or Homebrew.

```sh
# GitHub Releases
curl -L -o ~/bin/ghcp https://github.com/int128/ghcp/releases/download/v1.3.0/ghcp_linux_amd64

# Homebrew
brew tap int128/ghcp
brew install ghcp
```

You need to get a personal access token from the [settings](https://github.com/settings/tokens) and set it to the `GITHUB_TOKEN` environment variable or `--token` option.


### Release to GitHub Pages

You can commit files to the `gh-pages` branch.

```sh
ghcp -u YOUR -r REPO -b gh-pages -m MESSAGE index.html index.css
```

You can commit files without parents, i.e. discarding history of the past commits.

```sh
ghcp -u YOUR -r REPO -b gh-pages --no-parent -m MESSAGE index.html index.css
```

### Release your Homebrew formula

You can release a Homebrew formula to a tap repository.

Create a repository with the prefix `homebrew-`, e.g. `homebrew-hello`.

Create a formula.
For example, the following script will generate a file `hello.rb`:

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

Commit the formula to the repository.

```sh
ghcp -u YOUR -r homebrew-hello -m v1.0.0 hello.rb
```

Now you can install the formula from the repository.

```sh
brew tap YOUR/hello
brew install hello
```

See also [Makefile](Makefile) because ghcp is released to [the tap repository](https://github.com/int128/homebrew-ghcp) by using ghcp self.

### Bump version string

You can change version string in files such as README or build script.

```sh
# substitute version string in files
sed -i -e "s/version '[0-9.]*'/version '$TAG'/g" README.md build.gradle

# commit the changes to a new branch
ghcp -u YOUR -r REPO -B bump-v1.1.0 -m v1.1.0 README.md build.gradle
```


### Usage

```
Usage: ghcp [options] [file or directory...]

Options:
      --api string          GitHub API v3 URL (v4 will be inferred) [$GITHUB_API]
  -b, --branch string       Update the branch (default: default branch of repository)
      --debug               Show debug logs
  -C, --directory string    Change to directory before copy
      --dry-run             Upload files but do not update the branch actually
  -m, --message string      Commit message (mandatory)
  -B, --new-branch string   Create a branch
      --no-file-mode        Ignore executable bit of file and treat as 0644
      --no-parent           Create a commit without a parent
  -u, --owner string        GitHub repository owner (mandatory)
      --parent string       Create a commit from the parent branch or tag (default: default branch of repository)
  -r, --repo string         GitHub repository name (mandatory)
      --token string        GitHub API token [$GITHUB_TOKEN]
```

Author and comitter of a commit are set to the login user, that depends on the token.

It does not create a new commit if the branch has same files.
Therefore it prevents an empty commit.

It does not read the current Git config and Git state.
You need to always set owner and name of a repository.

### Examples

To update the default branch of repository `YOUR/REPO` with files in directory `dist/`:

```sh
ghcp -u YOUR -r REPO -m MESSAGE dist/
```

To update branch `gh-pages` of repository `YOUR/REPO` with files in directory `dist/`:

```sh
ghcp -u YOUR -r REPO -b gh-pages -m MESSAGE dist/
```

To create a new branch `topic` from the default branch on repository `YOUR/REPO` with files in directory `dist/`:

```sh
ghcp -u YOUR -r REPO -B topic -m MESSAGE dist/
```

To create a new branch `topic` from branch `develop` on repository `YOUR/REPO` with files in directory `dist/`:

```sh
ghcp -u YOUR -r REPO -B topic --parent develop -m MESSAGE dist/
```

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
