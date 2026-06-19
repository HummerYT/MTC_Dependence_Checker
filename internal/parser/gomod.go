package parser

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/mod/modfile"
)

// ModInfo содержит данные из go.mod
type ModInfo struct {
	ModuleName string
	GoVersion  string
	Requires   []Dependency
}

type Dependency struct {
	Path    string
	Version string
}

// ParseGoMod читает и парсит go.mod из указанной папки репозитория
func ParseGoMod(repoDir string) (*ModInfo, error) {
	gomodPath := filepath.Join(repoDir, "go.mod")

	data, err := os.ReadFile(gomodPath)
	if err != nil {
		return nil, fmt.Errorf("go.mod не найден в %s: %w", repoDir, err)
	}

	f, err := modfile.Parse(gomodPath, data, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга go.mod: %w", err)
	}

	info := &ModInfo{
		ModuleName: f.Module.Mod.Path,
		GoVersion:  f.Go.Version,
	}

	for _, req := range f.Require {
		if req.Indirect {
			continue // пропускаем косвенные зависимости для чистоты вывода
		}
		info.Requires = append(info.Requires, Dependency{
			Path:    req.Mod.Path,
			Version: req.Mod.Version,
		})
	}

	return info, nil
}
