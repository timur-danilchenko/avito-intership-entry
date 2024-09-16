# Стажировка Авито. Тестовое задание

## Задание
Полноe условие задачи лежит [здесь](api/readme.md)

Авито — большая компания, в рамках которой пользователи не только продают/покупают товары и услуги, но и предоставляют помощь крупному бизнесу и предприятиям.

Поэтому ребята из Авито решили сделать сервис, который позволит бизнесу создать тендер на оказание каких-либо услуг. А пользователи/другие бизнесы будут предлагать свои выгодные условия для получения данного тендера.

Помогите ребятам из Авито реализовать новое HTTP API!

## Создание и настройка базы данных
### Установка
Для установки postgres использовал следующую команду, так как предоставленное окружение Ubuntu.
```
sudo apt-get install postgresql-all
```
### Миграции
В рамках разработки были созданы следующие сущности

**Тендер**
```
CREATE TYPE tender_status_type AS ENUM (
    'Created',
    'Published',
    'Closed'
);

CREATE TYPE tender_service_type AS ENUM (
    'Construction',
    'Delivery',
    'Manufacture'
);

CREATE TABLE IF NOT EXISTS tender (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description VARCHAR(500) NOT NULL,
    service_type tender_service_type,
    status tender_status_type DEFAULT 'Created',
    organization_id UUID REFERENCES organization(id) ON DELETE CASCADE,
    version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

**Предложения**
```
CREATE TYPE bid_author_type AS ENUM (
    'Organization',
    'User'
);

CREATE TYPE bid_status_type AS ENUM (
    'Created',
    'Published',
    'Canceled',
    'Approved',
    'Rejected'
);

CREATE TABLE IF NOT EXISTS bid (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description VARCHAR(500) NOT NULL,
    status bid_status_type DEFAULT 'Created',
    tender_id UUID REFERENCES tender(id) ON DELETE CASCADE,
    author_type bid_author_type,
    author_id UUID NOT NULL,
    version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS review (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    bid_id UUID REFERENCES bid(id) ON DELETE CASCADE,
    description TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)
```

## Использование
### Docker
Для удобства в проекте был описан ```Dockerfile```, а также ```docker-compose.yml``` что позволяет автоматически настраить и запустить сервер и базу данных в Docker-контейнерах

Для этого необходимо выполнить следующие инструкции
1. Собрать образ Docker-контейнеров
```
docker-compose build
```
2. Запустить контейнеры
```
docker-compose up
```

После запуска сервера, сервер будет доступен по адресу ```SERVER_ADDRESS``` описаному в ```.env``` файле


### Локально
Перед началом необходимо объявить ```.env``` файл в котором надо добавить такие поля как ```SERVER_ADDRESS```. ```POSTGRES_USERNAME```. ```POSTGRES_PASSWORD```. ```POSTGRES_HOST```. ```POSTGRES_PORT```. ```POSTGRES_DATABASE```. ```SSLMODE```. ```POSTGRES_JDBC_URL```. ```POSTGRES_CONN```

Для удобства есть ```env_template```, в котором описаны необходимые поля

```Makefile``` содержит в себе правила для удобства взаимодействия с сервисом

Список комвнд ```Makefile```
```
all             # запускает процесс в фоновом режиме, предварительно выполнив setup, migrate-up, start
setup           # устанавливает необходимые зависимости проекта
start           # запускает сервер
migrate-up      # выполняет все миграции, до последней
migrate-down    # откатывает миграции на предыдущее состояние
drop-db         # стирание всех данных баззы данных
dbshell         # подключение к оболочке базы данных
```

Пошаговая инструкия запуска

1. Копирование и установка значений в ```.env```
```
cp env_template .env
```
2. Подготовка зависимостей и инструментов 
```
make setup
```
3. Выполнение миграций
```
make migrate-up
```
4. Запуск сервера
```
make start
```

После запуска сервера, сервер будет доступен по адресу ```SERVER_ADDRESS``` описаному в ```.env``` файле
