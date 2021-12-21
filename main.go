package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/google/go-github/github"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type config struct {
	dockerHost  string
	githubToken string
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
func initConfig() *config {
	return &config{
		dockerHost:  getEnv("DOCKER_HOST", ""),
		githubToken: getEnv("GITHUB_TOKEN", ""),
	}
}
func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("no .env file found")
	}
}
func handleWebhook(w http.ResponseWriter, r *http.Request) {
	payload, err := github.ValidatePayload(r, []byte(""))
	if err != nil {
		log.Error("error validating request body: err=%s\n", err)
		return
	}
	defer r.Body.Close()

	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		log.Error("could not parse webhook: err=%s\n", err)
		return
	}
	switch e := event.(type) {
	case *github.PullRequestEvent:
		if e.GetAction() == "closed" {
			handlePullRequestClosedEvent(e)
		}
	default:
		log.Error("unknown event type %s\n", github.WebHookType(r))
		return
	}
}

func handlePullRequestClosedEvent(GithubEvent *github.PullRequestEvent) {
	labels := filters.NewArgs(
		filters.KeyValuePair{"projectName", replaceDotAndUppercase(GithubEvent.GetRepo().GetName())},
		filters.KeyValuePair{"pullRequestNumber", strconv.Itoa(GithubEvent.GetPullRequest().GetNumber())}
	)
	err := stopContainer(labels)
	if _, err != nil {
		return err
	}
}

func replaceDotAndLowercase(repoName string) string {
	return strings.ToLower(strings.Replace(repoName, ".", "_", -1))
}

func createPullRequestComment(GithubEvent *github.PullRequestEvent) (error) {
	comment := "The container linked to this pull request has been stopped (%v)", GithubEvent.GetPullRequest().GetHTMLURL() + "."
	if err := cli.Activity.CreateComment(ctx, GithubEvent.GetRepo().GetOwner().GetLogin(), GithubEvent.GetRepo().GetName(), GithubEvent.GetNumber(), &github.IssueComment{Body: &comment}); 
	err != nil {
		log.Error("error creating comment: err=%s\n", err)
		return err
	}
	return nil
}

func findContainerByLabel(labels filters.Args) ([]types.Container, error) ) { 
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation(), client.WithHost("unix:///tmp/docker.sock"))
	if err != nil {
		return _, err
	}
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{Filters: args})
	if err != nil {
		return _, err
	}
	log.Debug("Containers: %v", containers)
	return containers, nil
}

func stopContainer(labels filters.Args) (listContainersStopped, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation(), client.WithHost("unix:///tmp/docker.sock"))
	if err != nil {
		 return _, err
	}
	containers, err := findContainerByLabel(labels)
	if err != nil {
		return _, err
	}
	for _, container := range containers {
		if err := cli.ContainerStop(ctx, container.ID, nil); err != nil {
			return _, err
		}
		log.Debug("Stopping container %v...", container.ID[:10])
	}
	return listContainersStopped, nil
}

func main() {
	log.Debug("Server started..")
	http.HandleFunc("/webhook", handleWebhook)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
