# labels

A command line tool to add and remove Github labels.

## Installation

- Ensure you have a working `golang` installation
- Clone this repo
- Run `go install` from this repo
- Ensure a Github token is stored within your git config (see https://github.com/blog/180-local-github-config)

## Example Usage

Add labels to a repo

```sh
labels add jimmyhillis/labels label:color label:color
```

Remove all labels from a repo

```sh
labels remove jimmyhillis/labels
```
