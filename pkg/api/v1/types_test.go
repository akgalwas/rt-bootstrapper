package v1_test

import (
	"strings"
	"testing"

	v1 "github.com/kyma-project/rt-bootstrapper/pkg/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	tcs := []struct {
		name     string
		val      string
		expected v1.Config
	}{
		{
			name: "happy path",
			val: `{ 
  "imagePullSecretName": "ipsn2",
  "imagePullSecretNamespace": "ipsns2",
  "secretSyncInterval": "10m",
  "overrides": { "rn2": "orn2" }
}`,
			expected: v1.Config{
				Overrides: map[string]string{
					"rn2": "orn2",
				},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			r := strings.NewReader(tc.val)
			actual, err := v1.NewConfig(r)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, *actual)
		})
	}
}
