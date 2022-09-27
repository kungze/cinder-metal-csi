#!/usr/bin/bash

set -x
set -e

set -o nounset
set -o errexit
set -o pipefail

git clone --single-branch --branch v1.10.1 https://github.com/rook/rook.git
cd rook/deploy/examples
kubectl create -f crds.yaml -f common.yaml -f operator.yaml
kubectl create -f cluster-test.yaml

get_pod_cmd=(kubectl --namespace rook-ceph get pod --no-headers)
timeout=450
start_time="${SECONDS}"
while [[ $((SECONDS - start_time)) -lt $timeout ]]; do
    pod="$("${get_pod_cmd[@]}" --selector=app=rook-ceph-osd-prepare --output custom-columns=NAME:.metadata.name,PHASE:status.phase | awk 'FNR <= 1')"
    if echo "$pod" | grep 'Running\|Succeeded\|Failed'; then break; fi
    echo 'waiting for at least one osd prepare pod to be running or finished'
    sleep 5
done

kubectl -n rook-ceph wait pods -l app=rook-ceph-mon --for condition=Ready --timeout=600s

wget https://get.helm.sh/helm-v3.10.0-rc.1-linux-amd64.tar.gz
tar -zxvf helm-v3.10.0-rc.1-linux-amd64.tar.gz
sudo mv linux-amd64/helm /usr/local/bin/helm

docker pull registry.aliyuncs.com/kolla-helm/ubuntu-source-kolla-toolbox:yoga
docker pull registry.aliyuncs.com/kolla-helm/ubuntu-source-keystone:yoga
docker pull registry.aliyuncs.com/kolla-helm/ubuntu-source-cinder-backup:yoga
docker pull registry.aliyuncs.com/kolla-helm/ubuntu-source-cinder-volume:yoga
docker pull registry.aliyuncs.com/kolla-helm/ubuntu-source-glance-api:yoga
docker pull registry.aliyuncs.com/kolla-helm/ubuntu-source-kolla-toolbox:yoga
docker pull registry.aliyuncs.com/kolla-helm/ubuntu-source-cinder-api:yoga

kubectl create namespace openstack
helm repo add kungze  https://charts.kungze.net
helm install openstack-password kungze/password --namespace openstack
helm install openstack-dependency kungze/openstack-dep --namespace openstack --set mariadb.primary.persistence.enabled=false --set rabbitmq.persistence.enabled=false --wait --timeout 2400s
helm install openstack-keystone kungze/keystone --namespace openstack --wait --timeout 2400s
helm install openstack-glance kungze/glance --set ceph.replicatedSize=1 --namespace openstack --wait --timeout 2400s
helm install openstack-cinder kungze/cinder --set lvm.enabled=false --set ceph.failureDomain=osd --set ceph.replicatedSize=1 --set ceph.backup.enabled=false --namespace openstack --wait --timeout 2400s
