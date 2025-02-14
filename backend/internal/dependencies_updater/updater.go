package dependenciesupdater

import (
	"fmt"
	"strings"

	"github.com/wojcikp/deps-dev-assignment/backend/internal/database"
	dependenciesloader "github.com/wojcikp/deps-dev-assignment/backend/internal/dependencies_loader"
)

type Updater struct {
	loader *dependenciesloader.Loader
	db     *database.SQLiteDB
}

func NewDependenciesUpdater(loader *dependenciesloader.Loader, db *database.SQLiteDB) *Updater {
	return &Updater{loader, db}
}

func (u *Updater) UpdateDependencies() ([]string, error) {
	dependenciesToUpdate, err := u.FindDependenciesToUpdate()
	if err != nil {
		return []string{}, fmt.Errorf("update dependencies failed due to an error: %w", err)
	}

	for _, dependency := range dependenciesToUpdate {
		url := "https://api.deps.dev/v3/projects/" + strings.ReplaceAll(dependency, "/", "%2F")
		newDetails, err := u.loader.FetchDependencyDetails(url)
		if err != nil {
			return []string{}, fmt.Errorf("update dependencies failed due to an error: %w", err)
		}
		if err = u.db.UpdateDependencyDetails(newDetails); err != nil {
			return []string{}, fmt.Errorf("update dependencies failed due to an error: %w", err)
		}
	}

	for _, node := range u.loader.Dependencies.Nodes {
		if err := u.db.UpdateVersionKeys(node.VersionKey.Name, node.VersionKey.Version); err != nil {
			return []string{}, fmt.Errorf("update dependencies failed due to an error: %w", err)
		}
	}

	return dependenciesToUpdate, nil
}

func (u *Updater) FindDependenciesToUpdate() ([]string, error) {
	dependenciesToUpdate := []string{}
	dbDependenciesVersions, err := u.db.GetVersionKeys()
	if err != nil {
		return []string{}, err
	}
	err = u.loader.FetchDepsDevDependencies()
	if err != nil {
		return []string{}, err
	}

	for _, dbDependency := range dbDependenciesVersions {
		if u.checkVersion(dbDependency) {
			dependenciesToUpdate = append(dependenciesToUpdate, dbDependency.Name)
		}
	}

	return dependenciesToUpdate, err
}

func (u *Updater) checkVersion(dbDependency dependenciesloader.VersionKey) bool {
	for _, node := range u.loader.Dependencies.Nodes {
		if dbDependency.Name == node.VersionKey.Name {
			return dbDependency.Version != node.VersionKey.Version
		}
	}
	return false
}
