package renderer_test

import (
	"flag"
	"os"
	"testing"

	"github.com/jppop/dmn2md/internal/parser"
	"github.com/jppop/dmn2md/internal/renderer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var update = flag.Bool("update", false, "update golden files")

func TestRenderLoanApprovalDoc(t *testing.T) {
	defs, err := parser.ParseFile("../../testdata/dmn/loan-approval.dmn")
	require.NoError(t, err)

	actual := renderer.Render(defs, renderer.ModeDoc)

	goldenPath := "../../testdata/expected/loan-approval.md"

	if *update {
		err := os.WriteFile(goldenPath, []byte(actual), 0644)
		require.NoError(t, err)
	}

	expected, err := os.ReadFile(goldenPath)
	require.NoError(t, err)

	assert.Equal(t, string(expected), actual)
}

func TestRenderLoanApprovalSkill(t *testing.T) {
	defs, err := parser.ParseFile("../../testdata/dmn/loan-approval.dmn")
	require.NoError(t, err)

	actual := renderer.Render(defs, renderer.ModeSkill)

	goldenPath := "../../testdata/expected/loan-approval-skill.md"

	if *update {
		err := os.WriteFile(goldenPath, []byte(actual), 0644)
		require.NoError(t, err)
	}

	expected, err := os.ReadFile(goldenPath)
	require.NoError(t, err)

	assert.Equal(t, string(expected), actual)
}

func TestRenderLoanApprovalRule(t *testing.T) {
	defs, err := parser.ParseFile("../../testdata/dmn/loan-approval.dmn")
	require.NoError(t, err)

	actual := renderer.Render(defs, renderer.ModeRule)

	goldenPath := "../../testdata/expected/loan-approval-rule.md"

	if *update {
		err := os.WriteFile(goldenPath, []byte(actual), 0644)
		require.NoError(t, err)
	}

	expected, err := os.ReadFile(goldenPath)
	require.NoError(t, err)

	assert.Equal(t, string(expected), actual)
}
