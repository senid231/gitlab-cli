package main

import (
	"fmt"
	"github.com/senid231/goopt"
	"github.com/xanzy/go-gitlab"
	"os"
	"strings"
)

// gitlab-ci create mr -P namespace/project -s my-branch -t master -a senid231 -A "GITLAB_ACCESS_TOKEN" -T "some title" -U "https://gitlab.com/api/v3"

const EXIT_CODE_ARGS = 1
const EXIT_CODE_API = 2

var token = goopt.StringWithLabel([]string{"-A", "--access-token"}, "", "<token>", "private access token")
var project = goopt.StringWithLabel([]string{"-P", "--project"}, "", "<name>", "namespaced name of project")
var srcBranch = goopt.StringWithLabel([]string{"-s", "--src-branch"}, "", "<branch>", "namespaced name of project")
var dstBranch = goopt.StringWithLabel([]string{"-t", "--target-branch"}, "", "<branch>", "target branch")
var assignee = goopt.StringWithLabel([]string{"-a", "--assignee"}, "", "<username>", "assignee username")
var title = goopt.StringWithLabel([]string{"-T", "--title"}, "", "<text>", "title of MR")
var baseUrl = goopt.StringWithLabel([]string{"-U", "--url"}, "", "<baseUrl>", "base URL for gitlab API")

func parseOpts() {
	goopt.Description = func() string {
		return "Console gitlab client.\nUse v3 API"
	}
	goopt.Version = "1.0.0"
	goopt.Summary = "gitlab-cli [options] <method> <resource> [options]"
	goopt.Author = "Denis Talakevich <senid231@gmail.com>"
	goopt.ExtraUsage = "Method:\n  create, c\tcreates resource\nResource:\n  projects\n  mr, merge-request\tMerge Request\n"
	goopt.Parse(nil)
}

func checkApiError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(EXIT_CODE_API)
	}
}

func createMergeRequest() {
	api := gitlab.NewClient(nil, *token)
	api.SetBaseURL(*baseUrl)

	usersOpts := gitlab.ListUsersOptions{Username: assignee}
	usersList, _, err := api.Users.ListUsers(&usersOpts)
	checkApiError(err)
	if len(usersList) == 0 {
		fmt.Printf("Can't found assignee %#v", assignee)
		os.Exit(EXIT_CODE_API)
	}
	assigneeID := usersList[0].ID

	mrOpts := gitlab.CreateMergeRequestOptions{
		SourceBranch: srcBranch,
		TargetBranch: dstBranch,
		Title:        title,
		AssigneeID:   &assigneeID,
	}
	mr, _, err := api.MergeRequests.CreateMergeRequest(*project, &mrOpts)
	checkApiError(err)
	webUrl := strings.Replace(api.BaseURL().String(), "/api/v3/", "", 1)
	fmt.Printf("%s/%s/merge_requests/%v\n", webUrl, *project, mr.IID)
}

func create(resource string) {
	switch resource {
	case "mr", "merge-request":
		createMergeRequest()
	default:
		fmt.Println(goopt.Help())
		fmt.Printf("invalid resource %s for method create\n", resource)
		os.Exit(EXIT_CODE_ARGS)
	}
}

func main() {
	//fmt.Printf("os.Args: %#v\n", os.Args)
	parseOpts()
	//fmt.Printf("token: %#v\nprojectID: %#v\ngoopt.Args: %#v\n", *token, *projectID, goopt.Args)
	if len(goopt.Args) != 2 {
		fmt.Println(goopt.Help(), "invalid quantity of args")
	}
	method := goopt.Args[0]
	switch method {
	case "create", "c":
		create(goopt.Args[1])
	default:
		fmt.Println(goopt.Help())
		fmt.Printf("invalid method %s", method)
		os.Exit(EXIT_CODE_ARGS)
	}
}
