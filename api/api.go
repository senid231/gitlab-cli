package api

import (
	"fmt"
	"os"
	"strings"

	"github.com/xanzy/go-gitlab"
)

type Opts struct {
	Token       *string
	Project     *string
	SrcBranch   *string
	DstBranch   *string
	Assignee    *string
	Title       *string
	BaseUrl     *string
	ProjectPath *string
}

func failErr(str string, err error) {
	failStr(str + "\n" + err.Error())
}

func failStrf(format string, x ...interface{}) {
	failStr(fmt.Sprintf(format, x...))
}

func failStr(str string) {
	fmt.Println("Error: " + str)
	os.Exit(1)
}

func CreateMergeRequest(opts *Opts) string {
	cli := setupApi(*opts.BaseUrl, *opts.Token)
	return createMergeRequest(cli, opts)
}

func findUser(api *gitlab.Client, username *string) *gitlab.User {
	usersOpts := gitlab.ListUsersOptions{Username: username}
	usersList, _, err := api.Users.ListUsers(&usersOpts)
	if err != nil {
		failErr("Can't fetch gitlab users", err)
	}
	if len(usersList) == 0 {
		failStrf("Can't found assignee %#v", *username)
	}
	return usersList[0]
}

func setupApi(url string, token string) *gitlab.Client {
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

func createMergeRequest(cli *gitlab.Client, opts *Opts) string {
	if *opts.Assignee == "" {
		failStr("assignee name required")
	}
	if *opts.Project == "" {
		failStr("project name required")
	}

	assigneeID := findUser(cli, opts.Assignee).ID

	mrOpts := gitlab.CreateMergeRequestOptions{
		SourceBranch: opts.SrcBranch,
		TargetBranch: opts.DstBranch,
		Title:        opts.Title,
		AssigneeID:   &assigneeID,
	}
	mr, _, err := cli.MergeRequests.CreateMergeRequest(*opts.Project, &mrOpts)
	if err != nil {
		failErr("Can't create merge request", err)
	}
	webUrl := strings.Replace(cli.BaseURL().String(), "/api/v3/", "", 1)
	return fmt.Sprintf("%s/%s/merge_requests/%v", webUrl, *opts.Project, mr.IID)
}
