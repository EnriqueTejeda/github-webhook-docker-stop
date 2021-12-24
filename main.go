package main

import (
	"context"
	"fmt"
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
	host               string
	port               string
	dockerHost         string
	githubClientSecret string
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
func initConfig() *config {
	return &config{
		host:               getEnv("HOST", "0.0.0.0"),
		port:               getEnv("PORT", "8080"),
		dockerHost:         getEnv("DOCKER_HOST", "unix:///var/run/docker.sock"),
		githubClientSecret: getEnv("GITHUB_CLIENT_SECRET", "123456"),
	}
}
func init() {
	if err := godotenv.Load(); err != nil {
		log.Warn("no .env file found")
	}
}
func handleWebhook(w http.ResponseWriter, r *http.Request) {
	payload, err := github.ValidatePayload(r, []byte(initConfig().githubClientSecret))
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
			err := handlePullRequestClosedEvent(e)
			if err != nil {
				log.Error("error handling pull request closed event: err=%s\n", err)
			}
		}
	default:
		log.Error("unknown event type %s\n", github.WebHookType(r))
		return
	}
}

func handlePullRequestClosedEvent(GithubEvent *github.PullRequestEvent) error {
	labels := filters.NewArgs(
		filters.KeyValuePair{"label", fmt.Sprintf("projectName=%v", sanitize(GithubEvent.GetRepo().GetName()))},
		filters.KeyValuePair{"label", fmt.Sprintf("pullRequestNumber=%v", strconv.Itoa(GithubEvent.GetPullRequest().GetNumber()))},
	)
	_, err := stopContainer(labels)
	if err != nil {
		return err
	}
	return nil
}

func sanitize(repoName string) string {
	replacer := strings.NewReplacer(",", "!", "?", "/", "-", "")
	return strings.ToLower(replacer.Replace(repoName))
}

func findContainerByLabel(labels filters.Args) ([]types.Container, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation(), client.WithHost(initConfig().dockerHost))
	if err != nil {
		return nil, err
	}
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{Filters: labels})
	if err != nil {
		return nil, err
	}
	return containers, nil
}

func stopContainer(labels filters.Args) ([]types.Container, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation(), client.WithHost(initConfig().dockerHost))
	if err != nil {
		return nil, err
	}
	containers, err := findContainerByLabel(labels)
	if err != nil {
		return nil, err
	}
	for _, container := range containers {
		if err := cli.ContainerStop(ctx, container.ID, nil); err != nil {
			return nil, err
		}
		log.Infof("stopping container %v (%v)...", container.Names[0], container.ID[:10])
	}
	return containers, nil
}

func main() {
	log.Infof("server started on port %v:%v..", initConfig().host, initConfig().port)
	http.HandleFunc("/webhook", handleWebhook)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%v:%v", initConfig().host, initConfig().port), nil))
}
