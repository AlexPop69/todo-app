# todo-app
Аналог приложения Microsoft Todo на Golang

## Описание
Нужно сделать backend с базовым набором функционала, с использованием `REST API`.

В приложении должны быть реализованы:
- Регистрация и авторизация
- Возможность создания, просмотра и удаления списка задач
- Возможность добавления, изменения и удаления отдельных задач

В API будет следующие операции:
- Со списками:
    - создать список;
    - получить список;
    - изменить список;
    - удалить список;
- С задачами:
    - добавить задачу;
    - получить задачу;
    - изменить задачу;
    - удалить задачу 

#### Схема архитектуры проекта

<a href="https://imgbb.com/"><img src="https://i.ibb.co/6F8zXcf/image.png" alt="image"></a>

## Техническое задание 
Чтобы упорядочить работу, выполнение проекта разбито на шаги.

#### Шаг 1. Создать и запустить веб-сервер
- Создать веб-сервер, который будет слушать определённый порт.
- Определить будущие эндпоинты

#### Шаг 2. Спроектировать и создать базу данных
Нужно создать базу данных для хранения задач с использованием **Postgresql**. Также я использовал библиотеку [GOOSE](https://github.com/pressly/goose) для миграций.

Таблицы БД:
- для хранения пользователей (`users`)
- для хранение списков (`todo_list`)
- для хранения элементов списка (`todo_item`)
- таблица, связывающая пользователей и их списки (`users_lists`)
- таблица, связывающая списки и задачи (`list_items`)

## Стек:
<img src="https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/go/go-original.svg" 
    title="golang" width="50" height="50"/>&nbsp;
<img src="https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/postgresql/postgresql-original.svg"                 
    title="postgresql" width="50" height="50"/>&nbsp;
<img src="https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/docker/docker-original.svg"          
    title="docker" width="50" height="50"/>&nbsp;
<img src="https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/git/git-original.svg"          
    title="git" width="50" height="50"/>&nbsp;

<!--
#### Шаг 2. Проектируем и создаём базу данных
- ✅ Создать базу данных для хранения задач с использованием **Poistgresql**
- ✅ Реализовать при запуске сервера проверку, существует ли в директории приложения файл `scheduler.db`. Если его нет, следует создать базу данных с таблицей `scheduler`.
- ✅ Создать в таблицу `scheduler` с полями: **id**, **date**, **title**, **comment**, **repeat**
- ✅ ❄️ Реализовать возможность определять путь к файлу базы данных через переменную окружения `TODO_DBFILE`

-->

