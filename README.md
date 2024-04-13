# Сервис баннеров

Сервис баннеров реализует [API](https://github.com/vvauijij/avito-tech-assignment/blob/main/server/api/openapi.yml) для пользователей и админов:

- пользователь имеет возможность получать закэшированную/актуальную информацию о содержимом баннера по фиче и тэгу
- админ имеет возможность создавать/обновлять/удалять баннер и получать/удалять баннеры c фильтрацией по фиче и/или тегу

Архитектура сервиса описана в [C4 диаграмме](https://github.com/vvauijij/avito-tech-assignment/blob/main/docs/containers.puml):

- сервер предоставляет функциональность сервиса для пользователей и админов через REST API
- для хранения баннеров используется MongoDB
- для кэширования запросов пользователей используется Redis

Сущность баннера описана в [диаграмме](https://github.com/vvauijij/avito-tech-assignment/blob/main/docs/entities.puml):

- один баннер может быть связан только с одной фичей и несколькими тегами
- содержимое баннера представляет собой JSON-документ неопределенной структуры
- пользователи не имеют доступа к выключенным баннерам, админы имеют доступ ко всем баннерам

## Аутентификация

Для аутентификации используются JWT-токены, сгенерированные с помощью [пары публичного и приватного ключей](https://github.com/vvauijij/avito-tech-assignment/blob/main/secrets):
- для верификации прав доступа (пользователь/админ) сервер должен иметь доступ к публичному ключу
- для генерации токенов в тестах используется приватный ключ

## Тестирование
В [docker-compose.yml](https://github.com/vvauijij/avito-tech-assignment/blob/main/docker-compose.yml) описана конфигурация сервиса и тестов: `make test` собирает и поднимает сервис, запускает E2E тесты. Тестирование поддерживается в [CI](https://github.com/vvauijij/avito-tech-assignment/actions).

Для изоляции E2E тестов используется ручка `/test_clean_up`, очищающая кэш и хранилище баннеров - доступна только в тестовом окружении.

- [тесты аутентификации](https://github.com/vvauijij/avito-tech-assignment/blob/main/tests/auth_test.go) покрывают сценарии запросов с некорректными правами доступа. Для генерации токенов используется приватный ключ

- [тесты пользовательского функционала](https://github.com/vvauijij/avito-tech-assignment/blob/main/tests/user_test.go) покрывают сценарии запросов получения содержимого баннеров разной структуры, получения закэшированной и актуальной информации о содержимом баннеров, получения содержимого неактивных баннеров

- [тесты админского функционала](https://github.com/vvauijij/avito-tech-assignment/blob/main/tests/admin_test.go) покрывают сценарии запросов создания/обновления/удаления баннера, получения баннеров c фильтрацией по фиче и/или тегу, количеству, отступу и асинхронного удаления баннеров c фильтрацией по фиче и/или тегу

## Примеры

Для удобства локальной отправки запросов к сервису предусмотрены следующие механизмы, доступные только в тестовом окружении:
- DELETE-запрос к ручке `/test_clean_up` очищает кэш и хранилище баннеров
- токен `test_token` дает права доступа админа без необходимости генерации токена с использованием приватного ключа 

Собираем и поднимаем сервис:
```bash
make run
```

Очищаем кэш и хранилище баннеров:
```bash
curl -X 'DELETE' \
  'http://localhost:8080/test_clean_up'
```

Создаем баннеры с произвольным содержимым:
```bash
curl -X 'POST' \
  'http://localhost:8080/banner' \
  -H 'accept: application/json' \
  -H 'token: test_token' \
  -H 'Content-Type: application/json' \
  -d '{
  "tag_ids": [
    0, 1, 2
  ],
  "feature_id": 0,
  "content": {
    "title": "first banner",
    "text": "text",
    "number": 0,
    "object": {
        "text": "text",
        "number": 0
    }
  },
  "is_active": true
}'
```
```bash
{"banner_id":6853104346068991556}
```


```bash
curl -X 'POST' \
  'http://localhost:8080/banner' \
  -H 'accept: application/json' \
  -H 'token: test_token' \
  -H 'Content-Type: application/json' \
  -d '{
  "tag_ids": [
    2, 3, 4
  ],
  "feature_id": 1,
  "content": {
    "title": "second banner",
    "array": [2, 3, 4]
  },
  "is_active": true
}'
```
```bash
{"banner_id":119881502863611628}
```

Получаем содержимое первого баннера:
```bash
curl -X 'GET' \
  'http://localhost:8080/user_banner?tag_id=0&feature_id=0' \
  -H 'accept: application/json' \
  -H 'token: test_token'
```
```bash
{"number":0,"object":{"number":0,"text":"text"},"text":"text","title":"first banner"}
```

Обновляем первый баннер:
```bash
curl -X 'PATCH' \
  'http://localhost:8080/banner/6853104346068991556' \
  -H 'accept: */*' \
  -H 'token: test_token' \
  -H 'Content-Type: application/json' \
  -d '{
  "tag_ids": null,
  "feature_id": null,
  "content": {
    "title": "updated first banner"
  },
  "is_active": null
}'
```

Получаем закэшированное содержимое первого баннера:
```bash
curl -X 'GET' \
  'http://localhost:8080/user_banner?tag_id=0&feature_id=0' \
  -H 'accept: application/json' \
  -H 'token: test_token'
```
```bash
{"number":0,"object":{"number":0,"text":"text"},"text":"text","title":"first banner"}
```

Получаем актуальное содержимое первого баннера:
```bash
curl -X 'GET' \
  'http://localhost:8080/user_banner?tag_id=0&feature_id=0&use_last_revision=true' \
  -H 'accept: application/json' \
  -H 'token: test_token'
```
```bash
{"title":"updated first banner"}
```

Получаем содержимое баннеров по фильтру:
```bash
curl -X 'GET' \
  'http://localhost:8080/banner?tag_id=2' \
  -H 'accept: application/json' \
  -H 'token: test_token'
```
```bash
[{"banner_id":6853104346068991556,"feature_id":0,"tag_ids":[0,1,2],"content":{"title":"updated first banner"},"is_active":true,"created_at":"2024-04-13T23:09:20Z","updated_at":"2024-04-13T23:15:36Z"},{"banner_id":119881502863611628,"feature_id":1,"tag_ids":[2,3,4],"content":{"array":[2,3,4],"title":"second banner"},"is_active":true,"created_at":"2024-04-13T23:09:45Z","updated_at":"2024-04-13T23:09:45Z"}]
```

Удаляем второй баннер:
```bash
curl -X 'DELETE' \
  'http://localhost:8080/banner/119881502863611628' \
  -H 'accept: */*' \
  -H 'token: test_token'
```

Размещаем запрос на удаление баннеров по фильтру:
```bash
curl -X 'DELETE' \
  'http://localhost:8080/banner?tag_id=2' \
  -H 'accept: */*' \
  -H 'token: test_token'
```

Полная доступная функциональность описана в [API](https://github.com/vvauijij/avito-tech-assignment/blob/main/server/api/openapi.yml).