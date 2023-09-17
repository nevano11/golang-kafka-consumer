# Обработчик для очереди Kafka
## Задание
* [x] Сервис слушает очередь Apache Kafka FIO, присылавшую объекты Fio
* [x] В случае некорректных данных сервис отправляет сообщение в очередь FIOFAILED
* [x] Обогатить данные возрастом, полом и национальностью из сторонних API
* [x] Полученное сообщение сохранить в Postgres, поднимающимся миграциями
* [x] Выставить CRUD REST методы
* [ ] Выставить GraphQl методы аналогично с REST
* [ ] Предусмотреть кэширование данных в redis
* [x] Покрыть код логами
* [ ] Покрыть бизнес-логику unit-тестами
* [ ] Вынести все конфигурационные данные в .env

## Команды для kafka и migration
Kafka topics

    docker compose exec broker kafka-topics --create --topic FIO --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1
    docker compose exec broker kafka-topics --create --topic FIOFAILED --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1
Migrations

    migrate -database postgres://postgres:password@localhost:5002/humans?sslmode=disable -path db/migrations up
    migrate -database postgres://postgres:password@localhost:5002/humans?sslmode=disable -path db/migrations down
## Сервис-писатель kafka
    https://github.com/nevano11/golang-kafka-producer
    