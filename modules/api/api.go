package api

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/Mrpye/go_k8_helm"
	"github.com/Mrpye/golib/encrypt"
	_ "github.com/Mrpye/helm-api/docs"
	"github.com/Mrpye/helm-api/modules/body_types"
	"github.com/Mrpye/helm-api/modules/git"
	"github.com/Mrpye/helm-api/modules/template"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var web_port string
var web_ip string
var web_default_context string
var web_config_path string

func remapConfig(config_file string) string {
	if config_file == "" {
		return ""
	}
	//****************************************
	// remap the ConfigName to the config_path
	//****************************************
	parts := strings.Split(config_file, "/")

	if strings.Contains(config_file, "://") {
		return config_file
	}
	config_path := ""
	if len(parts) < 2 {
		config_path = path.Join("config", config_file)
	}

	//*********************************************
	//check if the ConfigName has a .json extension
	//*********************************************
	if !strings.HasSuffix(config_file, ".json") {
		config_path = config_path + ".json"
	}
	return config_path
}

// @Summary get the config for helm chart
// @ID get-helm-chart-config
// @Produce json
// @Param request body body_types.GetPayload.request true true "query params"
// @Success 200 {object}  body_types.InstallUpgradeRequest.response
// @Failure 400 {string}  string "error"
// @Router /get_config [post]
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func postGetConfig(c *gin.Context) {
	var importRequest body_types.GetPayload

	// Call BindJSON to bind the received JSON to
	if err := c.BindJSON(&importRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, "Bad payload")
		return
	}

	//****************************************
	// remap the ConfigName to the config_path
	//****************************************
	config_path := remapConfig(importRequest.ConfigName)
	answer_path := remapConfig(importRequest.AnswerFile)

	token := ""
	if val, ok := c.Request.Header["Authorization"]; ok {
		token = strings.ReplaceAll(val[0], "Bearer ", "")
	}

	//***************
	//Load the config
	//***************
	obj, err := git.LoadHelmConfig(config_path, answer_path, importRequest.Params, importRequest.ReleaseName, importRequest.Namespace, token)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
	} else {
		c.IndentedJSON(http.StatusCreated, obj)
	}

}

// @Summary get the config for helm chart and installs
// @ID get-helm-chart-config-install
// @Produce json
// @Param request body body_types.GetPayload.request true true "query params"
// @Success 200 {object}  body_types.InstallUpgradeRequest.response
// @Failure 400 {object}  body_types.InstallUpgradeRequest.response
// @Router /get_config_install [post]
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func postGetConfigInstall(c *gin.Context) {
	var importRequest body_types.GetPayload

	// Call BindJSON to bind the received JSON to
	if err := c.BindJSON(&importRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, "Bad payload")
		return
	}

	//****************************************
	// remap the ConfigName to the config_path
	//****************************************
	config_path := remapConfig(importRequest.ConfigName)
	answer_path := remapConfig(importRequest.AnswerFile)

	token := ""
	if val, ok := c.Request.Header["Authorization"]; ok {
		token = strings.ReplaceAll(val[0], "Bearer ", "")
	}

	//***************
	//Load the config
	//***************
	obj, err := git.LoadHelmConfig(config_path, answer_path, importRequest.Params, importRequest.ReleaseName, importRequest.Namespace, token)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
	}
	//*******************
	//Perform the install
	//*******************
	helm, err := go_k8_helm.CreateK8KubeConfig(web_default_context, web_config_path)
	if err != nil {
		obj.Result = err.Error()
		c.IndentedJSON(http.StatusBadRequest, obj)
		return
	}

	//****************
	//parse the config
	//****************
	config, err := template.ParseInterfaceMap(*obj, obj.Config)
	if err != nil {
		obj.Result = err.Error()
		c.IndentedJSON(http.StatusBadRequest, obj)
		return
	}
	err = helm.DeployHelmChart(obj.Chart, obj.ReleaseName, obj.Namespace, config)

	if err != nil {
		obj.Result = err.Error()
		c.IndentedJSON(http.StatusBadRequest, obj)
	} else {
		obj.Result = "success"
		c.IndentedJSON(http.StatusCreated, obj)
	}

}

// @Summary get the config for helm chart and installs
// @ID get-helm-chart-config-upgrade
// @Produce json
// @Param request body body_types.GetPayload.request true true "query params"
// @Success 200 {object}  body_types.InstallUpgradeRequest.response
// @Failure 400 {object}  body_types.InstallUpgradeRequest.response
// @Router /get_config_upgrade [post]
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func postGetConfigUpgrade(c *gin.Context) {
	var importRequest body_types.GetPayload

	// Call BindJSON to bind the received JSON to
	if err := c.BindJSON(&importRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, "Bad payload")
		return
	}

	//****************************************
	// remap the ConfigName to the config_path
	//****************************************
	config_path := remapConfig(importRequest.ConfigName)
	answer_path := remapConfig(importRequest.AnswerFile)

	token := ""
	if val, ok := c.Request.Header["Authorization"]; ok {
		token = strings.ReplaceAll(val[0], "Bearer ", "")
	}

	//***************
	//Load the config
	//***************
	obj, err := git.LoadHelmConfig(config_path, answer_path, importRequest.Params, importRequest.ReleaseName, importRequest.Namespace, token)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
	}
	//*******************
	//Perform the install
	//*******************
	helm, err := go_k8_helm.CreateK8KubeConfig(web_default_context, web_config_path)
	if err != nil {
		obj.Result = err.Error()
		c.IndentedJSON(http.StatusBadRequest, obj)
		return
	}
	//****************
	//parse the config
	//****************
	config, err := template.ParseInterfaceMap(*obj, obj.Config)
	if err != nil {
		obj.Result = err.Error()
		c.IndentedJSON(http.StatusBadRequest, obj)
		return
	}
	err = helm.UpgradeHelmChart(obj.Chart, obj.ReleaseName, obj.Namespace, config)

	if err != nil {
		obj.Result = err.Error()
		c.IndentedJSON(http.StatusBadRequest, obj)
	} else {
		obj.Result = "success"
		c.IndentedJSON(http.StatusCreated, obj)
	}

}

// @Summary add a helm chart repo
// @ID add-helm-chart-repo
// @Produce json
// @Param request body body_types.ImportChartRepo.request true "query params"
// @Success 200 {string}  string "charts Repo Added"
// @Failure 400 {string}  string "error"
// @Router /add_repo [post]
func postAddRepo(c *gin.Context) {
	var importRequest body_types.ImportChartRepo

	// Call BindJSON to bind the received JSON to
	if err := c.BindJSON(&importRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, "Bad payload")
		return
	}
	user := ""
	password := ""
	if val, ok := c.Request.Header["Authorization"]; ok {
		token := strings.Split(val[0], " ")[1]
		dec, err := encrypt.Base64DecString(token)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
			return
		}
		parts := strings.Split(dec, ":")

		user = parts[0]
		password = parts[1]
	}

	helm, err := go_k8_helm.CreateK8KubeConfig(web_default_context, web_config_path)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}
	if importRequest.Repo != "" {
		err := helm.RepoAdd(importRequest.RepoName, importRequest.Repo, user, password)
		if err != nil {
			if !strings.Contains(err.Error(), "already exists") {
				c.IndentedJSON(http.StatusBadRequest, err.Error())
				return
			}
		}
		err = helm.RepoUpdate()
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
			return
		}
	}

	c.IndentedJSON(http.StatusCreated, importRequest.RepoName+" Repo Added")

}

// @Summary Install a helm chart
// @ID install-helm-chart
// @Produce json
// @Param request body body_types.InstallUpgradeRequest.request true "query params"
// @Success 200 {object}  body_types.InstallUpgradeRequest.response
// @Failure 400 {object}  body_types.InstallUpgradeRequest.response
// @Router /install [post]
func postInstall(c *gin.Context) {
	var importRequest body_types.InstallUpgradeRequest

	// Call BindJSON to bind the received JSON to
	if err := c.BindJSON(&importRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, "Bad payload")
		return
	}

	helm, err := go_k8_helm.CreateK8KubeConfig(web_default_context, web_config_path)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}
	//****************
	//parse the config
	//****************
	config, err := template.ParseInterfaceMap(importRequest, importRequest.Config)
	if err != nil {
		importRequest.Result = err.Error()
		c.IndentedJSON(http.StatusBadRequest, importRequest)
		return
	}
	err = helm.DeployHelmChart(importRequest.Chart, importRequest.ReleaseName, importRequest.Namespace, config)

	if err != nil {
		importRequest.Result = err.Error()
		c.IndentedJSON(http.StatusBadRequest, importRequest)
	} else {
		importRequest.Result = "success"
		c.IndentedJSON(http.StatusCreated, importRequest)
	}

}

// @Summary Upgrade a helm chart
// @ID upgrade-helm-chart
// @Produce json
// @Param request body body_types.InstallUpgradeRequest.request true "query params"
// @Success 200 {object}  body_types.InstallUpgradeRequest.response
// @Failure 400 {object}  body_types.InstallUpgradeRequest.response
// @Router /upgrade [post]
func postUpgrade(c *gin.Context) {
	var importRequest body_types.InstallUpgradeRequest

	// Call BindJSON to bind the received JSON to
	if err := c.BindJSON(&importRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, "Bad payload")
		return
	}

	helm, err := go_k8_helm.CreateK8KubeConfig(web_default_context, web_config_path)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}
	err = helm.UpgradeHelmChart(importRequest.Chart, importRequest.ReleaseName, importRequest.Namespace, importRequest.Config)
	if err != nil {
		importRequest.Result = err.Error()
		c.IndentedJSON(http.StatusBadRequest, importRequest)
	} else {
		importRequest.Result = "success"
		c.IndentedJSON(http.StatusCreated, importRequest)
	}

}

// @Summary uninstall a helm chart
// @ID uninstall-helm-chart
// @Produce json
// @Param request body body_types.UninstallChartRepo.request true "query params"
// @Success 200 {object} body_types.InstallUpgradeRequest.response
// @Failure 400 {object} body_types.InstallUpgradeRequest.response
// @Router /uninstall [post]
func postUninstall(c *gin.Context) {
	var importRequest body_types.UninstallChartRepo

	// Call BindJSON to bind the received JSON to
	if err := c.BindJSON(&importRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, "Bad payload")
		return
	}

	helm, err := go_k8_helm.CreateK8KubeConfig(web_default_context, web_config_path)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}
	err = helm.UninstallHelmChart(importRequest.ReleaseName, importRequest.Namespace)
	if err != nil {
		importRequest.Result = err.Error()
		c.IndentedJSON(http.StatusBadRequest, importRequest)
	} else {
		importRequest.Result = "success"
		c.IndentedJSON(http.StatusCreated, importRequest)
	}
}

// @Summary Create Namespace
// @ID create-namespace
// @Produce json
// @Param request body body_types.NamespaceChartRepo.request true "query params"
// @Success 200 {string}  string "namespace created"
// @Failure 400 {string}  string "error"
// @Router /create_ns [post]
func postCreateNS(c *gin.Context) {
	var importRequest body_types.NamespaceChartRepo

	// Call BindJSON to bind the received JSON to
	if err := c.BindJSON(&importRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, "Bad payload")
		return
	}

	helm, err := go_k8_helm.CreateK8KubeConfig(web_default_context, web_config_path)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}
	err = helm.CreateNS(importRequest.Namespace)
	if err != nil {
		if !strings.Contains(err.Error(), "cannot create default name space") {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
			return
		}
	}
	c.IndentedJSON(http.StatusCreated, importRequest.Namespace+" Created")
}

// @Summary Get Service IP
// @ID get-service-ip
// @Produce json
// @Param request body body_types.GetServiceIP.request true "query params"
// @Success 200 {object}  []body_types.ServiceDetails.response
// @Failure 400 {string}  string "error"
// @Router /get_ip [post]
func postGetServiceIP(c *gin.Context) {
	var importRequest body_types.GetServiceIP

	if err := c.BindJSON(&importRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, "Bad payload")
		return
	}

	helm, err := go_k8_helm.CreateK8KubeConfig(web_default_context, web_config_path)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
		return
	}
	results, err := helm.GetServiceIP(importRequest.Namespace, importRequest.ReleaseName)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err.Error())
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

// Function to start web server
func StartWebServer(ip string, port string, default_context string, config_path string) {
	//****************
	//Set the variable
	//****************
	web_ip = ip
	web_port = port
	web_default_context = default_context
	web_config_path = config_path

	//*****************
	//Set up the server
	//*****************
	fmt.Println("Starting Web-Server")
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	//*****************
	//Set up the routes
	//*****************
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.POST("/install", postInstall)
	router.POST("/uninstall", postUninstall)
	router.POST("/upgrade", postUpgrade)
	router.POST("/add_repo", postAddRepo)
	router.POST("/create_ns", postCreateNS)
	router.POST("/get_ip", postGetServiceIP)
	router.POST("/get_config", postGetConfig)
	router.POST("/get_config_install", postGetConfigInstall)
	router.POST("/get_config_upgrade", postGetConfigUpgrade)

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

	//****************
	//Start the server
	//****************
	fmt.Printf("Web-Server started %s:%s", web_ip, web_port)
	router.Run(fmt.Sprintf("%s:%s", web_ip, web_port))

}
