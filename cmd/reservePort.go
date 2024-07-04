package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var routerGroupName string

func init() {
	rootCmd.PersistentFlags().StringVar(&routerGroupName, "router-group-name", "", "Router Group Name")
	viper.BindPFlag("router-group-name", rootCmd.PersistentFlags().Lookup("router-group-name"))
	rootCmd.MarkPersistentFlagRequired("router-group-name")

	rootCmd.AddCommand(reservePort)
}

var reservePort = &cobra.Command{
	Use:   "reserve-port",
	Short: "Reserves next available TCP port and creates/updates a router group",
	Long: "Creates a router group with the given name and assign to it a free port from range (1024-1123) " +
		"If a router group with this name already exists, it will be updated with a free port from the specified range." +
		"The old port will be overwritten",
	Run: func(cmd *cobra.Command, args []string) {
		routerGroupName, _ := cmd.Flags().GetString("router-group-name")

		if routingAdapter == nil {
			fmt.Println("Routing Adapter is not initialized")
			os.Exit(1)
		}

		portRange := "1024-1123"

		tcpPort, err := routingAdapter.ReservePort(routerGroupName, portRange)
		if err != nil {
			fmt.Println("Unable to reserve TCP port", err)
			os.Exit(1)
		}
		fmt.Println("Reserved TCP port: ", tcpPort)
	},
}
