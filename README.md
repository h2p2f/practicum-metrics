##



## Practicum-metrics

эти приложения разработаны для курса Yandex.practicum - Go-developer

этот репозиторий содержит два приложения: сервер и агент

server - это http-сервер, который слушает порт 8080 (по умолчанию) и ожидает подключения клиента. После подключения он сохраняет метрики памяти клиента как в памяти, так и в файле или базе данных.

agent - это агент, который отправляет runtime метрики на сервер.

## История изменений:

### v.0.0.1

- iter 1: Разработан простой http-сервер, который слушает порт 8080 (по умолчанию) и ожидает подключения клиента в формате ```http://<server_addr>/update/<metric_type>/<metrics_name>/<value>```. После подключения он сохраняет метрики памяти клиента в памяти.
- iter 2: Разработан агент, который отправляет runtime метрики на сервер.
- iter 3: Разработан обработчик сервера для формата ```http://<server_addr>/value/<metric_type>/<metrics_name>/<value>```. Этот обработчик возвращает значение указанной метрики и ключа.
- iter 4: добавлена поддержка флагов параметров запуска для сервера и агента
- iter 5: добавлена поддержка переменных окружения для сервера и агента
- iter 6: добавлено middleware для логгирования событий сервера
- iter 7: добавлена обработка метрик в формате JSON для сервера
- iter 8: добавлены функции сжатия и распаковки gzip для сервера и агента
- iter 9: добавлено файловое хранилище для сервера
- iter 10-11: добавлено хранилище postgreSQL для сервера
- iter 12: добавлена пакетная отправка метрик для агента, сервер может обрабатывать метрики в пакетном режиме
- iter 13: добавлена обработка специфичных ошибок для сервера и агента
- iter 14: добавлен расчет хеша данных для сервера и агента
- iter 15: разработан пул воркеров для агента

### v.0.2.1

- полностью переработан код. Теперь inmemory хранилище и БД используют один и тот же интерфейс. Файловое хранилище работает только с inmemory хранилищем.
- Код стал чище и более понятным (субъективно).
- iter 16 - добавлены бенчмарк тесты для хендлеров сервера. Добавлены эндпоинты для ```pprof```, проанализирован код и выполнена его оптимизация. Diff находится в файле ```profiles/diff.pprof```
- iter 17 - код отформатирован при помощи ```gofmt``` и ```goimports```
- iter 18 - Написана документация для каждого пакета, для хендлеров - с примерами использования
- iter 19 - написан кастомный линтер, проверяющий использование ```os.Exit()``` в ```main``` функциях
- iter 20 - добавлены переменные для версионирования приложения при сборке
- iter 21 - добавлена возможность работы агента и сервера с шифрацией данных RSA ключами. Ключи генерируются ```cmd/server/crypokeygenerator/keygen``` Агент использует публичный RSA ключ, обработка данных сервером осуществяется middleware c приватным RSA ключом. Ответы сервера не шифруются. Ограничения: комбинированный режим не поддерживается, т.е. агент может отправлять данные только в зашифрованном виде, либо в открытом. Аналогично с сервером - он может принимать только зашифрованные данные, либо открытые.
- iter 22 - добавлена возможность обработки пользовательской конфигурации агента и сервера из файла в формате JSON. При запуске агента и сервера с флагом ```-config``` или ```-c``` будет использована конфигурация из файла. 
- iter 23 - подключена обработка сигналов ОС терминации приложения для корректного завершения работы агента и сервера. 
- --out of scope-- для отладки сервера интерфейсы работы с хранилищем в хендлнрах обернуты при помощи ```gowrap```, написан новый темплейт ```templates/gowrap/zap``` для логгера от Uber

--------------------

these apps developed for Yandex.practicum on the Go-developer course

this repo contains two apps: server and agent

server is an http server that listens on port 8080 (default) and waits for a client connection. Once connected, it stores the client's memory metrics both in memory and in a file or database.

agent - is an agent that sends runtime metrics to the server.

## Changelog:

### v.0.0.1

- iter 1: Developed simply http server that listens on port 8080 (default) and waits for a client connection in ```http://<server_addr>/update/<metric_type>/<metrics_name>/<value>``` format. Once connected, it stores the client's memory metrics in memory.
- iter 2: Developed agent that sends runtime metrics to the server.
- iter 3: Developed server's handler for ```http://<server_addr>/value/<metric_type>/<metrics_name>/<value>``` format. This handler returns the value of the specified metric and key.
- iter 4: added support flags startup parameters for server and agent
- iter 5: added support environment variables for server and agent
- iter 6: added logging middleware for server
- iter 7: added handle metrics in JSON format for server
- iter 8: added compress and decompress functions for server and agent
- iter 9: added file storage for server
- iter 10-11: added postgreSQL storage for server
- iter 12: added batch sending of metrics for agent and server can handle metrics in batch mode
- iter 13: added handle specific errors for server and agent
- iter 14: added data hash calculation for server and agent
- iter 15: developed worker pool for agent

### v.0.2.1

- fully refactored code. Now inmemory storage and db uses the same interface. File storage working only with inmemory storage.
- Code became cleaner and more understandable (subjectively).
- iter 16 - added benchmark tests for server handlers. Added endpoints for pprof, analyzed the code and optimized it. Diff is in the file ```profiles/diff.pprof```
- iter 17 - code formatted with ```gofmt``` and ```goimports```
- iter 18 - Documentation written for each package, for handlers - with examples of use
- iter 19 - written custom linter that checks the use of ```os.Exit()``` in ```main``` functions
- iter 20 - added variables for application versioning at build
- iter 21 - added the ability for the agent and server to work with data encryption using RSA keys. Key will be generated by ```cmd/server/crypokeygenerator/keygen```. The agent uses a public RSA key, data processing by the server is carried out by middleware with a private RSA key. Server responses are not encrypted. Limitations: combined mode is not supported, i.e. the agent can only send data in encrypted or open form. Similarly with the server - it can only accept encrypted data or open data.
- iter 22 - added the ability to process user configuration of the agent and server from a file in JSON format. When the agent and server are started with the ```-config``` or ```-c``` flag, the configuration from the file will be used.
- iter 23 - connected processing of OS termination signals for correct termination of the agent and server.
- --out of scope-- for debugging the server, the interfaces for working with the storage in the handlers are wrapped using ```gowrap```, a new template ```templates/gowrap/zap``` is written for the Uber logger



