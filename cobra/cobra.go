package cobra

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func Cobra() {
	Execute()
}

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "A complex CLI application",
	Long:  "A complex CLI application with multiple commands using Cobra in Golang",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(calcCmd)
}
