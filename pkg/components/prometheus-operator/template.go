// Copyright 2020 The Lokomotive Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package prometheus

const chartValuesTmpl = `
global:
  rbac:
    pspAnnotations:
      seccomp.security.alpha.kubernetes.io/allowedProfileNames: 'docker/default,runtime/default'
      seccomp.security.alpha.kubernetes.io/defaultProfileName:  'docker/default'

alertmanager:
{{.AlertManagerConfig}}
  alertmanagerSpec:
    retention: {{.AlertManagerRetention}}
    externalUrl: {{.AlertManagerExternalURL}}
    {{ if .AlertManagerNodeSelector }}
    nodeSelector:
      {{ range $key, $value := .AlertManagerNodeSelector }}
      {{ $key }}: {{ $value }}
      {{ end }}
    {{ end }}
    storage:
      volumeClaimTemplate:
        # This is done to reduce the name length of PVC that is autogenerated if metadata.Name is
        # not provided. More info: https://github.com/coreos/prometheus-operator/issues/2865.
        metadata:
          name: data
        spec:
          {{ if .StorageClass }}
          storageClassName: {{.StorageClass}}
          {{ end }}
          accessModes: ["ReadWriteOnce"]
          resources:
            requests:
              storage: "{{.AlertManagerStorageSize}}"

grafana:
  plugins: "grafana-piechart-panel"
  testFramework:
    enabled: false
  sidecar:
    dashboards:
      searchNamespace: ALL
      provider:
        foldersFromFilesStructure: true
  rbac:
    pspUseAppArmor: false
  adminPassword: {{.Grafana.AdminPassword}}
  {{- if .Grafana.SecretEnv }}
  envRenderSecret:
    {{ range $key, $value := .Grafana.SecretEnv }}
    {{ $key }}: {{ $value }}
    {{ end }}
  {{- end }}
  {{ if .Grafana.Ingress }}
  ingress:
    enabled: true
    annotations:
      kubernetes.io/ingress.class: {{.Grafana.Ingress.Class}}
      kubernetes.io/tls-acme: "true"
      cert-manager.io/cluster-issuer: {{.Grafana.Ingress.CertManagerClusterIssuer}}
    hosts:
    - {{ .Grafana.Ingress.Host }}
    tls:
    - hosts:
      - {{ .Grafana.Ingress.Host }}
      secretName: {{ .Grafana.Ingress.Host }}-tls
  grafana.ini:
    server:
      root_url: https://{{ .Grafana.Ingress.Host }}
  {{ end }}

kubeEtcd:
  enabled: {{.Monitor.Etcd}}
prometheus-node-exporter:
  service: {}
{{ if (or .PrometheusOperatorNodeSelector .DisableWebhooks) }}
prometheusOperator:
  {{- if .DisableWebhooks }}
  tlsProxy:
    enabled: false
  admissionWebhooks:
    enabled: false
  {{- end }}
  {{- if .PrometheusOperatorNodeSelector }}
  nodeSelector:
    {{ range $key, $value := .PrometheusOperatorNodeSelector }}
    {{ $key }}: {{ $value }}
    {{ end }}
  {{- end }}
{{ end }}
prometheus:
  {{ if .Prometheus.Ingress }}
  ingress:
    enabled: true
    annotations:
      kubernetes.io/ingress.class: {{.Prometheus.Ingress.Class}}
      kubernetes.io/tls-acme: "true"
      cert-manager.io/cluster-issuer: {{.Prometheus.Ingress.CertManagerClusterIssuer}}
    hosts:
    - {{ .Prometheus.Ingress.Host }}
    tls:
    - hosts:
      - {{ .Prometheus.Ingress.Host }}
      secretName: {{ .Prometheus.Ingress.Host }}-tls
  {{ end }}
  prometheusSpec:
    {{ if .Prometheus.ExternalURL }}
    externalUrl: {{ .Prometheus.ExternalURL }}
    {{ else if .Prometheus.Ingress }}
    externalUrl: https://{{.Prometheus.Ingress.Host}}
    {{ end }}
    {{ if .Prometheus.NodeSelector }}
    nodeSelector:
      {{ range $key, $value := .Prometheus.NodeSelector }}
      {{ $key }}: {{ $value }}
      {{ end }}
    {{ end }}
    {{ if .Prometheus.ExternalLabels }}
    externalLabels:
      {{ range $key, $value := .Prometheus.ExternalLabels}}
      {{ $key }}: {{ $value }}
      {{ end }}
    {{ end }}
    retention: {{.Prometheus.MetricsRetention}}
    serviceMonitorSelectorNilUsesHelmValues: {{.Prometheus.WatchLabeledServiceMonitors}}
    ruleSelectorNilUsesHelmValues: {{.Prometheus.WatchLabeledPrometheusRules}}
    storageSpec:
      volumeClaimTemplate:
        metadata:
          name: data
        spec:
          {{ if .StorageClass }}
          storageClassName: {{.StorageClass}}
          {{ end }}
          accessModes: ["ReadWriteOnce"]
          resources:
            requests:
              storage: "{{.Prometheus.StorageSize}}"

kubeControllerManager:
  enabled: {{.Monitor.KubeControllerManager}}
  service:
    selector:
      k8s-app: kube-controller-manager
      tier: control-plane

coreDns:
  service:
    selector:
      {{ range $k, $v := .CoreDNS.Selector }}
      {{ $k }}: "{{ $v }}"
      {{- end }}

kubeScheduler:
  enabled: {{.Monitor.KubeScheduler}}
  service:
    selector:
      k8s-app: kube-scheduler
      tier: control-plane

kube-state-metrics:
  podSecurityPolicy:
    annotations:
      seccomp.security.alpha.kubernetes.io/allowedProfileNames: 'docker/default,runtime/default'
      seccomp.security.alpha.kubernetes.io/defaultProfileName:  'docker/default'

kubeProxy:
  enabled: {{.Monitor.KubeProxy}}

kubelet:
  enabled: {{.Monitor.Kubelet}}
`
