package renderer

import (
	"fmt"
	"strings"

	"github.com/jppop/dmn2md/internal/model"
)

func renderRule(defs *model.Definitions) string {
	var sb strings.Builder
	writeRuleHeader(&sb, defs)
	writeRuleDecisions(&sb, defs)
	writeRuleConstraints(&sb, defs)
	return sb.String()
}

func writeRuleHeader(sb *strings.Builder, defs *model.Definitions) {
	fmt.Fprintf(sb, "# %s\n\n", defs.Name)
	fmt.Fprintf(sb, "Apply these rules whenever evaluating a %s.\n\n", defs.Name)
}

func writeRuleDecisions(sb *strings.Builder, defs *model.Definitions) {
	sorted := topoSort(defs)
	for _, d := range sorted {
		if d.DecisionTable != nil {
			writeRuleDecisionTable(sb, d, defs)
		} else if d.LiteralExpression != nil {
			writeRuleLiteralExpression(sb, d, defs)
		}
	}
}

func writeRuleDecisionTable(sb *strings.Builder, d model.Decision, defs *model.Definitions) {
	dt := d.DecisionTable

	fmt.Fprintf(sb, "## %s\n\n", d.Name)

	if parts := skillDepNames(d, defs); len(parts) > 0 {
		fmt.Fprintf(sb, "_Requires_: %s\n\n", strings.Join(parts, ", "))
	}

	outputLabels := make([]string, 0, len(dt.Outputs))
	for _, o := range dt.Outputs {
		outputLabels = append(outputLabels, fmt.Sprintf("`%s`", outputVarName(o)))
	}
	fmt.Fprintf(sb, "Determines %s. %s:\n\n", strings.Join(outputLabels, ", "), hitPolicyRuleIntro(dt.HitPolicy))

	for _, rule := range dt.Rules {
		conditions := ruleConditions(rule, dt)
		outputs := ruleOutputs(rule, dt)
		line := fmt.Sprintf("- %s → %s", conditions, outputs)
		if rule.Description != "" {
			line += fmt.Sprintf(" _%s_", rule.Description)
		}
		fmt.Fprintln(sb, line)
	}
	fmt.Fprintln(sb)
}

func writeRuleLiteralExpression(sb *strings.Builder, d model.Decision, defs *model.Definitions) {
	fmt.Fprintf(sb, "## %s\n\n", d.Name)

	if parts := skillDepNames(d, defs); len(parts) > 0 {
		fmt.Fprintf(sb, "_Requires_: %s\n\n", strings.Join(parts, ", "))
	}

	if d.Variable != nil {
		fmt.Fprintf(sb, "Computes `%s` (`%s`):\n\n", d.Variable.Name, d.Variable.TypeRef)
	}

	fmt.Fprintf(sb, "```feel\n")
	fmt.Fprintf(sb, "%s\n", strings.TrimSpace(d.LiteralExpression.Text))
	fmt.Fprintf(sb, "```\n\n")
}

func writeRuleConstraints(sb *strings.Builder, defs *model.Definitions) {
	// Collect dependency order as a constraint
	sorted := topoSort(defs)
	if len(sorted) < 2 {
		return
	}

	fmt.Fprintf(sb, "## Constraints\n\n")

	// Emit evaluation order
	names := make([]string, 0, len(sorted))
	for _, d := range sorted {
		names = append(names, fmt.Sprintf("`%s`", d.Name))
	}
	fmt.Fprintf(sb, "- Always evaluate decisions in this order: %s.\n", strings.Join(names, " → "))
	fmt.Fprintf(sb, "- Always state which rule determined each intermediate value.\n")
}

// ruleConditions formats the input entries of a rule as a readable condition string.
func ruleConditions(rule model.Rule, dt *model.DecisionTable) string {
	parts := make([]string, 0, len(rule.InputEntries))
	for i, entry := range rule.InputEntries {
		if i >= len(dt.Inputs) {
			break
		}
		text := entry.Text
		if text == "-" {
			continue // wildcard — skip for readability
		}
		varName := dt.Inputs[i].InputExpression.Text
		parts = append(parts, fmt.Sprintf("`%s` %s", varName, text))
	}
	if len(parts) == 0 {
		return "_(any)_"
	}
	return strings.Join(parts, ", ")
}

// hitPolicyRuleIntro returns the imperative intro sentence for a hit policy.
func hitPolicyRuleIntro(hp string) string {
	switch hp {
	case "FIRST":
		return "Apply the first matching rule"
	case "ANY":
		return "Any matching rule applies (all must produce the same output)"
	case "COLLECT":
		return "Collect the outputs of all matching rules"
	case "RULE ORDER":
		return "Apply all matching rules in order"
	default:
		return "Exactly one rule must match"
	}
}

// ruleOutputs formats the output entries of a rule as a readable result string.
func ruleOutputs(rule model.Rule, dt *model.DecisionTable) string {
	parts := make([]string, 0, len(rule.OutputEntries))
	for i, entry := range rule.OutputEntries {
		if i >= len(dt.Outputs) {
			break
		}
		varName := outputVarName(dt.Outputs[i])
		parts = append(parts, fmt.Sprintf("`%s` = %s", varName, entry.Text))
	}
	return strings.Join(parts, ", ")
}
