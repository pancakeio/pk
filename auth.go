package main

import (
  "fmt"

  "github.com/howeyc/gopass"
  "pk/api"
)

// Reads oauth tokens from disk if exist, else prompt for login
func authorize(force bool) {
  conf, err := getRc()
  client, err = api.NewPKClient(conf.URL)

  if force || err != nil {
    var username, password string
    fmt.Printf("Email: ")
    fmt.Scanln(&username)
    fmt.Printf("Password: ")
    password = string(gopass.GetPasswd())
    err = client.Authorize(username, password) // should return json on unauthorized
    if err != nil {
      fmt.Println("Authorization error:", err)
      return
    }

    conf.User = username
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
