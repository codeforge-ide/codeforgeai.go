package githubmodels

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/codeforge-ide/codeforgeai.go/config"
)

const githubModelsCatalogURL = "https://models.github.ai/catalog/models"

type ModelCatalogEntry struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	// Add more fields as needed
}

type ModelCatalogResponse struct {
	Models []ModelCatalogEntry `json:"models"`
}

// FetchModelCatalog fetches the model catalog from GitHub Models API.
func FetchModelCatalog(token string) ([]ModelCatalogEntry, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("GET", githubModelsCatalogURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("github models catalog error: %s", string(b))
	}
	var catalog ModelCatalogResponse
	if err := json.NewDecoder(resp.Body).Decode(&catalog); err != nil {
		return nil, err
	}
	return catalog.Models, nil
}

// UpdateConfigWithModelList updates the config with the fetched model list.
func UpdateConfigWithModelList(cfg *config.Config, models []ModelCatalogEntry) bool {
	var names []string
	for _, m := range models {
		names = append(names, m.ID)
	}
	joined := strings.Join(names, ",")
	if cfg.GithubModelsList != joined {
		cfg.GithubModelsList = joined
		return true
	}
	return false
}
