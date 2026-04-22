package parser

import (
	"strings"
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
	assert.Equal(t, "http://example.com/loan-approval", defs.Namespace)
	require.Len(t, defs.Decisions, 3)

	// ── Decision 1: Applicant Risk Rating ──────────────────────────────
	risk := defs.Decisions[0]
	assert.Equal(t, "applicantRiskRating", risk.ID)
	assert.Equal(t, "Applicant Risk Rating", risk.Name)
	assert.Empty(t, risk.InformationRequirements)
	assert.Nil(t, risk.Variable)
	assert.Nil(t, risk.LiteralExpression)

	require.NotNil(t, risk.DecisionTable)
	dt := risk.DecisionTable
	assert.Equal(t, "FIRST", dt.HitPolicy)
	assert.Len(t, dt.Inputs, 3)
	assert.Len(t, dt.Outputs, 1)
	assert.Len(t, dt.Rules, 4)

	assert.Equal(t, "age", dt.Inputs[0].InputExpression.Text)
	assert.Equal(t, "integer", dt.Inputs[0].InputExpression.TypeRef)
	assert.Equal(t, "employmentStatus", dt.Inputs[1].InputExpression.Text)
	assert.Equal(t, "creditScore", dt.Inputs[2].InputExpression.Text)

	assert.Equal(t, "Risk Rating", dt.Outputs[0].Label)
	assert.Equal(t, "riskRating", dt.Outputs[0].Name)
	assert.Equal(t, "string", dt.Outputs[0].TypeRef)

	assert.Equal(t, "Minor applicant — always high risk", dt.Rules[0].Description)
	assert.Equal(t, "< 18", dt.Rules[0].InputEntries[0].Text)
	assert.Equal(t, "-", dt.Rules[0].InputEntries[1].Text)
	assert.Equal(t, `"high"`, dt.Rules[0].OutputEntries[0].Text)

	// ── Decision 2: Loan Eligibility ────────────────────────────────────
	elig := defs.Decisions[1]
	assert.Equal(t, "loanEligibility", elig.ID)
	assert.Equal(t, "Loan Eligibility", elig.Name)
	require.Len(t, elig.InformationRequirements, 1)
	assert.Equal(t, "#applicantRiskRating", elig.InformationRequirements[0].RequiredDecision.Href)
	assert.Nil(t, elig.Variable)
	assert.Nil(t, elig.LiteralExpression)

	require.NotNil(t, elig.DecisionTable)
	dt2 := elig.DecisionTable
	assert.Equal(t, "UNIQUE", dt2.HitPolicy)
	assert.Len(t, dt2.Inputs, 3)
	assert.Equal(t, "isEligible", dt2.Outputs[0].Name)
	assert.Equal(t, "boolean", dt2.Outputs[0].TypeRef)
	assert.Len(t, dt2.Rules, 4)

	// ── Decision 3: Final Decision ──────────────────────────────────────
	final := defs.Decisions[2]
	assert.Equal(t, "finalDecision", final.ID)
	assert.Equal(t, "Final Decision", final.Name)
	require.Len(t, final.InformationRequirements, 2)
	assert.Equal(t, "#loanEligibility", final.InformationRequirements[0].RequiredDecision.Href)
	assert.Equal(t, "#applicantRiskRating", final.InformationRequirements[1].RequiredDecision.Href)
	assert.Nil(t, final.DecisionTable)

	require.NotNil(t, final.Variable)
	assert.Equal(t, "finalDecision", final.Variable.Name)
	assert.Equal(t, "string", final.Variable.TypeRef)

	require.NotNil(t, final.LiteralExpression)
	text := strings.TrimSpace(final.LiteralExpression.Text)
	assert.Contains(t, text, "isEligible")
	assert.Contains(t, text, "riskRating")
	assert.Contains(t, text, "approved with conditions")
	// multiline: text must contain actual newlines
	assert.Contains(t, text, "\n")
}
