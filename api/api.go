// Package api provides access to gitlab HTTP API, YAML config and GIT repo info
package api

import (
	"log"
	"net/http"
	"strings"

	"fmt"
	"github.com/xanzy/go-gitlab"
)

// Opts contains options that api uses
type Opts struct {
	Token       *string
	Project     *string
	ForkProject *string
	SrcBranch   *string
	DstBranch   *string
	Assignee    *string
	Title       *string
	BaseURL     *string
	ProjectPath *string
	Description *string
	Debug       *bool
}

// CreateMergeRequest creates merge request according to given options
func CreateMergeRequest(opts *Opts) string {
	cli := setupAPI(*opts.BaseURL, *opts.Token)
	return createMergeRequest(cli, opts)
}

func findUser(api *gitlab.Client, username *string, debug bool) *gitlab.User {
	usersOpts := gitlab.ListUsersOptions{Username: username}
	usersList, _, err := api.Users.ListUsers(&usersOpts, debugFunc(debug))
	if err != nil {
		log.Fatalf("Can't fetch gitlab users\n%s\n", err)
	}
	if len(usersList) == 0 {
		log.Fatalf("Can't found assignee %#v", *username)
	}
	debugObject(usersList[0])
	return usersList[0]
}

func findProject(cli *gitlab.Client, projectName *string, debug bool) *gitlab.Project {
	project, _, err := cli.Projects.GetProject(*projectName, debugFunc(debug))
	if err != nil {
		log.Fatalf("Can't create merge request\n%s\n", err)
	}
	debugObject(project)
	return project
}

func setupAPI(url string, token string) *gitlab.Client {
	cli := gitlab.NewClient(nil, token)
	if strings.HasSuffix(url, "/") {
		url = strings.TrimSuffix(url, "/")
	}
	if !strings.HasSuffix(url, "/api/v3") {
		url += "/api/v3"
	}
	cli.SetBaseURL(url)
	return cli
}

func debugFunc(debug bool) gitlab.OptionFunc {
	if debug {
		return debugRequestURL
	}
	return func(r *http.Request) error { return nil }
}

func debugRequestURL(r *http.Request) error {
	fmt.Printf("Debug: %s\n", r.URL.String())
	return nil
}

func debugObject(o interface{}) {
	fmt.Printf("Debug: %#v\n", o)
}

func createMergeRequest(cli *gitlab.Client, opts *Opts) string {
	if *opts.Assignee == "" {
		log.Fatalln("assignee name required")
	}
	if *opts.Project == "" {
		log.Fatalln("project name required")
	}
	if *opts.ForkProject == "" {
		opts.ForkProject = opts.Project
	}

	assigneeID := findUser(cli, opts.Assignee, *opts.Debug).ID
	targetProjectID := findProject(cli, opts.Project, *opts.Debug).ID
	sourceProjectID := findProject(cli, opts.ForkProject, *opts.Debug).ID

	mrOpts := gitlab.CreateMergeRequestOptions{
		SourceBranch:    opts.SrcBranch,
		TargetBranch:    opts.DstBranch,
		Title:           opts.Title,
		AssigneeID:      &assigneeID,
		TargetProjectID: &targetProjectID,
		Description:     opts.Description,
	}

	mr, _, err := cli.MergeRequests.CreateMergeRequest(sourceProjectID, &mrOpts, debugFunc(*opts.Debug))
	if err != nil {
		log.Fatalf("Can't create merge request\n%s\n", err)
	}
	debugObject(mr)
	return mr.WebURL
}
