package helm

import (
	"fmt"
	"os"

	"github.com/go-logr/logr"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
)

type HelmHelper struct {
	cfg *action.Configuration
}

func NewHelmHelper(log logr.Logger) (*HelmHelper, error) {
	settings := cli.New()
	actionConfig := new(action.Configuration)

	if err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), debug(log)); err != nil {
		return nil, err
	}

	return &HelmHelper{
		cfg: actionConfig,
	}, nil
}

func debug(log logr.Logger) func(format string, v ...interface{}) {
	return func(format string, v ...interface{}) {
		log.V(1).Info(fmt.Sprintf(format, v...))
	}
}

func (helm *HelmHelper) GetHelmRelease(ns string) ([]*release.Release, error) {
	client := action.NewList(helm.cfg)
	client.AllNamespaces = true
	client.SetStateMask()

	return client.Run()
}

func (helm *HelmHelper) GetChart(client *action.Install, name string) (*chart.Chart, error) {
	settings := cli.New()

	chartPath, err := client.ChartPathOptions.LocateChart(name, settings)
	if err != nil {
		return nil, err
	}

	return loader.Load(chartPath)
}

func (helm *HelmHelper) InstallChart(client *action.Install, chart *chart.Chart) (*release.Release, error) {
	return client.Run(chart, nil)
}

func (helm *HelmHelper) NewHelminstall(ns string, name string, repo string) (*action.Install, error) {
	client := action.NewInstall(helm.cfg)
	client.Namespace = ns
	client.ReleaseName = name
	client.RepoURL = repo

	return client, nil
}

// type HelmChart struct {
// 	Name        string
// 	Repo        string
// 	ReleaseName string
// 	Namespace   string
// 	Values      map[string]interface{}
// }

// func (chart *HelmChart) Install(log logr.Logger) error {
// 	settings := cli.New()
// 	actionConfig := new(action.Configuration)

// 	if err := actionConfig.Init(settings.RESTClientGetter(), chart.Namespace, os.Getenv("HELM_DRIVER"), func(format string, v ...interface{}) {
// 		log.V(1).Info(fmt.Sprintf(format, v...))
// 	}); err != nil {
// 		return err
// 	}

// 	client := action.NewInstall(actionConfig)
// 	client.Namespace = chart.Namespace
// 	client.ReleaseName = chart.ReleaseName
// 	client.RepoURL = chart.Repo

// 	chartPath, err := client.ChartPathOptions.LocateChart(chart.Name, settings)
// 	if err != nil {
// 		return err
// 	}

// 	chart, err := loader.Load(chartPath)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = client.Run(chart, nil)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
