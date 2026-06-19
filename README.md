# Go Dependency Analyzer

CLI для анализа зависимостей Go-репозитория.

## Входные данные
URL публичного Git-репозитория с go.mod

## Выходные данные
- Имя модуля
- Версия Go
- Список прямых зависимостей доступных для обновления

## Запуск

```bash
go run main.go <git-repo-url>
```

## Пример запуска

```bash
go run main.go https://github.com/HummerYT/Avito
```

## Сборка

```bash
go build -o depcheck .
./depcheck https://github.com/HummerYT/Avito
```
