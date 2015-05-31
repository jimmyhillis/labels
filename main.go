package main

import (
    "os"
    "os/exec"
    "fmt"
    "log"
    "strings"
    "errors"
    "golang.org/x/oauth2"
    "github.com/google/go-github/github"
    "github.com/codegangsta/cli"
)

func main() {
    app := cli.NewApp()
    app.Name = "labels"
    app.Usage = "Mange Github labels from the command line"
    app.Action = func(c *cli.Context) {
        fmt.Println("Yeah okay")
    }
    app.Commands = []cli.Command{
        {
            Name:      "add",
            Aliases:   []string{"a"},
            Usage:     "add labels to a github repo",
            Action: func(c *cli.Context) {
                println("added task: ", c.Args().First())
                githubrepo := c.Args().First()
                owner, repo, err := ParseRepoArgument(githubrepo)
                if err != nil {
                    log.Fatal(err)
                    os.Exit(1)
                }
                AddLabels(owner, repo)
            },
        },
        {
            Name: "delete",
            Aliases: []string{"d"},
            Usage: "delete all labels from a github repo",
            Action: func(c *cli.Context) {
                println("delete: ", c.Args().First())
                githubrepo := c.Args().First()
                owner, repo, err := ParseRepoArgument(githubrepo)
                if err != nil {
                    log.Fatal(err)
                    os.Exit(1)
                }
                DeleteCurrentLabels(owner, repo)
            },
        },
    }
    app.Run(os.Args)
}

func ParseRepoArgument(repo string) (string, string, error) {
    config := strings.Split(repo, "/")
    fmt.Println(config)
    if len(config) != 2 {
        err := errors.New("Github repo not provided")
        return "", "", err
    }
    return config[0], config[1], nil
}

type tokenSource struct {
  token *oauth2.Token
}

func GithubToken() (string, error) {
    out, err := exec.Command("git", "config", "--global", "github.token").Output()
    if err != nil {
      return "", err
    }
    return string(out), nil
}

// add Token() method to satisfy oauth2.TokenSource interface
func (t *tokenSource) Token() (*oauth2.Token, error){
  return t.token, nil
}

func GithubClient() (*github.Client) {
    token, err := GithubToken()
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
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
    labels := []Label{{ Name: "Ready For Review", Color: "24c0eb" }}
    for _, label := range labels {
        gh_label := github.Label{Name: &label.Name, Color: &label.Color}
        client.Issues.CreateLabel(owner, repo, &gh_label)
        fmt.Println("Added '", *gh_label.Name, "' to", owner, repo)
    }
}

func CurrentRepos(owner string) {
    client := GithubClient()
    repos, _, err := client.Repositories.ListByOrg(owner, nil)
    if err != nil {
        fmt.Println(err)
        return
    }
    for _, repo := range repos {
        fmt.Println(*repo.Name)
    }
}
