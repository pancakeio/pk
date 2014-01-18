package api

import (
  "encoding/json"
  "io/ioutil"
  "net/http"
  "net/url"
)

func (pk *PKClient) UploadKey(name, key string) error {
  resp, err := pk.postForm("keys", url.Values{
    "name": {name},
    "key":  {key},
  })
  if err != nil {
    return err
  }

  r, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return err
  }

  if resp.StatusCode != http.StatusOK {
    return &APIError{resp.StatusCode, string(r)}
  }

  return nil
}

type ListKeysResponse struct {
  Keys []struct {
    Name        string `json:"name"`
    Fingerprint string `json:"fingerprint"`
    Preview     string `json:"preview"`
  } `json:"keys"`
}

func (pk *PKClient) ListKeys() (*ListKeysResponse, error) {
  resp, err := pk.get("keys")
  if err != nil {
    return nil, err
  }

  r, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return nil, err
  }
  if resp.StatusCode != http.StatusOK {
    return nil, &APIError{resp.StatusCode, string(r)}
  }

  out := new(ListKeysResponse)
  json.Unmarshal(r, out)
  return out, nil
}

func (pk *PKClient) DeleteKey(fingerprint string) (bool, error) {
  resp, err := pk.delete("keys", url.Values{"fingerprint": {fingerprint}})
  if err != nil {
    return false, err
  }

  if resp.StatusCode != http.StatusOK {
    r, err := ioutil.ReadAll(resp.Body)
    if err != nil {
      return false, err
    }

    return false, &APIError{resp.StatusCode, string(r)}
  }

  return true, nil
}
