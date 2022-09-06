#!/bin/bash
set -ex

export OS_AUTH_TYPE=noauth
export OS_PROJECT_ID=admin
export OS_VOLUME_API_VERSION=3.10
export CINDER_ENDPOINT=http://cinder-api.{{ .Release.Namespace }}.svc.{{ .Values.cluster_domain_suffix }}:8776/v3

{{- if .Values.ceph.enabled }}
cinder type-create {{ .Values.ceph.volume_type }}
{{- end }}

{{- if .Values.lvm.enabled }}
cinder type-create {{ .Values.lvm.volume_type }}
{{- end }}
