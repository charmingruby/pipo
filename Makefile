.PHONY: all setup run clean

run-compose:
	@echo "Starting services..."
	docker-compose up -d

run-cluster: create-cluster apply
	
create-cluster:
	@echo "Checking if cluster exists..."
	@if ! kind get clusters | grep -q "pipo-cluster"; then \
		echo "Creating new cluster..."; \
		kind create cluster --config ./infra/k8s/config/kind.yml --name=pipo-cluster; \
	else \
		echo "Cluster 'pipo-cluster' already exists. Skipping creation."; \
	fi

apply:
	
	@echo "Applying metrics-server..."
	kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
	@echo "Applying application config manifests..."
	kubectl apply -k ./infra/k8s/config
	@echo "Applying application shared manifests..."
	kubectl apply -k ./infra/k8s/shared
	
	@echo "Waiting for ingress-nginx controller pod to be ready..."
	$(MAKE) wait-for-ingress
	
	
	@echo "Applying application refinery manifests..."
	kubectl apply -k ./infra/k8s/refinery

wait-for-ingress:
	kubectl wait --namespace ingress-nginx \
		--for=condition=Ready pod \
		--selector=app.kubernetes.io/component=controller \
		--timeout=120s

enable-metrics-server:
	kubectl patch deployment metrics-server -n kube-system \
		--type='json' \
		-p='[{"op": "add", "path": "/spec/template/spec/containers/0/args/-", "value": "--kubelet-insecure-tls"}]'
	

delete-cluster:
	@echo "Deleting cluster..."
	kind delete cluster --name=pipo-cluster


# -
# Docker
# -

.PHONY: run-compose
run-compose:
	@echo "Starting services..."
	docker-compose up -d

# -
# Kubernetes
# -

.PHONY: run-cluster
run-cluster: create-cluster apply

.PHONY: create-cluster
create-cluster:
	@echo "Checking if cluster exists..."
	@if ! kind get clusters | grep -q "pipo-cluster"; then \
		echo "Creating new cluster..."; \
		kind create cluster --config ./infra/k8s/config/kind.yml --name=pipo-cluster; \
	else \
		echo "Cluster 'pipo-cluster' already exists. Skipping creation."; \
	fi

.PHONY: apply
apply:
	@echo "Applying metrics-server..."
	kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml	

	@echo "Applying application config manifests..."
	kubectl apply -k ./infra/k8s/config
	
	@echo "Applying application shared manifests..."
	kubectl apply -k ./infra/k8s/shared
	
	@echo "Waiting for ingress-nginx controller pod to be ready..."
	$(MAKE) wait-for-ingress	

	@echo "Enabling metrics-server..."
	$(MAKE) enable-metrics-server

	@echo "Applying application refinery manifests..."
	kubectl apply -k ./infra/k8s/refinery

.PHONY: wait-for-ingress
wait-for-ingress:
	kubectl wait --namespace ingress-nginx \
		--for=condition=Ready pod \
		--selector=app.kubernetes.io/component=controller \
		--timeout=120s

.PHONY: enable-metrics-server
enable-metrics-server:
	kubectl patch deployment metrics-server -n kube-system \
		--type='json' \
		-p='[{"op": "add", "path": "/spec/template/spec/containers/0/args/-", "value": "--kubelet-insecure-tls"}]'
	
.PHONY: delete-cluster
delete-cluster:
	@echo "Deleting cluster..."
	kind delete cluster --name=pipo-cluster
