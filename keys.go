package main

import (
  "bytes"
  "crypto/md5"
  "errors"
  "flag"
  "fmt"
  "io/ioutil"
  "os"
  "os/exec"
  "path/filepath"
  "strings"

  "code.google.com/p/go.crypto/ssh"
)

var sshPubKeyPath = flag.String("--key-path", "", "path to an ssh-key to upload")

func upload(k ssh.PublicKey) {
  keyStr := string(ssh.MarshalAuthorizedKey(k))
  fmt.Println(strings.TrimSpace(keyStr))
}

func fingerprint(k ssh.PublicKey) []byte {
  w := md5.New()
  w.Write(ssh.MarshalPublicKey(k))
  return w.Sum(nil)
}

func findSSHKeys() ([]byte, error) {
  if *sshPubKeyPath != "" {
    return sshReadPubKey(*sshPubKeyPath)
  }

  candidateKeys := make(map[string]string)

  // key from id_rsa.pub
  keyBytes, err := sshReadPubKey(filepath.Join(homePath(), ".ssh", "id_rsa.pub"))
  if err == nil {
    key, comment, _, _, ok := ssh.ParseAuthorizedKey(keyBytes)
    if ok {
      candidateKeys[string(ssh.MarshalPublicKey(key))] = comment
    }
  }

  // keys from ssh-add
  out, err := exec.Command("ssh-add", "-L").Output()
  sshAddKeys := strings.TrimSpace(string(out))
  if err == nil && sshAddKeys != "" {
    for _, k := range strings.Split(sshAddKeys, "\n") {
      key, comment, _, _, ok := ssh.ParseAuthorizedKey([]byte(k))
      if ok {
        candidateKeys[string(ssh.MarshalPublicKey(key))] = comment
      }
    }
  }

  numCandidateKeys := len(candidateKeys)
  if numCandidateKeys == 0 {
    return nil, errors.New("No ssh keys found.")
  }

  i := 1
  keyLst := make([]ssh.PublicKey, len(candidateKeys))
  for key, comment := range candidateKeys {
    pubKey, _, _ := ssh.ParsePublicKey([]byte(key))
    keyLst[i-1] = pubKey

    k := strings.TrimSpace(string(ssh.MarshalAuthorizedKey(pubKey)))
    fmt.Printf("[ %d ] %s...%s %s\n", i, k[0:24], k[len(k)-24:], comment)
    i += 1
  }

  choice := -1
  fmt.Printf("Pick a key to use [1-%d]: ", numCandidateKeys)
  fmt.Scanf("%d", &choice)
  if choice < 1 || choice > numCandidateKeys {
    return nil, errors.New("wat")
  }

  return ssh.MarshalAuthorizedKey(keyLst[choice-1]), nil
}

func sshReadPubKey(s string) ([]byte, error) {
  f, err := os.Open(filepath.FromSlash(s))
  if err != nil {
    return nil, err
  }

  key, err := ioutil.ReadAll(f)
  if err != nil {
    return nil, err
  }

  if bytes.Contains(key, []byte("PRIVATE")) {
    return nil, privKeyError(s)
  }

  return key, nil
}

type privKeyError string

func (e privKeyError) Error() string {
  return "appears to be a private key: " + string(e)
}
