package gohmage

import (
  "net/http"
  "net/url"
  "fmt"
  "io/ioutil"
  "errors"
  "github.com/antonholmquist/jason"
)

type Client struct {
  Ohmage_URL string
  Ohmage_Client string
  Auth_Token string
  Is_Authenticated bool
}

func NewClient( ohmage_url, client string ) *Client {
  // todo: check trailing slash remove it if present
  fmt.Println("")
  return &Client {
    Ohmage_URL: ohmage_url,
    Ohmage_Client: client,
    Auth_Token: "",
    Is_Authenticated: false,
  }
}

func ( client *Client ) UserAuthToken( username string, password string ) bool {
  parameters := url.Values{ }
  parameters.Set( "user", username )
  parameters.Set( "password", password )
  response, err := client.request( "post", "/user/auth_token", parameters )
  if( err != nil ) {
    // todo : handle call-specific errors: auth success or failure
    return false;
  }
  result, err := response.GetString( "result" )
  token, err := response.GetString( "token" )
  if result == "success" && token != "" {
    client.Auth_Token = token
    client.Is_Authenticated = true
    return true;
  }
  return false;
}

func ( client *Client ) request( method string, endpoint string, parameters url.Values ) (*jason.Object, error) {
  var resp *http.Response
  var err error
  parameters.Set( "client", client.Ohmage_Client )
  if method == "post" {
    resp, err = http.PostForm( client.Ohmage_URL + endpoint, parameters )
  } else {
    resp, err = http.Get( client.Ohmage_URL + endpoint )
  }
  defer resp.Body.Close( )
  if( err != nil ) {
    return nil, errors.New( "An HTTP protocol error occurred." )
  }
  body, err := ioutil.ReadAll( resp.Body )
  if( err != nil ) {
    return nil, errors.New( "An error occurred while reading the HTTP response" )
  }
  response, err := jason.NewObjectFromBytes( body )
  if( err != nil ) {
    return nil, errors.New( "An error occurred while parsing the response as JSON" )
  }
  return response, nil
}