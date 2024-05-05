// Package grpc provides basic gRPC functions.
package grpc

import (
	"sort"
	"strings"

	"github.com/forest33/warthog/business/entity"
)

func (c *Client) sortServicesByName(services []*entity.Service) {
	sort.Slice(services, func(i, j int) bool {
		if strings.HasPrefix(services[j].Name, entity.ReflectionServicePrefix) {
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
