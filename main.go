// @title helm-api
// @version 1.0
// @description Helm-api is a CLI application written in Golang that gives the ability to perform Install, Uninstall and Upgrade Helm Charts via Rest API endpoint. The application can be run as a stand alone application or deployed as a Container. Also for convenience there is the  ability to create namespaces and retrieve service IPs of the deployed application. GitHub repository at https://github.com/Mrpye/helm-api

// @contact.url https://github.com/Mrpye/helm-api

// @license.name Apache 2.0 licensed
// @license.url https://github.com/Mrpye/helm-api/blob/main/LICENSE

// @host localhost:8080
// @BasePath /
// @query.collection.format multi
package main

import "github.com/Mrpye/helm-api/cmd"

func main() {
	cmd.Execute()
}
