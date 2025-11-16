# HTTP транспорт

- **handlers/** - обработчики запросов (team, pull request, user, stats)
- **dto/** - Data Transfer Objects (запросы/ответы API)
- **router/** - настройка маршрутов Gin

Поддерживает:
- REST API для управления командами и PR
- JSON валидацию через Gin binding
- Единый формат ошибок