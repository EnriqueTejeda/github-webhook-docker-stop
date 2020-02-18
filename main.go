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
		log.Print("No .env file found")
	}
}
func handleWebhook(w http.ResponseWriter, r *http.Request) {
	payload, err := github.ValidatePayload(r, []byte(""))
	if err != nil {
		log.Printf("error validating request body: err=%s\n", err)
		return
	}
	defer r.Body.Close()

	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		log.Printf("could not parse webhook: err=%s\n", err)
		return
	}
	log.Printf("Received a webhook event!")

	switch e := event.(type) {
	case *github.PullRequestEvent:
		if e.GetAction() == "closed" {
			stopDocker(e)
		}
	default:
		log.Printf("unknown event type %s\n", github.WebHookType(r))
		return
	}
}
func replaceDotAndUppercase(repoName string) string {
	return strings.ToUpper(strings.Replace(repoName, ".", "_", -1))
}
func stopDocker(GithubEvent *github.PullRequestEvent) {
	var nameContainer = "preview-" + replaceDotAndUppercase(GithubEvent.GetRepo().GetName()) + "-PR-" + strconv.Itoa(GithubEvent.GetPullRequest().GetNumber())
	args := filters.NewArgs(filters.KeyValuePair{"name", string(nameContainer)})
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation(), client.WithHost("unix:///tmp/docker.sock"))
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{Filters: args})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		if err := cli.ContainerStop(ctx, container.ID, nil); err != nil {
			panic(err)
		}
		createComment(GithubEvent)
		log.Printf("Stopping container %v...", container.ID[:10])
	}

}
func main() {

	log.Println("Server started..")
	http.HandleFunc("/webhook", handleWebhook)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
