package body_types

type InstallUpgradeRequest struct {
	Chart       string                 `json:"chart"`
	ReleaseName string                 `json:"release_name"`
	Namespace   string                 `json:"namespace"`
	Params      map[string]string      `json:"params"`
	Config      map[string]interface{} `json:"config"`
}

type ImportChartRepo struct {
	Repo     string `json:"repo"`
	RepoName string `json:"repo_name"`
}

type UninstallChartRepo struct {
	ReleaseName string `json:"release_name"`
	Namespace   string `json:"namespace"`
}

type GetServiceIP struct {
	ReleaseName string `json:"release_name"`
	Namespace   string `json:"namespace"`
}

type NamespaceChartRepo struct {
	Namespace string `json:"namespace"`
}

type ServiceDetails struct {
	ServiceName string `json:"service_name" yaml:"service_name"`
	ServiceType string `json:"service_type" yaml:"service_type"`
	IP          string `json:"ip" yaml:"ip"`
	Port        int32  `json:"port" yaml:"port"`
}

type GetPayload struct {
	ConfigName  string            `json:"config"`
	AnswerFile  string            `json:"answer_file"`
	ReleaseName string            `json:"release_name"`
	Namespace   string            `json:"namespace"`
	Params      map[string]string `json:"params"`
}
