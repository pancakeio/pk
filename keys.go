package main

import (
  "bytes"
  "crypto/md5"
  "errors"
  "fmt"
  "io/ioutil"
  "os"
  "os/exec"
  "path/filepath"
  "strings"

  "code.google.com/p/go.crypto/ssh"
)

var errNotKey = errors.New("not a key")

type errPrivKey string

func (e errPrivKey) Error() string {
  return "appears to be a private key: " + string(e)
}

// Return SSH key md5 fingerprint
func fingerprint(k ssh.PublicKey) []byte {
  w := md5.New()
  w.Write(ssh.MarshalPublicKey(k))
  return w.Sum(nil)
}

// Read SSH public key bytes from path
func sshReadPubKey(path string) (ssh.PublicKey, string, error) {
  f, err := os.Open(filepath.FromSlash(path))
  if err != nil {
    return nil, "", err
  }

  keyBytes, err := ioutil.ReadAll(f)
  if err != nil {
    return nil, "", err
  }

  if bytes.Contains(keyBytes, []byte("PRIVATE")) {
    return nil, "", errPrivKey(path)
  }

  key, comment, _, _, ok := ssh.ParseAuthorizedKey(keyBytes)
  if !ok {
    return nil, "", errNotKey
  }

  return key, comment, nil
}

// Find SSH keys on the local file system
func findSSHKeys() ([]byte, error) {
  if argSSHPubKeyPath != "" {
    key, _, err := sshReadPubKey(argSSHPubKeyPath)
    return ssh.MarshalPublicKey(key), err
  }

  candidateKeys := make(map[string]string)

  // get key from id_rsa.pub
  key, comment, err := sshReadPubKey(filepath.Join(homePath(), ".ssh", "id_rsa.pub"))
  if err == nil {
    candidateKeys[string(ssh.MarshalPublicKey(key))] = comment
  }

  // get keys from ssh-add
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

  if len(candidateKeys) == 0 {
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

  choice, err := pick("key", len(candidateKeys))
  if err != nil {
    return nil, err
  }

  return ssh.MarshalAuthorizedKey(keyLst[choice]), nil
}
