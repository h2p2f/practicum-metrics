# cmd/agent

Этот программный код реализует агента, который отправляет runtime метрики на сервер.

## Параметры запуска

Агент поддерживает следующие параметры запуска:
- -r (env: REPORT_INTERVAL) - интервал отправки метрик на сервер (по умолчанию 10 секунд)
- -p (env: POOL_INTERVAL) - интервал сбора метрик (по умолчанию 2 секунды)
- -a (env: SERVER_ADDRESS) - адрес сервера (по умолчанию http://localhost:8080)
- -k (env: KEY) - ключ, при наличии которого вычисляется хеш отправляемых данных (по умолчанию пустая строка)
- -l (env: RATE_LIMIT) - ограничение на количество воркеров при отправке метрик (по умолчанию 2). Если параметр не указан явно - используется пактная отправка метрик без пула воркеров.
- -crypto-key (env: CRYPTO_KEY) - путь к ключу для шифрования данных
- -с ( -config, env: CONFIG) - путь к конфигурационному файлу (по умолчанию ./config/config.json)
-----------------

This code implements an agent that sends runtime metrics to the server.

## Launch parameters

The agent supports the following launch parameters:
- -r (env: REPORT_INTERVAL) - interval for sending metrics to the server (default 10 seconds)
- -p (env: POOL_INTERVAL) - interval for collecting metrics (default 2 seconds)
- -a (env: SERVER_ADDRESS) - server address (default http://localhost:8080)
- -k (env: KEY) - key, if present, the hash of the sent data is calculated (default empty string)
- -l (env: RATE_LIMIT) - limit on the number of workers when sending metrics (default 2) If the parameter is not specified explicitly, batch sending of metrics without a worker pool is used.
- -crypto-key (env: CRYPTO_KEY) - path to the key for encrypting data
- -с ( -config, env: CONFIG) - path to the configuration file (default ./config/config.json)
