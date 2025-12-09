# TODO List (Планировщик задач)

## Описание проекта
Веб-сервер для управления задачами (TODO-лист) с возможностью добавления, редактирования, удаления и отметки выполнения задач. Поддерживает повторяющиеся задачи с правилами повторения.

## Задания со звёздочкой
 Выполнено: настройка порта через TODO_PORT и пути к БД через TODO_DBFILE.
 Не выполнено: расширенные правила повторения (w, m), поиск задач, аутентификация.

## Инструкция по запуску кода локально
1. `go run main.go`
2. Открыть: http://localhost:7540
3. Опционально: `export TODO_PORT=9000` или `export TODO_DBFILE=./my.db`

## Запуск тестов
- Все тесты: `go test ./tests`
- Конкретные: `go test -run ^TestAddTask$ ./tests` и т.д.
- Настройки в tests/settings.go: Port, DBFile, FullNextDate, Search.

## Docker
Docker конфигурация не реализована.

## Структура проекта

**main.go** - запуск сервера, обработка переменных окружения TODO_PORT и TODO_DBFILE

**pkg/api/**
- **api.go** - регистрация всех API маршрутов
- **addtask.go** - POST /api/task (добавление задачи, функция checkDate)
- **task.go** - GET/PUT/DELETE /api/task (получение, обновление, удаление)
- **taskdone.go** - POST /api/task/done (отметка выполнения задачи)
- **tasks.go** - GET /api/tasks (список задач)
- **nextdate_handler.go** - GET /api/nextdate
- **nextDate.go** - функция NextDate() для расчёта следующих дат

**pkg/db/**
- **db.go** - инициализация SQLite базы, создание таблицы scheduler
- **task.go** - функции работы с БД: AddTask, GetTask, UpdateTask, DeleteTask, UpdateDate
- **errors.go** - константы ошибок (ErrTaskNotFound и др.)

**tests/** - тесты для проверки API
**web/** - статические файлы фронтенда (HTML, CSS, JS)