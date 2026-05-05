package services

import "csbbrokerpakdatabricks/acceptance-tests/helpers/servicekeys"

func (s *ServiceInstance) CreateServiceKey() *servicekeys.ServiceKey {
	return servicekeys.Create(s.Name)
}
