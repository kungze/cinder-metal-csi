#!/bin/bash
set -ex
touch /var/lib/cinder/cinder.sqlite
kolla_start
