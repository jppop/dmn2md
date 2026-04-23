# Getting started: a rule from a DMN file

This guide shows how to convert a DMN decision model into a **Claude Code rule** — a standing behavioral instruction that Claude applies automatically whenever it works with relevant files, without being explicitly invoked.

**What you will build**: a rule that instructs Claude to evaluate every loan application row in a spreadsheet folder using the exact decision logic from `loan-approval.dmn`.

---

## Rule vs. Skill — which one to use?

| | Rule | Skill |
| --- | --- | --- |
| **When loaded** | Automatically, every session (or when matching files are opened) | On demand — you type `/skill-name` or Claude decides it is relevant |
| **Best for** | Always-on constraints, evaluation criteria, domain knowledge that must apply silently | Repeatable workflows you invoke explicitly |
| **This guide** | ✓ | See [skill-from-dmn.md](skill-from-dmn.md) |

A rule is the right choice here because we want Claude to *know* the loan approval logic before it even opens a file — not be asked to evaluate each application manually.

---

## Prerequisites

- `dmn2md` installed:
  ```bash
  go install github.com/jppop/dmn2md/cmd/dmn2md@latest
  ```
- Python 3 with `openpyxl` if you work with `.xlsx` files:
  ```bash
  pip install openpyxl
  ```

---

## Step 1 — Generate the rule file

```bash
mkdir -p .claude/rules
dmn2md --mode rule loan-approval.dmn > .claude/rules/loan-approval.md
```

The generated `.claude/rules/loan-approval.md` contains:

```markdown
# Loan Approval

Apply these rules whenever evaluating a Loan Approval.

## Applicant Risk Rating

Determines `riskRating`. Apply the first matching rule:

- `age` < 18 → `riskRating` = "high" _Minor applicant — always high risk_
- `age` >= 18, `employmentStatus` "employed","self-employed", `creditScore` >= 750 → `riskRating` = "low" _Excellent credit score_
- `age` >= 18, `employmentStatus` "employed","self-employed", `creditScore` [600..750) → `riskRating` = "medium" _Good profile_
- `age` >= 18 → `riskRating` = "high" _Poor credit or unemployed_

## Loan Eligibility

_Requires_: `Applicant Risk Rating`

Determines `isEligible`. Exactly one rule must match:

- `riskRating` "high" → `isEligible` = false _High risk — always rejected_
- `riskRating` "low", `loanAmount` <= 100000, `loanDuration` <= 240 → `isEligible` = true _Low risk — any amount and duration accepted_
- `riskRating` "medium", `loanAmount` <= 50000, `loanDuration` <= 120 → `isEligible` = true _Medium risk — limited amount and duration_
- _(any)_ → `isEligible` = false _All other cases rejected_

## Final Decision

_Requires_: `Loan Eligibility`, `Applicant Risk Rating`

Computes `finalDecision` (`string`):

```feel
if isEligible and riskRating = "low" then "approved"
else if isEligible and riskRating = "medium" then "approved with conditions"
else "rejected"
```

## Constraints

- Always evaluate decisions in this order: `Applicant Risk Rating` → `Loan Eligibility` → `Final Decision`.
- Always state which rule determined each intermediate value.
```

---

## Step 2 — Scope the rule to application files (optional but recommended)

By default a rule loads at every session start. To load it only when Claude works with files in your `applications/` folder, add a `paths` frontmatter:

```bash
# prepend frontmatter to the generated file
cat > .claude/rules/loan-approval.md << 'EOF'
---
paths:
  - "applications/**"
---

EOF
dmn2md --mode rule loan-approval.dmn >> .claude/rules/loan-approval.md
```

With `paths` set, Claude loads the rule automatically the moment it reads any file under `applications/` — and only then.

---

## Step 3 — Set up the example project

This example processes a folder of Excel spreadsheets. Each sheet contains loan applicant rows; Claude evaluates each one and writes the results back as new columns.

### Project layout

```
loan-processor/
├── .claude/
│   └── rules/
│       └── loan-approval.md       ← generated rule (commit this)
├── applications/
│   ├── batch-2024-q1.xlsx         ← input spreadsheets
│   └── batch-2024-q2.xlsx
├── results/                       ← Claude writes enriched files here
└── scripts/
    └── read_xlsx.py               ← helper: dump xlsx rows as JSON
```

### Input spreadsheet format

Each `.xlsx` file has one sheet with these columns (column names must match exactly):

| applicant_id | age | employmentStatus | creditScore | loanAmount | loanDuration |
| --- | --- | --- | --- | --- | --- |
| A001 | 25 | employed | 680 | 40000 | 84 |
| A002 | 17 | student | 500 | 10000 | 36 |
| A003 | 35 | self-employed | 760 | 95000 | 180 |
| A004 | 42 | unemployed | 590 | 20000 | 60 |

### Helper script — `scripts/read_xlsx.py`

Claude cannot open binary `.xlsx` files directly. This small script dumps a sheet to JSON so Claude can read the rows:

```python
import json, sys
import openpyxl

def xlsx_to_json(path: str) -> list[dict]:
    wb = openpyxl.load_workbook(path, read_only=True, data_only=True)
    ws = wb.active
    rows = list(ws.iter_rows(values_only=True))
    if not rows:
        return []
    headers = [str(h) for h in rows[0]]
    return [dict(zip(headers, row)) for row in rows[1:]]

if __name__ == "__main__":
    print(json.dumps(xlsx_to_json(sys.argv[1]), indent=2, default=str))
```

---

## Step 4 — Ask Claude to evaluate the applications

With the rule loaded, open Claude Code in the project and ask:

```
Evaluate all loan applications in applications/batch-2024-q1.xlsx.
For each row, apply the Loan Approval decision and add three columns:
riskRating, isEligible, and finalDecision.
Write the result to results/batch-2024-q1.xlsx.
```

Because the rule is scoped to `applications/**`, Claude loads it the moment it touches the input file. It then processes each row step by step:

```
Reading applications/batch-2024-q1.xlsx…
Running scripts/read_xlsx.py to extract rows.

Evaluating 4 applicants using the Loan Approval rule:

A001 — age=25, employed, score=680, amount=40000, duration=84
  Step 1 Applicant Risk Rating (FIRST): rule 3 matches → riskRating = "medium"
  Step 2 Loan Eligibility (UNIQUE): rule 3 matches (medium, ≤50000, ≤120 ✓) → isEligible = true
  Step 3 Final Decision: isEligible=true, riskRating="medium" → "approved with conditions"

A002 — age=17, …
  Step 1 Applicant Risk Rating (FIRST): rule 1 matches (age < 18) → riskRating = "high"
  Step 2 Loan Eligibility (UNIQUE): rule 1 matches → isEligible = false
  Step 3 Final Decision: isEligible=false → "rejected"

A003 — age=35, self-employed, score=760, amount=95000, duration=180
  Step 1 Applicant Risk Rating (FIRST): rule 2 matches (≥18, self-employed, ≥750) → riskRating = "low"
  Step 2 Loan Eligibility (UNIQUE): rule 2 matches (low, ≤100000 ✓, ≤240 ✓) → isEligible = true
  Step 3 Final Decision: isEligible=true, riskRating="low" → "approved"

A004 — age=42, unemployed, score=590, amount=20000, duration=60
  Step 1 Applicant Risk Rating (FIRST): rule 4 matches (≥18, any) → riskRating = "high"
  Step 2 Loan Eligibility (UNIQUE): rule 1 matches → isEligible = false
  Step 3 Final Decision: isEligible=false → "rejected"

Writing results/batch-2024-q1.xlsx with columns riskRating, isEligible, finalDecision added.
```

### Output spreadsheet

| applicant_id | age | employmentStatus | creditScore | loanAmount | loanDuration | riskRating | isEligible | finalDecision |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| A001 | 25 | employed | 680 | 40000 | 84 | medium | true | approved with conditions |
| A002 | 17 | student | 500 | 10000 | 36 | high | false | rejected |
| A003 | 35 | self-employed | 760 | 95000 | 180 | low | true | approved |
| A004 | 42 | unemployed | 590 | 20000 | 60 | high | false | rejected |

---

## Step 5 — Keep the rule in sync

When the DMN changes, regenerate:

```bash
# if you added paths frontmatter manually, preserve it
dmn2md --mode rule loan-approval.dmn > /tmp/rule-body.md
# then re-prepend your frontmatter as in Step 2
```

Commit both the `.dmn` source and `.claude/rules/loan-approval.md` together.

---

## Why this works

The rule file is pure text — a list of conditions and outcomes derived directly from the DMN source. Claude treats it as standing knowledge: it knows the decision logic before you ask a single question. You never have to paste the decision table into the chat or re-explain the criteria for each batch.

When the DMN model is the single source of truth and `dmn2md` regenerates the rule on each change, the knowledge Claude applies is always consistent with what the business analysts designed.

---

## Next step

For interactive, on-demand evaluation (one applicant at a time with guided input), see the skill version in [skill-from-dmn.md](skill-from-dmn.md). Rules and skills complement each other: use the rule to ensure every batch run is consistent, and the skill when a user wants to walk through a single case manually.
