# SPEC — Office App (Phase 1: platform core + Leave)

approved-by: owner on 2026-09-25

**Detailed design (source of truth for rules, data model, security and operations):**
`docs/superpowers/specs/2026-09-25-office-app-leave-design.md` (revised 2026-09-25, commit 13f612d).
**Leave policy defaults:** `docs/hr-policy/leave-policy.md`.
This file tiers that design for the line; where they differ, the detailed design wins and this file is fixed.

## In one sentence (a user telling a friend)
"I apply for leave on my phone, pick my manager, and I can see my balance and exactly who approved it."

## Users
Primary: employees (~200–250) applying for leave and approvers acting on requests.
Secondary: HR admin managing employees, balances, rules, holidays and payroll exports.

## The first ten things a user does
| # | Action | Tier | Acceptance check (observable by a non-programmer) |
|---|---|---|---|
| 1 | Sign in with employee code or e-mail and password | Core | A seeded employee signs in on a phone-width browser and lands on Home; a wrong password shows a plain message |
| 2 | See my leave balances | Core | Home shows a card per active leave type with the ledger balance and "available" (balance − pending) |
| 3 | Apply for leave and pick a first approver | Core | Picking dates shows working days and "balance after" without reload; a 3-day CL shows "max 2 consecutive days" and cannot be submitted; a valid request appears as Pending |
| 4 | Approve & forward, approve as final, or reject | Core | Approver A forwards to B; B approves as final; the employee's CL balance drops by exactly the request's days; pressing Approve twice deducts once |
| 5 | See my request's status and history | Core | The request page shows every step (submitted, forwarded, approved) with who, when and comment; nobody can edit a step |
| 6 | HR sets up the company: employees (one-by-one and CSV), approver/HR flags, holidays, opening balances | Usable | HR imports a 20-row employee CSV with a dry-run preview, loads opening balances, adds holidays; each change is in the audit log |
| 7 | Balances credit themselves and close the year | Usable | With a fixed clock: monthly CL/EL credits appear on the 1st; a two-month outage is caught up without double credits; 1 April carries EL up to 30 and lapses CL/SL, and the old year closes to 0 even with a pending request |
| 8 | Cancel, reassign, record on behalf | Usable | Employee cancels a pending request; cancellation of approved future leave goes to the final approver and restores the balance; HR reassigns a stuck request and records LOP for an absent employee |
| 9 | HR changes a leave rule and everyone gets notified of their requests | Usable | A new CL rule with a future effective date leaves existing requests on the old version; the bell shows unread notifications for each request event |
| 10 | The app runs for the company | Launch | The public URL works over HTTPS; `/_/` and `/api/*` return 404 from the internet; a deploy with a broken binary rolls back by itself; a backup lands in R2 twice a day and a missed one alerts; HR downloads the monthly leave/LOP CSV |

## Tiers (the line builds one at a time)
| Tier | Actions | End-to-end check |
|---|---|---|
| Core | 1–5 (seeded data: sample employees, leave types, default rules, holidays, opening balances via a dev seed command; cookie login with CSRF protection; all day-count and submit checks of design §5.3) | `make e2e-core`: scripted run on a fresh DB — employee applies 2 CL days, A forwards to B, B approves final, balance 4 → 2, history has 3 steps; 3-day CL rejected with the policy message; double approve deducts once; two concurrent submits exceeding balance → exactly one succeeds. Plus a manual click-through on laptop and phone widths |
| Usable | 6–9 (HR admin screens, CSV imports, reconcile jobs, year close, cancellations, reassign, record-on-behalf, rule versions, notifications, audit log, deactivation) | `make e2e-usable`: fixed-clock scenario from go-live through a year close with an outage, cancellations and a rule change; every balance equals the ledger sum and the closed year is 0 |
| Launch | 10 (systemd, cloudflared config, tailscale serve, deploy with rollback, R2 backups + healthchecks.io, UptimeRobot, security headers/CSP, trusted proxy, rate limits, payroll CSV reports, attachment retention job, breach runbook) | On the office machine: the Launch acceptance check of action 10, plus a restore drill from R2 on the laptop |
| Later | E-mail via Microsoft 365, Comp-off claim form, Voucher, Advance, Procurement forms, team calendar, SSO | — |

## Not in version one (the "not yet" list, with a reason each)
- E-mail notifications — in-app is enough to start (owner's decision).
- Comp-off earning/use — needs its own claim form (design §2).
- Other forms — one form at a time; the core is built for them.
- Encashment payout calculation, payroll integration — CSV export covers payroll.
- Per-location holidays/weekly offs, restricted holidays — single location.
- Partial cancellation / early return — cancel and reapply.
- SSO, self sign-up, self-service password reset — HR resets passwords.

## Platforms
First: web, mobile-friendly (phone and laptop browsers). Then: nothing else planned.

## Success metrics
| Metric | Target | By |
|---|---|---|
| Employees who applied for leave in the app | ≥ 90 % of leave taken is recorded in the app | 2 months after go-live |
| Balance disputes raised with HR | 0 unexplained (every balance explained from its ledger) | Ongoing from go-live |
| Median time from submit to final decision | ≤ 1 working day | 1 month after go-live |

## Non-negotiables
- Free and open-source; backend only in Go; popular libraries over hand-written code.
- Change safety (design §9): history, ledger and audit log never edited; additive migrations; versioned rules.
- Personal data safeguards per India DPDP Act (design §8).
- Works on phone browsers. English only.

## Research cross-check
| Pain (RESEARCH.md) | Addressed by |
|---|---|
| 1 Balances by hand | Actions 2, 7 |
| 2 Untraceable approvals | Actions 4, 5 |
| 3 Limits not enforced | Action 3 (checks), 7 |
| 4 New app per form | Core module boundary (design §4); later forms in Later |
| 5 Changes break data | Action 9 (rule versions) + design §9 |
