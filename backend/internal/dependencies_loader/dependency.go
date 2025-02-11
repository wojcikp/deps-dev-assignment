package dependenciesloader

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type Loader struct {
	repositoryUrl string
	Dependencies  Dependencies
}

func NewDependenciesLoader(repositoryUrl string) *Loader {
	return &Loader{repositoryUrl: repositoryUrl}
}

func (l *Loader) FetchDepsDevDependencies() error {
	resp, err := http.Get(l.repositoryUrl)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-OK HTTP status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(body, &l.Dependencies); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	return nil
}

func (l *Loader) FetchDetailsForAllDependencies() []DependencyDetails {
	detailedDependencies := []DependencyDetails{}
	for _, dependency := range l.Dependencies.Nodes {
		url := "https://api.deps.dev/v3/projects/" + strings.ReplaceAll(dependency.VersionKey.Name, "/", "%2F")
		dependencyDetails, err := l.FetchDependencyDetails(url)
		if err != nil {
			log.Printf("failed to fetch details for dependency: %s due to an error: %v", dependency.VersionKey.Name, err)
			continue
		}
		detailedDependencies = append(detailedDependencies, dependencyDetails)
	}

	return detailedDependencies
}

func (l *Loader) FetchDependencyDetails(apiUrl string) (DependencyDetails, error) {
	resp, err := http.Get(apiUrl)
	if err != nil {
		return DependencyDetails{}, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return DependencyDetails{}, fmt.Errorf("received non-OK HTTP status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return DependencyDetails{}, fmt.Errorf("failed to read response body: %w", err)
	}

	var details DependencyDetails
	if err := json.Unmarshal(body, &details); err != nil {
		return DependencyDetails{}, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return details, nil
}
