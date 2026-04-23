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
