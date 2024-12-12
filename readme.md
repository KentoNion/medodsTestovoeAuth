# Фичи и особенности
1. Добавил конфиг файл, строка host_db игнорируется и перезаписывается докером через внешнюю среду при запуске приложения из докера. /authApp/config.yaml
2. /refresh выдаёт 2 токена, новый аксес и старый рефреш если тот прошёл проверку. Так же если ip изменился реализована моковая отправка уведомления.
3. Есть тесты на бд и сервис(основная логика проекта) и мини-тест на существование файла конфига.
4. Сервису не важно каким методом вы в него стучитесь (особенность)

Пример отправляемого запроса Authorize
http://localhost:8050/login?user_id=1&secret=123
Пример ответа:
{"access":"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzQxMTUwNzEsImlwIjoiMTI3LjAuMC4xOjU5MTc1IiwicmVmcmVzaCI6ImV5SmhiR2NpT2lKSVV6VXhNaUlzSW5SNWNDSTZJa3BYVkNKOS5leUpsZUhBaU9pSXlNREkxTFRBeExURXlWREl4T2pNM09qVXhMamMyTVRZMU9UUXJNRE02TURBaUxDSnBjQ0k2SWpFeU55NHdMakF1TVRvMU9URTNOU0lzSW5ObFkzSmxkQ0k2SWpFeU15SXNJblZ6WlhKZmFXUWlPaUl4SW4wLnQtREFRcHduWFAzbnJ4YkQ4T1gzaG9iZGNWQTFNc2FWeVJmM0otNUN3NllOUXJUWDRvYjR6SDRvaE4tOUJkSDdiSTZIY0ZzV0FGODQ2ZWJ5S3JaS09nIiwidXNlcl9pZCI6IjEifQ.CL8Jzz9uzxFhDQWRYvSF2BEGJHJLEMVkq7DJFmVogu52R_lFIz2DSwa2PeroGJ7ptvIyrxj1GGGNB7wga0U_gQ","refresh":"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDI1LTAxLTEyVDIxOjM3OjUxLjc2MTY1OTQrMDM6MDAiLCJpcCI6IjEyNy4wLjAuMTo1OTE3NSIsInNlY3JldCI6IjEyMyIsInVzZXJfaWQiOiIxIn0.t-DAQpwnXP3nrxbD8OX3hobdcVA1MsaVyRf3J-5Cw6YNQrTX4ob4zH4ohN-9BdH7bI6HcFsWAF846ebyKrZKOg"}

Пример отправляемого запроса refresh:
localhost:8050/refresh
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