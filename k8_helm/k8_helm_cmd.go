package k8_helm

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/theckman/go-flock"
	"gopkg.in/yaml.v2"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type SimpleRESTClientGetter struct {
	Namespace  string
	KubeConfig rest.Config
}

func NewRESTClientGetter(namespace string, kubeConfig rest.Config) *SimpleRESTClientGetter {
	return &SimpleRESTClientGetter{
		Namespace:  namespace,
		KubeConfig: kubeConfig,
	}
}

func (c *SimpleRESTClientGetter) ToRESTConfig() (*rest.Config, error) {
	return &c.KubeConfig, nil
}

func (c *SimpleRESTClientGetter) ToDiscoveryClient() (discovery.CachedDiscoveryInterface, error) {
	config, err := c.ToRESTConfig()
	if err != nil {
		return nil, err
	}

	// The more groups you have, the more discovery requests you need to make.
	// given 25 groups (our groups + a few custom conf) with one-ish version each, discovery needs to make 50 requests
	// double it just so we don't end up here again for a while.  This config is only used for discovery.
	config.Burst = 100

	discoveryClient, _ := discovery.NewDiscoveryClientForConfig(config)
	return memory.NewMemCacheClient(discoveryClient), nil
}

func (c *SimpleRESTClientGetter) ToRESTMapper() (meta.RESTMapper, error) {
	discoveryClient, err := c.ToDiscoveryClient()
	if err != nil {
		return nil, err
	}

	mapper := restmapper.NewDeferredDiscoveryRESTMapper(discoveryClient)
	expander := restmapper.NewShortcutExpander(mapper, discoveryClient)
	return expander, nil
}

func (c *SimpleRESTClientGetter) ToRawKubeConfigLoader() clientcmd.ClientConfig {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	// use the standard defaults for this client command
	// DEPRECATED: remove and replace with something more accurate
	loadingRules.DefaultClientConfig = &clientcmd.DefaultClientConfig

	overrides := &clientcmd.ConfigOverrides{ClusterDefaults: clientcmd.ClusterDefaults}
	overrides.Context.Namespace = c.Namespace

	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, overrides)
}

//Read the config for the connection to K8s cluster
func (m *K8) buildRestConfig() (*rest.Config, error) {
	var kube_config string

	if m.ConfigPath == "" {
		home := homedir.HomeDir()
		kube_config = path.Join(home, ".kube", "config") // flag.String("kubeconfig", path.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		if strings.HasSuffix(m.ConfigPath, "/") {
			kube_config = path.Join(m.ConfigPath, "config") //flag.String("kubeconfig", m.ConfigPath, "absolute path to the kubeconfig file")
		} else {
			kube_config = m.ConfigPath
		}
	}
	flag.Parse()

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	loadingRules.ExplicitPath = kube_config
	// if you want to change the loading rules (which files in which order), you can do so here
	var configOverrides clientcmd.ConfigOverrides

	if m.DefaultContext != "" {
		configOverrides = clientcmd.ConfigOverrides{
			CurrentContext: m.DefaultContext,
		}
	} else {
		configOverrides = clientcmd.ConfigOverrides{}
	}
	// if you want to change override values or bind them to flags, there are methods to help you
	//clientcmd.LoadFromFile(*kubeconfig)
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &configOverrides)
	config, err := kubeConfig.ClientConfig()

	//config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)

	if err != nil {
		return nil, fmt.Errorf("unable to load kube config %s with context %s", kube_config, m.DefaultContext)

	}
	return config, nil
}

func (m *K8) UninstallHelmChart(release_name string, namespace string) error {

	nameSpace := "default"
	if namespace != "" {
		nameSpace = namespace
	}
	log.Printf("Info: Uninstalling Chart from path: %s in namespace: %s\n", release_name, nameSpace)
	actionConfig := new(action.Configuration)
	// You can pass an empty string instead of settings.Namespace() to list
	// all namespaces
	config, err := m.buildRestConfig()
	if err != nil {
		return err
	}
	getter := NewRESTClientGetter(namespace, *config)

	if err := actionConfig.Init(getter, namespace,
		os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		log.Printf("%+v", err)
		return err
	}
	client := action.NewUninstall(actionConfig)
	rel, err := client.Run(release_name)
	if err != nil {
		return err
	}
	log.Print(rel)
	log.Printf("Info: Uninstalled Chart from path: %s in namespace: %s\n", release_name, namespace)
	return err
}

func (m *K8) DeployHelmChart(chart_path string, release_name string, namespace string, configs map[string]interface{}) error {

	chartPath := chart_path
	nameSpace := "default"
	if namespace != "" {
		nameSpace = namespace
	}

	releaseName := release_name
	log.Printf("Info: Installing Chart from path: %s in namespace: %s\n", release_name, nameSpace)
	settings := cli.New()

	actionConfig := new(action.Configuration)

	config, err := m.buildRestConfig()
	if err != nil {
		return err
	}
	getter := NewRESTClientGetter(nameSpace, *config)
	// You can pass an empty string instead of settings.Namespace() to list
	// all namespaces
	if err := actionConfig.Init(getter, nameSpace,
		os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		return err
	}

	client := action.NewInstall(actionConfig)
	client.Namespace = nameSpace
	client.ReleaseName = releaseName

	ch_path, err := client.LocateChart(chartPath, settings)
	if err != nil {
		return err
	}
	// load chart from the path
	chart, err := loader.Load(ch_path)
	if err != nil {
		return err
	}

	// install the chart here
	rel, err := client.Run(chart, configs)
	if err != nil {
		return err
	}

	log.Printf("Info: Installed Chart from path: %s in namespace: %s\n", rel.Name, rel.Namespace)
	// this will confirm the values set during installation
	//log.Println(rel.Config)
	return nil
}

func (m *K8) UpgradeHelmChart(chart_path string, release_name string, namespace string, configs map[string]interface{}) error {

	chartPath := chart_path
	nameSpace := "default"
	if namespace != "" {
		nameSpace = namespace
	}

	releaseName := release_name
	log.Printf("Info: Updating Chart from path: %s in namespace: %s\n", release_name, nameSpace)
	settings := cli.New()

	actionConfig := new(action.Configuration)

	config, err := m.buildRestConfig()
	if err != nil {
		return err
	}
	getter := NewRESTClientGetter(nameSpace, *config)
	// You can pass an empty string instead of settings.Namespace() to list
	// all namespaces
	if err := actionConfig.Init(getter, nameSpace,
		os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		return err
	}

	client := action.NewUpgrade(actionConfig)
	client.Namespace = nameSpace

	ch_path, err := client.LocateChart(chartPath, settings)
	if err != nil {
		return err
	}
	// load chart from the path
	chart, err := loader.Load(ch_path)
	if err != nil {
		return err
	}

	// install the chart here
	rel, err := client.Run(releaseName, chart, configs)
	if err != nil {
		return err
	}

	log.Printf("Info: Installed Chart from path: %s in namespace: %s\n", rel.Name, rel.Namespace)
	// this will confirm the values set during installation
	//log.Println(rel.Config)
	return nil
}

func (m *K8) RepoAdd(name string, url string, user string, password string) error {
	settings := cli.New()
	repoFile := settings.RepositoryConfig

	//Ensure the file directory exists as it is required for file locking
	err := os.MkdirAll(filepath.Dir(repoFile), os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err
	}

	// Acquire a file lock for process synchronization
	fileLock := flock.New(strings.Replace(repoFile, filepath.Ext(repoFile), ".lock", 1))
	lockCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	locked, err := fileLock.TryLockContext(lockCtx, time.Second)
	if err == nil && locked {
		defer fileLock.Unlock()
	}
	if err != nil {
		return err
	}

	b, err := ioutil.ReadFile(repoFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	var f repo.File
	if err := yaml.Unmarshal(b, &f); err != nil {
		return err
	}

	if f.Has(name) {

		return fmt.Errorf("repository name (%s) already exists", name)
	}

	c := repo.Entry{
		Name:     name,
		URL:      url,
		Username: user,
		Password: password,
		//PassCredentialsAll: true,
	}

	r, err := repo.NewChartRepository(&c, getter.All(settings))
	if err != nil {
		return err
	}
	_, err = r.DownloadIndexFile()
	if err != nil {
		err := errors.Wrapf(err, "looks like %q is not a valid chart repository or cannot be reached", url)
		return err
	}

	f.Update(&c)

	if err := f.WriteFile(repoFile, 0644); err != nil {
		return err
	}
	fmt.Printf("%q has been added to your repositories\n", name)
	return nil
}

// RepoUpdate updates charts for all helm repos
func (m *K8) RepoUpdate() error {
	settings := cli.New()
	repoFile := settings.RepositoryConfig

	f, err := repo.LoadFile(repoFile)
	if os.IsNotExist(errors.Cause(err)) || len(f.Repositories) == 0 {
		return errors.New("no repositories found. You must add one before updating")
	}
	var repos []*repo.ChartRepository
	for _, cfg := range f.Repositories {
		r, err := repo.NewChartRepository(cfg, getter.All(settings))
		if err != nil {
			return err
		}
		repos = append(repos, r)
	}

	fmt.Printf("Hang tight while we grab the latest from your chart repositories...\n")
	var wg sync.WaitGroup
	for _, re := range repos {
		wg.Add(1)
		go func(re *repo.ChartRepository) {
			defer wg.Done()
			if _, err := re.DownloadIndexFile(); err != nil {
				fmt.Printf("...Unable to get an update from the %q chart repository (%s):\n\t%s\n", re.Config.Name, re.Config.URL, err)
			} else {
				fmt.Printf("...Successfully got an update from the %q chart repository\n", re.Config.Name)
			}
		}(re)
	}
	wg.Wait()
	fmt.Printf("Update Complete. ⎈ Happy Helming!⎈\n")
	return nil
}
