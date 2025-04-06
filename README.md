# TTK-Organazier (backend)
## Architecture
```
├── cmd 
│    └── main.go
├── config    
│    └── local.yaml
├── internal
|    ├── controller
|    |     └── config.go
│    ├── controller
│    │    ├── article_controller.go
│    │    ├── user_controller.go
│    │    ├── task_controller.go
│    │    └── history_controller.go
|    ├── lib
|    |    └── jwt.go
|    ├── middleware
|    |    └── auth_middleware.go
|    |
│    ├── domain
│    │    ├── article.go
│    │    ├── user.go
│    │    ├── task.go
│    │    └── history.go
│    └── usecase
│         ├── article
│         ├── user
│         ├── task
│         └── history
└── storage
     └── prisma
            └── storage.go
```
### Основные компоненты:
- **/cmd**: Главное приложение и точка входа
- **/config**: Конфигурационные файлы
- **/internal**:
  - **controller**: Обработчики HTTP-запросов
  - **lib**: Общие утилиты и библиотеки
  - **middleware**: Промежуточное ПО для HTTP
  - **domain**: Бизнес-модели и структуры данных с интерфейсами
  - **usecase**: Слой бизнес-логики
- **/storage**: Реализация хранилища данных

Проект следует Clean Architecture и Domain-Driven Design принципам.

## Launch methods

env пример для подключения к supabase
```.env
DATABASE_URL="postgresql://postgres.nacuemduatmefatbmnoq:[YOUR-PASSWORD]@aws-0-eu-central-1.pooler.supabase.com:6543/postgres?pgbouncer=true"

DIRECT_URL="postgresql://postgres.nacuemduatmefatbmnoq:[YOUR-PASSWORD]@aws-0-eu-central-1.pooler.supabase.com:5432/postgres"
```

config .yaml
```yaml
env: "local"
token_ttl: 1h
app_secret: "YOUR_JWT_SECRET"
```

Скрипт запуска проекта
```bash
go run cmd/main.go --config=./config/local.yaml
```

# Run with docker
Стяните образ с DockerHub`a и запустите контейнер
```bash
docker pull c0dys/ttk-back:latest
sudo docker run -e CONFIG_PATH=/app/config/local.yaml -p 8080:8080  ttk-back
```