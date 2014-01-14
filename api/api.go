package api

import (
  "code.google.com/p/goauth2/oauth"
  "encoding/json"
  "fmt"
  "io/ioutil"
  "net/http"
  "net/url"
  "strings"
  "time"
)

type PKClient struct {
  BaseURL         *url.URL
  client          *http.Client
  ClientId        string
  AccessToken     string
  TokenExpiration time.Time
}

type oauthResponse struct {
  AccessToken      string `json:"access_token"`
  TokenType        string `json:"token_type"`
  ExpiresInSeconds int    `json:"expires_in"`

  Error            string `json:"error"`
  ErrorDescription string `json:"error_description"`
}

func NewPKClient(pancakeURL string) (*PKClient, error) {
  u, err := url.Parse(pancakeURL)
  return &PKClient{BaseURL: u}, err
}

// refresh?
func (pk *PKClient) Authorize(username, password string) (err error) {
  now := time.Now()
  u := *pk.BaseURL // copy
  u.Path = "/oauth/token"
  v := url.Values{
    "grant_type": {"password"},
    "username":   {username},
    "password":   {password},
    "client_id":  {pk.ClientId},
  }

  req, err := http.NewRequest("POST", u.String(), strings.NewReader(v.Encode()))
  if err != nil {
    return
  }

  req.Header.Set("Accept", "application/json")
  req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
  resp, err := http.DefaultClient.Do(req)
  if err != nil {
    return
  }

  if resp.StatusCode == http.StatusUnauthorized {
    return fmt.Errorf("Bad login.")
  }

  r, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return
  }

  oauthResp := new(oauthResponse)
  err = json.Unmarshal(r, oauthResp)
  if err != nil {
    return
  }

  if oauthResp.Error != "" {
    return fmt.Errorf("Error: %s", oauthResp.ErrorDescription)
  }
  if oauthResp.AccessToken == "" {
    return fmt.Errorf("Error: no access token returned.")
  }

  pk.AccessToken = oauthResp.AccessToken
  pk.TokenExpiration = now.Add(time.Duration(oauthResp.ExpiresInSeconds) * time.Second)
  pk.client = nil // clear client in case we have cached credentials

  return nil
}

func (p *PKClient) Client() *http.Client {
  if p.client == nil {
    transport := &oauth.Transport{
      Config: &oauth.Config{ClientId: p.ClientId},
      Token: &oauth.Token{
        AccessToken: p.AccessToken,
        Expiry:      p.TokenExpiration,
      },
    }
    p.client = transport.Client()
  }
  return p.client
}

func (p *PKClient) postForm(apiPath string, val url.Values) (*http.Response, error) {
  return p.Client().PostForm(p.BaseURL.String()+apiPath, val)
}

func (p *PKClient) get(apiPath string) (*http.Response, error) {
  return p.Client().Get(p.BaseURL.String() + apiPath)
}

func (p *PKClient) delete(apiPath string, val url.Values) (*http.Response, error) {
  req, err := http.NewRequest("DELETE", p.BaseURL.String()+apiPath, strings.NewReader(val.Encode()))
  if err != nil {
    return nil, err
  }
  req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

  return p.Client().Do(req)
}

type APIError struct {
  Code    int
  Message string
}

func (a *APIError) Error() string {
  return fmt.Sprintf("error %d: %s", a.Code, a.Message)
}
