# Copyright 2024 Flant JSC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

variable "clusterConfiguration" {
  type = any
}

variable "providerClusterConfiguration" {
  type = any
}

variable "nodeIndex" {
  type    = string
  default = 0
}

variable "cloudConfig" {
  type    = string
  default = ""
}

variable "clusterUUID" {
  type = string
}


locals {
  prefix                = var.clusterConfiguration.cloud.prefix
  namespace             = var.providerClusterConfiguration.provider.namespace
  master_instance_class = var.providerClusterConfiguration.masterNodeGroup.instanceClass

  vm_additional_labels      = lookup(local.master_instance_class.virtualMachine, "additionalLabels", {})
  vm_additional_annotations = lookup(local.master_instance_class.virtualMachine, "additionalAnnotations", {})
  vm_priority_class_name    = lookup(local.master_instance_class.virtualMachine, "priorityClassName", null)
  vm_node_selector          = lookup(local.master_instance_class.virtualMachine, "nodeSelector", null)
  vm_tolerations            = lookup(local.master_instance_class.virtualMachine, "tolerations", null)

  vm_cpu_cores         = local.master_instance_class.virtualMachine.cpu.cores
  vm_cpu_core_fraction = lookup(local.master_instance_class.virtualMachine.cpu, "coreFraction", "100%")
  vm_memory_size       = local.master_instance_class.virtualMachine.memory.size

  ip_addresses  = lookup(local.master_instance_class.virtualMachine, "ipAddresses", [])
  vm_ip_address = length(local.ip_addresses) > 0 ? local.ip_addresses[var.nodeIndex] : ""

  root_disk_size               = lookup(local.master_instance_class.rootDisk, "size", "20Gb")
  root_disk_storage_class = lookup(local.master_instance_class.rootDisk, "storageClass", null)
  root_disk_image_name         = local.master_instance_class.rootDisk.image.name
  root_disk_image_kind         = local.master_instance_class.rootDisk.image.kind

  etcd_disk_size               = lookup(local.master_instance_class.etcdDisk, "size", "20Gb")
  etcd_disk_storage_class = lookup(local.master_instance_class.etcdDisk, "storageClass", null)

  ssh_public_key = var.providerClusterConfiguration.sshPublicKey

  region = lookup(var.providerClusterConfiguration, "region", "")

  actual_zones = lookup(var.providerClusterConfiguration, "zones", [])
  zones        = lookup(var.providerClusterConfiguration.masterNodeGroup, "zones", null) != null ? tolist(setintersection(local.actual_zones, var.providerClusterConfiguration.masterNodeGroup["zones"])) : local.actual_zones
  zone         = length(local.actual_zones) > 0 ? element(local.zones, var.nodeIndex) : ""
}