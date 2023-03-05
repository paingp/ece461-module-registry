#!/bin/sh
set -xe
PROJECT_ID="principal-yen-379617"
COMPUTE_ENGINE_ZONE="us-east5-a"
CLUSTER_NAME="sec-privacy-cluster"
_SERVICE_NAME="Module Registry SA"

docker build -t gcr.io/${PROJECT_ID}/dev:v1 .
docker images
gcloud auth configure-docker
docker push gcr.io/${PROJECT_ID}/dev:v1
docker run --rm -p 8080:8080 gcr.io/${PROJECT_ID}/dev:v1
gcloud config set project principal-yen-379617
gcloud config set compute/zone us-east5-a
gcloud container clusters create sec-privacy-cluster --num-nodes=2
gcloud compute instances list kubectl create deployment sec-privacy-cluster --image=gcr.io/${PROJECT_ID}/dev:v1
kubectl get pods kubectl expose deployment sec-privacy-cluster --type=LoadBalancer --port 80 --target-port 8080
