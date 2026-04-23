# Getting started: your first skill from a DMN file

This guide walks you through converting a DMN decision model into a **Claude Code skill** — a reusable command that lets Claude evaluate your business logic interactively, step by step.

**What you will build**: a `/loan-approval` skill that Claude can run to evaluate a loan application following the exact decision chain defined in the DMN.

---

## Prerequisites

- `dmn2md` installed:
  ```bash
  go install github.com/jppop/dmn2md/cmd/dmn2md@latest
  ```
- A project already opened in Claude Code

---

## Step 1 — Generate the skill

Given `loan-approval.dmn`, run:

```bash
dmn2md --mode skill loan-approval.dmn
```

`dmn2md` reads the DMN and emits a Markdown document structured as a Claude Code skill: a YAML frontmatter block followed by the decision tables and literal expressions in topological (dependency-first) order.

Now create the skill directory and pipe the output in:

```bash
mkdir -p .claude/skills/loan-approval
dmn2md --mode skill loan-approval.dmn > .claude/skills/loan-approval/SKILL.md
```

---

## Step 2 — Understand the generated file

`.claude/skills/loan-approval/SKILL.md` contains:

```markdown
---
description: Evaluate Loan Approval given input values
---

# Loan Approval

Evaluate the complete decision chain. Ask the user for any missing required input before proceeding.

## Required inputs

| Variable | Type | Description |
| --- | --- | --- |
| `age` | `integer` | Age |
| `employmentStatus` | `string` | Employment Status |
| `creditScore` | `integer` | Credit Score |
| `loanAmount` | `integer` | Loan Amount |
| `loanDuration` | `integer` | Loan Duration (months) |

## Evaluation steps

### Step 1 — Applicant Risk Rating  _(Hit policy: FIRST — stop at first match)_

| # | `age` | `employmentStatus` | `creditScore` | → `riskRating` |
| --- | --- | --- | --- | --- |
| 1 | < 18 | * | * | "high" |
| 2 | >= 18 | "employed","self-employed" | >= 750 | "low" |
| 3 | >= 18 | "employed","self-employed" | [600..750) | "medium" |
| 4 | >= 18 | * | * | "high" |

### Step 2 — Loan Eligibility  _(Hit policy: UNIQUE — exactly one rule must match)_

**Depends on**: `Applicant Risk Rating`

| # | `riskRating` | `loanAmount` | `loanDuration` | → `isEligible` |
| --- | --- | --- | --- | --- |
| 1 | "high" | * | * | false |
| 2 | "low" | <= 100000 | <= 240 | true |
| 3 | "medium" | <= 50000 | <= 120 | true |
| 4 | * | * | * | false |

### Step 3 — Final Decision  _(Literal Expression)_

**Depends on**: `Loan Eligibility`, `Applicant Risk Rating`

**Output**: `finalDecision` (`string`)

```feel
if isEligible and riskRating = "low" then "approved"
else if isEligible and riskRating = "medium" then "approved with conditions"
else "rejected"
```

## Response format

For each step, show:

- For decision tables: the matched rule number and the computed output value
- For literal expressions: the evaluated result

Return the final result: `finalDecision`.
```

### Key parts explained

| Section | Role |
| --- | --- |
| `description:` frontmatter | The text Claude uses to decide when to load the skill automatically |
| **Required inputs** | Variables Claude must collect before evaluating — derived from DMN inputs that are not produced by any other decision |
| **Evaluation steps** | Each decision in dependency order; `*` means "any value" (DMN wildcard `-`) |
| Hit policy annotation | Tells Claude how many rules to apply: `FIRST` stops at the first match; `UNIQUE` expects exactly one |
| **Response format** | How Claude should present intermediate and final results |

---

## Step 3 — Refine the frontmatter (optional)

The generated `description` is functional but generic. You can enrich it so Claude loads the skill more precisely:

```yaml
---
description: >
  Evaluate a loan application following the Loan Approval DMN model.
  Use when the user provides applicant data (age, employment, credit score,
  loan amount, duration) and wants to know whether the loan should be approved.
when_to_use: >
  Triggered by phrases like "evaluate this loan", "check loan eligibility",
  or "run the loan approval decision".
argument-hint: "[age] [employmentStatus] [creditScore] [loanAmount] [loanDuration]"
---
```

The `description` and `when_to_use` fields are always loaded into context (truncated at 1 536 characters combined), so Claude can decide to invoke the skill without being explicitly asked. The full `SKILL.md` body is only loaded when the skill is actually invoked.

---

## Step 4 — Invoke the skill

In Claude Code, type:

```
/loan-approval
```

Claude prompts for the five required inputs, then evaluates each step in order:

```
> /loan-approval

I need a few values. Please provide:
- age (integer)
- employmentStatus (string — "employed", "self-employed", or other)
- creditScore (integer)
- loanAmount (integer)
- loanDuration (integer, in months)

---

User: age=25, employmentStatus="employed", creditScore=680,
      loanAmount=40000, loanDuration=84

---

**Step 1 — Applicant Risk Rating** (FIRST — stop at first match)
- Rule 1: age < 18 → no (25 ≥ 18)
- Rule 2: age ≥ 18, employed, creditScore ≥ 750 → no (680 < 750)
- Rule 3: age ≥ 18, employed, creditScore in [600, 750) → **match**
→ `riskRating` = "medium"

**Step 2 — Loan Eligibility** (UNIQUE — exactly one rule must match)
- Rule 1: riskRating = "high" → no
- Rule 2: riskRating = "low" → no
- Rule 3: riskRating = "medium", loanAmount ≤ 50 000 (✓), loanDuration ≤ 120 (✓ 84) → **match**
→ `isEligible` = true

**Step 3 — Final Decision** (Literal Expression)
- isEligible = true and riskRating = "medium"
→ `finalDecision` = **"approved with conditions"**
```

Claude can also invoke the skill automatically when it recognises the context — for example, when you paste applicant data and ask "is this loan approvable?".

---

## Step 5 — Keep the skill in sync

When the DMN model changes, regenerate with the same command:

```bash
dmn2md --mode skill loan-approval.dmn > .claude/skills/loan-approval/SKILL.md
```

Commit both the `.dmn` source and the generated `SKILL.md` together. The skill stays consistent with the model, and your team gets the updated version on pull.

---

## Project layout

```
your-project/
├── decisions/
│   └── loan-approval.dmn          ← source of truth
└── .claude/
    └── skills/
        └── loan-approval/
            └── SKILL.md           ← generated skill (commit this)
```

Skills in `.claude/skills/` are scoped to the project and visible to every Claude Code session opened in this directory. To share them with your team, commit `.claude/skills/` to your repository.

Personal skills (available across all your projects) go in `~/.claude/skills/`.

---

## Next step

Once the skill is working, you can also generate a **rule** from the same DMN. A rule embeds the decision logic directly into Claude's standing instructions so it applies automatically during a task — without an explicit invocation. See [rule-from-dmn.md](rule-from-dmn.md).
