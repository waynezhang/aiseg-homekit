package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = ""
var Revision = ""

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("aiseg v%s+%s\n", Version, Revision)
	},
}
