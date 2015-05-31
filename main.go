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
    app.Commands = []cli.Command{
        {
            Name: "add",
            Aliases: []string{"a"},
            Usage: "Adds labels to a Github repo",
            Action: func(c *cli.Context) {
                githubrepo := c.Args().First()
                owner, repo, err := ParseRepoArgument(githubrepo)
                if err != nil {
                    log.Fatal(err)
                    os.Exit(1)
                }
                labels, err := ParseLabelArguments(c.Args()[1:])
                if err != nil {
                    log.Fatal(err)
                    os.Exit(1)
                }
                AddLabels(owner, repo, labels)
            },
        },
        {
            Name: "delete",
            Aliases: []string{"d"},
            Usage: "Deletes labels from a github repo",
            Action: func(c *cli.Context) {
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

func ParseLabelArguments(labelopts []string) ([]Label, error) {
    labels := []Label{}
    if len(labelopts) == 0 {
        err := errors.New("No labels provided")
        return []Label{}, err
    }
    for _, label := range labelopts {
        parts := strings.Split(label, ":")
        if len(parts) != 2 {
            err := errors.New("Incorrect label format <name:color>")
            return []Label{}, err
        }
        labels = append(labels, Label{ Name: parts[0], Color: parts[1] })
    }
    return labels, nil
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

func AddLabels(owner string, repo string, labels []Label) {
    client := GithubClient()
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
