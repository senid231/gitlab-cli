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
help usage
```
Usage of gitlab-cli:
        api [options] <subcommand> <action> [options]
Options:
  -d <path>      --dir=<path>              path to project directory
  -A <token>     --access-token=<token>    private access token
  -P <name>      --project=<name>          namespaced name of project
  -s <branch>    --src-branch=<branch>     namespaced name of project
  -t <branch>    --target-branch=<branch>  target branch
  -a <username>  --assignee=<username>     assignee username
  -T <text>      --title=<text>            title of MR
  -U <baseUrl>   --url=<baseUrl>           base URL for gitlab API
  -h             --help                    Show usage message
  -v             --version                 Show version
Subcommands:
  mr            Merge Request manipulations
  Actions:
    create      creates merge request

```
example
```
$ gitlab-cli mr create -a senid231
https://gitlab.com/namespace/project/merge_requests/123
```

## Installation
If you has go environment you can do:
```
go install github.com/senid231/gitlab-cli
```
Or just extract binary to your system: 
[Linux](https://github.com/senid231/gitlab-cli/releases/download/v1.1.0/gitlab-cli_linux_amd64.tar.gz)
or
[OSX](https://github.com/senid231/gitlab-cli/releases/download/v1.1.0/gitlab-cli_darwin_amd64.tar.gz)
