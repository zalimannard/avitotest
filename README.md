# Тестовое задание Avito

- СУБД: PostgreSQL
- Библиотеки: cleanenv, chi, playground/validator/v10

# Запуск

## 1. Подготовка

Нужно создать сеть Docker и Volume для отчётов:

```shell
docker network create avitonetwork
docker volume create reports-volume
```

## 2. База данных

Можно запустить базу данных любым удобным образом, например, через docker-compose:

```yml
version: "3"

services:
  database:
    image: postgres:latest
    container_name: avitotest-db
    ports:
      - "5432:5432"
    volumes:
      - /var/lib/postgresql/data/
    environment:
      - POSTGRES_DB=avitotest-db
      - POSTGRES_USER=user
      - POSTGRES_PASSWORD=password
    restart: unless-stopped
    networks:
      - avitonetwork

networks:
  avitonetwork:
    external: true
```

Запускать командой

```shell
docker-compose up -d
```

Создать таблицы базы данных по приложенному скрипту. Схема данных:

![2023 08 31 23:38:55](https://github.com/zalimannard/avitotest/assets/90920120/5e3db849-bc70-4b1d-b58f-dd5197b08efe)

## 3. Приложение

```yml
version: '3'
services:
  backend:
    container_name: user-segmentation
    image: zalimannard/user-segmentation
    ports:
      - "8080:8080"
    environment:
      - ENV=dev
      - DB_URL=postgresql://user:password@database:5432/avitotest-db?sslmode=disable
      - HTTP_SERVER_URL=0.0.0.0:8080
      - REPORT_DIR=./reports
    restart: unless-stopped
    networks:
      - avitonetwork
    volumes:
      - ./reports/:/reports

networks:
  avitonetwork:
    external: true
```

Переменные среды:

- ENV - Тип запуска. Сделан только dev
- DB_URL - Путь к базе данных
- HTTP_SERVER_URL - По какому адресу находится сервис. Нужен для формирования ссылок на отчёты
- REPORT_DIR - Используемая папка для отчётов. Должна совпадать с созданной volume для сохранности отчётов

Можно пересобрать контейнер командой

```shell
docker-compose build
```

# Вопросы

---

Q: Удаление сегмента должно быть несмотря на наличие у пользователей?

A: Нет, каскадное удаление зло. В таком случае сначала пусть удалят сегмент у всех

---

Q: Реализовывать ли эндпоинты для пользователей?

A: Нет, идейно наше приложение отвечает только за сегменты. За пользователей отвечает сторонний сервис

---

Q: Использовать ли slug как первичные ключи?

A: Нет, они будут занимать много места и создадут проблемы в случае упомянутого в задании изменения slug. Делаем искуственные ключи

---

# Запросы и ответы

Создан Swagger-файл с подробным описанием. Здесь короткая версия успешных запросов

## Сегменты

### Создание сегмента

POST http://localhost:8080/api/segments

Request:
```json
{
    "slug": "A"
}
```

Response:
```json
{
    "status": "Ok",
    "id": 2
}
```

### Удаление сегмента

DELETE http://localhost:8080/api/segments?slug=A

Request: Empty

Response: Empty

## Сегменты пользователя

### Добавление сегментов пользователю

POST http://localhost:8080/api/users/1/segments

Request:
```json
{
    "slugs": ["A", "B"]
}
```

Response:
```json
{
    "status": "Ok"
}
```

### Удаление сегмента у пользователя

DELETE http://localhost:8080/api/users/4/segments

Request:
```json
{
    "slugs": ["A", "B"]
}
```

Response: Emtpy

### Получение зегментов пользователя

GET http://localhost:8080/api/users/1/segments

Request: Empty

Response:
```json
{
    "status": "Ok",
    "userId": 1,
    "slugs": [
        "A"
    ]
}
```

### Множественное присвоение сегментов (Доп. задание 3)

POST http://localhost:8080/api/users/segments/percent

Request:
```json
{
    "slug": "A",
    "percent": 50
}
```

Response:
```json
{
    "status": "Ok"
}
```

## Отчёты

### Генерация отчётов (Доп. задание 1)

GET http://localhost:8080/api/reports?year=2023&month=8

Request: Empty

Response:
```json
{
    "url": "http://0.0.0.0:8080/reports/report_1693511988283384182.csv"
}
```

## P.S.

Я первый раз пишу на Go приложение такого рода. Раньше я писал, в основном, на Java. Очень много нюансов, с которыми нужно было разобраться. Так, я не написал автоматические тесты, обойдясь ручным, потому что это тоже займёт время на изучение, а его не осталось. 

Основные идеи я брал из этого ролика: [https://www.youtube.com/watch?v=rCJvW2xgnk0](https://www.youtube.com/watch?v=rCJvW2xgnk0)

В чём-то я разбирался через поиск, в чём-то через ChatGPT, что-то мне подсказывал знакомый. В итоге получился вот такой первый блин комом с не самым лучшим выбором библиотек и в целом реализацией. В следующий раз будет лучше
