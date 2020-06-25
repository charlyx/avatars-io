package usage

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
)

type UsageSuite struct {
	suite.Suite
}

func (s *UsageSuite) Test_HandlerFunc() {
	assert := s.Assert()
	const twitterExampleLink = "https://avatars.charlyx.dev/twitter?username=charlyx"
	const gravatarExampleLink = "https://avatars.charlyx.dev/gravatar?email=mon@email"

	assert.HTTPStatusCode(HandlerFunc, "GET", "/", nil, http.StatusNotFound)
	assert.HTTPBodyContains(HandlerFunc, "GET", "/", nil, twitterExampleLink)
	assert.HTTPBodyContains(HandlerFunc, "GET", "/", nil, gravatarExampleLink)
}

func TestUsageSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(UsageSuite))
}
