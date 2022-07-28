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
- Реализован ratelimiter с помощью Redis. Ограничение на количество запросов регулируется в этом куске кода
    ```
    func NewHTTPHandler(storage storage.Storage, limiterFactory *ratelimit.Factory) *HTTPHandler {
        return &HTTPHandler{
            storage: storage,
            // POST 10 действия в 10 сек
            postLimit: limiterFactory.NewLimiter("post_url", 10*time.Second, 10),
            // GET 20 действий в минуту
            getLimit: limiterFactory.NewLimiter("get_url", 1*time.Minute, 20),
        }
   }
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
        - Значение времени ограничения количество запросов специально завышены, для наглядности работы
        - (cache и ratelimiter реализован на основе курса [Дизайн систем](http://wiki.cs.hse.ru/%D0%94%D0%B8%D0%B7%D0%B0%D0%B9%D0%BD_%D1%81%D0%B8%D1%81%D1%82%D0%B5%D0%BC_21/22) ФКН НИУ ВШЭ)
- Unit-тесты(mocking)
    - TestGetURL
    - TestPutURL
- проект структурирован на освнове [go standarts](https://github.com/golang-standards/project-layout)
    