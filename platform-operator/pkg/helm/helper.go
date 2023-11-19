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
	config   *action.Configuration
	settings *cli.EnvSettings
	log      logr.Logger
	repoFile *repo.File
}

func NewHelmHelper(log logr.Logger) (*HelmHelper, error) {
	settings := cli.New()

	config := new(action.Configuration)
	if err := config.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), debug(log)); err != nil {
		return nil, err
	}

	repoConfig := settings.RepositoryConfig
	repoFile, err := loadRepoFile(repoConfig)
	if err != nil {
		return nil, err
	}

	return &HelmHelper{
		config:   config,
		settings: settings,
		log:      log,
		repoFile: repoFile,
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

func (h *HelmHelper) AddRepo(cfg *repo.Entry) error {
	repo, err := repo.NewChartRepository(cfg, getter.All(h.settings))
	if err != nil {
		return err
	}

	if _, err := repo.DownloadIndexFile(); err != nil {
		return err
	}

	h.repoFile.Update(cfg)

	if err := h.syncRepoFile(); err != nil {
		return err
	}

	return nil
}

func (h *HelmHelper) syncRepoFile() error {
	return h.repoFile.WriteFile(h.settings.RepositoryConfig, 0644)
}

func (h *HelmHelper) InstallRelease(name, chart, namespace string, values map[string]interface{}) (*release.Release, error) {
	install := action.NewInstall(h.config)
	install.Namespace = namespace
	install.ReleaseName = name

	chartPath, err := install.ChartPathOptions.LocateChart(chart, h.settings)
	if err != nil {
		return nil, err
	}

	chartRequested, err := loader.Load(chartPath)
	if err != nil {
		return nil, err
	}

	return install.Run(chartRequested, values)
}

func (h *HelmHelper) GetRelease(name string) (*release.Release, error) {
	client := action.NewGet(h.config)

	return client.Run(name)
}

func (h *HelmHelper) IsDeployed(rel *release.Release) (bool, error) {
	if rel.Info.Status == release.StatusDeployed {
		return true, nil
	}

	return false, nil
}
