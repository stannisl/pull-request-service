# Pull request service

## кратко про проект


### стек проекта

- **Язык**: Go 1.19+
- **Фреймворк**: Gin (HTTP)
- **База данных**: PostgreSQL
- **Миграции**: Встроенные через embed
- **Тестирование**: testify + testcontainers

### структура

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
3. 

## Допольнительные задачи 

### E2E тестирование

```bash
make test-e2e
```

### Нагрузочное тестирование

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