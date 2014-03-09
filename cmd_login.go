package main

import (
  "fmt"
  "net/http"

  "pk/api"
)

var cmdLogin = &cmd{
  name: "login",
  run: func() error {
    // test auth, otherwise it will reauth
    didAuth := false
    _, err := client.ListProjects()
    if err != nil {
      switch err.(type) {
      case *api.APIError:
        if err.(*api.APIError).Code == http.StatusUnauthorized {
          authorize(true)
          didAuth = true
        } else {
          return err
        }
      default:
        return err
      }
    }

    conf, err := getRc()
    if err != nil {
      return err
    }

    fmt.Println("Logged in as", white(conf.User))
    if !didAuth && shouldContinue("Log in as different user?") {
      authorize(true)
    }
    return nil
  },
  usage: func() string {
    return "log in to pancake"
  },
}
