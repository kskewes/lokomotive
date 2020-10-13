# Required variables.
variable "ssh_keys" {
  type        = list(string)
  description = "List of SSH public keys for user `core`. Each element must be specified in a valid OpenSSH public key format, as defined in RFC 4253 Section 6.6, e.g. 'ssh-rsa AAAAB3N...'."
  default     = []
}

# Optional variables.
variable "count_index" {
  type        = number
  description = "Number of nodes to create."
  default     = 1
}

variable "cluster_dns_service_ip" {
  type        = string
  description = "IP address of cluster DNS Service. Passed to kubelet as --cluster_dns parameter."
  default     = "10.3.0.10"
}

variable "clc_snippets" {
  type        = list(string)
  description = "Extra CLC snippets to include in the configuration."
  default     = []
}

variable "kubelet_image_name" {
  type        = string
  description = "Source of kubelet Docker image"
  default     = "quay.io/poseidon/kubelet"
}

variable "kubelet_image_tag" {
  type        = string
  description = "Tag for kubelet Docker image"
  default     = "v1.19.3"
}

variable "kubelet_taints" {
  type    = map(string)
  default = {}
}

variable "kubelet_labels" {
  type    = map(string)
  default = {}
}

variable "cluster_domain_suffix" {
  type        = string
  description = "Cluster domain suffix. Passed to kubelet as --cluster_domain flag."
  default     = "cluster.local"
}

variable "host_dns_ip" {
  type        = string
  description = "IP address of DNS server to configure on the nodes."
  default     = "8.8.8.8"
}

variable "ca_cert" {
  description = "Kubernetes CA certificate needed in the kubeconfig file."
  type        = string
}

variable "apiserver" {
  description = "Apiserver private endpoint needed in the kubeconfig file."
  type        = string
}
