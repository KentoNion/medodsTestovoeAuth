# Фичи и особенности
1. Добавил конфиг файл, строка host_db игнорируется и перезаписывается докером через внешнюю среду при запуске приложения из докера. /authApp/config.yaml
2. /refresh выдаёт 2 токена, новый аксес и старый рефреш если тот прошёл проверку. Так же если ip изменился реализована моковая отправка уведомления.
3. Есть тесты на бд и сервис(основная логика проекта) и мини-тест на существование файла конфига.
4. /login и /regresh запросы необходимо осуществлять методом POST
5. Докер работает!
6. Если пользователь уже существует, не получится применить authorize с тем же самым GUID
7. Можно использовать любую string как GUID
8. функция рефреш обновляет рефреш токен как того требует задание, и не позволяет использовать рефреш токен более 1 раза (как я понял это задание)
9. Тело запросов игнорируется, можно было бы GUID и передачу токена реализовать через Json, но в задании это не обговорено. 
10. Server, Service, Postgres (db) и config обложены тестами

Пример отправляемого запроса Authorize POST:
localhost:8050/login?GUID=1
Пример ответа:
{"access":"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzQyOTc5NzEsImlwIjoiMTI3LjAuMC4xOjUyNTM4IiwicmVmcmVzaCI6ImV5SmhiR2NpT2lKSVV6STFOaUlzSW5SNWNDSTZJa3BYVkNKOS5leUpsZUhBaU9qRTNNelk1TnpZek5qRXNJbWx3SWpvaU1USTNMakF1TUM0eE9qVXlOVE00SWl3aWMyVmpjbVYwSWpvaVlqRmlZbUV3WVRJdFpHSmxNeTAwTUdGa0xXRmtORFF0T0RaaFlqTXpZV1ZrT1dOaElpd2lkWE5sY2w5cFpDSTZJalVpZlEuNDdOU1QwdjFtOVpqZDc5UExZdFhpTVhNWHhkZDc4YjhRckh5QWo5dHZyOCIsInVzZXJfaWQiOiI1In0.oODBUJdI9OvOPuOp9Or-VJJHGcamhE4l5yzJEqT1QVqsKZeoXStg5jwzTZMAegjnZtqaBrxQ_ZSZmef8ai7Vvw","refresh":"lqVWlmUS40N05TVDB2MW05WmpkNzlQTFl0WGlNWE1YeGRkNzhiOFFySHlBajl0dn"}

Пример отправляемого запроса refresh POST:
localhost:8050/refresh?refresh_token=lqVWlmUS40N05TVDB2MW05WmpkNzlQTFl0WGlNWE1YeGRkNzhiOFFySHlBajl0dn&GUID=1
Ответ будет состоять из двух токенов: новый access и нового refresh (см пункт 8 фичей и особенностей)

# Test task BackDev

Тестовое задание на позицию Junior Backend Developer

**Используемые технологии:**

- Go
- JWT
- PostgreSQL

**Задание:**

Написать часть сервиса аутентификации.

Два REST маршрута:

- Первый маршрут выдает пару Access, Refresh токенов для пользователя с идентификатором (GUID) указанным в параметре запроса
- Второй маршрут выполняет Refresh операцию на пару Access, Refresh токенов

**Требования:**

Access токен тип JWT, алгоритм SHA512, хранить в базе строго запрещено.

Refresh токен тип произвольный, формат передачи base64, хранится в базе исключительно в виде bcrypt хеша, должен быть защищен от изменения на стороне клиента и попыток повторного использования.

Access, Refresh токены обоюдно связаны, Refresh операцию для Access токена можно выполнить только тем Refresh токеном который был выдан вместе с ним.

Payload токенов должен содержать сведения об ip адресе клиента, которому он был выдан. В случае, если ip адрес изменился, при рефреш операции нужно послать email warning на почту юзера (для упрощения можно использовать моковые данные).

**Результат:**

Результат выполнения задания нужно предоставить в виде исходного кода на Github. Будет плюсом, если получится использовать Docker и покрыть код тестами.