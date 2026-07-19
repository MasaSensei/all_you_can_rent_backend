APP_NAME=rentos

docker-up:
	docker compose -f deployments/docker-compose.yml up -d

docker-down:
	docker compose -f deployments/docker-compose.yml down

docker-logs:
	docker compose -f deployments/docker-compose.yml logs -f

docker-restart:
	docker compose -f deployments/docker-compose.yml down
	docker compose -f deployments/docker-compose.yml up -d

docker-ps:
	docker compose -f deployments/docker-compose.yml ps