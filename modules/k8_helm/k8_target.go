package k8_helm

import (
	"fmt"
)

type K8 struct {
	DefaultContext string `json:"default_context" yaml:"default_context"`
	ConfigPath     string `json:"config_path" yaml:"config_path"`
	dry_run        bool
	verbose        bool
}

func (m *K8) Verbose() bool {
	return m.verbose
}
func (m *K8) SetVerbose(verbose bool) {
	m.verbose = verbose
}

func (m *K8) DryRun() bool {
	return m.dry_run
}
func (m *K8) SetDryRun(dry_run bool) {
	m.dry_run = dry_run
}

type K8Option func(*K8)

func OptionK8DefaultContext(default_context string) K8Option {
	return func(h *K8) {
		h.DefaultContext = default_context
	}
}

func OptionK8ConfigPath(config_path string) K8Option {
	return func(h *K8) {
		h.ConfigPath = config_path
	}
}

func (m *K8) Update(opts ...K8Option) {
	// Loop through each option
	for _, opt := range opts {
		// Call the option giving the instantiated
		opt(m)
	}
}

func (m *K8) String() string {
	return fmt.Sprintf("%s,%s", m.DefaultContext, m.ConfigPath)
}

//Create a instance of the docker registry
func CreateK8(default_context string, config_path string) *K8 {
	return &K8{
		DefaultContext: default_context,
		ConfigPath:     config_path,
	}
}

func CreateK8Options(opts ...K8Option) *K8 {
	obj := &K8{}
	obj.Update(opts...)
	return obj
}
