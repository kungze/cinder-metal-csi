#!/bin/bash -e

CEPH_CONFIG="/etc/ceph/ceph.conf"
MON_CONFIG="/etc/rook/mon-endpoints"
KEYRING_FILE="/etc/ceph/ceph.client.cinder.keyring"

# create a ceph config file in its default location so ceph/rados tools can be used
# without specifying any arguments
write_endpoints() {
  endpoints=$(cat ${MON_CONFIG})

  # filter out the mon names
  # external cluster can have numbers or hyphens in mon names, handling them in regex
  # shellcheck disable=SC2001
  mon_endpoints=$(echo "${endpoints}"| sed 's/[a-z0-9_-]\+=//g')

  DATE=$(date)
  echo "$DATE writing mon endpoints to ${CEPH_CONFIG}: ${endpoints}"
    cat <<EOF > ${CEPH_CONFIG}
[global]
mon_host = ${mon_endpoints}

[client.cinder]
keyring = ${KEYRING_FILE}
EOF
}

# create the keyring file
cat <<EOF > ${KEYRING_FILE}
[client.${CEPH_USERNAME}]
key = ${CEPH_SECRET}
EOF

# write the initial config file
write_endpoints
