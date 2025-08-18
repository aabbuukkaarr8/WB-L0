## L0 Orders Service

Микросервис для чтения/обработки заказов (orders):
- REST: GET /order/{order_uid}
- Kafka consumer: прием заказов из топика и сохранение в БД
- PostgreSQL с миграциями и пулом подключений
- Кэш LRU+TTL на уровне сервиса
- Retry/Backoff для БД и Kafka
- Метрики Prometheus (/metrics)
- OpenAPI (Swagger) в docs/openapi.yaml

### Архитектура
- internal/api: Gin-обработчики и роутинг
- internal/service: бизнес-логика, кэш, транзакции, retry
- internal/repository: доступ к БД (orders, delivery, payment, items)
- internal/kafka: consumer group на Sarama
- internal/db: обертка над database/sql, конфиг пула
- internal/config: загрузка config.yaml
- pkg/retry: экспоненциальный backoff с jitter
- pkg/metrics: метрики Prometheus
- pkg/validator: валидация входных данных (go-playground/validator)

### Требования
- Docker + Docker Compose (для Postgres/Kafka)
- Go 1.24+

### Быстрый старт
1) Поднять инфраструктуру:
```bash
docker compose up -d postgres zookeeper kafka kafka-ui
```

2) Запуск приложения:
```bash
go run cmd/L0/main.go -config config.yaml
```

3) Проверка здоровья:
```bash
curl -s http://localhost:8081/healthz
```

4) Просмотр метрик:
```bash
curl -s http://localhost:8081/metrics
```

### Конфигурация (config.yaml)
```yaml
logLevel:
  log_level: debug
server:
  port: :8081
db:
  host: localhost
  port: 5432
  login: orders_user
  password: orders_pass
  dbname: orders_db
  sslmode: disable
  max_open_conns: 25
  max_idle_conns: 10
  conn_max_lifetime: 30m
  conn_max_idle_time: 5m
kafka:
  broker: localhost:29092
  topic: orders
  group_id: order-group
```

### API
- GET /order/{order_uid}
  - 200: объект заказа
  - 404: {"error": "..."}
- OpenAPI: `docs/openapi.yaml` (откройте в Swagger Editor или подключите Swagger UI)

Пример:
```bash
curl -s http://localhost:8081/order/b563feb7b2b84b6test
```

### Kafka
- Брокер: `localhost:29092`
- Топик: `orders`
- Отправка сообщения (скрипт):
```bash
bash script/kafka-produce/produce_message.sh
```
- Сообщение валидируется (`pkg/validator`) и сохраняется в БД. Ошибки парсинга/валидации/сохранения логируются, потребление продолжится (с backoff).

### Кэш LRU+TTL
- Реализован в `internal/service/orders.Service`: ограничения по размеру и времени.
- `Get` сначала смотрит в кэш, затем в БД; `SaveOrder` кладет в кэш.

### Retry/Backoff
- БД: `pkg/retry` с экспоненциальным backoff + jitter, повторяются транзиентные ошибки Postgres.
- Kafka: при ошибке `Consume` — экспоненциальный backoff, метрика ошибок растет.

### Метрики
- Экспонируются на `GET /metrics` (Prometheus формат)
- Примеры метрик:
  - `orders_kafka_consume_errors_total`
  - `orders_db_retry_attempts_total`

### Миграции БД
- Миграции находятся в `db/migrations` и применяются через Docker (initdb) при первом старте Postgres из compose.
- Быстрая очистка (альтернатива): TRUNCATE таблиц с RESTART IDENTITY, если требуется только почистить данные.

### Тесты
```bash
go test ./...
```
- Repository: per-method тесты (`internal/repository/orders/*_test.go`)
- Service: `get_test.go`, `save_test.go`
- Handler: `get_test.go`

### Разработка
- Добавление новых валидаторов: `pkg/validator/custom.go`
- Расширение метрик: `pkg/metrics`
- Настройка пула БД: поля в `db` секции `config.yaml`

### Полезные ссылки
- Gin: https://github.com/gin-gonic/gin
- Sarama (Kafka): https://github.com/Shopify/sarama
- Prometheus client: https://github.com/prometheus/client_golang
- Swagger/OpenAPI: https://swagger.io/specification/
