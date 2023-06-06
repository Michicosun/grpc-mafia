# GRPC-Mafia

Реализация сервера и клиента для игры в мафию.

## Архитектура

Проект состоит их 5х частей:

- GRPC-Server - сервис, отвечающий за движок игры, реализован посредством передачи сообщений, использует stream-stream rpc вызов.
- Client - cli приложение, способное предоставлять доступ до предыдущих двух сервисов.
- RabbitMQ-chat - чат, основанный на брокере сообщений RabbitMQ
- Registry - HTTP REST-API сервис, отвечающий за регистрацию игроков и сбор статистики,
реботает на основе паттерна очерерь задач, поверх RabbitMQ
- Round-Tracker(aka. Tracker) - GraphQL сервис, позволяющий получить информацию о прошедших и текущих играх, а также прокомментировать их.

## Запуск Клиента

Команда запуска клиента:

```golang
path/grpc-mafia$ go run cmd/client/main.go
```

## Запуск Сервисов

Все Docker образы загружены в DockerHub, поэтому достаточно выполнить

```bash
path/grpc-mafia$ docker-compose up
```

В docker-compose поднимается сам сервис игры, rabbitmq, registry, а также 3 бота, которые сразу же подключаются к серверу, это сделано, чтобы можно было удобно протестировать возможности проекта.

## Build

Docker образы можно собрать по следующей схеме:

-   Сборка GRPC-Server'а
    ```bash
    path/grpc-mafia$ docker build . -f server/Dockerfile -t michicosun/mafia-server
    ```
-   Сборка Client'а
    ```bash
    path/grpc-mafia$ docker build . -f client/Dockerfile -t michicosun/mafia-bot
    ```
-   Сборка Registry
    ```bash
    path/grpc-mafia$ docker build . -f registry/Dockerfile -t michicosun/mafia-registry
    ```
-   Сборка Tracker'a
    ```bash
    path/grpc-mafia$ docker build . -f round-tracker/Dockerfile -t michicosun/mafia-tracker
    ```

## Registry API

У Registry есть 2 группы публичных endpoint'ов, доступных по 8080 порту:

- POST `/users/:login` - создает или обновляет информацию по пользователю, аргументы передаются, используя `multipart/form-data`
формат, пример:
  ```bash
    curl -X POST localhost:8080/users/michicosun -F mail=mr-robot@protonmail.ch -F avatar=@picture.jpg
  ```

- GET `/users/:login` - получить информацию по логину пользователя

- GET `/users/?logins=1,2,3` - получить информацию для группы пользователей, через запятую в logins необходимо указать login'ы интересующих пользователей

- DELETE `/users/:login` - удаляет информацию о пользователе, но не о его статистике

---

- POST `/pdf/:login` - создать запрос на рендер pdf документа для данного пользователя, ответом будет ссылка, по которой в будущем будет доступен отчет

- GET `/pdf/:filename` - получить очет по ссылке, полученной из предыдущего запроса

## Tracker API

API Трекера можно посмотреть в файле [schema.graphqls](./round-tracker/graph/schema.graphqls).

- AddComment - мутирующий метод, добавляющий комментарий к игре по ее id.
- GetRoundInfo - отдает информацию об игре
- ListRounds - отдает список из последних n игр, находящихся в состоянии state.

Также на порту 9090 будет доступна песочница, где можно удобно потестировать запросы.

## CLI

Client работает по протоколу автодополнения команд пользователя, возможные команды:

- login __username__ - сохраняет login пользователя, который используется всеми сервисами mafia.
- exit - закрывает клиента

---

- connect - подключается к игре и зависает, ожидая достаточное кол-во игроков.
- nothing - команда выхода из prepare стадии игры(до 1го дня).
- vote __username__ - основная команда игры, позволяет проголосовать как днем, так и ночью.
- message __[all, group]__ __text__ - отправляет сообщение всем или только группе, в которую вас определило.
- disconnect - закрывает сессию, но не клиента

---

- add-comment __id__, __comment__ - добавляет комментарий к игре с переданным id
- get-round-info __id__ - выводит всю информацию про игру с переданным id
- list-rounds __n__, __state__ - отдает список из последних n игр, находящихся в состоянии state.
