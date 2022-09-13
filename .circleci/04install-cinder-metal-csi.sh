#!/usr/bin/bash

set -x
set -e

sudo systemctl stop systemd-resolved.service
sudo systemctl disable systemd-resolved.service
sudo rm /etc/resolv.conf
sudo touch /etc/resolv.conf
sudo chmod a+w /etc/resolv.conf
sudo echo "nameserver 169.254.20.10" > /etc/resolv.conf
sudo echo "nameserver 8.8.8.8" >> /etc/resolv.conf

sudo apt install python3-openstackclient -y
export OS_USERNAME=$(kubectl get secret -n openstack openstack-keystone -o jsonpath="{.data.OS_USERNAME}" | base64 --decode)
export OS_PROJECT_DOMAIN_NAME=$(kubectl get secret -n openstack openstack-keystone -o jsonpath="{.data.OS_PROJECT_DOMAIN_NAME}" | base64 --decode)
export OS_USER_DOMAIN_NAME=$(kubectl get secret -n openstack openstack-keystone -o jsonpath="{.data.OS_USER_DOMAIN_NAME}" | base64 --decode)
export OS_PROJECT_NAME=$(kubectl get secret -n openstack openstack-keystone -o jsonpath="{.data.OS_PROJECT_NAME}" | base64 --decode)
export OS_REGION_NAME=$(kubectl get secret -n openstack openstack-keystone -o jsonpath="{.data.OS_REGION_NAME}" | base64 --decode)
export OS_PASSWORD=$(kubectl get secrets -n openstack openstack-password -o jsonpath="{.data.keystone-admin-password}" | base64 --decode)
export OS_AUTH_URL=$(kubectl get secret -n openstack openstack-keystone -o jsonpath="{.data.OS_CLUSTER_URL}" | base64 --decode)
export OS_INTERFACE=internal
openstack project create kubernetes
openstack user create --project kubernetes --password ChangeMe kubernetes
openstack role add --project kubernetes --user kubernetes member

cd  $HOME/project
export CEPH_CINDER_KEYRING=$(kubectl -n rook-ceph get secrets rook-ceph-client-cinder -o jsonpath="{.data.cinder}" | base64 --decode)
sed -i 's?__CEPH_CINDER_KEYRING__?'$CEPH_CINDER_KEYRING'?g' manifests/rbd-secrets.yaml

kubectl apply -f manifests/cinder-metal-csi-config.yaml
kubectl apply -f manifests/cinder-metal-csi-controllerplugin-rbac.yaml
kubectl apply -f manifests/cinder-metal-csi-controllerplugin.yaml
kubectl apply -f manifests/cinder-metal-csi-nodeplugin-rbac.yaml
kubectl apply -f manifests/cinder-metal-csi-nodeplugin.yaml
kubectl apply -f manifests/rbd-secrets.yaml
kubectl apply -f manifests/storageclass-rbd.yaml
kubectl -n kube-system wait pods -l app=cinder-metal-csi-controller-plugin --for condition=Ready --timeout=600s
kubectl -n kube-system wait pods -l app=cinder-metal-csi-node-plugin --for condition=Ready --timeout=600s
