package secrets

import (
	"errors"
	"fmt"
	"os"
)

type Env struct{}

func NewEnvClient() Env {
	return Env{}
}

func (Env) Get(key string) (string, error) {
	if key == "" {
		return "", errors.New("please provide a secret key")
	}

	value, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("provided key %s is unset", key)
	}

	return value, nil
}
