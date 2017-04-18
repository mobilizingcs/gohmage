package gohmage

import (
  "os"
  "testing"
  "github.com/stretchr/testify/assert"
)

var ohmage_url, ohmage_username, ohmage_password string
const ohmage_client = "gohmage-test"

func TestMain(m *testing.M) {
  ohmage_url = os.Getenv( "OHMAGE_URL" )
  ohmage_username = os.Getenv( "OHMAGE_USERNAME" )
  ohmage_password = os.Getenv( "OHMAGE_PASSWORD" )
}

func TestClient( t *testing.T ) {
  client := NewClient( ohmage_url, ohmage_client )
  assert.Equal( t, client.Ohmage_URL, ohmage_url, "URL should be set" )
  assert.Equal( t, client.Ohmage_Client, ohmage_client, "Client string should be set" )
}

func TestApi( t *testing.T ) {
  client := NewClient( ohmage_url, ohmage_client )
  t.Run( "UserAuthToken", client.testUserAuthToken )
}

func (client *Client) testUserAuthToken( t *testing.T ) {
  assert.Equal( t, client.UserAuthToken( ohmage_username, ohmage_password ), true, "Auth should succeed" )
  assert.Equal( t, client.Is_Authenticated, true, "Is_Authenticated should be true" )
  assert.NotEqual( t, client.Auth_Token, "", "Auth_Token should not be empty" )
}

