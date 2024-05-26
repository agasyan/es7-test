# Install Elasticsearch

Source: https://www.digitalocean.com/community/tutorials/how-to-install-and-configure-elasticsearch-on-ubuntu-22-04

```
curl -fsSL https://artifacts.elastic.co/GPG-KEY-elasticsearch | sudo gpg --dearmor -o /usr/share/keyrings/elastic.gpg

echo "deb [signed-by=/usr/share/keyrings/elastic.gpg] https://artifacts.elastic.co/packages/7.x/apt stable main" | sudo tee -a /etc/apt/sources.list.d/elastic-7.x.list

sudo apt update
```

After that install speicific version (7.7.1) and setup the cluster

```
sudo apt-get install elasticsearch=7.7.1

<<BEFORE CONTINUE EDIT the elasticsearch.yml>>
sudo nano /etc/elasticsearch/elasticsearch.yml

sudo systemctl start elasticsearch
sudo systemctl enable elasticsearch

<<ANY CHANGES>>
sudo systemctl restart elasticsearch
```

## Log

```
cd /var/log/elasticsearch/
```

## Setup ES Cluster
https://stackoverflow.com/questions/62856133/add-nodes-to-make-local-cluster-elasticsearch-7-8

### Master Setup
Things to note
- Need to add `cluster-name` on the yml
- Need to add `network.host: 0.0.0.0` on the yml to allow remote access from other hosts
- Need to setup node name `node.name: es-vm-master-1`
- Need to add role for master below for data only `node.master: false`
```
node.master: true
node.data: false
```
- Need to setup host and initial master node
```
discovery.seed_hosts:
  - 157.230.33.188
cluster.initial_master_nodes:
  - es-vm-master-1
```
- 

***

# Kube

## doctl
https://docs.digitalocean.com/reference/doctl/how-to/install/

For ctl and connect between kubectl and doctl

```
sudo snap install doctl
doctl auth init --context agas-main
doctl auth list
doctl auth switch --context agas-main
```


Kube pointer
```
kubectl get nodes
kubectl config get-contexts
kubectl config use-context do-sgp1-es7-kube-2
kubectl config delete-context do-sgp1-es7-kube-2
```

## Kube
https://github.com/zalando-incubator/es-operator
https://github.com/zalando-incubator/es-operator/blob/master/docs/GETTING_STARTED.md

### Setup and setup operator

#### Setup service acc

```
kubectl create namespace es-operator-demo
kubectl create serviceaccount es-operator -n es-operator-demo
kubectl get serviceaccount es-operator -n es-operator-demo

```

```
kubectl apply -f docs/cluster-roles.yaml
kubectl apply -f docs/zalando.org_elasticsearchdatasets.yaml
kubectl apply -f docs/zalando.org_elasticsearchmetricsets.yaml
kubectl apply -f es-operator-prod.yaml
kubectl -n es-operator-demo get pods
```

### Master
```
kubectl apply -f docs/es-master-prod.yaml
MASTER_POD=$(kubectl -n es-operator-demo get pods -l application=elasticsearch,role=master -o custom-columns=:metadata.name --no-headers | head -n 1)
kubectl -n es-operator-demo port-forward $MASTER_POD 9200
```

### Data
```
kubectl apply -f docs/es-data-prod.yaml
kubectl -n es-operator-demo get eds
kubectl -n es-operator-demo get sts
kubectl -n es-operator-demo get pods
```


*** 

# IGNORE BELOW
### v 7.7.1 - linux / mac
```
wget https://artifacts.elastic.co/downloads/elasticsearch/elasticsearch-7.7.1-linux-x86_64.tar.gz
wget https://artifacts.elastic.co/downloads/elasticsearch/elasticsearch-7.7.1-linux-x86_64.tar.gz.sha512
shasum -a 512 -c elasticsearch-7.7.1-linux-x86_64.tar.gz.sha512 
tar -xzf elasticsearch-7.7.1-linux-x86_64.tar.gz
cd elasticsearch-7.7.1/ 

```

### Networking
https://www.digitalocean.com/community/tutorials/initial-server-setup-with-ubuntu

Skip UFW If already using do firewall
```
sudo ufw status
sudo ufw enable
```