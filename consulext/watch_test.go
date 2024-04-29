package consulext

import (
	"testing"
)

func TestNewWatcher(t *testing.T) {
	address := "127.0.0.1:8500"
	RegisterServiceWatcher("friend.rpc", address, func(serviceName string, services []*ServiceInfo) {
		print("%#v", services)
	})
}
