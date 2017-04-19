package gohmage

import (
  "net/http"
  "net/url"
  "io/ioutil"
  "github.com/pkg/errors"
  "github.com/antonholmquist/jason"
)

type Client struct {
  Ohmage_URL string
  Ohmage_Client string
  Auth_Token string
  Is_Authenticated bool
  Username string
}

func NewClient( ohmage_url, client string ) *Client {
  // todo: check trailing slash remove it if present
  return &Client {
    Ohmage_URL: ohmage_url,
    Ohmage_Client: client,
    Auth_Token: "",
    Is_Authenticated: false,
    Username: "",
  }
}

func ( client *Client ) UserAuthToken( username string, password string ) ( bool, error ) {
  parameters := url.Values{ }
  parameters.Set( "user", username )
  parameters.Set( "password", password )
  response, err := client.request( "post", "/user/auth_token", parameters, false )
  if err != nil {
    return false, errors.Wrap( err, "Ohmage API authorization call failed" );
  }
  result, err := response.GetString( "result" )
  token, err := response.GetString( "token" )
  if result == "success" && token != "" {
    client.Auth_Token = token
    client.Is_Authenticated = true
    client.Username = username
    return true, nil;
  }
  return false, errors.New( "User authentication failed" )
}

func ( client *Client ) UserInfoRead( ) ( *jason.Object, error ) {
  response, err := client.request( "post", "/user_info/read", url.Values{ }, true )
  if err != nil {
    return nil, errors.Wrap( err, "Ohmage API error occurred" )
  }
  return response, nil
}

func ( client *Client ) request(  method string,
                                  endpoint string,
                                  parameters url.Values,
                                  is_protected_endpoint bool ) ( *jason.Object, error ) {
  var resp *http.Response
  var err error
  if is_protected_endpoint {
    if client.Is_Authenticated {
      parameters.Set( "auth_token", client.Auth_Token )
    } else {
      return nil, errors.New( "Authorization token not found" )
    }
  }
  parameters.Set( "client", client.Ohmage_Client )
  if method == "post" {
    resp, err = http.PostForm( client.Ohmage_URL + endpoint, parameters )
  } else {
    resp, err = http.Get( client.Ohmage_URL + endpoint )
  }
  defer resp.Body.Close( )
  if( err != nil ) {
    return nil, errors.New( "An HTTP protocol error occurred" )
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