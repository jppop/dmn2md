package renderer

import (
	"fmt"
	"strings"

	"github.com/jppop/dmn2md/internal/model"
)

func renderSkill(defs *model.Definitions) string {
	var sb strings.Builder
	writeSkillFrontmatter(&sb, defs)
	writeSkillHeader(&sb, defs)
	writeSkillInputs(&sb, defs)
	writeSkillSteps(&sb, defs)
	writeSkillResponseFormat(&sb, defs)
	return sb.String()
}

func writeSkillFrontmatter(sb *strings.Builder, defs *model.Definitions) {
	fmt.Fprintf(sb, "---\n")
	fmt.Fprintf(sb, "description: Evaluate %s given input values\n", defs.Name)
	fmt.Fprintf(sb, "---\n\n")
}

func writeSkillHeader(sb *strings.Builder, defs *model.Definitions) {
	fmt.Fprintf(sb, "# %s\n\n", defs.Name)
	fmt.Fprintf(sb, "Evaluate the complete decision chain. Ask the user for any missing required input before proceeding.\n\n")
}

// writeSkillInputs emits the table of external inputs (not produced by any decision).
func writeSkillInputs(sb *strings.Builder, defs *model.Definitions) {
	outputNames := collectOutputNames(defs)

	type inputEntry struct {
		name    string
		typeRef string
		label   string
	}

	seen := make(map[string]bool)
	var inputs []inputEntry

	for _, d := range defs.Decisions {
		if d.DecisionTable == nil {
			continue
		}
		for _, inp := range d.DecisionTable.Inputs {
			name := inp.InputExpression.Text
			if !seen[name] && !outputNames[name] {
				seen[name] = true
				inputs = append(inputs, inputEntry{name, inp.InputExpression.TypeRef, inp.Label})
			}
		}
	}

	if len(inputs) == 0 {
		return
	}

	fmt.Fprintf(sb, "## Required inputs\n\n")
	fmt.Fprintf(sb, "| Variable | Type | Description |\n")
	fmt.Fprintf(sb, "| --- | --- | --- |\n")
	for _, inp := range inputs {
		fmt.Fprintf(sb, "| `%s` | `%s` | %s |\n", inp.name, inp.typeRef, inp.label)
	}
	fmt.Fprintf(sb, "\n")
}

func writeSkillSteps(sb *strings.Builder, defs *model.Definitions) {
	sorted := topoSort(defs)
	fmt.Fprintf(sb, "## Evaluation steps\n\n")
	for i, d := range sorted {
		writeSkillStep(sb, i+1, d, defs)
	}
}

func writeSkillStep(sb *strings.Builder, stepNum int, d model.Decision, defs *model.Definitions) {
	if d.DecisionTable != nil {
		writeSkillDecisionTableStep(sb, stepNum, d, defs)
	} else if d.LiteralExpression != nil {
		writeSkillLiteralExpressionStep(sb, stepNum, d, defs)
	}
}

func writeSkillDecisionTableStep(sb *strings.Builder, stepNum int, d model.Decision, defs *model.Definitions) {
	dt := d.DecisionTable
	hp := dt.HitPolicy
	if hp == "" {
		hp = "UNIQUE"
	}
	fmt.Fprintf(sb, "### Step %d — %s  _(%s)_\n\n", stepNum, d.Name, hitPolicyNote(hp))

	if parts := skillDepNames(d, defs); len(parts) > 0 {
		fmt.Fprintf(sb, "**Depends on**: %s\n\n", strings.Join(parts, ", "))
	}

	sb.WriteString("| # |")
	for _, inp := range dt.Inputs {
		fmt.Fprintf(sb, " `%s` |", inp.InputExpression.Text)
	}
	for _, out := range dt.Outputs {
		fmt.Fprintf(sb, " → `%s` |", outputVarName(out))
	}
	sb.WriteString("\n")

	sb.WriteString("| --- |")
	for range dt.Inputs {
		sb.WriteString(" --- |")
	}
	for range dt.Outputs {
		sb.WriteString(" --- |")
	}
	sb.WriteString("\n")

	for i, rule := range dt.Rules {
		fmt.Fprintf(sb, "| %d |", i+1)
		for _, entry := range rule.InputEntries {
			text := entry.Text
			if text == "-" {
				text = "*"
			}
			fmt.Fprintf(sb, " %s |", text)
		}
		for _, entry := range rule.OutputEntries {
			fmt.Fprintf(sb, " %s |", entry.Text)
		}
		sb.WriteString("\n")
	}
	fmt.Fprintf(sb, "\n")
}

func writeSkillLiteralExpressionStep(sb *strings.Builder, stepNum int, d model.Decision, defs *model.Definitions) {
	fmt.Fprintf(sb, "### Step %d — %s  _(Literal Expression)_\n\n", stepNum, d.Name)

	if parts := skillDepNames(d, defs); len(parts) > 0 {
		fmt.Fprintf(sb, "**Depends on**: %s\n\n", strings.Join(parts, ", "))
	}

	if d.Variable != nil {
		fmt.Fprintf(sb, "**Output**: `%s` (`%s`)\n\n", d.Variable.Name, d.Variable.TypeRef)
	}

	fmt.Fprintf(sb, "```feel\n")
	fmt.Fprintf(sb, "%s\n", strings.TrimSpace(d.LiteralExpression.Text))
	fmt.Fprintf(sb, "```\n\n")
}

func writeSkillResponseFormat(sb *strings.Builder, defs *model.Definitions) {
	fmt.Fprintf(sb, "## Response format\n\n")
	fmt.Fprintf(sb, "For each step, show:\n\n")
	fmt.Fprintf(sb, "- For decision tables: the matched rule number and the computed output value\n")
	fmt.Fprintf(sb, "- For literal expressions: the evaluated result\n\n")

	sorted := topoSort(defs)
	if len(sorted) == 0 {
		return
	}
	last := sorted[len(sorted)-1]
	switch {
	case last.Variable != nil:
		fmt.Fprintf(sb, "Return the final result: `%s`.\n", last.Variable.Name)
	case last.DecisionTable != nil && len(last.DecisionTable.Outputs) > 0:
		parts := make([]string, 0, len(last.DecisionTable.Outputs))
		for _, o := range last.DecisionTable.Outputs {
			parts = append(parts, fmt.Sprintf("`%s`", outputVarName(o)))
		}
		fmt.Fprintf(sb, "Return the final result: %s.\n", strings.Join(parts, ", "))
	}
}

func hitPolicyNote(hp string) string {
	switch hp {
	case "FIRST":
		return "Hit policy: FIRST — stop at first match"
	case "UNIQUE":
		return "Hit policy: UNIQUE — exactly one rule must match"
	case "ANY":
		return "Hit policy: ANY — all matching rules must produce the same output"
	case "COLLECT":
		return "Hit policy: COLLECT — collect all matching outputs"
	case "RULE ORDER":
		return "Hit policy: RULE ORDER — apply all matching rules in order"
	default:
		return fmt.Sprintf("Hit policy: %s", hp)
	}
}

func skillDepNames(d model.Decision, defs *model.Definitions) []string {
	parts := make([]string, 0, len(d.InformationRequirements))
	for _, ir := range d.InformationRequirements {
		srcID := strings.TrimPrefix(ir.RequiredDecision.Href, "#")
		dep := findDecision(defs, srcID)
		if dep != nil {
			parts = append(parts, fmt.Sprintf("`%s`", dep.Name))
		}
	}
	return parts
}

// collectOutputNames returns the set of variable names produced by any decision.
func collectOutputNames(defs *model.Definitions) map[string]bool {
	names := make(map[string]bool)
	for _, d := range defs.Decisions {
		if d.DecisionTable != nil {
			for _, o := range d.DecisionTable.Outputs {
				names[outputVarName(o)] = true
			}
		}
		if d.Variable != nil {
			names[d.Variable.Name] = true
		}
	}
	return names
}
