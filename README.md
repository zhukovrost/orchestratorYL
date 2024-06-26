# Распределенный вычислитель арифметических выражений (оркестратор)

Финальное задание по второму спринту Яндекс Лицея (GoLang)

Этот сервис принимает на вход арифметические выражения и разбивает их на задачи. 
Другой сервис (*агент*) запрашивает у оркестратора эти задачи, считает их, возвращает результат. 
Оркестратор также обрабатывает результаты *агента*.
Также можно отслеживать статус и результат каждого выражения.

![SVG Image](service.svg)

## Данный сервис работает с агентом

Подробнее тут: https://github.com/zhukovrost/agentYL.git

## Установка

1. Клонируйте репозиторий:
    ```sh
    git clone https://github.com/zhukovrost/orchestratorYL.git
    ```

2. Перейдите в директорию проекта:
    ```sh
    cd orchestratorYL
    ```
   
3. **При необходимости** установите окружение (Linux):

   ```sh 
   export PATH=$PATH:/usr/local/go/bin
   ```

4. Установите зависимости:
    ```sh
    go mod tidy
    ```

## Запуск

Для запуска сервера выполните:

```sh
go run cmd/orchestrator/main.go
```

## Настройка (при необходимости)

Вы можете поменять время выполнения арифметических операций и отображение логов уровня debug.
Для этого можно поменять значения флагов. Например:

```sh
go run cmd/orchestrator/main.go -add=3000 -div=1000000 -debug
```

Подробнее про это:

```sh
go run cmd/orchestrator/main.go --help
```

## Инструкция по использованию

### 1. Добавление арифметического выражения

* Правильное выражение:

   Код ответа будет **201**. Также выражение добавится в очередь и будет иметь статус *pending*.
   ```sh
   curl http://localhost:8080/api/v1/calculate \
   --header 'Content-Type: application/json' \
   --data '{
     "id": "validExpression1",
     "expression": "2 + 2 * 2"
   }'
   ```
  
   Если у выражения не задан id, то он будет сгенерирован.
   ```sh
   curl http://localhost:8080/api/v1/calculate \
   --header 'Content-Type: application/json' \
   --data '{
     "expression": "100 - 5 * (40 - 23) + 3"
   }'
   ```

* Выражение с ошибкой (деление на ноль):

  Код ответа будет **422**. Также выражение добавится в очередь и будет иметь статус *invalid*.
   ```sh
   curl http://localhost:8080/api/v1/calculate \
   --header 'Content-Type: application/json' \
   --data '{
     "id": "invalidExpression1",
     "expression": "4/0 + 1"
   }'
   ```

   Но если добавить выражение, сразу по которому не скажешь, что там есть деление на ноль, то 
код ответа будет **201**. Также выражение добавится в очередь и будет иметь статус *pending*.
В процессе подсчёта выражения будет выявлена ошибка, и статус выражения обновится на *invalid*.
   ```sh
   curl http://localhost:8080/api/v1/calculate \
   --header 'Content-Type: application/json' \
   --data '{
     "id": "invalidExpression2",
     "expression": "4/(3 - 3) + 90"
   }'
   ```

* Выражение с ошибкой (некорректные входные данные):

   Некоторые выражение могут выдать ошибку **422**. Приведу пример несколько таких:

  * / 5 + 9 -- выражение начинается не с числа. **ВАЖНО!!!** Если вы хотите начать выражение с отрицательного числа, пишите (0 - число)
  * 15 * + 7 -- выражение содержит 2 идущих подряд операнда. **ВАЖНО!!!** Если вы хотите произвести операцию с отрицательным числом, пишите (0 - число)
  * 15 - (123 -- выражение содержит некорректные скобки
  * (11) + 7 -- выражение содержит бесполезные скобки
  * 83 - d + @l -- выражение содержит буквы или другие неиспользуемые символы
  * 6 5 -- выражение содержит пробел между цифрами без арифметического знака
      
* Выражение не содержит выражения:

  Код ответа будет **500**. Также выражение не добавится в очередь.
   ```sh
   curl http://localhost:8080/api/v1/calculate \
   --header 'Content-Type: application/json' \
   --data '{
     "id": "invalidExpression3",
     "expression": {"notExpression": true}
   }'
   ```   
  
* Выражение уже существует:

   Код ответа будет **422**. Также выражение не добавится в очередь.

   ```sh
   curl http://localhost:8080/api/v1/calculate \
   --header 'Content-Type: application/json' \
   --data '{
     "id": "validExpression1",
     "expression": "2 + 2"
   }'
   ```

### 2. Получение всех арифметических выражений

```sh 
curl http://localhost:8080/api/v1/expressions
```

Код ответа **200**. Ответом будет:
```json
[
   {
      "id":"invalidExpression1",
      "expression":"4/0 + 1",
      "result":0,
      "status":"invalid"
   },
   {
      "id":"invalidExpression2",
      "expression":"4/(3 - 3) + 90",
      "result":0,
      "status":"pending"
   },
   {
      "id":"validExpression1",
      "expression":"2 + 2 * 2",
      "result":0,
      "status":"pending"
   },
   {
      "id":"605702",
      "expression":"100 - 5 * (40 - 23) + 3",
      "result":0,
      "status":"pending"
   }
]

```

### 3. Получение арифметических выражений по ID

Дан пример получения выражения с ID **validExpression1**. Код ответа будет **200**:
```sh
curl http://localhost:8080/api/v1/expressions/validExpression1
```
И ответ:
```json
{
   "id":"validExpression1",
   "expression":"2 + 2 * 2",
   "result":0,
   "status":"pending"
}
```

Если запросить несуществующее выражение код ответа будет **404**. 

### 4. Получение задачи

С данной функцией работает **агент**. Он постоянно запрашивает задачу. Если их нет, то получает код **404**. 
Если есть, то код **200** и ответ.

Вот пример:

```sh
curl http://localhost:8080/internal/task
```

Ответ:

```json
{
   "task": {
      "id":0,
      "arg1":2,
      "arg2":2,
      "operation":"*",
      "operation_time":5000
   }
}
```

### 5. Установить результат задачи

С данной функцией работает **агент**. Он постоянно решает задачу и отправляет результат обратно. 
Если нет задачи с данным id, то получает код **404**. Если есть, то код **200** и ответ, при условии, что данные result валидны, иначе **422**.

```sh
curl http://localhost:8080/internal/task \
--header 'Content-Type: application/json' \
--data '{
  "id": 0,
  "result": 4
}'
```

## Структура проекта

```
orchestratorYL/
├── cmd/
│   └── orchestrator/
│       └── main.go
├── internal/
│   ├── app/
│   │   └── app.go
│   ├── config/
│   │   └── config.go
│   ├── handler/
│   │   └── handler.go
│   ├── router/
│   │   └── router.go
│   └── service/
│       └── service.go
├── pkg/
│   └── utils/
│       └── utils.go
├── .gitignore
├── go.mod
└── README.md
```
