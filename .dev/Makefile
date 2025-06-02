.PHONY: all setup run clean

all: run

run:
	@echo "Starting services..."
	docker-compose up -d --build

clean:
	@echo "Cleaning up..."
	docker-compose down -v
	@echo "Cleanup complete!"