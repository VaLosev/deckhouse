package reginjector

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/deckhouse/deckhouse/deckhouse-controller/pkg/apis/deckhouse.io/v1alpha1"
)

func TestOpenapiInjection(t *testing.T) {
	source := `
x-extend:
  schema: config-values.yaml
type: object
properties:
  internal:
    type: object
    default: {}
    properties:
      pythonVersions:
        type: array
        default: []
        items:
          type: string
  registry:
    type: object
    description: "System field, overwritten by Deckhouse. Don't use"
`

	sourceModule := &v1alpha1.ModuleSource{}
	sourceModule.Spec.Registry.Repo = "test.deckhouse.io/foo/bar"
	sourceModule.Spec.Registry.DockerCFG = "dGVzdG1lCg=="
	sourceModule.Spec.Registry.Scheme = "http"
	sourceModule.Spec.Registry.CA = "someCA"

	data, err := mutateOpenapiSchema([]byte(source), sourceModule)
	require.NoError(t, err)

	assert.YAMLEq(t, `
type: object
x-extend:
  schema: config-values.yaml
properties:
  registry:
    type: object
    default: {}
    properties:
      base:
        type: string
        default: test.deckhouse.io/foo/bar
      dockercfg:
        type: string
        default: dGVzdG1lCg==
      scheme:
        type: string
        default: http
      ca:
        type: string
        default: someCA
  internal:
    default: {}
    properties:
      pythonVersions:
        default: []
        items:
          type: string
        type: array
    type: object
`, string(data))
}
