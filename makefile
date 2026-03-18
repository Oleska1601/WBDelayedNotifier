up:
	docker compose -f docker-compose.yaml --env-file .env up -d --build

down:
	docker compose -f docker-compose.yaml down