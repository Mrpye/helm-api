package cmd

import (
	_ "github.com/Mrpye/helm-api/docs"
	"github.com/Mrpye/helm-api/modules/api"
	"github.com/spf13/cobra"
)

func Cmd_WebServer() *cobra.Command {
	// webserverCmd represents the webserver command
	var web_port string
	var web_ip string
	var web_default_context string
	var web_config_path string
	var cmd = &cobra.Command{
		Use:   "web",
		Short: "Start a API Web-Server",
		Long:  `Start a API Web-Server`,
		Run: func(cmd *cobra.Command, args []string) {
			api.StartWebServer(web_ip, web_port, web_default_context, web_config_path)
		},
	}
	cmd.Flags().StringVarP(&web_port, "port", "p", "8080", "The Port to listen on")
	cmd.Flags().StringVarP(&web_ip, "ip", "i", "localhost", "IP to listen on")
	cmd.Flags().StringVarP(&web_default_context, "context", "x", "", "K8s KubeConfig Context to use")
	cmd.Flags().StringVarP(&web_config_path, "config", "c", "", "Override KubeConfig path")

	return cmd
}
func init() {
	rootCmd.AddCommand(Cmd_WebServer())
}
