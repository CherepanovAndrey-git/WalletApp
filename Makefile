.PHONY: cdown cbuild cup clog restart install scratch help test-api swagger

cdown:
	docker-compose down

cbuild:
	docker-compose build

cup:
	docker-compose up -d

clog:
	docker ps && docker-compose logs -f app

restart: cdown cbuild cup
	@echo "Services restarted successfully!"

install: cbuild cup clog

test-api:
	@echo "Running API tests..."
	@./scripts/test_api.sh

prune:
	docker-compose down -v --rmi all && docker system prune -f
	

swagger:
	swag init -g cmd/server/main.go -o swagger/docs
	
help:
	@echo "Available commands:"
	@echo "  make install - Build and start services"
	@echo "  make cdown    - Docker Compose down"
	@echo "  make cbuild   - Docker Compose build"
	@echo "  make cup      - Docker Compose up, detached mode"
	@echo "  make clog     - Docker Compose logs"
	@echo "  make restart  - Restart services (down, build, up)"
	@echo "  make test-api - Run API tests"
	@echo "  make prune    - Full wipe and clean installation. Deletes containers, volumes and images, and system prune"

	
	