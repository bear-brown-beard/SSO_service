# SSO Service

Современный сервис единого входа (Single Sign-On) для веб-приложений, построенный на Go с использованием Clean Architecture и современных практик разработки.

## 🚀 Возможности

- **Аутентификация по SMS** - отправка и валидация кодов подтверждения
- **JWT токены** - безопасная аутентификация с поддержкой refresh токенов
- **Интеграция с Mindbox** - автоматическая регистрация и логин пользователей
- **Кеширование** - Redis для хранения кодов и токенов
- **Тестовые аккаунты** - поддержка тестовых номеров для разработки
- **REST API** - полноценный HTTP API для интеграции
- **Docker поддержка** - готовые контейнеры для развертывания

## 🛠️ Технологии

- **Go 1.23** - основной язык разработки
- **Chi Router** - HTTP роутер
- **QB Package** - SQL query builder для работы с БД
- **Redis** - кеширование и хранение токенов
- **MySQL** - основная база данных
- **JWT** - JSON Web Tokens для аутентификации
- **Docker** - контейнеризация
- **Structured Logging** - структурированное логирование

## 📦 Установка и запуск

1. **Клонировать репозиторий**
```bash
git clone https://github.com/bear-brown-beard/SSO_service.git
cd SSO_service
```

2. **Настроить переменные окружения**
```bash
cp example.env .env

3. **Запустить с Docker**
```bash
docker-compose up -d
```

4. **Или запустить локально**
```bash
go mod tidy
go run cmd/sso/main.go
```

## 🔒 Безопасность

- **Валидация телефонов** - проверка формата и нормализация
- **Временные коды** - коды верификации с TTL
- **JWT токены** - безопасная аутентификация
- **Черный список токенов** - блокировка скомпрометированных токенов
- **Rate limiting** - защита от брутфорс атак
- **TLS для внешних API** - защищенные запросы к SMSC и Mindbox

## 🚀 Развертывание

### Docker
```bash
docker build -t sso-service .

docker run -p 8080:8080 sso-service
```

### Docker Compose
```bash
docker-compose up -d
```

### Kubernetes
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: sso-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: sso-service
  template:
    metadata:
      labels:
        app: sso-service
    spec:
      containers:
      - name: sso-service
        image: sso-service:latest
        ports:
        - containerPort: 8080
```

---

**SSO Service** - современное решение для аутентификации пользователей с поддержкой SMS-верификации и интеграцией с внешними сервисами. 