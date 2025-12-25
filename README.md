# Idle Factory Game

a *Factorio* inspired idle game, made for the sole purpose of learning (and having some fun along the way). With this project I aim to apply concepts about Event Driven Architecture, Work Queues/Message Brokers (RabbitMQ) and Container Orchestration (Kubernetes)

# Build

first, make sure you have docker daemon, kubectl and minikube properly set up in your machine.

1. start minikube:
```sh
minikube start
```

2. set the environment to use minikube's docker daemon (if you are not on windows you can execute `minikube docker-env` and then run the command they provide):
```sh
& minikube -p minikube docker-env --shell powershell | Invoke-Expression
```

3. build the docker images inside minikube and make kubernetes deployments (if you are not on windows you can execute the commands from the script individually):
```sh
.\build-and-deploy.bat
```

# Usage

- you can check the logs from the state server by running the command:
    ```sh
    kubectl logs -f deployment/state-server
    ```

- you can check the rabbitmq dashboard by running the command, using the default credentials {guest:guest}:
    ```sh
    minikube service rabbitmq-service --url
    ```

- you can check kubernetes dashboard by running the command:
    ```sh
    minikube dashboard
    ```

After you are done you can `kubectl delete -f k8s.yaml --ignore-not-found=true` to stop the cluster.