package storagepool

import (
	"errors"
)

type healthChecker interface {
	populateHostMap() (map[string]string, error)
	deleteHost(string) error
}

func newHealthChecker(healthCheckBaseDir, healthCheckType, metadataUrl string) (healthChecker, error) {
	if healthCheckType == "file" {
		return &fileBasedHealthCheck{
			healthCheckBaseDir: healthCheckBaseDir,
		}, nil
	}
	if healthCheckType == "metadata" {
		return &metadataBasedHealthCheck{
			metadataUrl: metadataUrl,
		}, nil
	}
	return nil, errors.New("Unknown healthcheck type " + healthCheckType)
}
