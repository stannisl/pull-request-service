# Pull request service

## Кратко про проект

### Стек проекта

- **Язык**: Go 1.24
- **Фреймворк**: Gin (HTTP)
- **База данных**: PostgreSQL
- **Миграции**: Встроенные через embed
- **Тестирование**: testify + testcontainers

### Структура

```
./
├──cmd/service/ # Точка входа
├── internal/ # Внутренние модули
├── pkg/db/ # Утилиты работы с БД
└── tests/e2e/ # End-to-end тесты
```

## Запуск сервиса

Запустить сервис

```bash
docker compose up
```


## Вопросы и решения

1. Авторизация пользователей и 401 ответы. У нас в изначальном апи нет никакой авторизации и подтвеждения что действия
идут от автора, поэтому опущено
2. Обновление команды было опущено, потому что в openapi.yaml есть только создание, а при создании с таким же названием
мы имеем ошибку TEAM_EXISTS
3. Добавлена ручка GET /stats для взятия статистики по пользователям назначенными ревьюверами на OPEN PullRequest

## Допольнительные задачи 

### E2E тестирование

Запуск end-to-end тестов:
```bash
make test-e2e
```

#### Что проверяют E2E тесты

Создание команды (POST /team/add) и её получение (GET /team/get) \
Создание Pull Request (POST /pullRequest/create) \
Слияние Pull Request (POST /pullRequest/merge) \
Переназначение ревьювера (POST /pullRequest/reassign) \
Проверка ошибок: NotFound (404) и InvalidRequest (400) \
Проверка базовых метрик успешности и статусов HTTP

#### Пример успешного запуска

```bash
[GIN] 2025/11/17 - 04:33:16 | 200 |      89.155µs |       127.0.0.1 | GET      "/health"
2025/11/17 04:33:16 Application is ready!
=== RUN   TestE2ETestSuite/TestPullRequest_CreateMergeAndReassign
[GIN] 2025/11/17 - 04:33:16 | 201 |    6.161514ms |             ::1 | POST     "/team/add"
[GIN] 2025/11/17 - 04:33:16 | 201 |     4.67464ms |             ::1 | POST     "/pullRequest/create"
[GIN] 2025/11/17 - 04:33:16 | 200 |    2.584251ms |             ::1 | POST     "/pullRequest/merge"
=== RUN   TestE2ETestSuite/TestPullRequest_NotFound
2025/11/17 04:33:16 error getting pull request: entity not found in db
[GIN] 2025/11/17 - 04:33:16 | 404 |     524.929µs |             ::1 | POST     "/pullRequest/merge"
=== RUN   TestE2ETestSuite/TestPullRequest_ReassignReviewer
[GIN] 2025/11/17 - 04:33:16 | 201 |    3.761544ms |             ::1 | POST     "/team/add"
[GIN] 2025/11/17 - 04:33:16 | 201 |    2.224457ms |             ::1 | POST     "/pullRequest/create"
[GIN] 2025/11/17 - 04:33:16 | 201 |    3.122926ms |             ::1 | POST     "/pullRequest/reassign"
=== RUN   TestE2ETestSuite/TestTeam_CreateAndGet
[GIN] 2025/11/17 - 04:33:16 | 201 |    2.431098ms |             ::1 | POST     "/team/add"
[GIN] 2025/11/17 - 04:33:16 | 200 |     740.738µs |             ::1 | GET      "/team/get?team_name=backend-team"
=== RUN   TestE2ETestSuite/TestTeam_InvalidRequest
[GIN] 2025/11/17 - 04:33:16 | 400 |      17.743µs |             ::1 | POST     "/team/add"
=== RUN   TestE2ETestSuite/TestTeam_NotFound
2025/11/17 04:33:16 Error getting team: entity not found in db
[GIN] 2025/11/17 - 04:33:16 | 404 |     481.418µs |             ::1 | GET      "/team/get?team_name=nonexistent"
2025/11/17 04:33:16 Tearing down test suite...
2025/11/17 04:33:16 🐳 Stopping container: 0edb9ba08ae6
2025/11/17 04:33:17 ✅ Container stopped: 0edb9ba08ae6
2025/11/17 04:33:17 🐳 Terminating container: 0edb9ba08ae6
2025/11/17 04:33:17 🚫 Container terminated: 0edb9ba08ae6
2025/11/17 04:33:17 PostgreSQL container terminated
--- PASS: TestE2ETestSuite (5.59s)
    --- PASS: TestE2ETestSuite/TestPullRequest_CreateMergeAndReassign (0.02s)
    --- PASS: TestE2ETestSuite/TestPullRequest_NotFound (0.00s)
    --- PASS: TestE2ETestSuite/TestPullRequest_ReassignReviewer (0.01s)
    --- PASS: TestE2ETestSuite/TestTeam_CreateAndGet (0.00s)
    --- PASS: TestE2ETestSuite/TestTeam_InvalidRequest (0.00s)
    --- PASS: TestE2ETestSuite/TestTeam_NotFound (0.00s)
PASS
ok      github.com/stannisl/pull-request-service/tests/e2e      5.635s
```

### Нагрузочное тестирование

Запуск нагрузочных тестов (требуется k6):
```bash
make test-load
```

#### Целевые показатели

- RPS: 5 запросов/секунду
- Время ответа (SLI): 95-й перцентиль < 300 мс
- Доступность (SLI): 99.9% успешных запросов

#### Фактические результаты тестирования

| Метрика             | Целевое значение | Результат |
|---------------------|------------------|-----------|
| RPS                 | 5                | 28.06     |
| Время ответа (p95)  | < 300 мс         | 10.82 мс  | 
| Успешность запросов | 99.9%            | 100%      | 
| Ошибки              | < 0.1%           | 0%        | 

#### Сырой результат теста

```
  █ THRESHOLDS 

    errors
    ✓ 'rate<0.001' rate=0.00%

    http_req_duration
    ✓ 'p(95)<300' p(95)=10.82ms

    http_req_failed
    ✓ 'rate<0.001' rate=0.00%

    team_creation_duration
    ✓ 'p(95)<250' p(95)=13

    team_retrieval_duration
    ✓ 'p(95)<200' p(95)=6


  █ TOTAL RESULTS 

    checks_total.......: 35352   84.167066/s
    checks_succeeded...: 100.00% 35352 out of 35352
    checks_failed......: 0.00%   0 out of 35352

    ✓ team created successfully
    ✓ response time acceptable
    ✓ has correct content type
    ✓ team retrieved successfully
    ✓ team data is correct

    CUSTOM
    errors.........................: 0.00%  0 out of 0
    team_creation_duration.........: avg=6.773082 min=3        med=6        max=32      p(90)=11       p(95)=13      
    team_retrieval_duration........: avg=1.898337 min=0        med=2        max=9       p(90)=2        p(95)=6       

    HTTP
    http_req_duration..............: avg=4.18ms   min=577.43µs med=4.56ms   max=32.77ms p(90)=7.48ms   p(95)=10.82ms 
      { expected_response:true }...: avg=4.18ms   min=577.43µs med=4.56ms   max=32.77ms p(90)=7.48ms   p(95)=10.82ms 
    http_req_failed................: 0.00%  0 out of 11784
    http_reqs......................: 11784  28.055689/s

    EXECUTION
    iteration_duration.............: avg=361.72ms min=207.93ms med=361.77ms max=518.7ms p(90)=481.57ms p(95)=494.72ms
    iterations.....................: 5892   14.027844/s
    vus............................: 1      min=1          max=10
    vus_max........................: 10     min=10         max=10

    NETWORK
    data_received..................: 5.6 MB 13 kB/s
    data_sent......................: 3.5 MB 8.3 kB/s




running (7m00.0s), 00/10 VUs, 5892 complete and 0 interrupted iterations
default ✓ [======================================] 00/10 VUs  7m0s
The last command took 420.656 seconds.
```