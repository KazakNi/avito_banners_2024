# REST API-сервис для баннеров

### Решение:
1. Для запросов баннеров юзеров применил кеширование с таймингом горячего кеша 5 мин. Полагаю, что юзеров гораздо больше, чем админов.
2. Применил индексы на колонки, по которым ищутся баннеры и совершаются джоины тэгов.
3. Тэги хранятся в таблице Many To Many по отношению к баннерам, фичи денормализированно хранятся в одной таблице с баннерами, т. к. уникально определяют последние.

Для разворачивания сервиса локально необходимо клонировать репозиторий:
```
git clone https://github.com/KazakNi/avito_banners_2024.git
```

Перейти в директорию проекта

```
cd avito_banners_2024
```
Заполнить файл с конфигом

Запустить docker-compose

```
docker-compose up --build
```

Для регистрации пользователя необходимо направить запрос:

```curl
curl -X 'POST' \
  'https://localhost:8080/user/sign_up' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "username": "test",
  "password": "test"
  "email": "test@test.com"
}'
```
Далее получить токен по эндпоинту:
```curl
curl -X 'POST' \
  'https://localhost:8080/user/sign_in' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "username": "test",
  "password": "test"
  "email": "test@test.com"
}'
```
Выполнять дальнейшие опросы по ручкам с Bearer-токеном

Тест на получение баннера выполняется на отдельной таблице для тестирования:

```
go test -v /test
```

TODO:
- Yandex tank load testing
- Message broker for delayed tasks
- e2e others endpoints
