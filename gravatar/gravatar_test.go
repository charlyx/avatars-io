package gravatar

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

type GravatarSuite struct {
	suite.Suite
}

func (g *GravatarSuite) Test_HandlerFuncError() {
	assert := g.Assert()

	resp := getResponse(HandlerFunc, "/")
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	assert.Equal(http.StatusBadRequest, resp.StatusCode)
	assert.Contains(string(body), "You must specify email query parameter.")
}

func (g *GravatarSuite) Test_HandlerNoSize() {
	assert := g.Assert()

	resp := getResponse(HandlerFunc, "/?email=mon@email")
	url, err := resp.Location()

	assert.NoError(err)
	assert.Equal(http.StatusFound, resp.StatusCode)

	if assert.NotNil(url) {
		const expectedURL = "https://www.gravatar.com/avatar/b7121d6cc0de7b0560723f352ea29cf8?s=80"
		assert.Equal(expectedURL, url.String())
	}
}

func (g *GravatarSuite) Test_HandlerSize() {
	assert := g.Assert()

	resp := getResponse(HandlerFunc, "/?email=mon@email&size=200")
	url, err := resp.Location()

	assert.NoError(err)
	assert.Equal(http.StatusFound, resp.StatusCode)

	if assert.NotNil(url) {
		const expectedURL = "https://www.gravatar.com/avatar/b7121d6cc0de7b0560723f352ea29cf8?s=200"
		assert.Equal(expectedURL, url.String())
	}
}

func getResponse(handlerFunc http.HandlerFunc, target string) *http.Response {
	req := httptest.NewRequest("GET", target, nil)
	w := httptest.NewRecorder()
	handlerFunc(w, req)

	return w.Result()
}

func TestServerSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(GravatarSuite))
}
