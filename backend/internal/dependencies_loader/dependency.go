package dependenciesloader

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func (l *Loader) FetchDependencyDetails(apiUrl string) error {
	return nil
}
