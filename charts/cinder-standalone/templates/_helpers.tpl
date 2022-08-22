{{/* vim: set filetype=mustache: */}}

{{/*
Return the proper image name
{{ include "cinder.images.image" ( dict "registry" .Values.registry "namespace" .Values.namespace "repository" "kolla-toolbox" "tag" .Values.openstackTag) }}
*/}}
{{- define "cinder.images.image" -}}
{{- $registry := index . "registry" -}}
{{- $namespace := index . "namespace" -}}
{{- $repository := index . "repository" -}}
{{- $tag := index . "tag" -}}
{{ printf "%s/%s/%s:%s" $registry $namespace $repository $tag }}
{{- end -}}

{{/*
Return entrypoint image name
*/}}
{{- define "rabbitmq.image" -}}
{{ $repository := "ubuntu-source-rabbitmq" }}
{{- include "cinder.images.image" (dict "registry" .Values.imageRegistry "namespace" .Values.imageNamespace "repository" $repository "tag" .Values.openstackTag) }}
{{- end -}}

{{/*
Return entrypoint image name
*/}}
{{- define "kubernetes.entrypoint.image" -}}
{{ $repository := "kubernetes-entrypoint" }}
{{- include "cinder.images.image" (dict "registry" .Values.imageRegistry "namespace" .Values.imageNamespace "repository" $repository "tag" "v1.0.0") }}
{{- end -}}

{{/*
Return the proper cinder api image name
*/}}
{{- define "cinder.api.image" -}}
{{ $repository := "ubuntu-source-cinder-api" }}
{{- include "cinder.images.image" (dict "registry" .Values.imageRegistry "namespace" .Values.imageNamespace "repository" $repository "tag" .Values.openstackTag) }}
{{- end -}}

{{/*
Return the proper cinder volume image name
*/}}
{{- define "cinder.volume.image" -}}
{{ $repository := "ubuntu-source-cinder-volume" }}
{{- include "cinder.images.image" (dict "registry" .Values.imageRegistry "namespace" .Values.imageNamespace "repository" $repository "tag" .Values.openstackTag) }}
{{- end -}}

{{/*
Return the proper cinder backup image name
*/}}
{{- define "cinder.backup.image" -}}
{{ $repository := "ubuntu-source-cinder-backup" }}
{{- include "cinder.images.image" (dict "registry" .Values.imageRegistry "namespace" .Values.imageNamespace "repository" $repository "tag" .Values.openstackTag) }}
{{- end -}}

{{/*
Return the proper cinder scheduler image name
*/}}
{{- define "cinder.scheduler.image" -}}
{{ $repository := "ubuntu-source-cinder-scheduler" }}
{{- include "cinder.images.image" (dict "registry" .Values.imageRegistry "namespace" .Values.imageNamespace "repository" $repository "tag" .Values.openstackTag) }}
{{- end -}}


{{/*
Return the proper tgtd image name
*/}}
{{- define "kolla.tgtd.image" -}}
{{ $repository := "ubuntu-source-tgtd" }}
{{- include "cinder.images.image" (dict "registry" .Values.imageRegistry "namespace" .Values.imageNamespace "repository" $repository "tag" .Values.openstackTag) }}
{{- end -}}

{{/*
Return the proper cinder loop image name
*/}}
{{- define "cinder.loop.image" -}}
{{- include "cinder.images.image" (dict "registry" .Values.imageRegistry "namespace" .Values.imageNamespace "repository" "loop" "tag" "latest") }}
{{- end -}}

{{/*
Return the sync endpoint image name
*/}}
{{- define "cinder.kolla-toolbox.image" -}}
{{- $repository := "ubuntu-source-kolla-toolbox" -}}
{{- printf "%s/%s/%s:%s" .Values.imageRegistry .Values.imageNamespace $repository .Values.openstackTag }}
{{- end -}}

{{/*
Expand the name of the chart.
*/}}
{{- define "cinder.names.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "cinder.names.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Kubernetes standard labels
*/}}
{{- define "cinder.labels.standard" -}}
app.kubernetes.io/name: {{ include "cinder.names.name" . }}
helm.sh/chart: {{ include "cinder.names.chart" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}

{{/*
Labels to use on deploy.spec.selector.matchLabels and svc.spec.selector
*/}}
{{- define "cinder.labels.matchLabels" -}}
app.kubernetes.io/name: {{ include "cinder.names.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{- define "cinder.utils.template" -}}
{{- $name := index . 0 -}}
{{- $context := index . 1 -}}
{{- $last := base $context.Template.Name }}
{{- $wtf := $context.Template.Name | replace $last $name -}}
{{ include $wtf $context }}
{{- end -}}
