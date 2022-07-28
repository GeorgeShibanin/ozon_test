# Ozon-test
## Как запустить?
Запуск из корня
```
docker-compose up --build
```
## Описание:
- В качестве роутера используется [gorilla/mux](github.com/gorilla/mux)
- Пример POST запроса на сокращение изначальной ссылки
  ```
    curl -d '{"url": "http://google.com"}' -X POST http://localhost:8081/urls
  ```
  Ответ:
  ```
    {"shorturl":"HMR5f91Qth"}
  ```
- Пример GET запроса на выдачу изначальной ссылки по сокращённой
    ```
      curl -X GET  http://localhost:8080/HMR5f91Qth
    ```
  Ответ:
    ```
      {"url":"http://google.com"}
    ```

- Реализованы три типа хранения данных, которые выбираются при запуске приложения через переменную окружения STORAGE_MODE=
    - **in_memory**, где пары (shor_url, original_url) хранятся в структуре вида
      ```
      type inMemoryStore struct {
	      mutex sync.RWMutex
	      store map[storage.ShortedURL]storage.URL
      }
      ```
    - **postgres**, где пары (shor_url, original_url) хранятся в базе данных. База данных содержит одну таблицу links с полями id в качестве ключа и поле url
    - **redis**, где пары (shor_url, original_url) хранятся в базе данных с подключение к redis server,
      который реализует cache и позволяет ограничивать нагрузку на сервис.
        - Данные добавляются в cache при добавлении новой ссылки в базу,
          а так же данные добавляются в cache при get запросе, при условии что этой ссылки нет в cache.
        - Ограничение нагрузки на сервер, просиходит путём ограничения количества запросов. Запрос перестаёт обрабатываться
          если уже было много запросов этого типа с этим индентификатором
    - Unit-тесты
        - 1
        - 2
        - 3
    - проект структурирован на освнове [go standarts](https://github.com/golang-standards/project-layout)
    