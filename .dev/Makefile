.PHONY: all setup run clean

run-compose:
	@echo "Starting services..."
	docker-compose up -d

run-cluster: create-cluster apply
	
create-cluster:
	@echo "Checking if cluster exists..."
	@if ! kind get clusters | grep -q "pipo-cluster"; then \
		echo "Creating new cluster..."; \
		kind create cluster --config ../infra/k8s/config/kind.yml --name=pipo-cluster; \
	else \
		echo "Cluster 'pipo-cluster' already exists. Skipping creation."; \
	fi

delete-cluster:
	@echo "Deleting cluster..."
	kind delete cluster --name=pipo-cluster

apply:
	@echo "Applying application config manifests..."
	kubectl apply -k ../infra/k8s/config
	@echo "Applying application shared manifests..."
	kubectl apply -k ../infra/k8s/shared
	@echo "Waiting for ingress-nginx controller pod to be ready..."
	$(MAKE) wait-for-ingress
	@echo "Applying application refinery manifests..."
	kubectl apply -k ../infra/k8s/refinery

wait-for-ingress:
	@echo "Waiting for ingress-nginx controller pod to be ready..."
	kubectl wait --namespace ingress-nginx \
		--for=condition=Ready pod \
		--selector=app.kubernetes.io/component=controller \
		--timeout=120s
