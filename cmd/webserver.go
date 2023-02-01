package cmd

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"

	_ "github.com/Mrpye/helm-api/docs"
	"github.com/Mrpye/helm-api/k8_helm"
	"github.com/Mrpye/helm-api/lib"
	"github.com/Mrpye/helm-api/template"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var web_base_folder string
var web_default_context string
var web_config_path string

func DecodeB64(message string) (retour string) {
	base64Text := make([]byte, base64.StdEncoding.DecodedLen(len(message)))
	base64.StdEncoding.Decode(base64Text, []byte(message))
	fmt.Printf("base64: %s\n", base64Text)

	return strings.ReplaceAll(string(base64Text), "\x00", "")
}

// @Summary add a helm chart repo
// @ID add-helm-chart-repo
// @Produce json
// @Param request body lib.ImportChartRepo.request true "query params"
// @Success 200 {string}  string "charts Repo Added"
// @Failure 404 {string}  string "error"
// @Router /add_repo [post]
func postAddRepo(c *gin.Context) {
	var importRequest lib.ImportChartRepo

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

	helm := k8_helm.CreateK8(web_default_context, web_config_path)
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

// @Summary Install a helm chart
// @ID install-helm-chart
// @Produce json
// @Param request body lib.InstallUpgradeRequest.request true "query params"
// @Success 200 {string}  string "chart installed"
// @Failure 404 {string}  string "error"
// @Router /install [post]
func postInstall(c *gin.Context) {
	var importRequest lib.InstallUpgradeRequest

	// Call BindJSON to bind the received JSON to
	if err := c.BindJSON(&importRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, "Bad payload")
		return
	}

	helm := k8_helm.CreateK8(web_default_context, web_config_path)

	//****************
	//parse the config
	//****************
	config, err := template.ParseInterfaceMap(importRequest, importRequest.Config)
	if err != nil {
		c.IndentedJSON(400, err.Error())
		return
	}
	err = helm.DeployHelmChart(importRequest.Chart, importRequest.ReleaseName, importRequest.Namespace, config)

	if err != nil {
		c.IndentedJSON(400, err.Error())
	} else {
		c.IndentedJSON(http.StatusCreated, importRequest.ReleaseName+" Installed")
	}

}

// @Summary Upgrade a helm chart
// @ID upgrade-helm-chart
// @Produce json
// @Param request body lib.InstallUpgradeRequest.request true "query params"
// @Success 200 {string}  string "chart upgraded"
// @Failure 404 {string}  string "error"
// @Router /upgrade [post]
func postUpgrade(c *gin.Context) {
	var importRequest lib.InstallUpgradeRequest

	// Call BindJSON to bind the received JSON to
	if err := c.BindJSON(&importRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, "Bad payload")
		return
	}

	helm := k8_helm.CreateK8(web_default_context, web_config_path)

	err := helm.UpgradeHelmChart(importRequest.Chart, importRequest.ReleaseName, importRequest.Namespace, importRequest.Config)

	if err != nil {
		c.IndentedJSON(400, err.Error())
	} else {
		c.IndentedJSON(http.StatusCreated, importRequest.ReleaseName+" Installed")
	}

}

// @Summary uninstall a helm chart
// @ID uninstall-helm-chart
// @Produce json
// @Param request body lib.UninstallChartRepo.request true "query params"
// @Success 200 {string}  string "chart uninstalled"
// @Failure 404 {string}  string "error"
// @Router /uninstall [post]
func postUninstall(c *gin.Context) {
	var importRequest lib.UninstallChartRepo

	// Call BindJSON to bind the received JSON to
	if err := c.BindJSON(&importRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, "Bad payload")
		return
	}

	helm := k8_helm.CreateK8(web_default_context, web_config_path)
	err := helm.UninstallHelmChart(importRequest.ReleaseName, importRequest.Namespace)
	if err != nil {
		c.IndentedJSON(400, err.Error())
	} else {
		c.IndentedJSON(http.StatusCreated, importRequest.ReleaseName+" Uninstalled")
	}
}

// @Summary Create Namespace
// @ID create-namespace
// @Produce json
// @Param request body lib.NamespaceChartRepo.request true "query params"
// @Success 200 {string}  string "namespace created"
// @Failure 404 {string}  string "error"
// @Router /create_ns [post]
func postCreateNS(c *gin.Context) {
	var importRequest lib.NamespaceChartRepo

	// Call BindJSON to bind the received JSON to
	if err := c.BindJSON(&importRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, "Bad payload")
		return
	}

	helm := k8_helm.CreateK8(web_default_context, web_config_path)
	err := helm.CreateNS(importRequest.Namespace)
	if err != nil {
		c.IndentedJSON(400, err.Error())
	} else {
		c.IndentedJSON(http.StatusCreated, importRequest.Namespace+" Created")
	}
}

// @Summary Get Service IP
// @ID get-service-ip
// @Produce json
// @Param request body lib.GetServiceIP.request true "query params"
// @Success 200 {object}  []lib.ServiceDetails.response
// @Failure 404 {string}  string "error"
// @Router /get_ip [post]
func postGetServiceIP(c *gin.Context) {
	var importRequest lib.GetServiceIP

	// Call BindJSON to bind the received JSON to
	if err := c.BindJSON(&importRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, "Bad payload")
		return
	}

	helm := k8_helm.CreateK8(web_default_context, web_config_path)
	results, err := helm.GetServiceIP(importRequest.Namespace, importRequest.ReleaseName)
	if err != nil {
		c.IndentedJSON(400, err.Error())
	} else {
		c.IndentedJSON(http.StatusCreated, results)
	}
}

// @Summary Check API Endpoint
// @ID check-api-endpoint
// @Produce json
// @Success 200 {string}  string "ok"
// @Router / [get]
func getOK(c *gin.Context) {

	c.IndentedJSON(http.StatusCreated, "OK")

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
			router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
			router.POST("/install", postInstall)
			router.POST("/uninstall", postUninstall)
			router.POST("/upgrade", postUpgrade)
			router.POST("/add_repo", postAddRepo)
			router.POST("/create_ns", postCreateNS)
			router.POST("/get_ip", postGetServiceIP)

			router.GET("/", getOK)

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
				web_default_context = os.Getenv("WEB_DEFAULT_CONTEXT")
			}
			if os.Getenv("WEB_CONFIG_PATH") != "" {
				web_config_path = os.Getenv("WEB_CONFIG_PATH")
			}
			if os.Getenv("BASE_FOLDER") != "" {
				web_base_folder = os.Getenv("BASE_FOLDER")
			}

			router.Run(web_ip + ":" + web_port)
		},
	}
	cmd.Flags().StringVarP(&web_port, "port", "p", "8080", "The Port to listen on")
	cmd.Flags().StringVarP(&web_ip, "ip", "i", "localhost", "IP to listen on")
	cmd.Flags().StringVarP(&web_base_folder, "folder", "f", "", "The local helm chart folder")
	cmd.Flags().StringVarP(&web_default_context, "context", "x", "", "K8s KubeConfig Context to use")
	cmd.Flags().StringVarP(&web_config_path, "config", "c", "", "Override KubeConfig path")

	return cmd
}
func init() {
	rootCmd.AddCommand(Cmd_WebServer())
}
