package cobra

import (
	"fmt"
	"github.com/spf13/cobra"
)

var port int

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start a server",
	Long:  `Start a server and listen on the specified port.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Starting server on port %d\n", port)
		// 启动服务器的逻辑
	},
}

func init() {
	serveCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to start the server on")
}
