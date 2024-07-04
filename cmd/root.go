package cmd

import (
	"code.cloudfoundry.org/lager/v3"
	"fmt"
	"github.com/nitishSriv/svc-gateway-cli/routingAPI"
	"github.com/pivotal/service-gateway/routing"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"reflect"
)

var adapterParams routingAPI.AdapterParams
var routingAdapter *routing.RoutingAdapter

var rootCmd = &cobra.Command{
	Use:   "svcgwcli",
	Short: "CLI utility for pivotal-cf/service-gateway routingAPI API interaction for working with Router Groups",
	Long: `This CLI wraps around the Service Gateway github repository in Pivotal platform.
It provides convenience commands to run from cli for performing routingAPI API calls provided by the service gateway repo.
For example:

Reserve next available TCP port- 
./svcgwcli reserve-port`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Check if any required field is missing
		v := reflect.ValueOf(&adapterParams).Elem()
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			if field.String() == "" && v.Type().Field(i).Name != "UaaCA" {
				fmt.Printf("%s must be set via command line flags\n", v.Type().Field(i).Name)
				os.Exit(1)
			}
		}

		if adapterParams.TLSRoutingCredentials.ClientID == "" || adapterParams.TLSRoutingCredentials.ClientSecret == "" ||
			adapterParams.TLSRoutingCredentials.ClientKey == "" || adapterParams.TLSRoutingCredentials.ClientCert == "" {
			fmt.Print("[routing-client-id, routing-client-secret, routing-api-server-ca, routing-api-client-cert, routing-api-client-key] must be set via command line flags\n")
			os.Exit(1)
		}

		var logger lager.Logger
		var err error
		routingAdapter, err = routingAPI.BuildRoutingAdapter(adapterParams, logger)
		if err != nil {
			fmt.Println("Error building routingAPI adapter", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.PersistentFlags().StringVar(&adapterParams.RoutingAPIURL, "routing-api-url", "", "Routing API URL")
	viper.BindPFlag("routing-api-url", rootCmd.PersistentFlags().Lookup("routing-api-url"))
	rootCmd.MarkPersistentFlagRequired("routing-api-url")

	rootCmd.PersistentFlags().StringVar(&adapterParams.UaaURL, "uaa-url", "", "UAA URL")
	viper.BindPFlag("uaa-url", rootCmd.PersistentFlags().Lookup("uaa-url"))
	rootCmd.MarkPersistentFlagRequired("uaa-url")

	rootCmd.PersistentFlags().StringVar(&adapterParams.TLSRoutingCredentials.ClientID, "routing-client-id", "", "Routing client ID")
	viper.BindPFlag("routing-client-id", rootCmd.PersistentFlags().Lookup("routing-client-id"))
	rootCmd.MarkPersistentFlagRequired("routing-client-id")

	rootCmd.PersistentFlags().StringVar(&adapterParams.TLSRoutingCredentials.ClientSecret, "routing-client-secret", "", "Routing client secret")
	viper.BindPFlag("routing-client-secret", rootCmd.PersistentFlags().Lookup("routing-client-secret"))
	rootCmd.MarkPersistentFlagRequired("routing-client-secret")

	rootCmd.PersistentFlags().StringVar(&adapterParams.TLSRoutingCredentials.ServerCA, "routing-api-server-ca", "", "Routing API server CA")
	viper.BindPFlag("routing-api-server-ca", rootCmd.PersistentFlags().Lookup("routing-api-server-ca"))
	rootCmd.MarkPersistentFlagRequired("routing-api-server-ca")

	adapterParams.UaaCA = viper.GetString("routing-api-server-ca")

	rootCmd.PersistentFlags().StringVar(&adapterParams.TLSRoutingCredentials.ClientCert, "routing-api-client-cert", "", "Routing API client cert")
	viper.BindPFlag("routing-api-client-cert", rootCmd.PersistentFlags().Lookup("routing-api-client-cert"))
	rootCmd.MarkPersistentFlagRequired("routing-api-client-cert")

	rootCmd.PersistentFlags().StringVar(&adapterParams.TLSRoutingCredentials.ClientKey, "routing-api-client-key", "", "Routing API client key")
	viper.BindPFlag("routing-api-client-key", rootCmd.PersistentFlags().Lookup("routing-api-client-key"))
	rootCmd.MarkPersistentFlagRequired("routing-api-client-key")
}
