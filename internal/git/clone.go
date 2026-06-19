package git

import (
	"fmt"
	"os"

	gogit "github.com/go-git/go-git/v5"
)

// CloneRepo клонирует репозиторий во временную папку
// и возвращает путь к ней. Вызывающий код обязан удалить папку после использования.
func CloneRepo(url string) (string, error) {
	tmpDir, err := os.MkdirTemp("", "mtc-repo-*")
	if err != nil {
		return "", fmt.Errorf("не удалось создать временную папку: %w", err)
	}

	fmt.Printf("Клонирование %s ...\n", url)

	_, err = gogit.PlainClone(tmpDir, false, &gogit.CloneOptions{
		URL:      url,
		Depth:    1, // shallow clone — нам нужен только go.mod
		Progress: os.Stdout,
	})
	if err != nil {
		os.RemoveAll(tmpDir)
		return "", fmt.Errorf("ошибка клонирования: %w", err)
	}

	return tmpDir, nil
}
