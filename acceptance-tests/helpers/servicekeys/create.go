package servicekeys

import (
	"csbbrokerpakdatabricks/acceptance-tests/helpers/cf"
	"csbbrokerpakdatabricks/acceptance-tests/helpers/random"
)

type ServiceKey struct {
	name                string
	serviceInstanceName string
}

func Create(serviceInstanceName string) *ServiceKey {
	name := random.Name()
	cf.Run("create-service-key", serviceInstanceName, name)

	return &ServiceKey{
		name:                name,
		serviceInstanceName: serviceInstanceName,
	}
}
