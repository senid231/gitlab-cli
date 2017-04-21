// Package api provides access to gitlab HTTP API, YAML config and GIT repo info
package api

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"encoding/json"
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

const pathPrefix string = "/api/v3"

// sprintMR print MR details
func sprintMR(cli *gitlab.Client, mr *gitlab.MergeRequest, debug bool) string {
	source := findProject(cli, mr.SourceProjectID, debug)
	target := findProject(cli, mr.TargetProjectID, debug)
	return fmt.Sprintf("ID:\t\t%d\n\n"+
		"IID:\t\t%d\n\n"+
		"SourceBranch:\t%s\n\n"+
		"TargetBranch:\t%s\n\n"+
		"Author:\t\t%s\n\n"+
		"Assignee:\t%s\n\n"+
		"State:\t\t%s\n\n"+
		"WorkInProgress:\t%v\n\n"+
		"WebURL:\t\t%s\n\n"+
		"\nTitle:\n\n%s\n\n"+
		"\nDescription:\n\n%s\n",
		mr.ID,
		mr.IID,
		fmt.Sprintf("%s:%s", source.PathWithNamespace, mr.SourceBranch),
		fmt.Sprintf("%s:%s", target.PathWithNamespace, mr.TargetBranch),
		mr.Author.Username,
		mr.Assignee.Username,
		mr.State,
		mr.WorkInProgress,
		mr.WebURL,
		mr.Title,
		mr.Description)
}

// FindMergeRequest finds MR by iid
func FindMergeRequest(opts *Opts, mrID int) string {
	cli := setupAPI(*opts.BaseURL, *opts.Token)
	mrOpts := gitlab.ListMergeRequestsOptions{IID: &mrID}
	mrs, _, err := cli.MergeRequests.ListMergeRequests(*opts.Project, &mrOpts, debugFunc(*opts.Debug))
	if len(mrs) == 0 {

	}
	if err != nil {
		log.Fatalf("Can't find Merge Request\n%s\n", err)
	}
	return sprintMR(cli, mrs[0], *opts.Debug)
}

func prettyPrint(o interface{}, prefix string) {
	data, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		log.Printf("prettyPrint failed with error %q\n", err)
	}
	log.Printf("%s%s\n", prefix, data)
}

func findUser(api *gitlab.Client, username *string, debug bool) *gitlab.User {
	if debug {
		log.Printf("findUser with username %q\n", username)
	}
	usersOpts := gitlab.ListUsersOptions{Username: username}
	usersList, _, err := api.Users.ListUsers(&usersOpts, debugFunc(debug))
	if err != nil {
		log.Fatalf("Can't fetch gitlab users\n%s\n", err)
	}
	if len(usersList) == 0 {
		log.Fatalf("Can't find assignee %#v", *username)
	}
	if debug {
		prettyPrint(usersList[0], "Response:\n")
	}
	return usersList[0]
}

func findProject(cli *gitlab.Client, projectName interface{}, debug bool) *gitlab.Project {
	if debug {
		log.Printf("findProject with id or namespaced path %v\n", projectName)
	}
	project, _, err := cli.Projects.GetProject(projectName, debugFunc(debug))
	if err != nil {
		log.Fatalf("Can't find project\n%s\n", err)
	}
	if debug {
		prettyPrint(project, "Response:\n")
	}
	return project
}

func setupAPI(url string, token string) *gitlab.Client {
	cli := gitlab.NewClient(nil, token)
	if strings.HasSuffix(url, "/") {
		url = strings.TrimSuffix(url, "/")
	}
	if !strings.HasSuffix(url, pathPrefix) {
		url += pathPrefix
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
	fmt.Printf("RequestURL: %s\n", r.URL.String())
	return nil
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
	targetProjectID := findProject(cli, *opts.Project, *opts.Debug).ID
	sourceProjectID := findProject(cli, *opts.ForkProject, *opts.Debug).ID

	mrOpts := gitlab.CreateMergeRequestOptions{
		SourceBranch:    opts.SrcBranch,
		TargetBranch:    opts.DstBranch,
		Title:           opts.Title,
		AssigneeID:      &assigneeID,
		TargetProjectID: &targetProjectID,
		Description:     opts.Description,
	}

	if *opts.Debug {
		prettyPrint(mrOpts, fmt.Sprintf("create MR with sourceProjectID %d and options:\n", sourceProjectID))
	}
	mr, _, err := cli.MergeRequests.CreateMergeRequest(sourceProjectID, &mrOpts, debugFunc(*opts.Debug))
	if err != nil {
		log.Fatalf("Can't create merge request\n%s\n", err)
	}
	if *opts.Debug {
		prettyPrint(mr, "Response:\n")
	}
	return sprintMR(cli, mr, *opts.Debug)
}
