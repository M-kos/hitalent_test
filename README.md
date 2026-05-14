# API для управления организационной структурой компании: подразделения и сотрудники.

## Описание

Проект реализует REST API для работы с иерархической структурой подразделений и сотрудниками. Поддерживает создание, чтение, обновление и удаление подразделений и сотрудников, включая сложные операции вроде перемещения подразделений и каскадного удаления.

### Возможности
- Создание и редактирование подразделений с древовидной структурой
- Добавление сотрудников в подразделения
- Получение подразделения с иерархией и сотрудниками
- Перемещение подразделений (с проверкой циклов)
- Удаление подразделений с режимами `cascade` и `reassign`

## Запуск проекта

1. Установите [Docker](https://www.docker.com/products/docker-desktop) и [docker-compose](https://docs.docker.com/compose/install/)
2. Склонируйте репозиторий:
   ```bash
   git clone https://github.com/your-username/hitalent_test.git
   cd hitalent_test
   ```
3. Создайте файл `.env` в корне проекта с переменными окружения:
   ```env
   POSTGRES_USER=postgres
   POSTGRES_PASSWORD=postgres
   POSTGRES_NAME=department_db
   POSTGRES_PORT=5432
   POSTGRES_HOST=postgres
   SERVICE_PORT=8080
   LOG_LEVEL=debug
   GOOSE_DRIVER=postgres
   GOOSE_DBSTRING=host=postgres port=5432 user=postgres password=postgres dbname=department_db sslmode=disable
   GOOSE_TABLE=goose_db_version
   ```
4. Доступные команды (Taskfile)
    ```text
    Проект использует [Task](https://taskfile.dev/) для управления задачами. Доступные команды:
    
    Docker
    - `task docker:up` — поднять контейнеры
    - `task docker:up:local` — поднять только БД и миграции
    - `task docker:down` — остановить контейнеры
    
    Запуск сервиса
    - `task department:run` — запустить сервис локально
    - `task department:run:local:hot` — запустить с горячей перезагрузкой
    - `task department:build` — собрать бинарник
    
    Тесты
    - `task department:test` — запустить все тесты
    
    Миграции
    - `task migration:create <name>` — создать новую миграцию
    - `task migration:up` — применить миграции
    - `task migration:down` — откатить последнюю миграцию
    
    Линтинг и форматирование
    - `task lint` — запустить линтер
    - `task fmt` — отформатировать код
    ```
## API

### Создать подразделение
- **Метод:** `POST`
- **Эндпоинт:** `/departments/`
- **Тело запроса:**
  ```json
  {
    "name": "Engineering",
    "parent_id": 1
  }
  ```

### Создать сотрудника
- **Метод:** `POST`
- **Эндпоинт:** `/departments/{id}/employees/`
- **Тело запроса:**
  ```json
  {
    "full_name": "John Doe",
    "position": "Developer",
    "hired_at": "2023-01-01T00:00:00Z"
  }
  ```

### Получить подразделение
- **Метод:** `GET`
- **Эндпоинт:** `/departments/{id}`
- **Параметры:**
  - `depth` — глубина вложенности (1–5, по умолчанию 1)
  - `include_employees` — включать ли сотрудников (true/false, по умолчанию true)

### Обновить подразделение
- **Метод:** `PATCH`
- **Эндпоинт:** `/departments/{id}`
- **Тело запроса:**
  ```json
  {
    "name": "Engineering",
    "parent_id": 2
  }
  ```

### Удалить подразделение
- **Метод:** `DELETE`
- **Эндпоинт:** `/departments/{id}`
- **Параметры:**
  - `mode` — `cascade` или `reassign`
  - `reassign_to_department_id` — обязателен при `mode=reassign`

## Архитектура

Проект разделён на слои:
- `internal/domain` — модели
- `internal/repository` — работа с БД
- `internal/service` — бизнес-логика
- `internal/handler` — HTTP-обработчики
- `pkg/logger` — логгирование
