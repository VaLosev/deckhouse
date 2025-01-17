/*
Copyright 2022 Flant JSC

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

package configtools

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/flant/addon-operator/pkg/module_manager/models/modules"
	"github.com/flant/addon-operator/pkg/utils"
	"github.com/flant/addon-operator/pkg/values/validation"

	"github.com/deckhouse/deckhouse/deckhouse-controller/pkg/apis/deckhouse.io/v1alpha1"
	"github.com/deckhouse/deckhouse/go_lib/configtools/conversion"
)

// Validator is a validator for values in ModuleConfig.
type Validator struct {
	valuesValidator ValuesValidator
}

// ValuesValidator is a part of ValuesValidator from addon-operator with needed
// methods to validate config values.
type ValuesValidator interface {
	GetGlobal() *modules.GlobalModule
	GetModule(name string) *modules.BasicModule
}

func NewValidator(valuesValidator ValuesValidator) *Validator {
	return &Validator{
		valuesValidator: valuesValidator,
	}
}

type ValidationResult struct {
	Settings map[string]interface{}
	Version  int

	Warning string
	Error   string
}

// HasError return true if result has error.
// TODO(future) use regular error instead.
func (v ValidationResult) HasError() bool {
	return v.Error != ""
}

// validateCR checks if ModuleConfig resource is well-formed.
func (v *Validator) validateCR(config *v1alpha1.ModuleConfig) ValidationResult {
	result := ValidationResult{}

	if config.Spec.Version == 0 {
		// Resource is not valid when spec.settings are specified without version.
		if len(config.Spec.Settings) > 0 {
			result.Error = "spec.version is required when spec.settings are specified"
		}
		// Resource is valid without spec.version and spec.settings.
		return result
	}

	// Can run conversions and validations if spec.version and spec.settings are specified.
	if len(config.Spec.Settings) == 0 {
		// Warn about spec.version without spec.settings.
		result.Warning = "spec.version has no effect without spec.settings, defaults from the latest version of settings schema will be applied"
	}

	converter := conversion.Store().Get(config.GetName())
	latestVersion := converter.LatestVersion()

	// Check if version is unknown.
	if !converter.IsKnownVersion(config.Spec.Version) {
		prevVersionsMsg := concatIntList(converter.ListVersionsWithoutLatest())
		if prevVersionsMsg != "" {
			prevVersionsMsg = fmt.Sprintf(", or one of previous versions: %s", prevVersionsMsg)
		}

		msg := fmt.Sprintf("spec.version=%d is unsupported. Use latest version %d%s", config.Spec.Version, latestVersion, prevVersionsMsg)
		if hasVersionedSettings(config) {
			// Error if spec.settings are specified. Can't start conversions for such configuration.
			result.Error = msg
		} else {
			// Warning if there are no spec.settings.
			result.Warning = msg
		}
		return result
	}

	newVersion, newSettings, err := converter.ConvertToLatest(config.Spec.Version, config.Spec.Settings)
	if err != nil {
		result.Error = fmt.Sprintf("spec.settings conversion from version %d to %d: %v", config.Spec.Version, newVersion, err)
		return result
	}
	result.Settings = newSettings
	result.Version = newVersion

	if config.Spec.Version != latestVersion {
		result.Warning = fmt.Sprintf("spec.version=%d is obsolete. Please migrate spec.settings to the latest version %d", config.Spec.Version, latestVersion)
	}

	return result
}

// Validate checks ModuleConfig resource:
// - check if resource is well-formed
// - runs conversions for spec.settings if it`s needed
// - use OpenAPI schema defined in related config-values.yaml file to validate converted spec.settings.
// TODO(future) return cfg, error. Put cfg.Spec into result cfg.
func (v *Validator) Validate(config *v1alpha1.ModuleConfig) ValidationResult {
	result := v.validateCR(config)
	if result.HasError() {
		return result
	}

	if config.Spec.Enabled != nil && !(*config.Spec.Enabled) {
		return result
	}

	if module := v.valuesValidator.GetModule(config.Name); module == nil {
		return result
	}

	if err := v.validateSettings(config.GetName(), result.Settings); err != nil {
		convMsg := ""
		if config.Spec.Version != result.Version {
			convMsg = fmt.Sprintf(" converted to %d", result.Version)
		}
		result.Error = fmt.Sprintf("spec.settings are not valid (version %d%s): %v", config.Spec.Version, convMsg, cleanupMultilineError(err))
	}

	return result
}

// validateSettings uses module from ModuleManager instance to validate spec.settings.
// cfgName arg is a kebab-cased name of the ModuleConfig resource.
// cfgSettings is a content of spec.settings and can be nil if settings field wasn't set.
// (Note: cfgSettings map is a map with 'plain values', i.e. without camelCased module name as a root key).
func (v *Validator) validateSettings(configName string, configSettings map[string]interface{}) error {
	// Ignore empty validator.
	if v.valuesValidator == nil {
		return nil
	}

	// init cfg settings if it equals nil
	if configSettings == nil {
		configSettings = make(map[string]interface{})
	}

	valuesKey := valuesKeyFromObjectName(configName)
	var schemaStorage *validation.SchemaStorage
	if configName == "global" {
		schemaStorage = v.valuesValidator.GetGlobal().GetSchemaStorage()
	} else {
		module := v.valuesValidator.GetModule(configName)
		schemaStorage = module.GetSchemaStorage()
	}

	// Instantiate defaults from the OpenAPI schema.
	defaultSettings := make(map[string]interface{})
	s := schemaStorage.Schemas[validation.ConfigValuesSchema]
	if s != nil {
		validation.ApplyDefaults(defaultSettings, s)
	}

	// Merge defaults with passed settings as addon-operator will do.
	values := utils.MergeValues(
		utils.Values{valuesKey: defaultSettings},
		utils.Values{valuesKey: configSettings},
	)

	return schemaStorage.ValidateConfigValues(valuesKey, values)
}

func valuesKeyFromObjectName(name string) string {
	if name == "global" {
		return name
	}
	return utils.ModuleNameToValuesKey(name)
}

func concatIntList(items []int) string {
	var buf strings.Builder
	for i, item := range items {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(strconv.Itoa(item))
	}
	return buf.String()
}

func cleanupMultilineError(err error) string {
	if err == nil {
		return ""
	}
	parts := strings.Split(err.Error(), "\n")
	buf := strings.Builder{}
	for _, part := range parts {
		buf.WriteString(" ")
		buf.WriteString(strings.TrimSpace(part))
	}
	return buf.String()
}

func hasVersionedSettings(cfg *v1alpha1.ModuleConfig) bool {
	return cfg != nil && cfg.Spec.Version > 0 && cfg.Spec.Settings != nil
}
