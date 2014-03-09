package main

import (
  "bytes"
  "crypto/md5"
  "errors"
  "fmt"
  "io"
  "io/ioutil"
  "os"
  "os/exec"
  "path/filepath"
  "strings"

  "code.google.com/p/go.crypto/ssh"
)

var (
  idRsaPubPath = filepath.Join(homePath(), ".ssh", "id_rsa.pub")
  idDsaPubPath = filepath.Join(homePath(), ".ssh", "id_dsa.pub")
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
func getSSHKeys() map[string]string {
  candidateKeys := make(map[string]string)

  // get key from id_rsa.pub
  rsaKey, rsaComment, rsaErr := sshReadPubKey(idRsaPubPath)
  if rsaErr == nil {
    candidateKeys[string(ssh.MarshalPublicKey(rsaKey))] = rsaComment
  }

  // get key from id_dsa.pub
  dsaKey, dsaComment, dsaErr := sshReadPubKey(idDsaPubPath)
  if dsaErr == nil {
    candidateKeys[string(ssh.MarshalPublicKey(dsaKey))] = dsaComment
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

  return candidateKeys
}

func pickSSHKey(candidateKeys map[string]string, w io.Writer) (ssh.PublicKey, error) {
  i := 0
  keyLst := make([]ssh.PublicKey, len(candidateKeys))
  for key, comment := range candidateKeys {
    pubKey, _, ok := ssh.ParsePublicKey([]byte(key))
    if !ok {
      continue
    }
    keyLst[i] = pubKey

    k := strings.TrimSpace(string(ssh.MarshalAuthorizedKey(pubKey)))
    l := fmt.Sprintf("[ %d ] %s...%s %s\n", i+1, k[0:24], k[len(k)-24:], comment)
    w.Write([]byte(l))
    i += 1
  }

  if i == 0 {
    return nil, errors.New("No ssh keys found.")
  }

  choice, err := pick("key", i)
  if err != nil {
    return nil, err
  }

  return keyLst[choice], nil
}

func createSSHKey() {
  cmd := exec.Command(
    "ssh-keygen",
    "-t", "rsa",
    "-f", strings.TrimSuffix(idDsaPubPath, ".pub"),
    "-C", "created by pancake.io",
  )
  cmd.Stdin = os.Stdin
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  cmd.Run()
}
