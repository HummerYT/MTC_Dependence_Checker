package main

import (
	"fmt"
	"os"
	"strings"

	"MTC_task/internal/checker"
	gitclone "MTC_task/internal/git"
	"MTC_task/internal/parser"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "depcheck <git-repo-url>",
		Short: "Анализирует Go-зависимости указанного Git-репозитория",
		Args:  cobra.ExactArgs(1),
		RunE:  run,
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	repoURL := args[0]

	// 1. Клонируем репозиторий
	repoDir, err := gitclone.CloneRepo(repoURL)
	if err != nil {
		return err
	}
	defer os.RemoveAll(repoDir) // чистим за собой

	// 2. Парсим go.mod
	info, err := parser.ParseGoMod(repoDir)
	if err != nil {
		return err
	}

	// 3. Выводим основную информацию
	fmt.Println("\n========== Go Module Info ==========")
	fmt.Printf("  Модуль:     %s\n", info.ModuleName)
	fmt.Printf("  Версия Go:  %s\n", info.GoVersion)
	fmt.Printf("  Зависимостей (прямых): %d\n", len(info.Requires))

	// 4. Проверяем доступные обновления
	fmt.Println("\nПроверка обновлений зависимостей...")
	updates, err := checker.CheckUpdates(info.Requires)
	if err != nil {
		return err
	}

	// 5. Выводим результат
	fmt.Println("\n========== Доступные обновления ==========")
	if len(updates) == 0 {
		fmt.Println("  Все зависимости актуальны!")
	} else {
		fmt.Printf("  %-50s %-20s %s\n", "Пакет", "Текущая", "Последняя")
		fmt.Println("  " + strings.Repeat("-", 80))
		for _, u := range updates {
			fmt.Printf("  %-50s %-20s %s\n", u.Path, u.CurrentVersion, u.LatestVersion)
		}
	}

	return nil
}
