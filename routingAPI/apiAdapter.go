package routingAPI

import (
	"code.cloudfoundry.org/lager/v3"
	"github.com/pivotal/service-gateway/routing"
	"os"
)

type TLSRoutingCredentials struct {
	ClientID     string
	ClientSecret string
	ServerCA     string
	ClientCert   string
	ClientKey    string
}

type AdapterParams struct {
	RoutingAPIURL         string
	UaaURL                string
	UaaCA                 string
	TLSRoutingCredentials TLSRoutingCredentials
}

func BuildRoutingAdapter(routingAdapterParams AdapterParams, logger lager.Logger) (*routing.RoutingAdapter, error) {
	uaaCAFile, err := os.CreateTemp("", "uaa_ca")
	defer func(name string) {
		_ = os.Remove(name)
	}(uaaCAFile.Name())
	if err != nil {
		return nil, err
	}
	_, _ = uaaCAFile.Write([]byte(routingAdapterParams.UaaCA))

	routingAPIServerCAFile, err := os.CreateTemp("", "routingServer_ca")
	defer func(name string) {
		_ = os.Remove(name)
	}(routingAPIServerCAFile.Name())
	if err != nil {
		return nil, err
	}
	_, _ = routingAPIServerCAFile.Write([]byte(routingAdapterParams.TLSRoutingCredentials.ServerCA))

	routingClientCertFile, err := os.CreateTemp("", "routingClient_cert")
	defer func(name string) {
		_ = os.Remove(name)
	}(routingClientCertFile.Name())
	if err != nil {
		return nil, err
	}
	_, _ = routingClientCertFile.Write([]byte(routingAdapterParams.TLSRoutingCredentials.ClientCert))

	routingClientKeyFile, err := os.CreateTemp("", "routingClient_key")
	defer func(name string) {
		_ = os.Remove(name)
	}(routingClientKeyFile.Name())
	if err != nil {
		return nil, err
	}
	_, _ = routingClientKeyFile.Write([]byte(routingAdapterParams.TLSRoutingCredentials.ClientKey))

	routingAdapter, err := routing.NewRoutingAdapter(
		routingAdapterParams.RoutingAPIURL,
		routingAdapterParams.UaaURL,
		routingAdapterParams.TLSRoutingCredentials.ClientID,
		routingAdapterParams.TLSRoutingCredentials.ClientSecret,
		uaaCAFile.Name(),
		routingClientCertFile.Name(),
		routingClientKeyFile.Name(),
		routingAPIServerCAFile.Name(),
		logger,
	)
	if err != nil {
		return nil, err
	}
	return routingAdapter, nil
}
