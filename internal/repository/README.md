# Слой работы с базой данных

репозитории:
- **TeamRepository** - команды
- **PullRequestRepository** - пулл-реквесты и ревьюверы
- **UserRepository** - пользователи
- **StatsRepository** - статистика

особенности:
- Использует sqlx для работы с БД
- Поддержка транзакций через TransactionManager

вспомогательное:
- **dependencies** - структура со всеми сервисами