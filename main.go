package main

import (
	"fmt"
	"github.com/droundy/goopt"
	"github.com/senid231/gitlab-cli/api"
	"os"
)

func newOpts() *api.Opts {
	argOpts := api.Opts{}
	argOpts.Token = goopt.StringWithLabel([]string{"-A", "--access-token"}, "", "<token>", "private access token")
	argOpts.Project = goopt.StringWithLabel([]string{"-P", "--project"}, "", "<name>", "namespaced name of project")
	argOpts.SrcBranch = goopt.StringWithLabel([]string{"-s", "--src-branch"}, "", "<branch>", "namespaced name of project")
	argOpts.DstBranch = goopt.StringWithLabel([]string{"-t", "--target-branch"}, "", "<branch>", "target branch")
	argOpts.Assignee = goopt.StringWithLabel([]string{"-a", "--assignee"}, "", "<username>", "assignee username")
	argOpts.Title = goopt.StringWithLabel([]string{"-T", "--title"}, "", "<text>", "title of MR")
	argOpts.BaseUrl = goopt.StringWithLabel([]string{"-U", "--url"}, "https://gitlab.com", "<baseUrl>", "base URL for gitlab API")
	return &argOpts
}

func parseArgs() []string {
	goopt.Description = func() string {
		return "Console gitlab client.\n" +
			"Use v3 API"
	}
	goopt.Version = "1.0.0"
	goopt.Summary = "api [options] <subcommand> <action> [options]"
	goopt.Author = "Denis Talakevich <senid231@gmail.com>"
	goopt.ExtraUsage = "Subcommands:\n" +
		"  mr\t\tMerge Request manipulations\n" +
		"  Actions:\n" +
		"    create\tcreates merge request"

	goopt.Parse(nil)
	return goopt.Args
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
	opts := newOpts()
	args := parseArgs()
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
