## Запуск для локального тестирования
```
go get -u golang.org/x/net
go build -o http-server cmd/http-server/main.go; ./http-server
```

## Выполнить тестовый запрос
```
curl -v localhost:8090/parse -d '{"url":["http://localhost:8090/test"]}'
```
