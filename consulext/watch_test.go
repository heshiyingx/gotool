package consulext

import (
	"testing"
)

func TestNewWatcher(t *testing.T) {
	address := "192.168.20.246:8500"
	RegisterServiceWatcher("friend.rpc", address, "", func(serviceName string, services []*ServiceInfo) {
		for _, service := range services {
			println(" ", service.ServiceName, service.IP)

		}
	})
}
