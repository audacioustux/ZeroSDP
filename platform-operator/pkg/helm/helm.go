package helm

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-logr/logr"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/repo"
)

type HelmHelper struct {
	cfg  *action.Configuration
	cli  *cli.EnvSettings
	log  logr.Logger
	repo *repo.File
}

func NewHelmHelper(log logr.Logger) (*HelmHelper, error) {
	cli := cli.New()
	actionConfig := new(action.Configuration)

	if err := actionConfig.Init(cli.RESTClientGetter(), cli.Namespace(), os.Getenv("HELM_DRIVER"), debug(log)); err != nil {
		return nil, err
	}

	repoConfig := cli.RepositoryConfig
	repofile, err := loadRepoFile(repoConfig)
	if err != nil {
		return nil, err
	}

	return &HelmHelper{
		cfg:  actionConfig,
		cli:  cli,
		log:  log,
		repo: repofile,
	}, nil
}

func loadRepoFile(repoConfig string) (*repo.File, error) {
	if _, err := os.Stat(repoConfig); os.IsNotExist(err) {
		if err := createRepoFile(repoConfig); err != nil {
			return nil, err
		}
	}

	return repo.LoadFile(repoConfig)
}

func createRepoFile(repoConfig string) error {
	if err := os.MkdirAll(filepath.Dir(repoConfig), 0755); err != nil {
		return err
	}

	if _, err := os.Create(repoConfig); err != nil {
		return err
	}

	return nil
}

func debug(log logr.Logger) func(format string, v ...interface{}) {
	return func(format string, v ...interface{}) {
		log.V(1).Info(fmt.Sprintf(format, v...))
	}
}

func (h *HelmHelper) AddRepo(name, url string) error {
	chart := &repo.Entry{
		Name: name,
		URL:  url,
	}

	repo, err := repo.NewChartRepository(chart, getter.All(h.cli))
	if err != nil {
		return err
	}

	if _, err := repo.DownloadIndexFile(); err != nil {
		return err
	}

	h.repo.Update(chart)

	if err := h.syncRepoFile(); err != nil {
		return err
	}

	return nil
}

func (h *HelmHelper) syncRepoFile() error {
	return h.repo.WriteFile(h.cli.RepositoryConfig, 0644)
}

func (h *HelmHelper) InstallRelease(name, chart, namespace string, values map[string]interface{}) (*release.Release, error) {
	install := action.NewInstall(h.cfg)
	install.Namespace = namespace
	install.ReleaseName = name

	chartPath, err := install.ChartPathOptions.LocateChart(chart, h.cli)
	if err != nil {
		return nil, err
	}

	chartRequested, err := loader.Load(chartPath)
	if err != nil {
		return nil, err
	}

	return install.Run(chartRequested, values)
}
