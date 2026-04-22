package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseEligibilite(t *testing.T) {
	table, err := ParseFile("../../testdata/dmn/eligibilite.dmn")
	require.NoError(t, err)

	assert.Len(t, table.Inputs, 2)
	assert.Len(t, table.Outputs, 1)
	assert.Len(t, table.Rules, 3)
	assert.Equal(t, "FIRST", table.HitPolicy)
	assert.Equal(t, "Âge", table.Inputs[0].Label)
}
