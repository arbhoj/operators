—————————————————

Install
KUDO
    Install kudo-cli
                        brew tap kudobuilder/tap 
                        brew install kudo-cli
                              Or
                        go get github.com/kudobuilder/kudo
                        cd ~/go/src/github.com/kudobuilder/kudo/
                        make cli-install
    Deploy Kudo to Kubernetes cluster              
                    kubectl kudo init 

Operator SDK 
   SDK
                        RELEASE_VERSION=v0.12.0
                        cd /usr/local/bin  
                        wget https://github.com/operator-framework/operator-sdk/releases/download/${RELEASE_VERSION}/operator-sdk-${RELEASE_VERSION}-x86_64-linux-gnu
                        mv operator-sdk-${RELEASE_VERSION}-x86_64-linux-gnu operator-sdk
                        chmod +x operator-sdk
                        go mod tidy; go mod vendor
   OLM
                        curl -sL https://github.com/operator-framework/operator-lifecycle-manager/releases/download/${RELEASE_VERSION}/install.sh | bash -s 0.12.0


———————————


Build Operator
      Kudo 
           Create a directory and manually create parameters.yaml, operator.yaml and templates dir with templatized k8s manifests

      SDK
             Helm Create:
             operator-sdk new helm-new --api-version=helmnew.operator.com/v1alpha1 --kind=HelmNew --type=helm 

Ansible Create:
operator-sdk new ansible-new --api-version=ansiblenew.operator.com/v1alpha1 --kind=NginxAnsibleNew --type=ansible

Go Create:
operator-sdk new sdk-new --repo github.com/arbhoj/sdk-operator-new

cd sdk-new
 
operator-sdk add api --api-version=sdk-operator-new.com/v1alpha1 --kind=NginxSDKNew

operator-sdk add controller --api-version=sdk-operator-new.com/v1alpha1 --kind=NginxSDKNew

operator-sdk generate k8s # if <operator>_types.go is modified 

Build:

cd ansible-new
operator-sdk build arbhoj/ansible-operator-new:v1.0.5

Docker push arbhoj/ansible-operator-new:v1.0.5

              Update operator.yaml with the docker image 
              sed -i "" 's|REPLACE_IMAGE|arbhoj/ansible-operator-new:v1.0.5|g' deploy/operator.yaml 

——————————————

Package Operator

KUDO

cp -r my-operator-v3 my-operator-v4

vi my-operator-v4/operator.yaml (update version)

kubectl kudo package create my-operator-v4 --destination=repo/
      Package created: ~/operators/kudo-demo/repo/my-kudo-operator-1.0.4.tgz

kubectl kudo repo index repo --overwrite

docker run -d -v ~/operators/kudo-demo/repo:/usr/share/nginx/html -p 80:80 nginx:latest

kubectl kudo repo add local http://localhost
kubectl kudo repo list

curl http://localhost/index.yaml


SDK

cd ~/operators/helm-operator
•Create OperatorGroup CR for specifying the targetNamespace(s) for the operator (one time)
    kubectl apply -f OperatorGroup_For_OLM.yaml
•Create a CSV (Cluster Service Version - clusterserviceversion.yaml file created locally)
operator-sdk olm-catalog gen-csv --csv-version 1.0.1 --update-crds



search for placeholder and replace it with helmnginx (i.e. the namespace where the operator will be deployed)
sed -i "" 's|placeholder|helmnginx|g' deploy/olm-catalog/helm-operator/1.0.1/helm-operator.v1.0.1.clusterserviceversion.yaml 

Add description and displayName to all the customresourcedefinitions (open file search for owned and add the following to each) 
 vi deploy/olm-catalog/helm-operator/1.0.1/helm-operator.v1.0.1.clusterserviceversion.yaml 

description: This creates a sample nginx deployment
displayName: HelmAnsible
    
To Upgrade operator-sdk olm-catalog gen-csv --csv-version <version> --update-crds [--from-version <older-version>]

operator-sdk olm-catalog gen-csv --csv-version 1.0.2 --update-crds --from-version 1.0.1


Deploy Operator
 
KUDO

watch kubectl get all -n mykudo

kubectl kudo install my-kudo-operator --instance=mykudotest --version=1.0.2 --namespace=mykudo

kubectl get operator -n mykudo
kubectl get operatorversion -n mykudo
kubectl get instance -n mykudo

curl http://34.212.226.38:30080    

Additional instances of the same operator can be deployed by creating a “kudo.dev/v1beta1
Instance” GVK cr.

kubectl apply -f testinstance.yaml 



SDK

helm
kubectl apply -f deploy/service_account.yaml -n helmnginx
kubectl apply -f deploy/role.yaml -n helmnginx
kubectl apply -f deploy/role_binding.yaml -n helmnginx
kubectl apply -f deploy/olm-catalog/helm-operator/1.0.1/ -n helmnginx


ansible
kubectl apply -f deploy/service_account.yaml -n na
kubectl apply -f deploy/role.yaml -n na
kubectl apply -f deploy/role_binding.yaml -n na
kubectl apply -f deploy/olm-catalog/ansible-operator/1.0.1/ -n na

go
kubectl apply -f deploy/service_account.yaml -n sdk
kubectl apply -f deploy/role.yaml -n sdk
kubectl apply -f deploy/role_binding.yaml -n sdk
kubectl apply -f deploy/olm-catalog/sdk-operator/1.0.1/ -n sdk




Create Operator Instance

KUDO 
kubectl apply -f testinstance.yaml 
—————————————————————————————
apiVersion: kudo.dev/v1beta1
kind: Instance
metadata:
  labels:
    controller-tools.k8s.io: "1.0"
    kudo.dev/operator: my-kudo-operator
  name: new
  namespace: mykudo
spec:
  operatorVersion:
    name: my-kudo-operator-1.0.3
  parameters:
    defaultBody: From CR
    nodePort: "31000"
—————————————————————————————————

SDK

helm
watch kubectl get all -n helmnginx

kubectl apply -f deploy/crds/helm.example.com_v1alpha1_helmtest_cr.yaml

curl http://34.212.226.38:32750

ansible
watch kubectl get all -n na

kubectl apply -f deploy/crds/ansible.operator.com_v1alpha1_nginxansible_cr.yaml

curl http://34.212.226.38:30088

go
watch kubectl get all -n sdk

kubectl apply -f deploy/olm-catalog/sdk-operator/1.0.1/

curl http://34.212.226.38:30020

——————————————



Upgrade Operator

KUDO
kubectl kudo upgrade my-kudo-operator --instance=mykudotest --version=1.0.3 -n mykudo


SDK

helm
operator-sdk olm-catalog gen-csv --csv-version 1.0.2 --from-version 1.0.1 --update-crds
kubectl apply -f deploy/olm-catalog/ansible-operator/1.0.2/

ansible
kubectl apply -f deploy/olm-catalog/ansible-operator/1.0.3/ansible-operator.v1.0.3.clusterserviceversion.yaml
kubectl apply -f deploy/olm-catalog/ansible-operator/1.0.4/ansible-operator.v1.0.4.clusterserviceversion.yaml

sdk
kubectl apply -f deploy/olm-catalog/sdk-operator/1.0.2/

——————————


Update Operator Instance

KUDO
kubectl kudo update --instance=mykudotest --namespace=mykudo -p defaultBody="Test this" -p replicas=3

SDK

kubectl edit nginxsdk -n sdk

———————————
Delete an Object

KUDO
kubectl delete deploy mykudotest-my-nginx -n mykudo

The operator does nothing

SDK
kubectl delete deploy run-nginxansible -n na

The operator recreates the object as expected 

----------------------------------------------------------


Community Operators
KUDO

kubectl kudo repo list

Install zookeeper cluster
kubectl kudo install zookeeper --instance=zk -n kudokafka --repo=community

Install kafka cluster
kubectl kudo install kafka --instance=kafka --namespace=kudokafka -p ZOOKEEPER_URI="zk-zookeeper-0.zk-hs:2181,zk-zookeeper-1.zk-hs:2181,zk-zookeeper-2.zk-hs:2181" repo=community

kubectl get sts -n  kudokafka
kubectl get pods -n kudokdfka

SDK
show olm ui
http://localhost:9000

Kafka Cluster using ophub


kubectl get pods -n ophub




https://github.com/operator-framework/community-operators/tree/master/community-operators


—————————————

Day 2 Operations

KUDO
kubectl run --generator="run-pod/v1" kafka-cli --rm -it --image=taion809/kafka-cli -- watch /opt/kafka/bin/kafka-topics.sh --bootstrap-server kafka-kafka-0.kafka-svc.kudokafka.svc.cluster.local:9092 --list
kubectl kudo install kafka-topic --instance my-kafka-topic -n kudokafka -p topic=mynewkafkatopic



SDK
kubectl run --generator="run-pod/v1" kafka-cli-sdk --rm -it --image=taion809/kafka-cli -- watch /opt/kafka/bin/kafka-topics.sh --bootstrap-server my-cluster-kafka-0.my-cluster-kafka-brokers.ophub.svc.cluster.local:9092 --list
Now create new instances of topic from the OLM UI
Then delete some instances from OLM UI

————————————————————————


cleanup 
kubectl run --generator="run-pod/v1" kafka-cli --rm -it --image=taion809/kafka-cli -- /opt/kafka/bin/kafka-topics.sh --bootstrap-server kafka-kafka-0.kafka-svc.kudokafka.svc.cluster.local:9092 --delete --topic mytopic
kubectl run --generator="run-pod/v1" kafka-cli-sdk --rm -it --image=taion809/kafka-cli -- /opt/kafka/bin/kafka-topics.sh --bootstrap-server my-cluster-kafka-0.my-cluster-kafka-brokers.ophub.svc.cluster.local:9092 --delete --topic my-topic


Monitor Operator


Metrics:
https://github.com/strimzi/strimzi-kafka-operator/blob/master/metrics/examples/kafka/kafka-metrics.yaml
Similar config from KUDO Kafka operator:
https://github.com/kudobuilder/operators/blob/master/repository/kafka/operator/templates/metrics-config.yaml

Import the following Dashboard template into grafana

~/KUDO_Repo/operators/repository/kafka/docs/latest/resources/grafana-dashboard.json

KUDO
less ~/operators/kudo-demo/prometheus-service-monitor.yaml 

less KUDO_Repo/operators/repository/kafka/operator/templates/metrics-config.yaml

SDK
less ~/operators/ansible-operator/prometheus-service-monitor-http.yaml
less ~/operators/ansible-operator/prometheus-service-monitor.yaml
