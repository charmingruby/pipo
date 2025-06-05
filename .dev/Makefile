.PHONY: all setup run clean

run-compose:
	@echo "Starting services..."
	docker-compose up -d

run-cluster: create-cluster apply
	
create-cluster:
	@echo "Checking if cluster exists..."
	@if ! kind get clusters | grep -q "pipo-cluster"; then \
		echo "Creating new cluster..."; \
		kind create cluster --config ../infra/cluster/config/kind.yml --name=pipo-cluster; \
	else \
		echo "Cluster 'pipo-cluster' already exists. Skipping creation."; \
	fi

delete-cluster:
	@echo "Deleting cluster..."
	kind delete cluster --name=pipo-cluster

apply:
	@echo "Applying resources..."
	kubectl apply -k ../infra/cluster
