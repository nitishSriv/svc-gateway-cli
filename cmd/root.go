package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "svcgwcli",
	Short: "CLI utility for pivotal-cf/service-gateway routing API interaction for working with Router Groups",
	Long: `This CLI wraps around the Service Gateway github repository in Pivotal platform.
It provides convenience commands to run from cli for performing routing API calls provided by the service gateway repo.
For example:

Reserve next available TCP port- 
./svcgwcli reserve-port`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
