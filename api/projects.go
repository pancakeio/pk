package api

import (
  "encoding/json"
  "io/ioutil"
  "net/http"
  "net/url"
)

type CreateProjectResponse struct {
  Name string `json:"name"`
}

type ProjectKind string

const (
  STATIC_PROJECT  ProjectKind = "static"
  DEFAULT_PROJECT ProjectKind = "default"
  DROPBOX_PROJECT ProjectKind = "site"
)

func (pk *PKClient) CreateProject(kind ProjectKind) (*CreateProjectResponse, error) {
  resp, err := pk.postForm("projects", url.Values{"kind": {string(kind)}})
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

  out := new(CreateProjectResponse)
  json.Unmarshal([]byte(r), out)
  return out, nil
}

type ListProjectsResponse struct {
  Projects []struct {
    Name       string      `json:"name"`
    Kind       ProjectKind `json:"kind"`
    Subdomain  string      `json:"subdomain"`
    PancakeURL string      `json:"pancake_url"`
    RepoName   string      `json:"repo_name"`
  } `json:"projects"`
}

func (pk *PKClient) ListProjects() (*ListProjectsResponse, error) {
  resp, err := pk.get("projects")
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

  out := new(ListProjectsResponse)
  json.Unmarshal([]byte(r), out)
  return out, nil
}

func (pk *PKClient) DeleteProject(subdomain string) (bool, error) {
  resp, err := pk.delete("projects", url.Values{"subdomain": {subdomain}})
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
