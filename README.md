# gitlab-cli

## Description

Utility written on [Go](https://golang.org) that allows to manipulate gitlab
repositories through using it's API (currently using `v3`).

### Features
 * [X] Create MR
 * [X] set default `project-name`, `url` and `token` from config
 * [X] Set default `src-branch` to current git branch
 * [X] Set default `title` to `HEAD` commit message
 
## Config example
```yaml
token: PRIVATETOKEN
url: https://gitlab.com
project_name: namespace/project
```

## Usage
example
```bash
// without config
gitlab-cli mr create -P namespace/project -a senid231 -A "PRIVATETOKEN" -U "https://gitlab.com/api/v3"
// with config
gitlab-cli mr create -a senid231
```

## Installation
If you has go environment you can just do:
```bash
go install github.com/senid231/gitlab-cli
```
or just extract binary to your system: 
[Linux](https://github.com/senid231/gitlab-cli/releases/download/v1.1.0/gitlab-cli_linux_amd64.tar.gz)
or
[OSX](https://github.com/senid231/gitlab-cli/releases/download/v1.1.0/gitlab-cli_darwin_amd64.tar.gz)
