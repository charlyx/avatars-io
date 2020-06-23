package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/charlyx/avatars.io/twitter"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/suite"
)

type ServerSuite struct {
	suite.Suite
	Server *httptest.Server
	Client *http.Client
}

func (s *ServerSuite) SetupSuite() {
	httpmock.Activate()

	s.Server = httptest.NewServer(New("TWITTER_BEARER_TOKEN", nil))
	s.Client = &http.Client{
		Transport: &http.Transport{TLSHandshakeTimeout: 60 * time.Second},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

func (s *ServerSuite) SetupTest() {
	httpmock.Reset()
}

func (s *ServerSuite) TearDownSuite() {
	httpmock.DeactivateAndReset()
	s.Server.Close()
}

func (s *ServerSuite) TestUsage() {
	assert := s.Assert()

	res, err := s.Client.Get(s.Server.URL)
	assert.NoError(err)
	defer res.Body.Close()

	assert.Equal(http.StatusNotFound, res.StatusCode, "http status should be not found")

	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(err)

	twitterExampleLink := "https://avatars.charlyx.dev/twitter?username=charlyx"
	assert.Contains(string(body), twitterExampleLink)

	gravatarExampleLink := "https://avatars.charlyx.dev/gravatar?email=mon@email"
	assert.Contains(string(body), gravatarExampleLink)
}

func (s *ServerSuite) Test_GravatarUsage() {
	assert := s.Assert()

	gravatarURL := fmt.Sprintf("%s/gravatar", s.Server.URL)
	res, err := s.Client.Get(gravatarURL)
	assert.NoError(err)
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(err)

	url, err := res.Location()
	assert.EqualError(err, http.ErrNoLocation.Error())
	assert.Nil(url)

	assert.Equal(http.StatusBadRequest, res.StatusCode)
	assert.Contains(string(body), "You must specify email query parameter.")
}

func (s *ServerSuite) Test_GravatarDefault() {
	assert := s.Assert()

	gravatarURL := fmt.Sprintf("%s/gravatar?email=mon@email", s.Server.URL)
	res, err := s.Client.Get(gravatarURL)
	assert.NoError(err)
	defer res.Body.Close()

	assert.Equal(http.StatusFound, res.StatusCode)

	url, err := res.Location()
	assert.NoError(err)

	if assert.NotNil(url) {
		assert.Equal("https://www.gravatar.com/avatar/b7121d6cc0de7b0560723f352ea29cf8?s=80", url.String())
	}
}

func (s *ServerSuite) Test_GravatarSize() {
	assert := s.Assert()

	gravatarURL := fmt.Sprintf("%s/gravatar?email=mon@email&size=%d", s.Server.URL, 200)
	res, err := s.Client.Get(gravatarURL)
	assert.NoError(err)
	defer res.Body.Close()

	assert.Equal(http.StatusFound, res.StatusCode)

	url, err := res.Location()
	assert.NoError(err)

	if assert.NotNil(url) {
		assert.Equal("https://www.gravatar.com/avatar/b7121d6cc0de7b0560723f352ea29cf8?s=200", url.String())
	}
}

func (s *ServerSuite) Test_GravatarSizeShort() {
	assert := s.Assert()

	gravatarURL := fmt.Sprintf("%s/gravatar?email=mon@email&s=%d", s.Server.URL, 200)
	res, err := s.Client.Get(gravatarURL)
	assert.NoError(err)
	defer res.Body.Close()

	assert.Equal(http.StatusFound, res.StatusCode)

	url, err := res.Location()
	assert.NoError(err)
	assert.NotNil(url)

	if assert.NotNil(url) {
		assert.Equal("https://www.gravatar.com/avatar/b7121d6cc0de7b0560723f352ea29cf8?s=200", url.String())
	}
}

func (s *ServerSuite) Test_TwitterUsage() {
	assert := s.Assert()

	twitterURL := fmt.Sprintf("%s/twitter", s.Server.URL)
	res, err := s.Client.Get(twitterURL)
	assert.NoError(err)
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(err)

	url, err := res.Location()
	assert.EqualError(err, http.ErrNoLocation.Error())
	assert.Nil(url)

	assert.Equal(http.StatusBadRequest, res.StatusCode)
	assert.Contains(string(body), "You must specify username query parameter.")
}

func (s *ServerSuite) Test_TwitterNotFound() {
	assert := s.Assert()
	username := "NotFound"
	expectedQuery := fmt.Sprintf("screen_name=%s", username)

	httpmock.RegisterResponderWithQuery("GET", twitter.ShowURL, expectedQuery,
		httpmock.NewStringResponder(http.StatusNotFound,
			`{"errors":[{"code":50,"message":"User not found."}]}`),
	)

	twitterURL := fmt.Sprintf("%s/twitter?username=%s", s.Server.URL, username)
	res, err := s.Client.Get(twitterURL)
	assert.NoError(err)
	defer res.Body.Close()

	url, err := res.Location()
	assert.NoError(err)

	assert.Equal(http.StatusFound, res.StatusCode)

	if assert.NotNil(url) {
		assert.Equal(twitter.DefaultImageURL, url.String())
	}
}

func (s *ServerSuite) Test_TwitterFoundNormal() {
	assert := s.Assert()
	username := "charlyx"
	expectedQuery := fmt.Sprintf("screen_name=%s", username)

	expectedURL := "http://pbs.twimg.com/profile_images/1180040914695327744/qTSU9ZXI_normal.jpg"
	httpmock.RegisterResponderWithQuery("GET", twitter.ShowURL, expectedQuery,
		httpmock.NewStringResponder(http.StatusOK,
			fmt.Sprintf(`{"profile_image_url":"%s"}`, expectedURL)),
	)

	twitterURL := fmt.Sprintf("%s/twitter?username=%s", s.Server.URL, username)
	res, err := s.Client.Get(twitterURL)
	assert.NoError(err)
	defer res.Body.Close()

	url, err := res.Location()
	assert.NoError(err)

	assert.Equal(http.StatusFound, res.StatusCode)

	if assert.NotNil(url) {
		assert.Equal(expectedURL, url.String())
	}
}

func (s *ServerSuite) Test_TwitterFoundBigger() {
	assert := s.Assert()
	username := "charlyx"
	expectedQuery := fmt.Sprintf("screen_name=%s", username)

	expectedURL := "http://pbs.twimg.com/profile_images/1180040914695327744/qTSU9ZXI_bigger.jpg"
	httpmock.RegisterResponderWithQuery("GET", twitter.ShowURL, expectedQuery,
		httpmock.NewStringResponder(http.StatusOK,
			fmt.Sprintf(`{"profile_image_url":"%s"}`, expectedURL)),
	)

	twitterURL := fmt.Sprintf("%s/twitter?username=%s&size=bigger", s.Server.URL, username)
	res, err := s.Client.Get(twitterURL)
	assert.NoError(err)
	defer res.Body.Close()

	url, err := res.Location()
	assert.NoError(err)

	assert.Equal(http.StatusFound, res.StatusCode)

	if assert.NotNil(url) {
		assert.Equal(expectedURL, url.String())
	}
}

func (s *ServerSuite) Test_TwitterFoundMini() {
	assert := s.Assert()
	username := "charlyx"
	expectedQuery := fmt.Sprintf("screen_name=%s", username)

	expectedURL := "http://pbs.twimg.com/profile_images/1180040914695327744/qTSU9ZXI_mini.jpg"
	httpmock.RegisterResponderWithQuery("GET", twitter.ShowURL, expectedQuery,
		httpmock.NewStringResponder(http.StatusOK,
			fmt.Sprintf(`{"profile_image_url":"%s"}`, expectedURL)),
	)

	twitterURL := fmt.Sprintf("%s/twitter?username=%s&size=mini", s.Server.URL, username)
	res, err := s.Client.Get(twitterURL)
	assert.NoError(err)
	defer res.Body.Close()

	url, err := res.Location()
	assert.NoError(err)

	assert.Equal(http.StatusFound, res.StatusCode)

	if assert.NotNil(url) {
		assert.Equal(expectedURL, url.String())
	}
}

func (s *ServerSuite) Test_TwitterFoundOriginal() {
	assert := s.Assert()
	username := "charlyx"
	expectedQuery := fmt.Sprintf("screen_name=%s", username)

	expectedURL := "http://pbs.twimg.com/profile_images/1180040914695327744/qTSU9ZXI.jpg"
	httpmock.RegisterResponderWithQuery("GET", twitter.ShowURL, expectedQuery,
		httpmock.NewStringResponder(http.StatusOK,
			fmt.Sprintf(`{"profile_image_url":"%s"}`, expectedURL)),
	)

	twitterURL := fmt.Sprintf("%s/twitter?username=%s&size=original", s.Server.URL, username)
	res, err := s.Client.Get(twitterURL)
	assert.NoError(err)
	defer res.Body.Close()

	url, err := res.Location()
	assert.NoError(err)

	assert.Equal(http.StatusFound, res.StatusCode)

	if assert.NotNil(url) {
		assert.Equal(expectedURL, url.String())
	}
}

func (s *ServerSuite) Test_TwitterFoundUnknownSize() {
	assert := s.Assert()
	username := "charlyx"
	expectedQuery := fmt.Sprintf("screen_name=%s", username)

	expectedURL := "http://pbs.twimg.com/profile_images/1180040914695327744/qTSU9ZXI_normal.jpg"
	httpmock.RegisterResponderWithQuery("GET", twitter.ShowURL, expectedQuery,
		httpmock.NewStringResponder(http.StatusOK,
			fmt.Sprintf(`{"profile_image_url":"%s"}`, expectedURL)),
	)

	twitterURL := fmt.Sprintf("%s/twitter?username=%s&size=unknown", s.Server.URL, username)
	res, err := s.Client.Get(twitterURL)
	assert.NoError(err)
	defer res.Body.Close()

	url, err := res.Location()
	assert.NoError(err)

	assert.Equal(http.StatusFound, res.StatusCode)

	if assert.NotNil(url) {
		assert.Equal(expectedURL, url.String())
	}
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(ServerSuite))
}
