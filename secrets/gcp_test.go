package secrets

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type GCPSuite struct {
	suite.Suite
}

func (g *GCPSuite) Test_NewClientError() {
	assert := g.Assert()

	client, err := NewClient("")
	assert.EqualError(err, "projectID must not be empty.")
	assert.Nil(client)
}

func (g *GCPSuite) Test_NewClientSecretAccessor() {
	assert := g.Assert()
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "test.json")

	client, err := NewClient("projectID")
	assert.NoError(err)
	assert.Implements(new(SecretAccessor), client)
}

func TestServerSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(GCPSuite))
}
