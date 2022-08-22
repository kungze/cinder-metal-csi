[DEFAULT]
auth_strategy = noauth
log_dir = /var/log/kolla/cinder
transport_url = rabbit://{{ .Values.rabbitmq.username }}:{{ .Values.rabbitmq.password }}@rabbitmq.{{ .Release.Namespace }}.svc.{{ .Values.cluster_domain_suffix }}:5672

{{- if and .Values.lvm.enabled .Values.ceph.enabled }}
enabled_backends = {{ printf "%s,%s" (.Values.lvm.volume_type) (.Values.ceph.volume_type) }}
default_volume_type = {{ .Values.ceph.volume_type }}
{{- else if .Values.lvm.enabled }}
enabled_backends = {{ .Values.lvm.volume_type }}
default_volume_type = {{ .Values.lvm.volume_type }}
{{- else if .Values.ceph.enabled }}
enabled_backends = {{ .Values.ceph.volume_type }}
default_volume_type = {{ .Values.ceph.volume_type }}
{{- end }}

{{- if and .Values.ceph.enabled .Values.ceph.backup.enabled }}
backup_driver = cinder.backup.drivers.ceph.CephBackupDriver
backup_ceph_conf = /etc/ceph/ceph.conf
backup_ceph_user = {{ .Values.ceph.cephClientName }}
backup_ceph_chunk_size = 134217728
backup_ceph_pool = {{ .Values.ceph.backup.poolName }}
backup_ceph_stripe_unit = 0
backup_ceph_stripe_count = 0
restore_discard_excess_bytes = true
{{- end }}

{{- if .Values.ceph.enabled }}
random_select_backend = True
{{- end }}

[keystone_authtoken]
auth_type = none

[database]
connection = sqlite:////var/lib/cinder/cinder.sqlite

{{- if .Values.ceph.enabled }}
[{{ .Values.ceph.volume_type }}]
volume_driver = cinder.volume.drivers.rbd.RBDDriver
volume_backend_name = {{ .Values.ceph.volume_type }}
rbd_pool = {{ .Values.ceph.poolName }}
rbd_ceph_conf = /etc/ceph/ceph.conf
rbd_flatten_volume_from_snapshot = false
rbd_max_clone_depth = 5
rbd_store_chunk_size = 4
rados_connect_timeout = 5
rbd_user = {{ .Values.ceph.cephClientName }}
rbd_secret_uuid = rbd_secret_uuid_palceholder
report_discard_supported = True
image_upload_use_cinder_backend = False
{{- end }}

{{- if .Values.lvm.enabled }}
[{{ .Values.lvm.volume_type }}]
volume_backend_name = {{ .Values.lvm.volume_type }}
volume_group = {{ .Values.lvm.vg_name }}
volume_driver = cinder.volume.drivers.lvm.LVMVolumeDriver
volumes_dir = /var/lib/cinder/volumes
target_helper = {{ .Values.lvm.lvm_target_helper }}
target_protocol = iscsi
lvm_type = default
{{- end }}

[oslo_messaging_notifications]
{{- if .Values.enabled_notification }}
transport_url = rabbit://{{ .Values.rabbitmq.username }}:{{ .Values.rabbitmq.password }}@rabbitmq.{{ .Release.Namespace }}.svc.{{ .Values.cluster_domain_suffix }}:5672
driver = messagingv2
topics = notifications
{{- else }}
driver = noop
{{- end }}
