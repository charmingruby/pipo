.PHONY: all setup run clean

local:
	@echo "Starting services..."
	docker-compose up -d --build

cluster:
	@echo "Starting cluster..."
	kind create cluster --config ../infra/cluster/config/kind.yml --name=pipo-cluster
