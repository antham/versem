# Versem [![CircleCI](https://circleci.com/gh/antham/versem/tree/master.svg?style=svg)](https://circleci.com/gh/antham/versem/tree/master) [![codecov](https://codecov.io/gh/antham/versem/branch/master/graph/badge.svg)](https://codecov.io/gh/antham/versem) [![Go Report Card](https://goreportcard.com/badge/github.com/antham/versem)](https://goreportcard.com/report/github.com/antham/versem) [![GolangCI](https://golangci.com/badges/github.com/antham/versem.svg)](https://golangci.com) [![GoDoc](https://godoc.org/github.com/antham/versem?status.svg)](http://godoc.org/github.com/antham/versem) [![GitHub tag](https://img.shields.io/github/tag/antham/versem.svg)]()

Versem creates a semver git tag and a github release when merging a pull request according to the version label set on the repository.

---

- [Usage](#usage)
- [Setup](#setup)
- [Contribute](#contribute)

---

## Usage

```
Semver manager

Usage:
  versem [command]

Available Commands:
  help        Help about any command
  label       Manage pull request labels
  release     Manage release

Flags:
  -h, --help   help for versem

Use "versem [command] --help" for more information about a command.

```

You must define several environment variables : _GITHUB_OWNER_, _GITHUB_REPOSITORY_ and _GITHUB_TOKEN_

### label check [commitSha|pullRequestId]

Ensure a semver label is defined on a pull request or a commit that belong to a pull request, if not it exit with an error, if the commit is not tied to a pull request, it aborts without any errors.

### label create

Create semver labels (patch, minor, major) on a repository.

### release create [commitSha]

Create the semver tag using label version defined in pull request tied to the commit given as argument, if the commit is not tied to a pull request, it aborts without any errors.

## Setup

Download the binary from the release page according to your architecture : https://github.com/antham/versem/releases

## Contribute

If you want to add a new feature to chyle project, the best way is to open a ticket first to know exactly how to implement your changes in code.

### Setup

After cloning the repository you need to install vendors with `go mod vendor`
To test your changes locally you can run all tests with : `make test-all`.
