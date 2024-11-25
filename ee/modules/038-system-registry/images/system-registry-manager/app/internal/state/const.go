/*
Copyright 2024 Flant JSC
Licensed under the Deckhouse Platform Enterprise Edition (EE) license. See https://github.com/deckhouse/deckhouse/blob/main/ee/LICENSE
*/

package state

const (
	RegistryModuleName = "system-registry"

	LabelTypeKey             = "type"
	LabelModuleKey           = "module"
	LabelNodeSecretTypeValue = "node-secret"
	LabelHeritageKey         = "heritage"
	LabelHeritageValue       = "deckhouse"
	LabelNodeIsMasterKey     = "node-role.kubernetes.io/master"
	LabelManagedBy           = "app.kubernetes.io/managed-by"
)
