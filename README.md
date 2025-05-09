# VK Тестовое задание :: PubSub Service

[![Go Version](https://img.shields.io/badge/go-1.20+-00ADD8?style=flat-square&logo=go)](https://golang.org/dl/)
[![GitHub Actions](https://img.shields.io/github/actions/workflow/status/mirotvoretts/vk_test_task/go.yml?style=flat-square&logo=github-actions&label=build)](https://github.com/mirotvoretts/vk_test_task/actions)
[![License](https://img.shields.io/badge/license-MIT-blue?style=flat-square)](LICENSE)

## Как работает сервис

### Основные функции:
1. **Подписка на события** по ключу (Subscribe)
2. **Публикация событий** по ключу (Publish)
3. Graceful Shutdown с завершением активных подключений
4. Буферизация сообщений для асинхронной обработки

### Паттерны:
- Publisher-Subscriber
- Dependency Injection
- Graceful Shutdown

## Зависимости

### Необходимые инструменты:
1. Go 1.20+
2. protoc (protobuf compiler)
3. protoc-gen-go и protoc-gen-go-grpc

```bash
# Установка зависимостей
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## Сборка сервиса и запуск

1. Генерация gRPC кода

```bash
protoc --go_out=./pkg/proto --go-grpc_out=./pkg/proto proto/pubsub.proto
```

2. Сборка сервера

```bash
go build -o bin/server ./cmd/server
```

3. Запускаем

```bash
./bin/server
```

## Тестирование

```bash
# Запуск всех тестов
go test -v ./...

# Тесты конкретного пакета
go test -v ./subpub
go test -v ./server
```

