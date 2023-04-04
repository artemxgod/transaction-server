# Transaction server

## layout

- Хэндлеры находятся в `/internal/app/tnserver/server.go`
- В папке `/configs` можно настроить порт, ссылку на бд
- За запуск сервера отвечает функция `Start()` `internal/app/tnserver/tnserver.go`
- Так как для каждого пользователя нужна своя очередь я создавал уникальные топики в kafka для каждого из них