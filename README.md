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

es:
  vm: 
    host: "http://localhost:9200"
    index_name: "book"

nsq:
  cvm: 
    nsqd_address: "localhost:4150"
    publish_topic: "doc_index"
    consumer_name: "indexer_book_vm"
    timeout_ms: 1000
    num_of_consumers: 10
    in_flight: 10
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

nsq:
  nsqd_address: "localhost:4150"
  publish_topic: "doc_index"

es:
  vm: 
    host: "http://localhost:9200"
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

es:
  vm: 
    host: "http://localhost:9200"
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

2. create namespace
```
kubectl create namespace es-operator
```

