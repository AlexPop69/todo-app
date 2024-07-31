# todo-app
Аналог приложения Microsoft Todo на Golang

## Описание
Backend с базовым набором функционала, с использованием `REST API`.

В приложении реализованы:
- Регистрация и авторизация пользователей
- Возможность создания, просмотра и удаления списка задач
- Возможность добавления, изменения и удаления отдельных задач

В API имеются следующие операции:
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

## Запуск приложения в docker 🚀
В корне проекта ввести в консоли команду `make run`

## Техническое задание 📝
 - **Создать веб-сервер и определить будущие эндпоинты**
 - **Спроектировать и создать базу данных** для хранения задач с использованием **Postgresql**. Реализовать миграции, используя библиотеку [GOOSE](https://github.com/pressly/goose).
    Таблицы БД:
    - для хранения пользователей (`users`)
    - для хранение списков (`todo_list`)
    - для хранения элементов списка (`todo_item`)
    - таблица, связывающая пользователей и их списки (`users_lists`)
    - таблица, связывающая списки и задачи (`list_items`)
  
 - **Добавить регистрацию и аутентификацию пользователей**
 - **Реализовать действия со списком:**
   - создание списка
   - получение всех списков пользователя и конкретного списка по ID
   - удаление списка 
   - изменение (обновления) списка
 - **Реализовать действия с элементами списка:**
   - добавление элемента в список
   - получение всех элементов списка
   - получение элемента списка по ID
   - удаление элемента из списка 
   - изменение (обновление) элемента
 - **Завершение выполнения всех запросов перед завершением приложения (Graceful Shutdown)**
 - **Добавить тесты**
   - Тесты регистрации и авторизации
   - Тесты репозитория (БД)
 - **Создать и настроить запуск приложения в контейнере**
 - **Добавить Makefile** с командами:
    - Запуск приложения `run`
    - Тестирование приложения `test`
    - Применение миграции `migrate`
 - **Добавить документирование API при помощуи Swagger**
     - Использовал утилиту [swag](https://github.com/swaggo/swag).

## Стек:
<img src="https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/go/go-original.svg" 
    title="golang" width="50" height="50"/>&nbsp;
<img src="https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/postgresql/postgresql-original.svg"                 
    title="postgresql" width="50" height="50"/>&nbsp;
<img src="https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/docker/docker-original.svg"          
    title="docker" width="50" height="50"/>&nbsp;
<img src="https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/git/git-original.svg"          
    title="git" width="50" height="50"/>&nbsp;
<img src="https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/postman/postman-original.svg"          
    title="git" width="50" height="50"/>&nbsp;
<img src="https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/swagger/swagger-original.svg"          
    title="git" width="50" height="50"/>&nbsp;
<!--
-->