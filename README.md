# Manage Helm Deployment to Kubernetes via Rest API endpoint (helm-api) 

## Description
Helm-api is a CLI application written in Golang that gives the ability to perform Install, Uninstall and Upgrade of Helm Charts via Rest API endpoint. The application can be run as a stand alone application or deployed as a Container. Also for convenience there is the  ability to create namespaces and retrieve service IPs of the deployed application.

## When to use helm-api
- helm-api can be used when you need to automate the install and uninstall Helm Charts into a cluster. 
- Used as part of a CI/CD pipeline.


## Requirements
* go 1.8 [https://go.dev/doc/install](https://go.dev/doc/install) to run and install helm-api
* helm if you want to rebuild the helm package
* docker if you want to build or run the container image 
* KubeConfig for the cluster you wish to deploy to.
* Git if you wish to clone helm-api project
* Swag to update swagger document [https://github.com/swaggo/swag/cmd/swag](https://github.com/swaggo/swag/cmd/swag) 

## Project folders
Below is a description helm-api project folders and what they contain
|   Folder        | Description  | 
|-----------|---|
| charts    | Contains the helm chart for helm-api  |
| docs      | Contains the swagger documents |
| documents | Contains cli and api markdown files  |
| modules   | Contains helm-api modules and code  |
| config    | Contains Example payload config files  |
| cmd       | Contains code for helm-api CLI   |
|           |   |

## Installation and Basic usage
This will take you through the steps to install and get helm-api up and running.
<details>
<summary>1. Install</summary>

Once you have installed golang you can run the following command to install helm-api
```yaml
go install github.com/Mrpye/helm-api
```
</details>

<details>
<summary>2. Create a local helm chart repo folder</summary>

This is where you can put you helm charts that you wish to install. it is also possible to add remote helm chart repo via API.
```bash
    mkdir charts
```
</details>
<details>
<summary>3. Start the helm-api web server</summary>

This will run the web-server on port 8080 and we override the default context in out KubeConfig file.
```bash
    helm-api web --port 8080 --ip 0.0.0.0 --folder "./charts" --context "user@k8cluster"
```
</details>

For more instructions on using the helm-api CLI,
check out the CLI documentation [here](./documents/helm-api.md)



## Build and Run helm-api as a container
The following steps will take you through how to build and run helm-api as a container image.


<details>
<summary>1. Clone the repository</summary>

This will clone the helm-api project from github
```bash
# clone the project
git clone https://github.com/Mrpye/helm-api.git

# Change into the directory
cd helm-api
```
</details>

<details>
<summary>2. Build the helm-api container image</summary>

This will build the container image you will need docker installed to build.
```
sudo docker build . -t  helm-api:v1.0.0 -f Dockerfile
```
</details>

<details>
<summary>3. Run the container</summary>

This will run the helm-api container and expose the API endpoint on port 8080 and map the local chart folder to the container so that helm-api can access the local charts.
```
sudo docker run -d -p 8080:8080 --name=helm-api -v /host_path/charts:/go/bin/charts  --env=WEB_IP=0.0.0.0 -t helm-api:1.0.0
```
</details>

---

## Environment  variables
as well as parameters that you can pass to helm-api via the CLI you can also configure helm-api using environmental variables.
- BASE_FOLDER (set where the images will be exported)
- WEB_IP (set the listening ip address 0.0.0.0 allow from everywhere)
- WEB_PORT (set the port to listen on)
- WEB_DEFAULT_CONTEXT (set the default context from kube config)
- WEB_CONFIG_PATH (set the path to the kube config)

## helm-api Helm chart
This guide will show you how to build the helm chart package for helm-api, you will need to have helm installed to build the package.

<details>
<summary>1. Build the helm chart package for helm-api</summary>

```bash
# change into the chart directory
cd charts
# Package the helm-api chart
helm package helm-api

```

the helm chart package will be saved under the charts folder helm-api-0.1.0.tgz

</details>


<details>
<summary>2. Configure helm-api chart</summary>
below are main setting you may want to modify

```yaml
image:
  repository: helm-api
  pullPolicy: Always
  # Overrides the image tag whose default is the chart appVersion.
  tag: "v1.0.0"

#Set the container env values
WebServer:
  listenOn: "0.0.0.0"
  port: "8080"

#volume mount
volumeMount:
  # Create the volume mount
  create: false
  nfsIP: ""
  nfsPath: ""

```

</details>


<details>
<summary>3. How to configure and deploy using helm-api</summary>
Based on the configuration setting from previous step this is what a payload would look like when installing using helm-api


```bash
curl --location --request POST 'localhost:8080/install' \
--header 'Content-Type: application/json' \
--data-raw '{
    "chart":"charts/helm-api-0.1.0.tgz",
    "release_name":"helm-api-install",
    "namespace":"default",
    "params":{
        "repo":"docker.io",
        "nfsIP":"172.0.0.1",
        "nfsPath":"/data/nfs/charts/"
    },
    "config":{
        "image":{
            "repository": "{{.Params.repo}}/library/helm-api",
            "tag": "v1.0.0",
            "pullPolicy":"Always"
        },
        "volumeMount":{
            "create": true,
            "nfsIP":"{{.Params.nfsIP}}",
            "nfsPath":"{{.Params.nfsPath}}"
        }
    }
}
```
</details>


---

## Payload Tokens   

Sometime depending on the Helm Chart the configuration may be large and may contain repeating values. this can make it hard and error prone when you need to update values. So to help with this we can use Tokens.

In the previous example you would have seen {{.Params.repo}} this is a Token. When the payload is sent to helm-api it will get replaced with a value. 

<details>
<summary>Example of inserting values into config</summary>

In this example the value will get replaced with **repo** value that is set in the params property of the json payload.

```json
{
    ...
    "params":{
        "repo":"127.0.0.1",
        "nfsIP":"127.0.0.1",
        "nfsPath":"/data/nfs/charts/"
    },
    ...
}
```

we can define params with a key value pair and use these to populate the config using Token {{.Params.[param key]}}. 

For example to insert the repo ip into our configuration we can use **{{.Params.repo}}**

```json
"config":{
    "image":{
        "repository": "{{.Params.repo}}/library/helm-api",
        "tag": "v1.0.0",
        "pullPolicy":"Always"
    },
    "volumeMount":{
        "create": true,
        "nfsIP":"{{.Params.nfsIP}}",
        "nfsPath":"{{.Params.nfsPath}}"
    }
}
```
</details>

In fact you can access any part of the payload to use as value you can insert.

- **{{.Chart}}** get the chart
- **{{.ReleaseName}}** get the release name
- **{{.Namespace}}** get the namespace
- **{{.Params.[param key]}}** get the parameters using the key
- **{{.Config.[config path]}}** Get values from the config (be careful if accessing config with token as this will not be resolved)
	
## Function Token
Also there a functions that you can use to manipulate the data 

- **base64enc** [string] base64 encode a string
- **base64dec** [string] base64 decode a string
- **lc** [string] make string lowercase
- **uc** [string] make string uppercase
- **domain** [url string] get the domain or ip from a url
- **port_string** [url string] get the port of from a url
- **clean** [string] [replace] clean a string of spaces and special charts
- **concat** [string] [string] concatenate two strings together
- **replace** [string] [find] [replace]replace a value in a string

<details>
<summary>Example</summary> 

```json
"config":{
    "deploy_name":"{{lc .ReleaseName}}",
    "password":"{{base64enc .Params.password}}",
    "image":{
        "repository": "{{.Params.repo}}/library/helm-api",
        "tag": "v1.0.0",
        "pullPolicy":"Always"
    },
    "volumeMount":{
        "create": true,
        "nfsIP":"{{.Params.nfsIP}}",
        "nfsPath":"{{.Params.nfsPath}}"
    }
}
```
</details>

---

## Using the API
The following guide will show you how to make call to the various API endpoints with examples. 

- you can also refer to the API document [here](./documents/api.md) for more details. 
- Also there is a swagger interface you can access when helm-api web-server is running. 

```
http://localhost:8080/docs/index.html
```

## Example API calls

<details>
<summary>Add chart repo</summary>

``` bash
curl --location --request POST 'localhost:8080/add_repo' \
--header 'Authorization: Basic YWRtaW46cGFzc3dvcmQ=' \
--header 'Content-Type: application/json' \
--data-raw '{
    "repo":"https://gitlab.com/api/v4/projects/<project id>/packages/helm/stable",
    "repo_name":"chart-repo"
}'
```

## Payload

- repo (url to the chart repo)
- repo_name (name for the chart repo)

</details>

<details>
<summary>Get Helm Config from a GitLab repo</summary>

``` bash
curl --location --request POST 'localhost:8080/get_config' \
--header 'Authorization: Bearer 9ksddaS7B-Yp45kix-' \
--header 'Content-Type: application/json' \
--data-raw '{
    "config":"https://gitlab.com//api/v4/projects/3/repository/files/helm-api-0.1.0.json/raw?ref=main",
    "answer_file":"https://gitlab.com//api/v4/projects/3/repository/files/helm-api-0.1.0-answer-qa.json/raw?ref=main"
   
}'
```

## Payload

- config (name of the config.json file)
- answer_file (name of the answer.json file)

</details>

<details>
<summary>Get Helm Config from a GitLab repo and answer file from local repo</summary>

``` bash
curl --location --request POST 'localhost:8080/get_config' \
--header 'Authorization: Bearer 9ksddaS7B-Yp45kix-' \
--header 'Content-Type: application/json' \
--data-raw '{
    "config":"https://gitlab.com//api/v4/projects/3/repository/files/helm-api-0.1.0.json/raw?ref=main",
    "answer_file":"helm-api-0.1.0-answer-qa"
}'
```

## Payload

- config (name of the config.json file)
- answer_file (name of the answer.json file)

</details>

<details>
<summary>Get Helm Config from a GitHub repo</summary>

``` bash
curl --location --request POST 'localhost:8080/get_config' \
--header 'Authorization: Bearer 9ksddaS7B-Yp45kix-' \
--header 'Content-Type: application/json' \
--data-raw '{
    "config":"https://api.github.com/repos/Mrpye/helm-api/contents/helm-api-0.1.0.json",
    "answer_file":"https://api.github.com/repos/Mrpye/helm-api/contents/helm-api-0.1.0-answer-qa.json"
   
}'
```

## Payload

- config (name of the config.json file)
- answer_file (name of the answer.json file)

</details>

<details>
<summary>Get Helm Config and Answer from local folder</summary>

``` bash
curl --location --request POST 'localhost:8080/get_config' \
--header 'Content-Type: application/json' \
--data-raw '{
    "config":"helm-api-0.1.0",
    "answer_file":"helm-api-0.1.0-answer-qa"
}'
```

## Payload

- config (name of the config.json file)
- answer_file (name of the answer.json file)

</details>

<details>
<summary>Get Helm Config from local folder and pass params</summary>

``` bash
curl --location --request POST 'localhost:8080/get_config' \
--header 'Content-Type: application/json' \
--data-raw '{
    "config":"helm-api-0.1.0",
    "params":{
        "nfsIP": "172.16.20.10",
        "nfsPath": "/data/nfs/Bifrostv2/charts/",
        "repo": "172.19.2.15"
    }
}'
```

## Payload

- config (name of the config.json file)
- params (params to pass to the config)

</details>

<details>
<summary>Install Chart from local file</summary>

```bash
curl --location --request POST 'localhost:8080/install' \
--header 'Content-Type: application/json' \
--data-raw '{
    "chart":"/charts/demo-0.2.0.tgz",
    "release_name":"demo-test",
    "namespace":"default",
    "config":null
}'
```

## Payload

- chart (path to chart to install)
- release_name (the release name for the installed chart)
- namespace (the name space override leave blank to use default)
- config (config values to override the default values)
</details>


<details>
<summary>Install Chart from repo</summary>

```bash
curl --location --request POST 'localhost:8080/install' \
--header 'Content-Type: application/json' \
--data-raw '{
    "chart":"myrepo/demo",
    "release_name":"demo-test",
    "namespace":"default",
    "config":null
}'
```

## Payload

- chart (chart to install from repo)
- release_name (the release name for the installed chart)
- namespace (the name space override leave blank to use default)
- config (config values to override the default values)
</details>


<details>
<summary>Upgrade Chart</summary>

```bash
curl --location --request POST 'localhost:8080/upgrade' \
--header 'Content-Type: application/json' \
--data-raw '{
    "chart":"charts/demo",
    "release_name":"demo-test",
    "namespace":"default",
    "config":null
}'
```

## Payload

- chart (chart to install)
- release_name (the release name for the installed chart)
- namespace (the name space override leave blank to use default)
- config (config values to override the default values)
</details>


<details>
<summary>UnInstall Chart</summary>

```bash
curl --location --request POST 'localhost:8080/uninstall' \
--header 'Content-Type: application/json' \
--data-raw '{
    "release_name":"demo-test",
    "namespace":"default"
}'
```
## Payload
- release_name (the release name of the chart to uninstall)
- namespace (the name space override leave blank to use default)
</details>


<details>
<summary>Get Service IP</summary>

```bash
curl --location --request POST 'localhost:8080/get_ip' \
--header 'Content-Type: application/json' \
--data-raw '{
    "release_name":"demo-test(.*)",
    "namespace":"default"
}'
```
## Payload
- release_name (regex of the release name to display)
- namespace (the name space override leave blank to use default)
</details>

<details>
<summary>Create a Namespace</summary>

```bash
curl --location --request POST 'localhost:8080/create_ns' \
--header 'Content-Type: application/json' \
--data-raw '{
    "namespace":"default"
}'
```
## Payload
- namespace to create

</details>


<details>
<summary>Test the Web Server is Alive</summary>

```bash
curl --location --request GET 'localhost:8080/'
```
Return OK

</details>

---

## Update the swagger document
The code below shows you how to update the swagger API documents.

If you need more helm on using these tools please refer to the links below
- gin-swagge [https://github.com/swaggo/gin-swagger](https://github.com/swaggo/gin-swagger)
- swag [https://github.com/swaggo/swag](https://github.com/swaggo/swag)

<details>
<summary>1. Install swag</summary>

```bash
#Install swag
go install github.com/swaggo/swag/cmd/swag
```
</details>

<details>
<summary>2. Update APi document</summary>

```bash
#update the API document
swag init
```
</details>
<details>
<summary>3. Update the api.md</summary>

```bash
swagger generate markdown -f .\docs\swagger.json --output .\documents\api.md 
```
</details>
---

## To Do
- Nothing at the moment

--- 

## Main 3rd party Libraries

- [https://github.com/helm/helm](https://github.com/helm/helm)
- [https://github.com/swaggo/gin-swagger](https://github.com/swaggo/gin-swagger) 
- [https://github.com/gin-gonic/gin](https://github.com/gin-gonic/gin)
- [https://github.com/kubernetes/client-go](https://github.com/kubernetes/client-go)


## license
helm-api is Apache 2.0 licensed.