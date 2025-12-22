REM Build images in Minikube's Docker
minikube image build -f cmd/workerserver/drill/Dockerfile -t idle-factory-drill:latest .
minikube image build -f cmd/workerserver/furnace/Dockerfile -t idle-factory-furnace:latest .
minikube image build -f cmd/workerserver/lab/Dockerfile -t idle-factory-lab:latest .
minikube image build -f cmd/stateserver/Dockerfile -t idle-factory-stateserver:latest .

REM Deploy to Kubernetes
kubectl apply -f k8s.yaml
kubectl rollout restart deployment/idle-factory-state-server 
kubectl rollout restart deployment/idle-factory-drill-worker
kubectl rollout restart deployment/idle-factory-furnace-worker
kubectl rollout restart deployment/idle-factory-lab-worker
