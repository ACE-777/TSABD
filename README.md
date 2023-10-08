# TSABD

# Hw 1

```bash
cd hw_1
```
запуск servochka:
```bash
go run cmd/server/main.go
```

запуск тестов:
```bash
cd internal/server
go test
```

# Hw 2

```bash
cd hw_2
```
запуск servochka:
```bash
go run cmd/server/main.go
```

Журанлы транзакции находятся в ``` hw_2/internal/logs```, каждые 10 минут заводим новый журнал.
Снэпшот хранится в переменной ```snapshot```
