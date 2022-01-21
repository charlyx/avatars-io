package twitter

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/suite"
)

const imageURL = "https://pbs.twimg.com/profile_images/1180040914695327744/qTSU9ZXI_normal.jpg"

type TwitterSuite struct {
	suite.Suite
	secretAccessor *secretAccessorMock
}

func (s *TwitterSuite) SetupSuite() {
	httpmock.Activate()

	s.secretAccessor = &secretAccessorMock{
		MockGet: func(key string) (string, error) {
			return key, nil
		},
	}
}

func (s *TwitterSuite) SetupTest() {
	httpmock.Reset()
}

func (s *TwitterSuite) TearDownSuite() {
	httpmock.DeactivateAndReset()
}

func (s *TwitterSuite) Test_NewHandlerFuncError() {
	assert := s.Assert()
	secretAccessorMockErr := &secretAccessorMock{
		MockGet: func(key string) (string, error) {
			return "", errors.New(key)
		},
	}

	handler, err := NewHandlerFunc(secretAccessorMockErr)

	assert.Error(err, "could not get twitter token: TWITTER_BEARER_TOKEN")
	assert.Nil(handler)
}

func (s *TwitterSuite) Test_NewHandlerFunc() {
	assert := s.Assert()

	handler, err := NewHandlerFunc(s.secretAccessor)

	assert.NoError(err)
	assert.Implements(new(http.Handler), handler)
}

func (s *TwitterSuite) Test_NewHandlerFuncUsage() {
	assert := s.Assert()
	handlerFunc, _ := NewHandlerFunc(s.secretAccessor)

	resp := getResponse(handlerFunc, "/")
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	assert.Equal(http.StatusBadRequest, resp.StatusCode)
	assert.Contains(string(body), "You must specify username query parameter.")
}

func (s *TwitterSuite) Test_HandlerUserNotFound() {
	assert := s.Assert()
	username := "NotFound"
	mockTwitterShowAPIErr(username)

	handlerFunc, _ := NewHandlerFunc(s.secretAccessor)
	resp := getUserResponse(handlerFunc, username, "")

	url, err := resp.Location()
	assert.NoError(err)
	assert.Equal(http.StatusFound, resp.StatusCode)

	if assert.NotNil(url) {
		assert.Equal(DefaultImageURL, url.String())
	}
}

func (s *TwitterSuite) Test_HandlerUserFound() {
	assert := s.Assert()
	username := "charlyx"
	mockTwitterShowAPI(username)

	handlerFunc, _ := NewHandlerFunc(s.secretAccessor)
	resp := getUserResponse(handlerFunc, username, "")

	url, err := resp.Location()
	assert.NoError(err)
	assert.Equal(http.StatusFound, resp.StatusCode)

	if assert.NotNil(url) {
		assert.Equal(imageURL, url.String())
	}
}

func (s *TwitterSuite) Test_HandlerUserFoundBigger() {
	assert := s.Assert()
	const username = "charlyx"
	const expectedURL = "https://pbs.twimg.com/profile_images/1180040914695327744/qTSU9ZXI_bigger.jpg"
	mockTwitterShowAPI(username)

	handlerFunc, _ := NewHandlerFunc(s.secretAccessor)
	resp := getUserResponse(handlerFunc, username, "bigger")

	url, err := resp.Location()
	assert.NoError(err)
	assert.Equal(http.StatusFound, resp.StatusCode)

	if assert.NotNil(url) {
		assert.Equal(expectedURL, url.String())
	}
}

func (s *TwitterSuite) Test_HandlerUserFoundMini() {
	assert := s.Assert()
	const username = "charlyx"
	const expectedURL = "https://pbs.twimg.com/profile_images/1180040914695327744/qTSU9ZXI_mini.jpg"
	mockTwitterShowAPI(username)

	handlerFunc, _ := NewHandlerFunc(s.secretAccessor)
	resp := getUserResponse(handlerFunc, username, "mini")

	url, err := resp.Location()
	assert.NoError(err)
	assert.Equal(http.StatusFound, resp.StatusCode)

	if assert.NotNil(url) {
		assert.Equal(expectedURL, url.String())
	}
}

func (s *TwitterSuite) Test_HandlerUserFoundOriginal() {
	assert := s.Assert()
	const username = "charlyx"
	const expectedURL = "https://pbs.twimg.com/profile_images/1180040914695327744/qTSU9ZXI.jpg"
	mockTwitterShowAPI(username)

	handlerFunc, _ := NewHandlerFunc(s.secretAccessor)
	resp := getUserResponse(handlerFunc, username, "original")

	url, err := resp.Location()
	assert.NoError(err)
	assert.Equal(http.StatusFound, resp.StatusCode)

	if assert.NotNil(url) {
		assert.Equal(expectedURL, url.String())
	}
}

func (s *TwitterSuite) Test_HandlerUserFoundUnknownSize() {
	assert := s.Assert()
	const username = "charlyx"
	mockTwitterShowAPI(username)

	handlerFunc, _ := NewHandlerFunc(s.secretAccessor)
	resp := getUserResponse(handlerFunc, username, "unknown")

	url, err := resp.Location()
	assert.NoError(err)
	assert.Equal(http.StatusFound, resp.StatusCode)

	if assert.NotNil(url) {
		assert.Equal(imageURL, url.String())
	}
}

func TestServerSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(TwitterSuite))
}

type secretAccessorMock struct {
	MockGet func(key string) (string, error)
}

func (s *secretAccessorMock) Get(key string) (string, error) {
	return s.MockGet(key)
}

func mockTwitterShowAPI(username string) {
	query := fmt.Sprintf("screen_name=%s", username)
	resp := fmt.Sprintf(`{"profile_image_url_https":"%s"}`, imageURL)

	httpmock.RegisterResponderWithQuery("GET", ShowURL, query,
		httpmock.NewStringResponder(http.StatusOK, resp),
	)
}

func mockTwitterShowAPIErr(username string) {
	query := fmt.Sprintf("screen_name=%s", username)
	resp := `{"errors":[{"code":50,"message":"User not found."}]}`

	httpmock.RegisterResponderWithQuery("GET", ShowURL, query,
		httpmock.NewStringResponder(http.StatusNotFound, resp),
	)
}

func getUserResponse(handlerFunc http.HandlerFunc, username, size string) *http.Response {
	target := fmt.Sprintf("/?username=%s&size=%s", username, size)

	return getResponse(handlerFunc, target)
}

func getResponse(handlerFunc http.HandlerFunc, target string) *http.Response {
	req := httptest.NewRequest("GET", target, nil)
	w := httptest.NewRecorder()
	handlerFunc(w, req)

	return w.Result()
}
