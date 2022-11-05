package grpc

import (
	"sort"

	"github.com/Forest33/warthog/business/entity"
)

func (c *Client) sortServicesByName(services []*entity.Service) {
	sort.Slice(services, func(i, j int) bool {
		if services[j].Name == entity.ReflectionServiceFQN {
			return true
		}
		return services[i].Name < services[j].Name
	})
}

func (c *Client) sortMethodsByName(methods []*entity.Method) {
	sort.Slice(methods, func(i, j int) bool {
		return methods[i].Name < methods[j].Name
	})
}
