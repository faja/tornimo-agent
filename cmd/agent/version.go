package agent

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "print the version info",
		//Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%-8s : %s\n%-8s : %s\n%-8s : %s\n%-8s : %s\n",
				"version", versionInfoCli,
				"revision", versionInfoCommit,
				"date", versionInfoDate,
				"go", runtime.Version(),
			)
		},
	}
	versionInfoCli    = "no version provided"
	versionInfoCommit = "no revision provided"
	versionInfoDate   = "no date provided"
)
