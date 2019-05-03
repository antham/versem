# Versem [![CircleCI](https://circleci.com/gh/antham/versem/tree/master.svg?style=svg)](https://circleci.com/gh/antham/versem/tree/master) [![codecov](https://codecov.io/gh/antham/versem/branch/master/graph/badge.svg)](https://codecov.io/gh/antham/versem) [![Go Report Card](https://goreportcard.com/badge/github.com/antham/versem)](https://goreportcard.com/report/github.com/antham/versem) [![GolangCI](https://golangci.com/badges/github.com/antham/versem.svg)](https://golangci.com) [![GoDoc](https://godoc.org/github.com/antham/versem?status.svg)](http://godoc.org/github.com/antham/versem) [![GitHub tag](https://img.shields.io/github/tag/antham/versem.svg)]()

Versem creates a semver git tag and a github release when merging a pull request according to the version label set on the repository.

---

- [Usage](#usage)
- [Documentation](#documentation)
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

Create labels (patch, minor, major and norelease) on a repository.

### release create [commitSha]

Create the semver tag using label version defined in pull request tied to the commit given as argument, if the commit is not tied to a pull request, it aborts without any errors.

## Documentation

### Workflow

You will use the command `versem label create` to add semver labels to your repository manually.

You will use the command `versem label check [pullRequestId]` in your CI to ensure a version label is linked to a pull request, when a pull request is built.

When the pull request is merged, you will use the command `versem release create` in your CI to create the release according to the version label defined in the pull request and according to the previous semver tag created.

Have a look to [versem-circleci](https://github.com/antham/versem-circleci) to have a full example of how to use it in a CI.

### Recommended settings

You should force in a CI, a check to ensure every pull request are labelled properly like in the example above.

You should enable this setting in your github repository : `Require branches to be up to date before merging`, to be sure 2 pull requests are not merged in the same time and avoiding release creation mess.

### Label norelease

When your pull request is not intended to produce a new semver tag, it must be labelled with `norelease`, the CI will pass and will not produce any new release on merge.

### V version suffix or not

If you started to prefix your semver tag with a `v`, versem will automatically detect it and will create new versions following this convention, if not it will continue not adding ```v``` as a suffix.

When no tag exist yet, a ```v``` is added for the first tag created.

### Repository not following semver before

If you want to install versem on a repository that wasn't following semver convention before, you must first create a proper semver tag manually before settting it, to let versem be able to understand from where it should start to tag.

## Setup

Download the binary from the release page according to your architecture : https://github.com/antham/versem/releases

## Contribute

If you want to add a new feature to versem project, the best way is to open a ticket first to know exactly how to implement your changes in code.

### Setup

After cloning the repository you need to install vendors with `go mod vendor`
To test your changes locally you can run all tests with : `make test-all`.
