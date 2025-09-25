# Captcha Service

Современный сервис капчи с интерактивным слайдер-пазлом, построенный на Go и gRPC с поддержкой gRPC-Gateway для объединения HTTP и gRPC на одном порту.

## 🚀 Быстрый старт

### 🐳 Запуск через Docker Compose (ПРИОРИТЕТНЫЙ)

**Рекомендуемый способ запуска** - все сервисы автоматически настроены и готовы к работе:

```bash
# 1. Клонирование репозитория
git clone https://github.com/ristoch/captcha.git
cd captcha

# 2. Запуск всех сервисов одной командой
docker-compose up -d

# 3. Проверка статуса
docker-compose ps

# 4. Просмотр логов (все сервисы)
docker-compose logs -f

# 5. Просмотр логов конкретного сервиса
docker-compose logs -f captcha-service-1
docker-compose logs -f balancer
docker-compose logs -f balancer-proxy

# 6. Остановка всех сервисов
docker-compose down

# 7. Полная очистка (включая volumes)
docker-compose down --volumes --remove-orphans
```

**Что запускается автоматически:**
- ✅ **3 экземпляра captcha-service** (порты 38000-38002) с gRPC-Gateway
- ✅ **Balancer** (HTTP: 8080, gRPC: 9090) с автоматической регистрацией
- ✅ **Balancer-proxy** (порт 8081) для балансировки нагрузки
- ✅ **Demo-приложение** (порт 8082) для тестирования

### 🔧 Запуск вручную (для разработки)

**Требования:** Go 1.24+, protobuf, все зависимости

```bash
# 1. Установка зависимостей
go mod download

# 2. Генерация protobuf файлов
protoc --go_out=gen/proto --go_opt=paths=source_relative \
       --go-grpc_out=gen/proto --go-grpc_opt=paths=source_relative \
       proto/captcha/captcha.proto proto/balancer/balancer.proto

# 3. Запуск сервисов (в отдельных терминалах)
go run cmd/balancer/main.go          # Балансер
go run cmd/server/main.go            # Captcha-service
go run cmd/balancer-proxy/main.go    # Прокси-балансер  
go run cmd/demo/main.go              # Демо
```

## ✅ Быстрая проверка работоспособности

После запуска `docker-compose up -d` выполните:

```bash
# 1. Проверка всех сервисов (должно показать 3 captcha-service)
curl http://localhost:8080/api/services | jq .

# 2. Проверка демо-страницы
open http://localhost:8082

# 3. Проверка здоровья всех сервисов
curl http://localhost:8082/health          # Демо
curl http://localhost:8081/api/health      # Прокси-балансер
curl http://localhost:8080/health          # Балансер
curl http://localhost:38000/health         # Captcha-service 1
curl http://localhost:38001/health         # Captcha-service 2
curl http://localhost:38002/health         # Captcha-service 3
```

**Ожидаемый результат:**
- Все health-чекки должны вернуть `{"status":"ok"}`
- `/api/services` должен показать 3 зарегистрированных сервиса с уникальными UUID
- Демо-страница должна открыться и показать интерактивную капчу

## 🌐 Доступные URL

После запуска сервисы будут доступны по следующим адресам:

- **Демо-страница**: http://localhost:8082
- **Прокси-балансер**: http://localhost:8081
- **Балансер**: http://localhost:8080 (HTTP) / localhost:9090 (gRPC)
- **Сервисы капчи (gRPC-Gateway)**: 
  - http://localhost:38000 (экземпляр 1) - HTTP + gRPC на одном порту
  - http://localhost:38001 (экземпляр 2) - HTTP + gRPC на одном порту
  - http://localhost:38002 (экземпляр 3) - HTTP + gRPC на одном порту

## 🎮 Как это работает

1. **Пользователь** открывает демо-страницу
2. **Система** генерирует уникальную капчу-слайдер
3. **Пользователь** перемещает слайдер для решения
4. **Сервис** проверяет правильность в реальном времени

## 🛠 Технические детали

- **Язык**: Go 1.24+
- **Протокол**: gRPC + WebSocket + gRPC-Gateway
- **Производительность**: 100+ RPS генерации заданий
- **Память**: Оптимизированная для 10k+ активных задач
- **Интеграция**: Автоматическая регистрация в балансере с UUID
- **Архитектура**: Микросервисная с балансировкой нагрузки
- **Порты**: Динамическое выделение в диапазоне 38000-40000

## 📁 Структура проекта

```
├── cmd/                    # Точки входа сервисов
├── internal/              # Внутренняя логика
│   ├── service/           # Бизнес-логика
│   ├── transport/         # HTTP/gRPC обработчики
│   └── infrastructure/    # Инфраструктурные компоненты
├── templates/             # HTML шаблоны капчи
├── proto/                 # gRPC протоколы
└── tests/                 # Тесты производительности
```

## 🔧 Конфигурация

### Docker-конфигурация (основная)
Все настройки уже настроены в `docker-compose.yml`:

```yaml
# Основные порты
- "8080:8080"    # Balancer HTTP
- "8081:8081"    # Balancer-proxy
- "8082:8082"    # Demo
- "38000:38000"  # Captcha-service 1
- "38001:38001"  # Captcha-service 2  
- "38002:38002"  # Captcha-service 3

# Переменные окружения
PORT=38000-38002  # Порты для captcha-service
BALANCER_ADDRESS=balancer:9090  # Адрес балансера
```

### Локальная конфигурация
Основные настройки через переменные окружения:

```bash
# Порты
MIN_PORT=38000
MAX_PORT=40000

# Производительность
MAX_CHALLENGES=10000
COMPLEXITY_MEDIUM=50

# Безопасность
MAX_ATTEMPTS=3
BLOCK_DURATION_MINUTES=5
```

### Docker-отладка
```bash
# Вход в контейнер для отладки
docker-compose exec captcha-service-1 sh
docker-compose exec balancer sh

# Просмотр переменных окружения
docker-compose exec captcha-service-1 env

# Проверка сетевых подключений
docker network ls
docker network inspect captcha_captcha-network
```

## 🧪 Тестирование

### Docker-тестирование (рекомендуется)
```bash
# 1. Запуск тестов производительности в контейнере
docker-compose exec captcha-service-1 go test ./tests/integration/... -v

# 2. Проверка логов всех сервисов
docker-compose logs --tail=50

# 3. Перезапуск конкретного сервиса
docker-compose restart captcha-service-1

# 4. Пересборка и перезапуск всех сервисов
docker-compose build --no-cache
docker-compose up -d

# 5. Проверка использования ресурсов
docker stats
```

### Локальное тестирование
```bash
# Запуск тестов производительности
go test ./tests/integration/... -v

# Проверка здоровья сервисов
curl http://localhost:8082/health          # Демо
curl http://localhost:8081/api/health      # Прокси-балансер
curl http://localhost:8080/health          # Балансер
curl http://localhost:38000/health         # Сервис капчи 1

# Проверка списка сервисов
curl http://localhost:8080/api/services    # Все зарегистрированные сервисы
```

## 📊 Мониторинг

### Доступные эндпоинты:

**Демо (порт 8082):**
- `GET /` - главная страница (редирект на /demo)
- `GET /demo` - демо-страница с капчей
- `GET /health` - статус демо-сервиса
- `WebSocket /ws` - WebSocket для демо

**Прокси-балансер (порт 8081):**
- `GET /api/health` - статус прокси
- `GET /api/memory` - метрики памяти
- `GET /api/stats` - общая статистика
- `POST /api/services/add` - добавить сервис
- `DELETE /api/services/remove` - удалить сервис
- `POST /api/challenge` - создать капчу
- `POST /api/validate` - проверить решение
- `WebSocket /ws` - события в реальном времени

**Балансер (порт 8080):**
- `GET /health` - статус балансера
- `GET /api/health` - статус балансера (альтернативный)
- `GET /api/services` - список всех зарегистрированных сервисов

**Сервисы капчи (порты 38000-38002, gRPC-Gateway):**
- `GET /health` - статус сервиса
- `GET /memory` - метрики памяти
- `GET /stats` - статистика сервиса
- `POST /api/challenge` - создать капчу (HTTP)
- `POST /api/validate` - проверить решение (HTTP)
- `WebSocket /ws` - события в реальном времени
- **gRPC**: `NewChallenge`, `ValidateChallenge`, `MakeEventStream`

**Логи**: `logs/` директория

## 🔒 Безопасность

- Автоматическая блокировка при превышении попыток
- Graceful shutdown с сохранением состояния и корректной остановкой сервисов
- Бинарная упаковка событий для экономии трафика
- Валидация всех входящих данных
- Уникальные UUID для идентификации инстансов
- Изоляция сервисов через Docker контейнеры

## 🆕 Последние обновления

### v2.0.0 - gRPC-Gateway интеграция
- ✅ **gRPC-Gateway**: Объединение HTTP и gRPC на одном порту для каждого сервиса
- ✅ **UUID идентификация**: Замена временных ID на уникальные UUID
- ✅ **Корректная остановка**: Реализация метода `Stop()` с очисткой ресурсов
- ✅ **Три инстанса**: Гарантированная регистрация всех captcha-service
- ✅ **Оптимизация портов**: Использование диапазона 38000-40000

### Архитектурные улучшения
- Микросервисная архитектура с балансировкой нагрузки
- Автоматическая регистрация и обнаружение сервисов
- Graceful shutdown с сохранением состояния
- Мониторинг и метрики в реальном времени

## 📝 API

### Создание капчи
```bash
curl -X POST http://localhost:8081/api/challenge \
  -H "Content-Type: application/json" \
  -d '{"complexity": 50, "user_id": "test-user"}'
```

### Проверка решения
```bash
curl -X POST http://localhost:8081/api/validate \
  -H "Content-Type: application/json" \
  -d '{"challenge_id": "slider_123", "answer": {"x": 100, "y": 50}}'
```

### Мониторинг памяти
```bash
curl http://localhost:8081/api/memory
```

### Статистика сервисов
```bash
curl http://localhost:8081/api/stats
```

### Проверка gRPC-Gateway сервисов
```bash
# Проверка здоровья captcha-service через gRPC-Gateway
curl http://localhost:38000/health
curl http://localhost:38001/health
curl http://localhost:38002/health

# Создание капчи через gRPC-Gateway
curl -X POST http://localhost:38000/api/challenge \
  -H "Content-Type: application/json" \
  -d '{"complexity": 50, "user_id": "test-user"}'

# Проверка всех зарегистрированных сервисов
curl http://localhost:8080/api/services | jq .
```

### Мониторинг системы
```bash
# Проверка всех сервисов
curl http://localhost:8082/health          # Демо
curl http://localhost:8081/api/health      # Прокси-балансер  
curl http://localhost:8080/health          # Балансер
curl http://localhost:38000/health         # Captcha-service 1
curl http://localhost:38001/health         # Captcha-service 2
curl http://localhost:38002/health         # Captcha-service 3

# Метрики памяти
curl http://localhost:38000/memory
curl http://localhost:38001/memory
curl http://localhost:38002/memory
```

---
