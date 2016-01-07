package storagepool

type healthChecker interface {
	populateHostMap() (map[string]string, error)
	deleteHost(string) error
}

func newHealthChecker(metadataUrl string) (healthChecker, error) {
	return &metadataBasedHealthCheck{
		metadataUrl: metadataUrl,
	}, nil
}
