package main

import (
  "flag"
  "fmt"
  "strings"

  "code.google.com/p/go.crypto/ssh"
)

var cmdKeyAdd = &cmd{
  name: "add-key",
  run: func() error {
    s, _ := findSSHKeys()
    key, _, _, _, ok := ssh.ParseAuthorizedKey(s)
    if !ok {
      return fmt.Errorf("refusing to upload")
    }

    keyStr := string(ssh.MarshalAuthorizedKey(key))
    err := client.UploadKey("myfirstkey", strings.TrimSpace(keyStr))
    return err
  },
  flags: flag.NewFlagSet("add-key", flag.ExitOnError),
  usage: func() string {
    return "add an ssh key to your account"
  },
}

func init() {
  cmdKeyAdd.flags.StringVar(&sshPubKeyPath, "key-path", "", "path to an ssh-key to upload")
}

var cmdKeysList = &cmd{
  name: "list-keys",
  run: func() error {
    resp, err := client.ListKeys()
    if err != nil {
      return err
    }

    for _, key := range resp.Keys {
      fmt.Printf("%-20s  %s  %s\n", key.Name, key.Fingerprint, key.Preview)
    }

    return nil
  },
  usage: func() string {
    return "list added ssh keys"
  },
}

var cmdKeyRemove = &cmd{
  name: "remove-key",
  run: func() error {
    resp, err := client.ListKeys()
    if err != nil {
      return err
    }

    for i, key := range resp.Keys {
      fmt.Printf("[ %d ] %-20s %s\n", i+1, key.Name, key.Preview)
    }

    choice, err := pick("key", len(resp.Keys))
    if err != nil {
      return err
    }

    chosenKey := resp.Keys[choice]
    _, err = client.DeleteKey(chosenKey.Fingerprint)
    if err != nil {
      return err
    }

    fmt.Println("Removed key", chosenKey.Preview)
    return nil
  },
  usage: func() string {
    return "diassocate an ssh key"
  },
}
