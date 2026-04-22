package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseEligibilite(t *testing.T) {
	defs, err := ParseFile("../../testdata/dmn/eligibilite.dmn")
	require.NoError(t, err)

	require.Len(t, defs.Decisions, 1)
	table := defs.Decisions[0].DecisionTable
	require.NotNil(t, table)

	assert.Len(t, table.Inputs, 2)
	assert.Len(t, table.Outputs, 1)
	assert.Len(t, table.Rules, 3)
	assert.Equal(t, "FIRST", table.HitPolicy)
	assert.Equal(t, "Âge", table.Inputs[0].Label)
}

func TestParseLoanApproval(t *testing.T) {
	defs, err := ParseFile("../../testdata/dmn/loan-approval.dmn")
	require.NoError(t, err)

	assert.Equal(t, "Loan Approval", defs.Name)
	assert.Equal(t, "http://example.com/loan", defs.Namespace)
	require.Len(t, defs.Decisions, 3)

	eligibility := defs.Decisions[0]
	assert.Equal(t, "eligibility", eligibility.ID)
	assert.Equal(t, "Eligibility", eligibility.Name)
	assert.NotNil(t, eligibility.DecisionTable)
	assert.Nil(t, eligibility.LiteralExpression)
	assert.Empty(t, eligibility.InformationRequirements)

	riskScore := defs.Decisions[1]
	assert.Equal(t, "risk_score", riskScore.ID)
	require.Len(t, riskScore.InformationRequirements, 1)
	assert.Equal(t, "#eligibility", riskScore.InformationRequirements[0].RequiredDecision.Href)
	assert.NotNil(t, riskScore.DecisionTable)
	assert.Nil(t, riskScore.LiteralExpression)

	finalDecision := defs.Decisions[2]
	assert.Equal(t, "final_decision", finalDecision.ID)
	require.Len(t, finalDecision.InformationRequirements, 1)
	assert.Equal(t, "#risk_score", finalDecision.InformationRequirements[0].RequiredDecision.Href)
	assert.Nil(t, finalDecision.DecisionTable)
	require.NotNil(t, finalDecision.LiteralExpression)
	assert.Contains(t, finalDecision.LiteralExpression.Text, "riskLevel")
}
