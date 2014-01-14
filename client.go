package main

import (
  "code.google.com/p/go.crypto/ssh"
  "code.google.com/p/goauth2/oauth"
  "flag"
  "fmt"
  "net/http"
  "net/url"
  "pk/api"
  "strings"
)

var commands = map[string]func() error{
  "add-key": func() error {
    s, _ := findSSHKeys()
    key, _, _, _, ok := ssh.ParseAuthorizedKey(s)
    if !ok {
      return fmt.Errorf("refusing to upload")
    }

    keyStr := string(ssh.MarshalAuthorizedKey(key))
    err := client.UploadKey("myfirstkey", strings.TrimSpace(keyStr))
    return err

  },
  "list-keys": func() error {
    resp, err := client.ListKeys()
    if err != nil {
      return err
    }

    for _, key := range resp.Keys {
      fmt.Printf("%-20s  %s  %s\n", key.Name, key.Fingerprint, key.Preview)
    }

    return nil
  },
  "remove-key": func() error {
    resp, err := client.ListKeys()
    if err != nil {
      return err
    }

    // choose key to remove
    return nil
  },
  "create-project": func() error {
    resp, err := client.CreateProject()
    if err != nil {
      return err
    }
    fmt.Printf("Created new project %s.\n", resp.Name)
    return nil
  },
  "list-projects": func() error {
    resp, err := client.ListProjects()
    if err != nil {
      return err
    }

    for _, key := range resp.Projects {
      fmt.Printf("%-28s  %s  %s.git\n", key.Name, key.PancakeURL, key.RepoName)
    }

    return nil
  },
  "delete-project": func() error {
    return nil
  },
}

var client *api.PKClient

func main() {

  var w = flag.Bool("w", false, "prints list of commands")

  flag.Usage = func() {
    fmt.Println("Commands: ")
    for commandName, _ := range commands {
      fmt.Println(commandName)
    }
    fmt.Println()
    flag.PrintDefaults()
  }

  flag.Parse()
  if *w {
    for commandName, _ := range commands {
      fmt.Printf("%s ", commandName)
    }
    fmt.Println()
    return
  }

  var err error
  authorize(false)

  // save access token
  if flag.NArg() == 0 {
    flag.PrintDefaults()
    return
  }

  command, ok := commands[flag.Arg(0)]
  if !ok {
    flag.PrintDefaults()
    return
  }

  err = tryWithReauth(command)
  if err != nil {
    fmt.Printf("%s error: %s\n", flag.Arg(0), err)
  }
}

func authorize(force bool) {
  conf, err := getRc()
  client, err = api.NewPKClient(conf.URL)

  if force || err != nil {
    var username, password string
    fmt.Printf("Username: ")
    fmt.Scanln(&username)
    fmt.Printf("Password: ")
    fmt.Scanln(&password)                      // ugh shows password
    err = client.Authorize(username, password) // should return json on unauthorized
    if err != nil {
      fmt.Println("Authorization error:", err)
      return
    }

    conf.AccessToken = client.AccessToken
    conf.Expiration = client.TokenExpiration

    if err := conf.saveRc(); err != nil {
      panic(err)
    }
  } else {
    client.AccessToken = conf.AccessToken
    client.TokenExpiration = conf.Expiration
  }
  return
}

func tryWithReauth(f func() error) error {
  err := f()
  switch err.(type) {
  case *url.Error:
    switch err.(*url.Error).Err.(type) {
    case oauth.OAuthError:
      fmt.Println("Access token has expired; please log in again.")
      authorize(true)
      return f()
    }
  case *api.APIError:
    if err.(*api.APIError).Code == http.StatusUnauthorized {
      fmt.Println("Bad access token; please log in again.")
      authorize(true)
      return f()
    }
  }
  return err
}
