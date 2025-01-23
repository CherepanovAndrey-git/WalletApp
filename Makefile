.PHONY: cdown cbuild cup clog restart install scratch help 

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

prune:
	docker-compose down -v --rmi all && docker system prune -f
	@echo " run make clog to check the logs and container status"

help:
	@echo "Available commands:"
	@echo "  make cdown    - Stop and remove Docker Compose services"
	@echo "  make cbuild   - Build Docker Compose services"
	@echo "  make cup      - Start Docker Compose services in detached mode"
	@echo "  make clog     - Follow logs for the 'app' service"
	@echo "  make restart  - Restart services (down, build, up)"
	@echo "  make prune  - Restarts all from scratch, wipes containers, cash and images"