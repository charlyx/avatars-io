package secrets

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type EnvSuite struct {
	suite.Suite
}

func (g *EnvSuite) Test_Env_SecretAccessor() {
	assert := g.Assert()

	assert.Implements(new(SecretAccessor), Env{})
}

func (g *EnvSuite) Test_Get_NoKey() {
	assert := g.Assert()
	env := Env{}

	value, err := env.Get("")

	assert.EqualError(err, "please provide a secret key")
	assert.Empty(value)
}

func (g *EnvSuite) Test_Get_UnsetKey() {
	assert := g.Assert()
	env := Env{}

	value, err := env.Get("UNSET_KEY")

	assert.EqualError(err, "provided key UNSET_KEY is unset")
	assert.Empty(value)
}

func (g *EnvSuite) Test_Get_Key() {
	assert := g.Assert()
	env := Env{}
	os.Setenv("TEST_KEY", "secret value")

	value, err := env.Get("TEST_KEY")

	assert.NoError(err)
	assert.Equal("secret value", value)
}

func TestSecretEnvSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(EnvSuite))
}
