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
```

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