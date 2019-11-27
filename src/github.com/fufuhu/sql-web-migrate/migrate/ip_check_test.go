package migrate

import (
	"net"
	"os"
	"testing"
)

func TestIsAllowed(t *testing.T) {
	var networks AllowedNetworks

	_, network, _ := net.ParseCIDR("127.0.0.0/8")
	networks = append(networks, network)
	_, network, _ = net.ParseCIDR("128.0.0.0/8")
	networks = append(networks, network)

	ip := net.ParseIP("127.0.0.1")
	if !networks.IsAllowed(ip) {
		t.Fail()
	}
}

func TestIsAllowedFailure(t *testing.T) {
	var networks AllowedNetworks

	_, network, _ := net.ParseCIDR("127.0.0.0/8")
	networks = append(networks, network)
	_, network, _ = net.ParseCIDR("128.0.0.0/8")
	networks = append(networks, network)

	ip := net.ParseIP("8.8.8.8")
	if networks.IsAllowed(ip) {
		t.Fail()
	}
}

func TestGetAllowedNetworks(t *testing.T) {

	allowedNetworksString := []string{
		"127.0.0.0/8",
		"128.0.0.0/8",
	}
	os.Setenv(SQLMigrateAllowedNetworks, "127.0.0.0/8,128.0.0.0/8")

	allowedNetworks := NetworkConfig.AllowedNetworks()

	for index, network := range allowedNetworks {
		_, cidr, _ := net.ParseCIDR(allowedNetworksString[index])

		if cidr.Network() != network.Network() {
			t.Fail()
		}
	}
}
