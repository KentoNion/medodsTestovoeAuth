# Фичи и особенности
1. Добавил конфиг файл, строка host_db игнорируется и перезаписывается докером через внешнюю среду при запуске приложения из докера. /authApp/config.yaml
2. /refresh выдаёт 2 токена, новый аксес и старый рефреш если тот прошёл проверку. Так же если ip изменился реализована моковая отправка уведомления.
3. Есть тесты на бд и сервис(основная логика проекта) и мини-тест на существование файла конфига.
4. /login и /regresh запросы необходимо осуществлять методом POST
5. Докер работает!
6. Если пользователь уже существует, не получится применить authorize с тем же самым userID
7. Можно использовать любую string как userID (т.е можно написать никнейм к примеру)
8. функция рефреш обновляет рефреш токен как того требует задание, и не позволяет использовать рефреш токен более 1 раза (как я понял это задание)
9. Тело запросов игнорируется, можно было бы userID, secret и передачу токена реализовать через Json, но в задании это не обговорено.
10. По хорошему secret в соответсвии с заданием должен генерироваться с помощью uuid.New().String(), но это не позволяет реализовать тест authorize, тк каждый раз будет генерироваться новый токен, поэтому вместо uuid.New().String() я вынес генерацию secret на внешнюю сторону, это допущение здесь только для реализации теста. (строка 44 файла authApp/auth/service.go)

Пример отправляемого запроса Authorize POST:
http://localhost:8050/login?user_id=1&secret=789
Пример ответа:
{"access":"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzQxMTUwNzEsImlwIjoiMTI3LjAuMC4xOjU5MTc1IiwicmVmcmVzaCI6ImV5SmhiR2NpT2lKSVV6VXhNaUlzSW5SNWNDSTZJa3BYVkNKOS5leUpsZUhBaU9pSXlNREkxTFRBeExURXlWREl4T2pNM09qVXhMamMyTVRZMU9UUXJNRE02TURBaUxDSnBjQ0k2SWpFeU55NHdMakF1TVRvMU9URTNOU0lzSW5ObFkzSmxkQ0k2SWpFeU15SXNJblZ6WlhKZmFXUWlPaUl4SW4wLnQtREFRcHduWFAzbnJ4YkQ4T1gzaG9iZGNWQTFNc2FWeVJmM0otNUN3NllOUXJUWDRvYjR6SDRvaE4tOUJkSDdiSTZIY0ZzV0FGODQ2ZWJ5S3JaS09nIiwidXNlcl9pZCI6IjEifQ.CL8Jzz9uzxFhDQWRYvSF2BEGJHJLEMVkq7DJFmVogu52R_lFIz2DSwa2PeroGJ7ptvIyrxj1GGGNB7wga0U_gQ","refresh":"VpZlEuSkZMVFQwRnVQblJkMnI4SGxxNU1FdVZtVHhmWFEydEF0em9odzNoaVhCTQ"}

Пример отправляемого запроса refresh POST:
http://localhost:8050/refresh?refresh_token=VpZlEuSkZMVFQwRnVQblJkMnI4SGxxNU1FdVZtVHhmWFEydEF0em9odzNoaVhCTQ&user_id=1
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