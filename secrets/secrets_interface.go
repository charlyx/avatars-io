package secrets

type SecretAccessor interface {
	Get(key string) (string, error)
}
