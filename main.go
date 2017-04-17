package main

import (
	"fmt"
	"os"

	"github.com/droundy/goopt"
	"github.com/senid231/gitlab-cli/api"
)

func parseOpts() (*api.Opts, []string) {
	opts := &api.Opts{}
	opts.ProjectPath = goopt.StringWithLabel([]string{"-d", "--dir"}, "", "<path>", "path to project directory")
	opts.Token = goopt.StringWithLabel([]string{"-A", "--access-token"}, "", "<token>", "private access token")
	opts.Project = goopt.StringWithLabel([]string{"-P", "--project"}, "", "<name>", "namespaced name of project")
	opts.SrcBranch = goopt.StringWithLabel([]string{"-s", "--src-branch"}, "", "<branch>", "namespaced name of project")
	opts.DstBranch = goopt.StringWithLabel([]string{"-t", "--target-branch"}, "master", "<branch>", "target branch")
	opts.Assignee = goopt.StringWithLabel([]string{"-a", "--assignee"}, "", "<username>", "assignee username")
	opts.Title = goopt.StringWithLabel([]string{"-T", "--title"}, "", "<text>", "title of MR")
	opts.BaseUrl = goopt.StringWithLabel([]string{"-U", "--url"}, "", "<baseUrl>", "base URL for gitlab API")
	goopt.Description = func() string {
		return "Console gitlab client.\n" +
			"Use v3 API"
	}
	goopt.Version = "1.1.0"
	goopt.Summary = "api [options] <subcommand> <action> [options]"
	goopt.Author = "Denis Talakevich <senid231@gmail.com>"
	goopt.ExtraUsage = "Subcommands:\n" +
		"  mr\t\tMerge Request manipulations\n" +
		"  Actions:\n" +
		"    create\tcreates merge request"

	goopt.Parse(nil)
	args := goopt.Args
	return opts, args
}

func optErr(str string) {
	fmt.Println(goopt.Usage(), "\n"+str)
	os.Exit(1)
}

func optErrf(format string, x ...interface{}) {
	optErr(fmt.Sprintf(format, x...))
}

// gitlab-ci mr create -P namespace/project -s my-branch -t master -a senid231 -A "TOKEN" -T "some title" -U "https://gitlab.com/api/v3"

func main() {
	opts, args := parseOpts()
	conf, err := api.NewConfig(opts.ProjectPath)
	if err != nil {
		fmt.Printf("[warn] can't find config for %s\n%v\n", *opts.ProjectPath, err)
	}
	if err != nil && *opts.ProjectPath != "" {
		optErrf("Can't find config %s in %s", api.CONFIG_NAME, *opts.ProjectPath)
	}
	if *opts.ProjectPath == "" {
		*opts.ProjectPath = "./"
	}
	gitInfo, err := api.NewGitInfo(opts.ProjectPath)
	if err != nil {
		optErrf("Can't get git info: %v", err)
	}
	if *opts.SrcBranch == "" {
		*opts.SrcBranch = gitInfo.CurrentBranch
	}
	if *opts.Title == "" {
		*opts.Title = gitInfo.LastCommit
	}
	if *opts.Token == "" {
		*opts.Token = conf.Token
	}
	if *opts.Project == "" {
		*opts.Project = conf.ProjectName
	}
	if *opts.BaseUrl == "" {
		*opts.BaseUrl = conf.Url
	}

	if len(args) != 2 {
		optErr("invalid arguments")
	}
	switch {
	case args[0] == "mr" && args[1] == "create":
		fmt.Println(api.CreateMergeRequest(opts))
	default:
		optErrf("invalid subcommmand and/or action: %s %s", args[0], args[1])
	}
}
