package checker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"MTC_task/internal/parser"
)

// UpdateAvailable описывает зависимость с доступным обновлением
type UpdateAvailable struct {
	Path           string
	CurrentVersion string
	LatestVersion  string
}

var httpClient = &http.Client{Timeout: 10 * time.Second}

// CheckUpdates проверяет каждую зависимость через proxy.golang.org
func CheckUpdates(deps []parser.Dependency) ([]UpdateAvailable, error) {
	var updates []UpdateAvailable

	for _, dep := range deps {
		latest, err := fetchLatestVersion(dep.Path)
		if err != nil {
			// не прерываем — просто пропускаем эту зависимость
			fmt.Printf("  [warn] не удалось проверить %s: %v\n", dep.Path, err)
			continue
		}

		if latest != "" && latest != dep.Version && isNewer(latest, dep.Version) {
			updates = append(updates, UpdateAvailable{
				Path:           dep.Path,
				CurrentVersion: dep.Version,
				LatestVersion:  latest,
			})
		}
	}

	return updates, nil
}

// fetchLatestVersion запрашивает последнюю версию модуля через GOPROXY
func fetchLatestVersion(modulePath string) (string, error) {
	// proxy.golang.org возвращает JSON с полем Version
	url := fmt.Sprintf("https://proxy.golang.org/%s/@latest", modulePath)

	resp, err := httpClient.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	var result struct {
		Version string `json:"Version"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Version, nil
}

// isNewer — грубое сравнение: latest != current и не содержит "-" (не pre-release)
// Для продакшена лучше использовать semver-библиотеку
func isNewer(latest, current string) bool {
	// убираем "v" префикс
	l := strings.TrimPrefix(latest, "v")
	c := strings.TrimPrefix(current, "v")
	return l != c
}
