# cmd/server
Этот программный код реализует сервер, который слушает порт 8080 (по умолчанию) и ожидает подключения клиента. После подключения он сохраняет метрики памяти клиента как в памяти, так и в файле.

На главной странице пользователи могут увидеть список всех метрик, которые были сохранены в памяти.

## Параметры запуска
- -a (env: ADDRESS) - адрес сервера, по умолчанию localhost:8080
- -f (env: FILE_STORAGE_PATH) - путь к файлу, по умолчанию /tmp/metrics-db.json
- -i (env: STORE_INTERVAL) - интервал сохранения метрик в файл, по умолчанию 10 секунд
- -r (env: RESTORE) - флаг восстановления метрик из файла при запуске сервера, по умолчанию false
- -d (env: DATABASE_DSN) - параметр подключения к postgreSQL
- -k (env: KEY) - ключ для вычисления хеша ответов сервера
- -crypto-key (env: CRYPTO_KEY) - путь к ключу для шифрования данных
- -с ( -config, env: CONFIG) - путь к конфигурационному файлу (по умолчанию ./config/config.json)

При запуске сервер загружает все метрики из файла в память при работе с inmemory хранилищем или файлом, при работе с postgreSQL метрики хранятся в только в БД.

Есть два способа отправить метрики на сервер:

используя имя и значение в URL

используя тело JSON

Клиент отправляет метрики в следующем формате:

```
ID    string   `json:"id"`
MType string   `json:"type"`
Delta *int64   `json:"delta,omitempty"`
Value *float64 `json:"value,omitempty"`
```

Где:

- ID - уникальное имя метрики
- MType - тип метрики (счетчик или метрика)
- Delta - значение приращения счетчика
- Value - значение метрики

Сервер обрабатывает следующие запросы:

- POST "/update/{metric}/{key}/{value}" - обновляет метрику с заданным ключом и значением
- GET "/value/{metric}/{key}" - возвращает значение заданной метрики и ключа
- GET "/" - возвращает текущие значения всех метрик, сохраненных в памяти
- POST "/update/" - обновляет метрику с заданным телом JSON
- GET "/value/" - возвращает текущие значения заданной метрики в формате JSON.
- POST "/updates/" - обновляет метрики с заданным телом JSON в пакетном режиме.

-----------

This code implements a server that listens on port 8080 (by default) and waits for a client to connect. Once connected, it stores the client's memory metrics into both memory and a file.

On the main page, users can see a list of all metrics that were stored in memory.

## Launch parameters
- -a (env: ADDRESS) - server address, default localhost:8080
- -f (env: FILE_STORAGE_PATH) - path to the file, default /tmp/metrics-db.json
- -i (env: STORE_INTERVAL) - interval for saving metrics to a file, default 10 seconds
- -r (env: RESTORE) - flag to restore metrics from a file when the server starts, default false
- -d (env: DATABASE_DSN) - postgreSQL connection parameter
- -k (env: KEY) - key for calculating the hash of server responses
- -crypto-key (env: CRYPTO_KEY) - path to the key for encrypting data
- -с ( -config, env: CONFIG) - path to the configuration file (default ./config/config.json)

Upon start-up, the server loads all metrics from the file into memory when working with inmemory storage or a file, when working with postgreSQL, metrics are stored only in the database.

There are two ways to send metrics to the server:

by using name and value in the URL
by using a JSON body
The client sends metrics in the following format:

```
ID    string   `json:"id"`
MType string   `json:"type"`
Delta *int64   `json:"delta,omitempty"`
Value *float64 `json:"value,omitempty"`
```
Where:

- ID is the unique metric name
- MType is the metric type (counter or gauge)
- Delta is the counter increment value
- Value is the gauge value

The server handles the following requests:

- POST "/update/{metric}/{key}/{value}" - updates metric with the given key and value
- GET "/value/{metric}/{key}" - returns the value of the given metric and key
- GET "/" - returns current values of all metrics stored in memory
- POST "/update/" - updates metric with the given JSON body
- GET "/value/" - returns current values of the given metric in JSON format.
- POST "/updates/" - updates metrics with the given JSON body in batch mode.
