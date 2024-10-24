/*
Copyright 2023 Flant JSC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package module

const DefinitionFile = "module.yaml"

type DisableOptions struct {
	Confirmation bool   `yaml:"confirmation"`
	Message      string `yaml:"message"`
}

// DeckhouseModuleDefinition represent module configuration located in DefinitionFile
// TODO: rename to Definition because the word "module" will be in the package name, and "Deckhouse" will be in the package import path
type DeckhouseModuleDefinition struct {
	Name         string            `yaml:"name"`
	Weight       uint32            `yaml:"weight,omitempty"`
	Tags         []string          `yaml:"tags"`
	Stage        string            `yaml:"stage"`
	Description  string            `yaml:"description"`
	Requirements map[string]string `json:"requirements"`

	DisableOptions DisableOptions `yaml:"disable"`

	Path string `yaml:"-"`
}