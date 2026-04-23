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
