package cmd

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Mrpye/helm-api/helm"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var web_base_folder string
var web_default_context string
var web_config_path string

type ImportExportRequest struct {
	Repo        string                 `json:"repo"`
	RepoName    string                 `json:"repo_name"`
	Chart       string                 `json:"chart"`
	ReleaseName string                 `json:"release_name"`
	Namespace   string                 `json:"namespace"`
	Config      map[string]interface{} `json:"config"`
}

func DecodeB64(message string) (retour string) {
	base64Text := make([]byte, base64.StdEncoding.DecodedLen(len(message)))
	base64.StdEncoding.Decode(base64Text, []byte(message))
	fmt.Printf("base64: %s\n", base64Text)

	return strings.ReplaceAll(string(base64Text), "\x00", "")
}

func postAddRepo(c *gin.Context) {
	var importRequest ImportExportRequest

	// Call BindJSON to bind the received JSON to
	if err := c.BindJSON(&importRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, "Bad payload")
		return
	}
	user := ""
	password := ""
	if val, ok := c.Request.Header["Authorization"]; ok {
		token := strings.Split(val[0], " ")[1]
		dec := DecodeB64(token)
		parts := strings.Split(dec, ":")

		user = parts[0]
		password = parts[1]
	}

	helm := helm.CreateK8(web_default_context, web_config_path)
	if importRequest.Repo != "" {
		err := helm.RepoAdd(importRequest.RepoName, importRequest.Repo, user, password)
		if err != nil {
			c.IndentedJSON(400, err.Error())
			return
		}
		helm.RepoUpdate()
		if err != nil {
			c.IndentedJSON(400, err.Error())
			return
		}
	}

	c.IndentedJSON(http.StatusCreated, importRequest.RepoName+" Repo Added")

}

func postInstall(c *gin.Context) {
	var importRequest ImportExportRequest

	// Call BindJSON to bind the received JSON to
	if err := c.BindJSON(&importRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, "Bad payload")
		return
	}

	helm := helm.CreateK8(web_default_context, web_config_path)

	err := helm.DeployHelmChart(importRequest.Chart, importRequest.ReleaseName, importRequest.Namespace, importRequest.Config)

	if err != nil {
		c.IndentedJSON(400, err.Error())
	} else {
		c.IndentedJSON(http.StatusCreated, importRequest.ReleaseName+" Installed")
	}

}

func postUninstall(c *gin.Context) {
	var importRequest ImportExportRequest

	// Call BindJSON to bind the received JSON to
	if err := c.BindJSON(&importRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, "Bad payload")
		return
	}

	helm := helm.CreateK8(web_default_context, web_config_path)
	err := helm.UninstallHelmChart(importRequest.ReleaseName, importRequest.Namespace)
	if err != nil {
		c.IndentedJSON(400, err.Error())
	} else {
		c.IndentedJSON(http.StatusCreated, importRequest.ReleaseName+" Uninstalled")
	}

}

func Cmd_WebServer() *cobra.Command {
	// webserverCmd represents the webserver command
	var web_port string
	var web_ip string

	var cmd = &cobra.Command{
		Use:   "web",
		Short: "Start a API Web-Server",
		Long:  `Start a API Web-Server`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Starting Web-Server")

			router := gin.Default()

			router.POST("/install", postInstall)
			router.POST("/uninstall", postUninstall)
			router.POST("/addrepo", postAddRepo)
			//**********************************
			//Set up the environmental variables
			//**********************************
			if os.Getenv("WEB_IP") != "" {
				web_ip = os.Getenv("WEB_IP")
			}
			if os.Getenv("WEB_PORT") != "" {
				web_port = os.Getenv("WEB_PORT")
			}
			if os.Getenv("WEB_DEFAULT_CONTEXT") != "" {
				web_ip = os.Getenv("WEB_DEFAULT_CONTEXT")
			}
			if os.Getenv("WEB_CONFIG_PATH") != "" {
				web_port = os.Getenv("WEB_CONFIG_PATH")
			}
			if os.Getenv("BASE_FOLDER") != "" {
				web_base_folder = os.Getenv("BASE_FOLDER")
			}

			router.Run(web_ip + ":" + web_port)
		},
	}
	cmd.Flags().StringVarP(&web_port, "port", "p", "8080", "Listen on Port")
	cmd.Flags().StringVarP(&web_ip, "ip", "i", "localhost", "Listen on Ip")
	cmd.Flags().StringVarP(&web_base_folder, "folder", "f", "", "base export import folder")
	cmd.Flags().StringVarP(&web_default_context, "context", "x", "", "K8s Config context to use")
	cmd.Flags().StringVarP(&web_config_path, "config", "c", "", "override kube config path")

	return cmd
}
func init() {
	rootCmd.AddCommand(Cmd_WebServer())
}
