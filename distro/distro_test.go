package distro_test

import (
	"encoding/json"
	"testing"

	"github.com/jxsl13/osfacts/distro"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_detect(t *testing.T) {

	got, err := distro.Detect()
	assert.NoError(t, err)

	require.NotNil(t, got)

	require.NotEmpty(t, got.Arch)
	require.NotEmpty(t, got.Family)
	require.NotEmpty(t, got.Distribution)
	require.NotEmpty(t, got.Version)
	require.NoError(t, err)

	data, err := json.MarshalIndent(got, "", " ")
	require.NoError(t, err)

	t.Logf("OS_INFO: %s", string(data))
}
