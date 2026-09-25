# GLOSSARY — one meaning per word

Use these words in code, UI and stories. Code identifiers in `backticks`. Full rules: design spec §5–§7.

## People
| Term | Meaning | Code |
|---|---|---|
| Employee | Any user of the app; one `users` record, identified by employee code | `users`, `employee_code` |
| Approver | An employee allowed to act on requests | `users.can_approve` |
| HR admin | An employee who manages employees, rules, holidays, balances | `users.is_hr_admin` |
| Requester | The employee a request belongs to | `requests.requester` |
| Current approver | The one approver who may act on a pending request now | `requests.current_approver` |
| Final approver | The approver who approved as final | `requests.final_approver` |
| Active | Can sign in; inactive employees are never deleted | `users.active` |

## Requests and approvals (any form)
| Term | Meaning | Code |
|---|---|---|
| Request | One submitted form of any type | `requests` |
| Form type | Which form a request is (`leave` first) | `requests.form_type` |
| Status | `pending`, `approved`, `rejected`, `cancelled`, `cancel_requested` | `requests.status` |
| Action | submit, forward, approveFinal, reject, cancel, requestCancel, approveCancel, declineCancel, hrCancel, reassign, recordOnBehalf | transition table in `internal/core/approvals` |
| Approval step | One recorded action on a request; never edited | `approval_steps` |
| Forward | Approve and send to another approver (max 5 per request) | step `forwarded` |
| Approve as final | Close the request as approved; leave balance is debited now | step `approved_final` |
| Record on behalf | HR creates an already-approved request for an employee | step `recorded_by_hr` |
| Reassign | HR changes the current approver | step `reassigned` |
| Form hooks | What a form plugs into approvals: Validate, OnFinalApproved, OnCancelledAfterApproval, Summary | `approvals` hooks interface |
| Error / warning | Error blocks submit or approval; warning is shown to the approver only (e.g. short notice) | `Validate` returns both |

## Leave
| Term | Meaning | Code |
|---|---|---|
| Leave type | CL, SL, EL, ML, PL, BL, MRL, CO (inactive), LOP | `leave_types.code` |
| Leave rule | The settings of a leave type from an effective date; a change is a new version | `leave_rules`, `effective_from` |
| Leave request | The leave form data of a request | `leave_requests` |
| Session | `full`, `first_half`, `second_half` of a day | `from_session`, `to_session` |
| Days | Days a leave request counts (working or calendar days per rule, halves = 0.5) | `leave_requests.days` |
| Working day | Not a weekly off and not a holiday | `internal/core/calendar` |
| Weekly off | Sundays and the 2nd and 4th Saturdays (setting) | `settings` weekly-off pattern |
| Holiday | A company holiday date | `holidays` |
| Leave year | 1 April – 31 March, written `2026-27` | `leave_year` |
| Ledger entry | One signed change to a balance; never edited | `leave_ledger` |
| Entry type | `opening`, `credit`, `debit`, `reversal`, `carry_forward`, `lapse`, `adjustment` | `leave_ledger.entry_type` |
| Balance | Sum of ledger entries for employee + type + leave year | — |
| Available | Balance minus days of that type in pending requests | — |
| Credit | Days added by the monthly or yearly job | entry `credit`, `period_key` |
| Period key | Makes each automatic entry unique (`2026-11`, `2026-27`, `2026-27-close`) | `leave_ledger.period_key` |
| Accrual start period | First month the monthly job may credit (go-live month is in opening balances) | setting `accrual_start_period` |
| Year close | Carry forward EL up to the cap and lapse the rest on 1 April | entries `carry_forward`, `lapse` |
| LOP | Loss of pay leave: no balance, capped per year | leave type `LOP`, rule `yearly_cap` |

## Platform
| Term | Meaning | Code |
|---|---|---|
| Module | A folder that registers itself; core modules or a form | `internal/core/*`, `internal/forms/*` |
| Part | One feature file's registration inside a module | `addPart` |
| Card | A block a module adds to the home page | `home.AddCard` |
| Notification | An in-app message about a request | `notifications` (Usable tier) |
| Audit log | Record of every HR/admin change; never edited | `audit_log` (Usable tier) |
