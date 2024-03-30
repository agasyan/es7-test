# es7-test

## Running
`sysctl -w vm.max_map_count=262144` -> in windows so es can start

Index name for local is `book`

## Sample config
Refer to `*.sample` file

sample curl:
```
curl --request GET \
  --url 'http://localhost:8082/get/vm?title=brothers' \
  --header 'User-Agent: insomnium/0.2.3-a'
```

indexhandler
```
curl --request POST \
  --url http://localhost:8080/random-action-es \
  --header 'Content-Type: application/json' \
  --data '{
	"count": 5
}'
```

## Sample Command Load test
Indexer
```

```

Fetcher
```

```

### Running

Docker
```
sudo sysctl -w vm.max_map_count=262144
docker-compose up
```

Minikube
```
minikube start
kubectl port-forward es-master-0 9200:9200 -n es-operator-demo
```

# Install

## Kubectl
```
curl https://storage.googleapis.com/kubernetes-release/release/stable.txt > ./stable.txt
export KUBECTL_VERSION=$(cat stable.txt)
curl -LO https://storage.googleapis.com/kubernetes-release/release/$KUBECTL_VERSION/bin/linux/amd64/kubectl
chmod +x ./kubectl
sudo mv ./kubectl /usr/local/bin/kubectl
mkdir -p ~/.kube
ln -sf "/mnt/c/users/$USER/.kube/config" ~/.kube/config
rm ./stable.txt
```

Minikube -> https://gist.github.com/wholroyd/748e09ca0b78897750791172b2abb051

## ES operator

1. clone repo
```
git clone https://github.com/zalando-incubator/es-operator.git
cd es-operator
```

# Useful command 

## minikube

| Command                                           | Description                                                       |
| ------------------------------------------------- | ----------------------------------------------------------------- |
| `minikube start` | starting minikube |
| `minikube stop` | stop minikube pods|


## kubectl

| Command                                           | Description                                                       |
| ------------------------------------------------- | ----------------------------------------------------------------- |
| `kubectl port-forward es-master-0 9200:9200 -n es-operator-demo` | forward minikube to localhost |
| `kubectl -n es-operator-demo get pods` | check pods|
| `kubectl apply -f docs/elasticsearchdataset-simple.yaml` | apply yaml|
| `kubectl get deployments -n es-operator-demo` | check deployments|
| `kubectl -n es-operator-demo logs es-data-simple-0 -c elasticsearch` | check logs |
| `kubectl get services --all-namespaces` | check all services |
| `kubectl describe service es-http -n es-operator-demo` | describe service |
| `kubectl scale statefulset/es-data-simple --replicas=2 -n es-operator-demo` | scale


https://onenr.io/08jqZkP3xQl