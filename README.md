# GRPC-Mafia

Реализация сервера и клиента для игры в мафию.

## Архитектура

Всего существует 3 сущности:

- GRPC-Server - сервис, отвечающий за движок игры, реализован посредством передачи сообщений, использует stream-stream rpc вызов.
- Chat-Coordinator - сервис, сохраняющий информацию о чат группах, а также занимающийся пересылкой сообщений по схеме one to all.
- Client - cli приложение, способное предоставлять доступ до предыдущих двух сервисов.

## Запуск

Все Docker образы загружены в DockerHub, поэтому достаточно выполнить

```bash
path/grpc-mafia$ docker-compose up
```

Chat-Coordinator по умолчанию прослушивает UDP порт 20500, а GRPC-Server 9000, по которому доступен GRPC сервис игры. Все значения портов/хостов подключения можно менять через параметры окружения.

В docker-compose поднимается сам сервис игры, чат, а также 3 бота, которые сразу же подключаются к серверу, это сделано, чтобы можно было удобно протестировать возможности проекта.

Команда запуска клиента:

```golang
path/grpc-mafia$ go run cmd/client/main.go
```

## Build

Docker образы можно собрать по следующей схеме:

-   Сборка GRPC-Server'а
    ```bash
    path/grpc-mafia$ docker build . -f server/Dockerfile -t michicosun/mafia-server
    ```
-   Сборка Chat-Coordiantor'а
    ```bash
    path/grpc-mafia$ docker build . -f chat/Dockerfile -t michicosun/chat-coordinator
    ```
-   Сборка Client'а
    ```bash
    path/grpc-mafia$ docker build . -f client/Dockerfile -t michicosun/mafia-bot
    ```

## CLI

Client работает по протоколу автодополнения команд пользователя, возможные команды:

- connect __username__ - подключается к игре и зависает, ожидая достаточное кол-во игроков.
- nothing - команда выхода из prepare стадии игры(до 1го дня).
- vote __username__ - основная команда игры, позволяет проголосовать как днем, так и ночью.
- message __[all, group]__ __text__ - отправляет сообщение всем или только группе, в которую вас определило.
- disconnect - закрывает сессию, но не клиента
- exit - закрывает клиента

