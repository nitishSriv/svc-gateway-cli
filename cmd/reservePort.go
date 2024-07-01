package cmd

import (
	"code.cloudfoundry.org/lager/v3"
	"encoding/json"
	"fmt"
	"github.com/pivotal/service-gateway/routing"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

type RoutingParams struct {
	RoutingAPIURL        string
	UAAURL               string
	RoutingClientID      string
	RoutingClientSecret  string
	UaaCA                string
	RoutingApiServerCA   string
	RoutingApiClientCert string
	RoutingApiClientKey  string
}

const routerGroupNamePrefix = "postgres-rg-"

var routingParams RoutingParams
var serviceInstanceId string

func init() {
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.svcgwcli.yaml)")
	rootCmd.PersistentFlags().StringVar(&routingParams.RoutingAPIURL, "routing-api-url", "", "Routing API URL")
	viper.BindPFlag("routing-api-url", rootCmd.PersistentFlags().Lookup("routing-api-url"))
	rootCmd.MarkPersistentFlagRequired("routing-api-url")

	rootCmd.PersistentFlags().StringVar(&routingParams.UAAURL, "uaa-url", "", "UAA URL")
	viper.BindPFlag("uaa-url", rootCmd.PersistentFlags().Lookup("uaa-url"))
	rootCmd.MarkPersistentFlagRequired("uaa-url")

	rootCmd.PersistentFlags().StringVar(&routingParams.RoutingClientID, "routing-client-id", "", "Routing client ID")
	viper.BindPFlag("routing-client-id", rootCmd.PersistentFlags().Lookup("routing-client-id"))
	rootCmd.MarkPersistentFlagRequired("routing-client-id")

	rootCmd.PersistentFlags().StringVar(&routingParams.RoutingClientSecret, "routing-client-secret", "", "Routing client secret")
	viper.BindPFlag("routing-client-secret", rootCmd.PersistentFlags().Lookup("routing-client-secret"))
	rootCmd.MarkPersistentFlagRequired("routing-client-secret")

	rootCmd.PersistentFlags().StringVar(&routingParams.UaaCA, "uaa-ca", "", "UAA CA")
	viper.BindPFlag("uaa-ca", rootCmd.PersistentFlags().Lookup("uaa-ca"))
	rootCmd.MarkPersistentFlagRequired("uaa-ca")

	rootCmd.PersistentFlags().StringVar(&routingParams.RoutingApiServerCA, "routing-api-server-ca", "", "Routing API server CA")
	viper.BindPFlag("routing-api-server-ca", rootCmd.PersistentFlags().Lookup("routing-api-server-ca"))
	rootCmd.MarkPersistentFlagRequired("routing-api-server-ca")

	rootCmd.PersistentFlags().StringVar(&routingParams.RoutingApiClientCert, "routing-api-client-cert", "", "Routing API client cert")
	viper.BindPFlag("routing-api-client-cert", rootCmd.PersistentFlags().Lookup("routing-api-client-cert"))
	rootCmd.MarkPersistentFlagRequired("routing-api-client-cert")

	rootCmd.PersistentFlags().StringVar(&routingParams.RoutingApiClientKey, "routing-api-client-key", "", "Routing API client key")
	viper.BindPFlag("routing-api-client-key", rootCmd.PersistentFlags().Lookup("routing-api-client-key"))
	rootCmd.MarkPersistentFlagRequired("routing-api-client-key")

	rootCmd.PersistentFlags().StringVar(&serviceInstanceId, "service-instance-id", "", "Service instance ID")
	viper.BindPFlag("service-instance-id", rootCmd.PersistentFlags().Lookup("service-instance-id"))
	rootCmd.MarkPersistentFlagRequired("service-instance-id")

	rootCmd.AddCommand(reservePort)
}

var reservePort = &cobra.Command{
	Use:   "reserve-port",
	Short: "Reserves next available TCP port",
	Run: func(cmd *cobra.Command, args []string) {
		var routingParams RoutingParams
		routingParams.RoutingAPIURL, _ = cmd.Flags().GetString("routing-api-url")
		routingParams.UAAURL, _ = cmd.Flags().GetString("uaa-url")
		routingParams.RoutingClientID, _ = cmd.Flags().GetString("routing-client-id")
		routingParams.RoutingClientSecret, _ = cmd.Flags().GetString("routing-client-secret")
		routingParams.UaaCA, _ = cmd.Flags().GetString("uaa-ca")
		routingParams.RoutingApiServerCA, _ = cmd.Flags().GetString("routing-api-server-ca")
		routingParams.RoutingApiClientCert, _ = cmd.Flags().GetString("routing-api-client-cert")
		routingParams.RoutingApiClientKey, _ = cmd.Flags().GetString("routing-api-client-key")
		serviceInstanceId, _ := cmd.Flags().GetString("service-instance-id")

		outputFormat, _ := cmd.Flags().GetString("output")

		var logger lager.Logger

		// Call the actual function from service-gateway
		routingAdapter, err := buildRoutingAdapter(routingParams, logger)
		if err != nil {
			return
		}

		routerGroupName := routerGroupNamePrefix + serviceInstanceId
		portRange := "1000-1200"

		tcpPort, err := routingAdapter.ReservePort(routerGroupName, portRange)
		if err != nil {
			return
		}

		switch outputFormat {
		case "json":
			jsonOutput, err := json.Marshal(tcpPort)
			if err != nil {
				fmt.Println("Error formatting output as JSON:", err)
				return
			}
			fmt.Println(string(jsonOutput))
		default:
			fmt.Println("Available TCP port:", tcpPort)
		}
	},
}

func buildRoutingAdapter(routingParams RoutingParams, logger lager.Logger) (*routing.RoutingAdapter, error) {
	uaaCAFile, err := os.CreateTemp("", "uaa_ca")
	defer func(name string) {
		_ = os.Remove(name)
	}(uaaCAFile.Name())
	if err != nil {
		return nil, err
	}
	_, _ = uaaCAFile.Write([]byte(routingParams.UaaCA))

	routingAPIServerCAFile, err := os.CreateTemp("", "routingServer_ca")
	defer func(name string) {
		_ = os.Remove(name)
	}(routingAPIServerCAFile.Name())
	if err != nil {
		return nil, err
	}
	_, _ = routingAPIServerCAFile.Write([]byte(routingParams.RoutingApiServerCA))

	routingClientCertFile, err := os.CreateTemp("", "routingClient_cert")
	defer func(name string) {
		_ = os.Remove(name)
	}(routingClientCertFile.Name())
	if err != nil {
		return nil, err
	}
	_, _ = routingClientCertFile.Write([]byte(routingParams.RoutingApiClientCert))

	routingClientKeyFile, err := os.CreateTemp("", "routingClient_key")
	defer func(name string) {
		_ = os.Remove(name)
	}(routingClientKeyFile.Name())
	if err != nil {
		return nil, err
	}
	_, _ = routingClientKeyFile.Write([]byte(routingParams.RoutingApiClientKey))

	routingAdapter, err := routing.NewRoutingAdapter(
		routingParams.RoutingAPIURL,
		routingParams.UAAURL,
		routingParams.RoutingClientID,
		routingParams.RoutingClientSecret,
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
