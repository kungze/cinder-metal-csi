#!/usr/bin/bash

set -x
set -e

sudo apt install qemu-utils tgt open-iscsi

sudo mkdir -p /data
sudo qemu-img create -f raw /data/rook-ceph-disk01.img 200G
sudo qemu-img create -f raw /data/rook-ceph-disk02.img 200G
sudo qemu-img create -f raw /data/rook-ceph-disk03.img 200G

sudo tgtadm --lld iscsi --op new --mode target --tid 1 -T iqn.cinder-metal-csi.kungze.net:localnode
sudo tgtadm --lld iscsi --op new --mode logicalunit --tid 1 --lun 1 -b /data/rook-ceph-disk01.img
sudo tgtadm --lld iscsi --op new --mode logicalunit --tid 1 --lun 2 -b /data/rook-ceph-disk02.img
sudo tgtadm --lld iscsi --op new --mode logicalunit --tid 1 --lun 3 -b /data/rook-ceph-disk03.img
sudo tgtadm --lld iscsi --mode target --op bind --tid 1 --initiator-address ALL

sudo iscsiadm -m discovery -t st -p 127.0.0.1
sudo iscsiadm -m node -T iqn.cinder-metal-csi.kungze.net:localnode --login
