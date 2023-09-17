docker compose exec broker kafka-topics --create --topic FIO --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1
docker compose exec broker kafka-topics --create --topic FIOFAILED --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1

migrate -database postgres://postgres:password@localhost:5002/humans?sslmode=disable -path db/migrations up
migrate -database postgres://postgres:password@localhost:5002/humans?sslmode=disable -path db/migrations down