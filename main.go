// gitlab-cli is a command line tool that uses gitlab api
package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/droundy/goopt"
	"github.com/senid231/gitlab-cli/api"
)

func parseOpts() (*api.Opts, []string) {
	opts := &api.Opts{}
	opts.ProjectPath = goopt.StringWithLabel([]string{"-d", "--dir"}, "", "<path>", "path to project directory")
	opts.Token = goopt.StringWithLabel([]string{"-A", "--access-token"}, "", "<token>", "private access token")
	opts.ForkProject = goopt.StringWithLabel([]string{"-F", "--fork-project"}, "", "<name>", "namespaced path to fork project")
	opts.Project = goopt.StringWithLabel([]string{"-P", "--project"}, "", "<name>", "namespaced path to main project")
	opts.SrcBranch = goopt.StringWithLabel([]string{"-s", "--src-branch"}, "", "<branch>", "namespaced name of project")
	opts.DstBranch = goopt.StringWithLabel([]string{"-t", "--target-branch"}, "master", "<branch>", "target branch")
	opts.Assignee = goopt.StringWithLabel([]string{"-a", "--assignee"}, "", "<username>", "assignee username")
	opts.Title = goopt.StringWithLabel([]string{"-T", "--title"}, "", "<text>", "title of MR")
	opts.BaseURL = goopt.StringWithLabel([]string{"-U", "--url"}, "", "<baseUrl>", "base URL for gitlab API")
	opts.Description = goopt.StringWithLabel([]string{"-D", "--description"}, "", "<text>", "MR description")
	opts.Debug = goopt.Flag([]string{"--debug"}, []string{}, "debug mode", "")

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

func main() {
	opts, args := parseOpts()
	conf, err := api.NewConfig(opts.ProjectPath)
	if err != nil {
		log.Printf("Warning! Can't find config for %s\n%v\n", *opts.ProjectPath, err)
	}
	if err != nil && *opts.ProjectPath != "" {
		log.Fatalf("Can't find config %s in %s", api.ConfigName, *opts.ProjectPath)
	}
	if *opts.ProjectPath == "" {
		*opts.ProjectPath = "./"
	}
	gitInfo, err := api.NewGitInfo(opts.ProjectPath)
	if err != nil {
		log.Printf("Can't get git info: %v", err)
	}
	if *opts.SrcBranch == "" {
		*opts.SrcBranch = gitInfo.CurrentBranch
	}
	commitParts := strings.SplitN(gitInfo.LastCommit, "\n", 2)
	if *opts.Title == "" {
		*opts.Title = commitParts[0]
	}
	if *opts.Description == "" && len(commitParts) > 1 {
		*opts.Description = commitParts[1]
	}
	if *opts.Token == "" {
		*opts.Token = conf.Token
	}
	if *opts.Project == "" {
		*opts.Project = conf.ProjectName
	}
	if *opts.ForkProject == "" {
		*opts.ForkProject = conf.ForkProjectName
	}
	if *opts.BaseURL == "" {
		*opts.BaseURL = conf.URL
	}

	if len(args) != 2 {
		log.Print("invalid arguments")
	}
	switch {
	case args[0] == "mr" && args[1] == "create":
		fmt.Println(api.CreateMergeRequest(opts))
	default:
		log.Printf("invalid subcommmand and/or action: %s %s", args[0], args[1])
	}
}
