// Copyright (c) 2021-2022 北京渠成软件有限公司(Beijing Qucheng Software Co., Ltd. www.qucheng.com) All rights reserved.
// Use of this source code is covered by the following dual licenses:
// (1) Z PUBLIC LICENSE 1.2 (ZPL 1.2)
// (2) Affero General Public License 3.0 (AGPL 3.0)
// license that can be found in the LICENSE file.

package helm

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/easysoft/qcadmin/internal/pkg/util/kutil"
	"github.com/easysoft/qcadmin/internal/pkg/util/log"
	"github.com/gofrs/flock"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"helm.sh/helm/v3/cmd/helm/search"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/helmpath"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/repo"
	"helm.sh/helm/v3/pkg/storage/driver"
)

const (
	helmDriver = "secrets"
)

func nolog(format string, v ...interface{}) {}

type Config struct {
	Namespace string
}

type Client struct {
	actionConfig *action.Configuration
	Namespace    string
	settings     *cli.EnvSettings
}

func NewClient(config *Config) (*Client, error) {
	settings := cli.New()
	client := &Client{}
	settings.SetNamespace(config.Namespace)
	client.settings = settings
	actionConfig := &action.Configuration{}
	if err := actionConfig.Init(settings.RESTClientGetter(), config.Namespace, helmDriver, nolog); err != nil {
		return nil, err
	}
	client.actionConfig = actionConfig
	client.Namespace = config.Namespace
	return client, nil
}

func (c Client) List(limit, offset int, pattern string) ([]*release.Release, int, error) {
	client := action.NewList(c.actionConfig)
	if c.Namespace == "" {
		client.AllNamespaces = true
		client.All = true
	}
	client.SetStateMask()
	list, err := client.Run()
	if err != nil {
		return nil, 0, err
	}
	// TODO limit & offset
	// client.Limit = limit
	// client.Offset = offset - 1
	if pattern != "" {
		client.Filter = pattern
	}
	result, err := client.Run()
	if err != nil {
		return nil, 0, err
	}

	return result, len(list), nil
}

func (c Client) GetDetail(name string) (*release.Release, error) {
	client := action.NewGet(c.actionConfig)
	result, err := client.Run(name)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c Client) Upgrade(name, repoName, chartName, chartVersion string, values map[string]interface{}) (*release.Release, error) {
	repos, err := c.ListRepo()
	if err != nil {
		return nil, err
	}
	var rp *repo.Entry
	for _, r := range repos {
		if r.Name == repoName {
			rp = r
		}
	}
	if rp == nil {
		return nil, errors.New("get chart detail failed, repo not found")
	}

	histClient := action.NewHistory(c.actionConfig)
	histClient.Max = 1
	if _, err := histClient.Run(name); err == driver.ErrReleaseNotFound {
		// If a release does not exist, install it.
		return c.Install(name, repoName, chartName, chartVersion, values)
	}

	client := action.NewUpgrade(c.actionConfig)
	client.Namespace = c.Namespace
	client.RepoURL = rp.URL
	client.Username = rp.Username
	client.Password = rp.Password
	client.DryRun = false
	client.ChartPathOptions.InsecureSkipTLSverify = true
	if len(chartVersion) > 0 {
		client.ChartPathOptions.Version = chartVersion
	}
	p, err := client.ChartPathOptions.LocateChart(chartName, c.settings)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("locate chart %s failed: %v", chartName, err))
	}
	ct, err := loader.Load(p)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("load chart %s failed: %v", chartName, err))

	}
	// TODO 获取之前参数，并且更新参数
	release, err := client.Run(name, ct, values)
	if err != nil {
		return release, errors.Wrap(err, fmt.Sprintf("upgrade tool %s with chart %s failed: %v", name, chartName, err))
	}
	return release, nil
}

func (c Client) Install(name, repoName, chartName, chartVersion string, values map[string]interface{}) (*release.Release, error) {
	repos, err := c.ListRepo()
	if err != nil {
		return nil, err
	}
	var rp *repo.Entry
	for _, r := range repos {
		if r.Name == repoName {
			rp = r
		}
	}
	if rp == nil {
		return nil, errors.New("get chart detail failed, repo not found")
	}
	client := action.NewInstall(c.actionConfig)
	client.ReleaseName = name
	client.Namespace = c.Namespace
	client.RepoURL = rp.URL
	client.Username = rp.Username
	client.Password = rp.Password
	client.ChartPathOptions.InsecureSkipTLSverify = true
	if len(chartVersion) != 0 {
		client.ChartPathOptions.Version = chartVersion
	}
	p, err := client.ChartPathOptions.LocateChart(chartName, c.settings)
	if err != nil {
		return nil, fmt.Errorf("locate chart %s failed: %v", chartName, err)
	}
	ct, err := loader.Load(p)
	if err != nil {
		return nil, fmt.Errorf("load chart %s failed: %v", chartName, err)
	}
	re, err := client.Run(ct, values)
	if err != nil {
		return re, errors.Wrap(err, fmt.Sprintf("install %s with chart %s failed: %v", name, chartName, err))
	}
	return re, nil
}

func (c Client) UpdateRepo() error {
	if !kutil.NeedCacheHelmFile() {
		return nil
	}
	repoFile := c.settings.RepositoryConfig
	repoCache := c.settings.RepositoryCache
	f, err := repo.LoadFile(repoFile)
	if err != nil {
		return fmt.Errorf("load file of repo %s failed: %v", repoFile, err)
	}
	var rps []*repo.ChartRepository
	for _, cfg := range f.Repositories {
		r, err := repo.NewChartRepository(cfg, getter.All(c.settings))
		if err != nil {
			return err
		}
		if repoCache != "" {
			r.CachePath = repoCache
		}
		rps = append(rps, r)
	}
	updateCharts(rps)
	return nil
}

func updateCharts(repos []*repo.ChartRepository) {
	log.Flog.Debug("Hang tight while we grab the latest from your chart repositories...")
	var wg sync.WaitGroup
	for _, re := range repos {
		wg.Add(1)
		go func(re *repo.ChartRepository) {
			defer wg.Done()
			if _, err := re.DownloadIndexFile(); err != nil {
				log.Flog.Errorf("...Unable to get an update from the %q chart repository (%s):\n\t%s", re.Config.Name, re.Config.URL, err)
			} else {
				log.Flog.Debugf("...Successfully got an update from the %q chart repository", re.Config.Name)
			}
		}(re)
	}
	wg.Wait()
	log.Flog.Debug("Update Complete. ⎈ Happy Helming!⎈ ")
}

func (c Client) ListRepo() ([]*repo.Entry, error) {
	var repos []*repo.Entry
	f, err := repo.LoadFile(c.settings.RepositoryConfig)
	if err != nil {
		return repos, err
	}
	return f.Repositories, nil
}

func (c Client) ListCharts(repoName, pattern string) ([]*search.Result, error) {
	repos, err := c.ListRepo()
	if err != nil {
		return nil, fmt.Errorf("list chart failed: %v", err)
	}
	i := search.NewIndex()
	for _, re := range repos {
		if repoName != "KRepoAll" {
			if repoName != re.Name {
				continue
			}
		}
		path := filepath.Join(c.settings.RepositoryCache, helmpath.CacheIndexFile(re.Name))
		indexFile, err := repo.LoadIndexFile(path)
		if err != nil {
			return nil, fmt.Errorf("list chart failed: %v", err)
		}
		i.AddRepo(re.Name, indexFile, true)
	}
	var result []*search.Result
	if pattern != "" {
		result = i.SearchLiteral(pattern, 1)
	} else {
		result = i.All()
	}
	search.SortScore(result)
	return result, nil
}

func (c Client) AddRepo(name, url, username, password string) error {
	repoFile := c.settings.RepositoryConfig

	err := os.MkdirAll(filepath.Dir(repoFile), os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err
	}

	fileLock := flock.New(strings.Replace(repoFile, filepath.Ext(repoFile), ".lock", 1))
	lockCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	locked, err := fileLock.TryLockContext(lockCtx, time.Second)
	if err == nil && locked {
		defer func() {
			if err := fileLock.Unlock(); err != nil {
				log.Flog.Errorf("addRepo fileLock.Unlock failed, error: %s", err.Error())
			}
		}()
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
		return errors.Errorf("repository name (%s) already exists, please specify a different name", name)
	}

	e := repo.Entry{
		Name:                  name,
		URL:                   url,
		Username:              username,
		Password:              password,
		InsecureSkipTLSverify: true,
	}

	r, err := repo.NewChartRepository(&e, getter.All(c.settings))
	if err != nil {
		return err
	}
	r.CachePath = c.settings.RepositoryCache
	if _, err := r.DownloadIndexFile(); err != nil {
		return errors.Wrapf(err, "looks like %q is not a valid chart repository or cannot be reached", url)
	}
	f.Update(&e)
	if err := f.WriteFile(repoFile, 0644); err != nil {
		return err
	}
	return nil
}

func (c Client) GetCharts(repoName, name string) ([]*search.Result, error) {
	charts, err := c.ListCharts(repoName, name)
	if err != nil {
		return nil, err
	}
	var result []*search.Result
	for _, chart := range charts {
		if chart.Chart.Name == name {
			result = append(result, chart)
		}
	}
	return result, nil
}

func (c Client) GetLastCharts(repoName, name string) ([]*search.Result, error) {
	res, err := c.GetCharts(repoName, name)
	if err != nil {
		return nil, err
	}
	// stable >0.0.0
	// alpha, beta, and release candidate releases >0.0.0-0
	constraint, err := semver.NewConstraint(">0.0.0-0")
	if err != nil {
		return res, errors.Wrap(err, "an invalid version/constraint format")
	}
	data := res[:0]
	foundNames := map[string]bool{}
	for _, r := range res {
		if foundNames[r.Name] {
			continue
		}
		v, err := semver.NewVersion(r.Chart.Version)
		if err != nil {
			continue
		}
		if constraint.Check(v) {
			data = append(data, r)
			foundNames[r.Name] = true
		}
	}
	return data, nil
}

func (c Client) Uninstall(name string) (*release.UninstallReleaseResponse, error) {
	client := action.NewUninstall(c.actionConfig)
	res, err := client.Run(name)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("uninstall tool %s failed: %v", name, err))
	}
	return res, nil
}
