# Links Checker Service

Веб-сервис для проверки доступности интернет-ресурсов с генерацией отчетов в PDF.

## Архитектура

Проект построен по принципам чистой архитектуры и SOLID:

### Структура пакетов

- `config/` - загрузка конфигурационного файла
- `model/` - доменные модели и интерфейсы
- `repository/` - слой данных (репозиторий)
- `service/` - бизнес-логика
- `handler/` - HTTP обработчики
- `pdf/` - генерация PDF отчетов


### Паттерны

1. **Repository Pattern** - абстракция над хранилищем данных
2. **Service Layer** - инкапсуляция бизнес-логики
3. **Worker Pool** - обработка задач в параллельных воркерах
4. **Semaphore** - ограничение в кол-ве одновременно выполняющихся горутин
5. **Graceful Shutdown** - корректная обработка остановки сервера


## Конфигурация

Конфигурационный файл local.yaml находится в папке config/ 

```yaml
http_server:
  port: "8080"

storage:
  filepath: "./storage/data.json"

app:
  worker_count: 2
```

## Запуск

```bash
go mod tidy
go run cmd/server/main.go
```

## API Endpoints

### POST /check
Проверка доступности ссылок

**Request:**
```json
{
    "links": ["https://mydrop.io", "https://digimetr.com"]
}
```

### POST /report
Генерация pdf файла

**Request:**
```json
{
    "links_num": [1, 2]
}
```

## Тестирование

```bash
curl -X POST http://localhost:8080/check \
  -H "Content-Type: application/json" \
  -d '{"links": ["https://mydrop.io", "https://digimetr.com"]]}'
```

```bash
curl -X POST http://localhost:8080/report \
  -H "Content-Type: application/json" \
  -d '{"links_num": [1,2]}' \
  --output report.pdf
```