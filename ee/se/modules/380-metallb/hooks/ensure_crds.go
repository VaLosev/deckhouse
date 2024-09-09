/*
Copyright 2024 Flant JSC
Licensed under the Deckhouse Platform Enterprise Edition (EE) license.
See https://github.com/deckhouse/deckhouse/blob/main/ee/LICENSE
*/

package hooks

import (
	"github.com/deckhouse/deckhouse/go_lib/hooks/ensure_crds"
)

var _ = ensure_crds.RegisterEnsureCRDsHook("/deckhouse/modules/380-metallb/crds/*.yaml")
var _ = ensure_crds.RegisterEnsureCRDsHook("/deckhouse/modules/380-metallb/crds/internal/*.yaml")
var _ = ensure_crds.RegisterEnsureCRDsHook("/deckhouse/modules/380-metallb/crds/native/*.yaml")
