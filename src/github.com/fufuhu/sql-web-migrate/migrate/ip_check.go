package migrate

import (
	"net"
	"net/http"
	"os"
	"strings"

	"go.uber.org/zap"
)

// AllowedNetworks 許可されたネットワークを利用
type AllowedNetworks []*net.IPNet

// IsAllowed 特定のIPアドレスが許可されたネットワークに含まれているかを確認する
func (networks *AllowedNetworks) IsAllowed(ip net.IP) bool {
	for _, network := range *networks {
		if network.Contains(ip) {
			return true
		}
	}
	return false
}

// NetworkConfigStruct Struct to get NetworkConfig
type NetworkConfigStruct struct {
	AllowedNetworks func() AllowedNetworks
}

// NetworkConfig Global Instance of NetworkConfigStruct
var NetworkConfig NetworkConfigStruct

const (
	// SQLMigrateAllowedNetworks 許可されたIPネットワークを表す環境変数(SQL_MIGRATE_ALLOWD_NETWORKS)
	SQLMigrateAllowedNetworks = "SQL_MIGRATE_ALLOWED_NETWORKS"
)

// GetAllowedNetworks 許可されたネットワークのリストを取得する
func GetAllowedNetworks() AllowedNetworks {
	var networks AllowedNetworks

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	value := os.Getenv(SQLMigrateAllowedNetworks)

	networkStrings := strings.Split(value, ",")

	for _, networkString := range networkStrings {
		_, network, err := net.ParseCIDR(networkString)
		if err != nil {
			logger.Error(
				"Network address parse error",
				zap.Error(err),
			)
			break
		}
		networks = append(networks, network)
	}

	return networks
}

const (
	// XForwardedFor X-Forwarded-For constant for IDE
	XForwardedFor = "X-Forwarded-For"
	// RemoteAddr RemoteAddr constant for IDE
	RemoteAddr = "RemoteAddr"
)

func getRemoteAddr(r *http.Request) net.IP {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	remote := r.RemoteAddr
	ip := net.ParseIP(remote)
	logger.Info(
		"Checking HTTP Request Header, RemoteAddr",
		zap.String("RemoteAddr", ip.String()),
	)

	return ip
}

func getXForwardedFor(r *http.Request) []net.IP {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	addresses := r.Header.Get(XForwardedFor)
	logger.Info(
		"Checking HTTP Request Header, X-Forwarded-For",
		zap.String(XForwardedFor, addresses),
	)

	var ips []net.IP
	for _, address := range strings.Split(addresses, " ") {
		ips = append(ips, net.ParseIP(address))
	}
	return ips
}

func init() {
	NetworkConfig = NetworkConfigStruct{
		AllowedNetworks: GetAllowedNetworks,
	}
}
