package main

import (
  "flag"
  "fmt"

  "pk/api"
)

var argProjectCreateDropbox bool
var cmdProjectCreate = &cmd{
  name: "create-project",
  run: func() error {
    kind := api.STATIC_PROJECT
    if argProjectCreateDropbox {
      kind = api.DROPBOX_PROJECT
    }

    resp, err := client.CreateProject(kind)
    if err != nil {
      return err
    }
    fmt.Printf("Created new project %s.\n", white(resp.Name))
    fmt.Println("Run the following in your project repo to get started:")
    fmt.Println()
    fmt.Printf("  git remote add pk git@build.pancake.io:%s.git\n", resp.Name)
    fmt.Println()
    fmt.Println("To deploy:")
    fmt.Println()
    fmt.Println("  git push pk master")
    fmt.Println()
    return nil
  },
  flags: flag.NewFlagSet("create-project", flag.ExitOnError),
  usage: func() string {
    return "create new pancake.io project"
  },
}

func init() {
  cmdProjectCreate.flags.BoolVar(&argProjectCreateDropbox, "dropbox", false, "creates a classic dropbox-based Pancake project (default: git-based project)")
}

var cmdProjectsList = &cmd{
  name: "list-projects",
  run: func() error {
    resp, err := client.ListProjects()
    if err != nil {
      return err
    }

    for _, key := range resp.Projects {
      if key.RepoName != "" {
        fmt.Printf("%-28s  %s  %s.git\n", key.Name, key.PancakeURL, key.RepoName)
      } else {
        fmt.Printf("%-28s  %s\n", key.Name, key.PancakeURL)
      }
    }

    return nil
  },
  usage: func() string {
    return "list your projects"
  },
}

var cmdProjectDelete = &cmd{
  name: "delete-project",
  run: func() error {
    resp, err := client.ListProjects()
    if err != nil {
      return err
    }

    var defaultProject int

    for i, project := range resp.Projects {
      if project.Kind == api.DEFAULT_PROJECT {
        defaultProject = i
        break
      }
    }

    p := resp.Projects
    p[defaultProject] = p[len(p)-1]
    p = p[0 : len(p)-1]

    for i, key := range p {
      fmt.Printf("[ %2d ] %-28s  %s\n", i+1, key.Name, key.PancakeURL)
    }
    choice, err := pick("project to delete", len(p))
    if err != nil {
      return err
    }

    _, err = client.DeleteProject(p[choice].Subdomain)
    if err != nil {
      return err
    }

    fmt.Println("Deleted project", p[choice].Name)
    return nil
  },
  usage: func() string {
    return "delete a project"
  },
}
