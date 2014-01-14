package main

import (
  "encoding/json"
  "io/ioutil"
  "os"
  "os/user"
  "path/filepath"
  "time"
)

type PKConfig struct {
  URL         string    "json:`api_url`"
  AccessToken string    "json:`access_token`"
  Expiration  time.Time "json:`token_expiraion`"
}

func getRc() (rc *PKConfig, err error) {
  rc = new(PKConfig)
  rc.URL = "https://pancake.io/v1/"
  rcBytes, err := ioutil.ReadFile(rcPath())
  if err != nil {
    return
  }

  err = json.Unmarshal(rcBytes, rc)
  return rc, err
}

func (p *PKConfig) saveRc() error {
  rcBytes, err := json.Marshal(p)
  if err != nil {
    return err
  }

  return ioutil.WriteFile(rcPath(), rcBytes, os.ModePerm)
}

func rcPath() string {
  return filepath.Join(homePath(), ".pk")
}

func homePath() string {
  u, err := user.Current()
  if err != nil {
    panic("couldn't determine user: " + err.Error())
  }
  return u.HomeDir
}
