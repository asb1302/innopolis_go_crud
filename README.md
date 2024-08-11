# innopolis_go_crud

## task

- дополните сервис эндпоинтами для пагинационного чтения рецептов, за максимальное количество получаемых записей в 1
  запросе примем 10;

## Инструкция по работе с приложением

### Актуальные данные

* Актуальный адрес прода: https://innopolisgocrud-production.up.railway.app/ping

### Получение рецепта по ID

```http request
GET https://innopolisgocrud-production.up.railway.app/?id=<укажите ID рецепта>
Authorization: <укажите токен>
```

### Добавление или обновление рецепта

```http request
POST https://innopolisgocrud-production.up.railway.app
Authorization: <укажите токен>
Content-Type: application/json

{
  "id": "<укажите ID рецепта>",
  "name": "<укажите название рецепта>",
  "ingredients": [
    {
      "name": "<укажите название ингредиента>",
      "quantity": "<укажите количество ингредиента>"
    }
  ],
  "steps": [
    "<укажите шаг приготовления>"
  ]
}
```

### Удаление рецепта по ID

```http request
DELETE https://innopolisgocrud-production.up.railway.app/?id=<укажите ID рецепта>
Authorization: <укажите токен>
```


### Получение рецептов с пагинацией

```http request
GET https://innopolisgocrud-production.up.railway.app/?page=<номер страницы>&limit=<количество рецептов на странице>
Authorization: <укажите токен
```

### Получение количества всех рецептов

```http request
GET https://innopolisgocrud-production.up.railway.app/count
Authorization: <укажите токен>
```

## Локальная работа с проектом

### Требования

- Docker
- Docker Compose
- Make

### Основные команды Makefile

#### Сборка и запуск проекта

Для сборки и запуска проекта выполните следующую команду:

```sh
make
```

#### Очистка проекта

Для остановки и удаления контейнеров, образов, томов и зависших контейнеров выполните следующую команду:

```sh
make clean
```

#### Проверка доступности сервиса

Для проверки доступности сервиса выполните следующую команду:

```sh
make check
```

#### Примечание

Убедитесь, что вы настроили переменные окружения(.env) перед запуском команд.
