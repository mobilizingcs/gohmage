package gohmage

import (
  "os"
  "testing"
  "fmt"
  "github.com/stretchr/testify/assert"
)

var ohmage_url, ohmage_username, ohmage_password string
const ohmage_client = "gohmage-test"

func TestMain( m *testing.M ) {
  ohmage_url = os.Getenv( "OHMAGE_URL" )
  ohmage_username = os.Getenv( "OHMAGE_USERNAME" )
  ohmage_password = os.Getenv( "OHMAGE_PASSWORD" )
  if ohmage_url == "" || ohmage_username == "" || ohmage_password == "" {
    fmt.Println( "Ohmage server environment variables not set. Cannot proceed with tests." )
    os.Exit( 1 )
  } else {
    os.Exit( m.Run( ) )
  }
}

func TestClient( t *testing.T ) {
  client := NewClient( ohmage_url, ohmage_client )
  assert.Equal( t, client.Ohmage_URL, ohmage_url, "URL should be set" )
  assert.Equal( t, client.Ohmage_Client, ohmage_client, "Client string should be set" )
}

func TestApi( t *testing.T ) {
  client := NewClient( ohmage_url, ohmage_client )
  t.Run( "UserAuthToken", client.testUserAuthToken )
  t.Run( "UserInfoRead", client.testUserInfoRead )
}

func ( client *Client ) testUserAuthToken( t *testing.T ) {
  assert.Equal( t, client.UserAuthToken( ohmage_username, ohmage_password ), true, "Auth should succeed" )
  assert.Equal( t, client.Is_Authenticated, true, "Is_Authenticated should be true" )
  assert.NotEqual( t, client.Auth_Token, "", "Auth_Token should not be empty" )
}

func ( client *Client ) testUserInfoRead( t *testing.T ) {
  response, err := client.UserInfoRead( )
  if err != nil {
    t.Error( "Failed to read JSON response" )
  }
  classes, err := response.GetObject( "data", ohmage_username, "classes" )
  if err != nil {
    t.Error( "Failed to find the classes attribute in JSON response" )
  }
  public_class, err := classes.GetString( "urn:class:public" )
  if err != nil {
    t.Error( "Failed to find the public class attribute in the JSON response" )
  }
  assert.Equal( t, public_class, "Public Class", "Public Class URN should match" )
}