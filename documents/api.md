


# helm-api
Helm-api is a CLI application written in Golang that gives the ability to perform Install, Uninstall and Upgrade Helm Charts via Rest API endpoint. The application can be run as a stand alone application or deployed as a Container. Also for convenience there is the  ability to create namespaces and retrieve service IPs of the deployed application. GitHub repository at https://github.com/Mrpye/helm-api
  

## Informations

### Version

1.0

### License

[Apache 2.0 licensed](https://github.com/Mrpye/helm-api/blob/main/LICENSE)

### Contact

  https://github.com/Mrpye/helm-api

## Content negotiation

### URI Schemes
  * http

### Consumes
  * application/json

### Produces
  * application/json

## All endpoints

###  operations

| Method  | URI     | Name   | Summary |
|---------|---------|--------|---------|
| POST | /add_repo | [add helm chart repo](#add-helm-chart-repo) | add a helm chart repo |
| GET | / | [check api endpoint](#check-api-endpoint) | Check API Endpoint |
| POST | /create_ns | [create namespace](#create-namespace) | Create Namespace |
| POST | /get_config | [get helm chart config](#get-helm-chart-config) | get the config for helm chart |
| POST | /get_ip | [get service ip](#get-service-ip) | Get Service IP |
| POST | /install | [install helm chart](#install-helm-chart) | Install a helm chart |
| POST | /uninstall | [uninstall helm chart](#uninstall-helm-chart) | uninstall a helm chart |
| POST | /upgrade | [upgrade helm chart](#upgrade-helm-chart) | Upgrade a helm chart |
  


## Paths

### <span id="add-helm-chart-repo"></span> add a helm chart repo (*add-helm-chart-repo*)

```
POST /add_repo
```

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| request | `body` | [BodyTypesImportChartRepo](#body-types-import-chart-repo) | `models.BodyTypesImportChartRepo` | | ✓ | | query params |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#add-helm-chart-repo-200) | OK | charts Repo Added |  | [schema](#add-helm-chart-repo-200-schema) |
| [404](#add-helm-chart-repo-404) | Not Found | error |  | [schema](#add-helm-chart-repo-404-schema) |

#### Responses


##### <span id="add-helm-chart-repo-200"></span> 200 - charts Repo Added
Status: OK

###### <span id="add-helm-chart-repo-200-schema"></span> Schema
   
  



##### <span id="add-helm-chart-repo-404"></span> 404 - error
Status: Not Found

###### <span id="add-helm-chart-repo-404-schema"></span> Schema
   
  



### <span id="check-api-endpoint"></span> Check API Endpoint (*check-api-endpoint*)

```
GET /
```

#### Produces
  * application/json

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#check-api-endpoint-200) | OK | ok |  | [schema](#check-api-endpoint-200-schema) |

#### Responses


##### <span id="check-api-endpoint-200"></span> 200 - ok
Status: OK

###### <span id="check-api-endpoint-200-schema"></span> Schema
   
  



### <span id="create-namespace"></span> Create Namespace (*create-namespace*)

```
POST /create_ns
```

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| request | `body` | [BodyTypesNamespaceChartRepo](#body-types-namespace-chart-repo) | `models.BodyTypesNamespaceChartRepo` | | ✓ | | query params |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#create-namespace-200) | OK | namespace created |  | [schema](#create-namespace-200-schema) |
| [404](#create-namespace-404) | Not Found | error |  | [schema](#create-namespace-404-schema) |

#### Responses


##### <span id="create-namespace-200"></span> 200 - namespace created
Status: OK

###### <span id="create-namespace-200-schema"></span> Schema
   
  



##### <span id="create-namespace-404"></span> 404 - error
Status: Not Found

###### <span id="create-namespace-404-schema"></span> Schema
   
  



### <span id="get-helm-chart-config"></span> get the config for helm chart (*get-helm-chart-config*)

```
POST /get_config
```

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| request | `body` | [BodyTypesGetPayload](#body-types-get-payload) | `models.BodyTypesGetPayload` | | ✓ | | query params |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-helm-chart-config-200) | OK | OK |  | [schema](#get-helm-chart-config-200-schema) |
| [404](#get-helm-chart-config-404) | Not Found | error |  | [schema](#get-helm-chart-config-404-schema) |

#### Responses


##### <span id="get-helm-chart-config-200"></span> 200 - OK
Status: OK

###### <span id="get-helm-chart-config-200-schema"></span> Schema
   
  

[][BodyTypesInstallUpgradeRequest](#body-types-install-upgrade-request)

##### <span id="get-helm-chart-config-404"></span> 404 - error
Status: Not Found

###### <span id="get-helm-chart-config-404-schema"></span> Schema
   
  



### <span id="get-service-ip"></span> Get Service IP (*get-service-ip*)

```
POST /get_ip
```

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| request | `body` | [BodyTypesGetServiceIP](#body-types-get-service-ip) | `models.BodyTypesGetServiceIP` | | ✓ | | query params |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-service-ip-200) | OK | OK |  | [schema](#get-service-ip-200-schema) |
| [404](#get-service-ip-404) | Not Found | error |  | [schema](#get-service-ip-404-schema) |

#### Responses


##### <span id="get-service-ip-200"></span> 200 - OK
Status: OK

###### <span id="get-service-ip-200-schema"></span> Schema
   
  

[][BodyTypesServiceDetails](#body-types-service-details)

##### <span id="get-service-ip-404"></span> 404 - error
Status: Not Found

###### <span id="get-service-ip-404-schema"></span> Schema
   
  



### <span id="install-helm-chart"></span> Install a helm chart (*install-helm-chart*)

```
POST /install
```

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| request | `body` | [BodyTypesInstallUpgradeRequest](#body-types-install-upgrade-request) | `models.BodyTypesInstallUpgradeRequest` | | ✓ | | query params |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#install-helm-chart-200) | OK | chart installed |  | [schema](#install-helm-chart-200-schema) |
| [404](#install-helm-chart-404) | Not Found | error |  | [schema](#install-helm-chart-404-schema) |

#### Responses


##### <span id="install-helm-chart-200"></span> 200 - chart installed
Status: OK

###### <span id="install-helm-chart-200-schema"></span> Schema
   
  



##### <span id="install-helm-chart-404"></span> 404 - error
Status: Not Found

###### <span id="install-helm-chart-404-schema"></span> Schema
   
  



### <span id="uninstall-helm-chart"></span> uninstall a helm chart (*uninstall-helm-chart*)

```
POST /uninstall
```

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| request | `body` | [BodyTypesUninstallChartRepo](#body-types-uninstall-chart-repo) | `models.BodyTypesUninstallChartRepo` | | ✓ | | query params |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#uninstall-helm-chart-200) | OK | chart uninstalled |  | [schema](#uninstall-helm-chart-200-schema) |
| [404](#uninstall-helm-chart-404) | Not Found | error |  | [schema](#uninstall-helm-chart-404-schema) |

#### Responses


##### <span id="uninstall-helm-chart-200"></span> 200 - chart uninstalled
Status: OK

###### <span id="uninstall-helm-chart-200-schema"></span> Schema
   
  



##### <span id="uninstall-helm-chart-404"></span> 404 - error
Status: Not Found

###### <span id="uninstall-helm-chart-404-schema"></span> Schema
   
  



### <span id="upgrade-helm-chart"></span> Upgrade a helm chart (*upgrade-helm-chart*)

```
POST /upgrade
```

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| request | `body` | [BodyTypesInstallUpgradeRequest](#body-types-install-upgrade-request) | `models.BodyTypesInstallUpgradeRequest` | | ✓ | | query params |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#upgrade-helm-chart-200) | OK | chart upgraded |  | [schema](#upgrade-helm-chart-200-schema) |
| [404](#upgrade-helm-chart-404) | Not Found | error |  | [schema](#upgrade-helm-chart-404-schema) |

#### Responses


##### <span id="upgrade-helm-chart-200"></span> 200 - chart upgraded
Status: OK

###### <span id="upgrade-helm-chart-200-schema"></span> Schema
   
  



##### <span id="upgrade-helm-chart-404"></span> 404 - error
Status: Not Found

###### <span id="upgrade-helm-chart-404-schema"></span> Schema
   
  



## Models

### <span id="body-types-get-payload"></span> body_types.GetPayload


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| answer_file | string| `string` |  | |  |  |
| config | string| `string` |  | |  |  |
| params | map of string| `map[string]string` |  | |  |  |



### <span id="body-types-get-service-ip"></span> body_types.GetServiceIP


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| namespace | string| `string` |  | |  |  |
| release_name | string| `string` |  | |  |  |



### <span id="body-types-import-chart-repo"></span> body_types.ImportChartRepo


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| repo | string| `string` |  | |  |  |
| repo_name | string| `string` |  | |  |  |



### <span id="body-types-install-upgrade-request"></span> body_types.InstallUpgradeRequest


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| chart | string| `string` |  | |  |  |
| config | [interface{}](#interface)| `interface{}` |  | |  |  |
| namespace | string| `string` |  | |  |  |
| params | map of string| `map[string]string` |  | |  |  |
| release_name | string| `string` |  | |  |  |



### <span id="body-types-namespace-chart-repo"></span> body_types.NamespaceChartRepo


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| namespace | string| `string` |  | |  |  |



### <span id="body-types-service-details"></span> body_types.ServiceDetails


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| ip | string| `string` |  | |  |  |
| port | integer| `int64` |  | |  |  |
| service_name | string| `string` |  | |  |  |
| service_type | string| `string` |  | |  |  |



### <span id="body-types-uninstall-chart-repo"></span> body_types.UninstallChartRepo


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| namespace | string| `string` |  | |  |  |
| release_name | string| `string` |  | |  |  |


