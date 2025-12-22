# Idle Factory Game

a *Factorio* inspired idle game, made for the sole purpose of learning (and having some fun along the way). With this project I aim to apply concepts about Event Driven Architecture, Work Queues/Message Brokers (RabbitMQ) and Container Orchestration (Kubernetes)

# Build

to build this application you first need to ensure you can run kubectl and minikube.

firt start minikube
minikube start

then, set the environment to use minikube's docker daemon
& minikube -p minikube docker-env --shell powershell | Invoke-Expression

to see the logs irl
kubectl logs -f deployment/state-server