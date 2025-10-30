docker-compose --env-file .env up -d

goose -dir ./internal/database/repo/migrations -s create add_notifications_index sql


goose -dir ./internal/database/repo/migrations -s create new_parsings sql  

docker build -t htmlparser -f deploy/Dockerfile . 
docker-compose -f ./deploy/docker-compose.yml up -d


go run github.com/vektra/mockery/v2@latest --name=CacheInterface --dir=./internal/usecase/interface.go