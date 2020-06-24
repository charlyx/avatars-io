package server

import (
	"errors"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ServerSuite struct {
	suite.Suite
}

func (s *ServerSuite) Test_ServerError() {
	assert := s.Assert()
	os.Setenv("PROJECT_ID", "")
	secretAccessorMockErr := &secretAccessorMock{
		MockGet: func(key string) (string, error) {
			return "", errors.New(key)
		},
	}

	server, err := New(secretAccessorMockErr)

	assert.EqualError(err, "could not get twitter token: TWITTER_BEARER_TOKEN")
	assert.Nil(server)
}

func (s *ServerSuite) TestUsage() {
	assert := s.Assert()
	os.Setenv("PROJECT_ID", "avatars-io")
	secretAccessorMock := &secretAccessorMock{
		MockGet: func(key string) (string, error) {
			return key, nil
		},
	}
	const twitterExampleLink = "https://avatars.charlyx.dev/twitter?username=charlyx"
	const gravatarExampleLink = "https://avatars.charlyx.dev/gravatar?email=mon@email"

	server, err := New(secretAccessorMock)

	assert.NoError(err)
	assert.Implements(new(http.Handler), server)
	assert.HTTPStatusCode(server.ServeHTTP, "GET", "/", nil, http.StatusNotFound)
	assert.HTTPBodyContains(server.ServeHTTP, "GET", "/", nil, twitterExampleLink)
	assert.HTTPBodyContains(server.ServeHTTP, "GET", "/", nil, gravatarExampleLink)
}

func TestServerSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ServerSuite))
}

type secretAccessorMock struct {
	MockGet func(key string) (string, error)
}

func (s *secretAccessorMock) Get(key string) (string, error) {
	return s.MockGet(key)
}
