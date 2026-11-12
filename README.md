# Semver-Manager

[![Build](https://github.com/blueforge-labs/semver-manager/actions/workflows/release.yaml/badge.svg)](https://github.com/blueforge-labs/semver-manager/actions/workflows/release.yaml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.27.4-00ADD8.svg)](https://golang.org/)
[![Go Report Card](https://goreportcard.com/badge/github.com/blueforge-labs/semver-manager)](https://goreportcard.com/report/github.com/blueforge-labs/semver-manager)

A CLI tool for managing [Semantic Versioning 2.0.0](https://semver.org) compliant versions — increment, filter, and fetch version tags from GitHub repositories.

## Table of Contents

- [Quick Start](#quick-start)
- [Installation](#installation)
- [Commands](#commands)
  - [increment](#increment)
  - [filter](#filter)
  - [fetch](#fetch)
- [Contributing](#contributing)
- [License](#license)

## Quick Start

```bash
# Increment a patch version from nothing (defaults to 0.0.1)
smgr increment --level patch
# → 0.0.1

# Increment major from existing versions
smgr increment --level major --source-versions "0.0.0,1.0.0,0.1.0"
# → 2.0.0

# Filter versions by stream and pick the highest
smgr filter --versions "1.0.0 2.0.0 1.1.0" --stream "1.*.*" --highest
# → 1.1.0

# Fetch tags from a GitHub repository
smgr fetch -o blueforge-labs -r semver-manager -t "$GITHUB_TOKEN"
```

## Installation

### Precompiled Binaries (Linux only)

1. Download the binary from the [Releases](https://github.com/blueforge-labs/semver-manager/releases) page.
2. Move the `smgr` binary to a directory in your `PATH`.

### Building from Source

Prerequisites: [Go](https://golang.org/dl/) 1.27.4+

```bash
go build -o smgr ./src/cmd/smgr/
```

### Install with Just

From a checkout, install to `~/.local/bin` without elevated privileges (requires Go 1.27.4+ and just):

```bash
just install
```

`just install` also builds the binary. Ensure `~/.local/bin` is on your `PATH`. Set `PREFIX` to change the install prefix and `DESTDIR` to stage the installation, for example `PREFIX=/usr DESTDIR=/tmp/staging just install`. A system-wide prefix requires write access to its `bin` directory. Remove the installed binary manually when no longer needed.

### Install with Go

Install the latest version from GitHub with Go 1.27.4+ to your user-local bin directory:

```bash
GOBIN="$HOME/.local/bin" go install github.com/blueforge-labs/semver-manager/src/cmd/smgr@latest
```

Ensure `~/.local/bin` is on your `PATH`. This installs only for the current user and does not require elevated privileges.

### Homebrew

Tap this repository and install the development version from its formula:

```bash
brew tap blueforge-labs/semver-manager https://github.com/blueforge-labs/semver-manager.git
brew install --HEAD blueforge-labs/semver-manager/smgr
```

The explicit Git URL lets this repository serve as a tap without renaming it to use Homebrew's `homebrew-` convention. This builds the upstream `main` branch from source. Remove it with `brew uninstall smgr`, then `brew untap blueforge-labs/semver-manager`.

### Arch Linux

From a checkout, build the development package from the upstream `main` branch and install it with pacman:

```bash
cd packaging/arch
makepkg -si
```

This requires `base-devel`, Git and Go; `makepkg -si` can install the declared build dependencies with pacman. Remove the package with `sudo pacman -R semver-manager-git`. The PKGBUILD is not yet published to the AUR, so `paru -S semver-manager-git` will only work after it is published there.

### Go Module

Other Go projects can depend on this repository with `go get github.com/blueforge-labs/semver-manager@latest`. Library packages are under `src/`, for example `github.com/blueforge-labs/semver-manager/src/pkg/increment`.

### Go Client

The `src/client` package provides a client facade for fetching, filtering, incrementing, validating, and pushing version tags. A prefix is applied to tag names and can be any namespace, such as `v` or `release/`.

```go
package main

import (
  "log"
  "os"

  "github.com/blueforge-labs/semver-manager/src/client"
  "github.com/blueforge-labs/semver-manager/src/models"
  "github.com/blueforge-labs/semver-manager/src/pkg/filter"
)

func main() {
  api := client.New(client.Config{
    Owner:      "blueforge-labs",
    Repository: "semver-manager",
    Token:      os.Getenv("GITHUB_TOKEN"),
    Prefix:     "v",
  })
  versions, err := api.Fetch()
  if err != nil {
    log.Fatal(err)
  }
  versions, err = api.Filter(versions, filter.Highest())
  if err != nil {
    log.Fatal(err)
  }
  stream, err := models.ParseVersionPattern("*.*.*")
  if err != nil {
    log.Fatal(err)
  }
  next, err := api.Increment(versions, stream, models.Patch)
  if err != nil {
    log.Fatal(err)
  }
  name := api.FormatRelease(next) // v<next version>
  valid, err := api.Validate(name, "strict")
  if err != nil {
    log.Fatal(err)
  }
  if !valid {
    log.Fatal("incremented version is invalid")
  }
  log.Println(name)
}
```

`Push(ctx, version, target)` creates a Git tag only (not a GitHub Release). Its target may be a commit SHA or GitHub ref; an empty target uses the local repository's `HEAD`. The optional `LocalPath` config selects the local repository directory.

## Commands

All commands support a `--dry-run` flag. Use `-v 1` to enable debug logging.

Use the global `--prefix` option to read and emit tags in a namespace. For example, `smgr --prefix release/ filter --versions "release/1.0.0 release/2.0.0" --highest` prints `release/2.0.0`, and `smgr --prefix v increment --source-versions "v1.2.3"` prints `v1.2.4`.

### increment

Increment a version number (MAJOR.MINOR.PATCH) with optional pre-release support. Defaults to `0.0.1` if no source versions are provided.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--level` | `-l` | `patch` | Increment level: `major`, `minor`, `patch` (defaults to `patch` if `--target-stream` not specified) |
| `--target-stream` | `-t` | | Target stream pattern, e.g. `1.2.*` or `*.*.*-alpha.*` |
| `--source-versions` | `-s` | | Comma-separated source versions, e.g. `"0.0.0,1.0.0,1.1.0"` |

**Examples:**

```bash
# Default patch increment (no source → 0.0.1)
smgr increment

# Major increment from existing versions
smgr increment --level major --source-versions "0.0.0,1.0.0,0.1.0"
# → 2.0.0

# Pre-release increment targeting an alpha stream
smgr increment --level minor --source-versions "0.0.0,1.0.0,0.1.0" --target-stream "*.*.*-alpha.*"
# → 1.1.0-alpha.0
```

### filter

Filter a list of versions using stream patterns and/or select the highest match.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--versions` | `-V` | | Space-separated version list to filter |
| `--stream` | `-s` | | Stream pattern using `*` wildcards for any identifier |
| `--highest` | `-H` | `false` | Return only the highest version after filtering |

**Examples:**

```bash
# Filter by stream
smgr filter --versions "1.0.0 2.0.0 1.1.0" --stream "1.*.*"
# → 1.0.0 1.1.0

# Highest across all versions
smgr filter --versions "1.0.0 2.0.0 1.1.0" --highest
# → 2.0.0

# Combined: highest in a stream
smgr filter --versions "1.0.0 2.0.0 1.1.0" --stream "1.*.*" --highest
# → 1.1.0
```

<details>
<summary>Stream pattern reference</summary>

Stream patterns use `*` as a wildcard for any identifier. The absence of an identifier (or wildcard) means "no match", except for build metadata which matches anything when not specified.

| Pattern | Input versions | Result |
|---------|---------------|--------|
| `1.*.*` | 1.1.1, 2.1.1, 1.1.1+build01, 1.1.1-alpha | 1.1.1, 1.1.1+build01 |
| `*.*.*+AMD` | 1.1.1+AMD, 1.1.2, 1.1.1-alpha+AMD | 1.1.1+AMD |
| `1.0.0-Beta.*` | 1.0.0-Alpha.0, 1.0.0-Beta.0, 1.0.0-Beta.1 | 1.0.0-Beta.0, 1.0.0-Beta.1 |
| `1.0.0-Beta` | 0.1.0-Alpha, 0.1.0-Beta, 1.0.0-Beta | 1.0.0-Beta |
| `1.0.0-Beta.*` | 1.0.0-Beta.Alpha.0, 1.0.0-Beta, 1.0.0-Beta.Alpha | 1.0.0-Beta.Alpha.0, 1.0.0-Beta, 1.0.0-Beta.Alpha |
| `1.0.0-*.Beta.*` | 1.0.0-0.Alpha.0, 1.0.0-Beta.0, 1.0.0-Alpha.Beta.1 | 1.0.0-Beta.0, 1.0.0-Alpha.Beta.1 |

</details>

### fetch

Fetch semantic version tags from a GitHub repository. Automatically chains with all `filter` flags.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--owner` | `-o` | | Repository owner or organization |
| `--repo` | `-r` | | Repository name |
| `--token` | `-t` | | GitHub access token |
| `--platform` | `-p` | `github` | Platform to fetch from (currently: `github`) |
| `--stream` | `-s` | | *(from filter)* Stream pattern |
| `--highest` | `-H` | `false` | *(from filter)* Return only the highest version |
| `--versions` | `-V` | | *(from filter)* Additional versions to merge with fetched results |

**Examples:**

```bash
# Fetch all semver tags from a repo
smgr fetch -o blueforge-labs -r semver-manager -t "$GITHUB_TOKEN"

# Fetch and filter to highest in a stream
smgr fetch -o blueforge-labs -r semver-manager -t "$GITHUB_TOKEN" --stream "1.*.*" --highest
```

## Contributing

Contributions are welcomed! Please read the [Contributing Guidelines](CONTRIBUTING.md) to get started. By participating in this project, you agree to abide by the [Code of Conduct](CODE_OF_CONDUCT.md).

## License

Semver-Manager is released under the [MIT License](LICENSE).
