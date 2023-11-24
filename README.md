# es7-test

## Running
`sysctl -w vm.max_map_count=262144` -> in windows so es can start

Index name for local is `book`

## Sample config
indexconsumer
```
server:
  port: 8081
  env: "development"
  app_name: "indexconsumer"

nr: <CAN BE OMITTED>
  license_key: <<FILL>>
  acc_id: <<FILL>>
  buffer_time_ms: 250
  buffer_size: 100

es:
  vm: 
    host: "http://localhost:9210"
    index_name: "book"
  kube: 
    host: "http://localhost:9210"
    index_name: "book"

nsq:
  cvm: 
    nsqd_address: "localhost:4150"
    publish_topic: "doc_index"
    consumer_name: "indexer_book_vm"
    timeout_ms: 1000
    num_of_consumers: 50
    in_flight: 50
  ckube: 
    nsqd_address: "localhost:4150"
    publish_topic: "doc_index"
    consumer_name: "indexer_book_kube"
    timeout_ms: 1000
    num_of_consumers: 50
    in_flight: 50
```


indexhandlerapp
```
server:
  port: 8080
  env: "development"
  app_name: "indexapphandler"

nr: <CAN BE OMITTED>
  license_key: <<FILL>>
  acc_id: <<FILL>>
  buffer_time_ms: 250
  buffer_size: 100

nsq:
  nsqd_address: "localhost:4150"
  publish_topic: "doc_index"

es:
  vm: 
    host: "http://localhost:9210"
    index_name: "book"
  kube: 
    host: "http://localhost:9210"
    index_name: "book"
```

fetchhandler
```
server:
  port: 8082
  env: "development"
  app_name: "fetchhandler"

nr: <CAN BE OMITTED>
  license_key: <<FILL>>
  acc_id: <<FILL>>
  buffer_time_ms: 250
  buffer_size: 100

es:
  vm: 
    host: "http://localhost:9210"
    index_name: "book"
  kube: 
    host: "http://localhost:9210"
    index_name: "book"
handler:
  vm: 
    timeout_ms: 1000
```

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
| `kubectl port-forward es-master-0 9200:9200 -n es-operator-demo` | forward minikube |
| `kubectl -n es-operator-demo get pods` | check pods|
| `kubectl apply -f docs/elasticsearchdataset-simple.yaml` | apply yaml|
| `kubectl get deployments -n es-operator-demo` | check deployments|
| `kubectl -n es-operator-demo logs es-data-simple-0 -c elasticsearch` | check logs |
| `kubectl get services --all-namespaces` | check all services |
| `kubectl describe service es-http -n es-operator-demo` | describe service |
| `kubectl scale statefulset/es-data-simple --replicas=2 -n es-operator-demo` | scale


https://onenr.io/08jqZkP3xQl