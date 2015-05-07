package main

import (
    "fmt"
    "golang.org/x/oauth2"
    "github.com/google/go-github/github"
)

func main() {
    fmt.Println("Updating labels")
    DeleteCurrentLabels("XXX", "XXX")
    AddLabels("XXX", "XXX")
}

type tokenSource struct {
  token *oauth2.Token
}

// add Token() method to satisfy oauth2.TokenSource interface
func (t *tokenSource) Token() (*oauth2.Token, error){
  return t.token, nil
}

func GithubClient() (*github.Client) {
    token := "XXX"
    ts := &tokenSource{
      &oauth2.Token{ AccessToken: token} ,
    }
    tc := oauth2.NewClient(oauth2.NoContext, ts)
    client := github.NewClient(tc)
    return client
}

func CurrentLabels(owner string, repo string) []github.Label {
    client := GithubClient()
    labels, _, err := client.Issues.ListLabels(owner, repo, nil)
    if err != nil {
        fmt.Println(err)
        return labels
    }
    for _, label := range labels {
        fmt.Println(*label.Name, *label.Color)
    }
    return labels
}

func DeleteCurrentLabels(owner string, repo string) {
    client := GithubClient()
    labels, _, err := client.Issues.ListLabels(owner, repo, nil)
    if err != nil {
        fmt.Println(err)
        return
    }
    for _, label := range labels {
        fmt.Println("Deleting", *label.Name, "from", owner, repo)
        client.Issues.DeleteLabel(owner, repo, *label.Name)
    }
}

type Label struct {
    Name string
    Color string
}

func AddLabels(owner string, repo string) {
    client := GithubClient()
    labels := []Label{{ Name: "XXX", Color: "XXX" },
                      { Name: "XXX", Color: "XXX" }}
    for _, label := range labels {
        gh_label := github.Label{Name: &label.Name, Color: &label.Color}
        client.Issues.CreateLabel(owner, repo, &gh_label)
    }
}

func CurrentRepos() {
    client := GithubClient()
    repos, _, err := client.Repositories.ListByOrg("meerkats", nil)
    if err != nil {
        fmt.Println(err)
        return
    }
    for _, repo := range repos {
        fmt.Println(*repo.Name)
    }
}
