package agent

import (
	"github.com/spf13/cobra"
)

var (
	RootCmd = &cobra.Command{
		Use: "tornimo-agent [command]",
	}

	configFile string
)

func init() {
	RootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "/etc/tornimo/config.yaml", "config file")
}
