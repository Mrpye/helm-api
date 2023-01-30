# helm-api
Rest API Web Server to manage the install and uninstall of Hel Charts

# Description
The solution enables you to easily install and uninstall Helm Charts via REST API.

# When to use helm-api
helm-api can be used when you want to automate the install and uninstall Helm Charts. Also can be used as part of your CI/CD pipeline.

## Requirements
* you will need to install go 1.8 [https://go.dev/doc/install](https://go.dev/doc/install) to run and install helm-api

## Installation

```yaml
go install github.com/Mrpye/helm-api
```

## Run as a container
1. Clone the repository

```
git clone https://github.com/Mrpye/helm-api.git
```

2. Build the container as API endpoint
```
sudo docker build . -t  helm-api:v1.0.0 -f Dockerfile
```
3. Run the container
```
sudo docker run -d -p 8080:8080 --name=helm-api --restart always  -v /host_path/charts:/go/bin/charts  --env=WEB_IP=0.0.0.0 --env=WEB_PORT=8080 -t helm-api:1.0.0
```

### Environment  variables
- BASE_FOLDER (set where the images will be exported)
- WEB_IP (set the listening ip address 0.0.0.0 allow from everywhere)
- WEB_PORT (set the port to listen on)
- WEB_DEFAULT_CONTEXT (set the default context from kube config)
- WEB_CONFIG_PATH (set the path to the kube config)


## How to use helm-api CLI
Check out the CLI documentation [here](./documents/helm-api.md)


# Using the API
## Run web server
```bash
    helm-api.md web -f /charts -p 8080 -i 0.0.0.0
```

## Add chart repo
``` bash
curl --location --request POST 'localhost:8080/addrepo' \
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

---

## Install Chart
```bash
curl --location --request POST 'localhost:8080/install' \
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
---
## UnInstall Chart
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
--

# To Do
- Need to add the upgrade chart option

# 3rd party Libraries
https://github.com/helm/helm

# license
helm-api is Apache 2.0 licensed.