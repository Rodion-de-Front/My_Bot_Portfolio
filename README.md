# Портфолио бот на языке Go
Привет! Это мой портфолио бот, написанный на языке Go. Бот разработан для автоматизации определенных задач и демонстрации моих навыков в программировании.

## Описание проекта
Этот проект представляет собой пример использования языка Go для создания бота, способного выполнять различные задачи. Бот имеет модульную структуру, что позволяет легко добавлять новые функциональные возможности.

## Требования
Для запуска этого проекта вам понадобится следующее:

Go (версия X.X.X): Ссылка на установку: https://golang.org/
Docker (версия X.X.X): Ссылка на установку: https://www.docker.com/
### Установка и запуск
Для начала создайте `.env` файл с `TELEGRAM_API` в папке `backend`

Установите зависимости проекта с помощью go mod:

```
go mod download
go run .
```
Или вы можете собрать исполняемый файл и запустить его:

```
go build -o bot
./bot
```
---------
Docker

Для удобства развертывания и запуска бота предоставляются Docker-контейнеры. В проекте имеется docker-compose, который позволяет создать контейнер для бота.

Чтобы создать Docker-образ и запустить контейнер, выполните следующие действия:

Убедитесь, что у вас установлен Docker.

Перейдите в директорию проекта и запустите docker-compose:

```
docker-compose up
```
Ваш бот должен быть успешно запущен внутри контейнера Docker.