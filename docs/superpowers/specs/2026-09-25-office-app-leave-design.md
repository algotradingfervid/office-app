# Office App — Phase 1 (Platform Core + Leave) Design

**Date:** 2026-09-25
**Status:** Draft for review
**Related:** `docs/hr-policy/leave-policy.md` (sample leave policy — source of the default leave rules)

---

## 1. Goal and success criteria

An in-house web app for ~200–250 employees that replaces paper/e-mail forms and grows **one form at a time**.
Phase 1 delivers the shared platform core plus the **Leave** form with balances. Later forms (Voucher,
Advance, Procurement, …) reuse the same core.

Phase 1 is done when:

1. HR admin can create employees, mark approvers and HR admins, and load opening leave balances.
2. An employee can apply for leave, see live working-day count and policy errors, and pick a first approver.
3. Approvers can approve & forward, approve as final, or reject; final approval deducts the balance.
4. Balances are credited automatically (monthly CL/EL, yearly SL) and rolled over at year end.
5. Every action is recorded in an uneditable history; every balance is explainable from its ledger.
6. HR can change leave rules with an effective date without affecting past or in-flight requests.
7. The app runs on the office Ubuntu machine, reachable publicly via Cloudflare Tunnel, with nightly
   off-site backups and a one-command deploy with automatic rollback.

## 2. Decisions and constraints

| Topic | Decision |
|-------|----------|
| Cost / licensing | Free and open-source only |
| Simplicity | Very simple; ERPNext-class systems rejected as too complex |
| Language | **Backend entirely in Go.** No non-Go backend services |
| Reuse | Use popular, well-maintained libraries wherever one fits; hand-write only business rules |
| Platform | **PocketBase** used as a Go framework (embedded SQLite, auth, files, cron, backups, admin UI) |
| Front end | Go `html/template` + **htmx** + **Alpine.js** + **Pico.css** (+ **Tom Select** only if the approver list gets long), served as static assets from the same binary. No JS build step |
| Login | PocketBase email + password; accounts created by HR admin (no self sign-up) |
| Approvals | Employee picks the first approver from users with *Can approve*; each approver forwards, approves as final, or rejects |
| Notifications | In-app only (Phase 1). Microsoft 365 e-mail later, from the same events |
| Hosting | Office Ubuntu machine; admin access via Tailscale; public access via Cloudflare Tunnel (domain TBD by owner) |
| Locale | India; time zone Asia/Kolkata; leave year April–March |
| Working week | Mon–Sat; Sundays and 2nd & 4th Saturdays are off (configurable) |

### Out of scope for Phase 1
E-mail notifications, other forms (Voucher/Advance/Procurement), payroll integration, leave encashment
payout calculation, team leave calendar, mobile app, SSO with Microsoft 365, self sign-up.

### Assumptions
- Three roles: **Employee** (everyone), **Approver** (flag), **HR admin** (flag). A user can hold several.
- Only sick leave beyond N days and similar rules need attachments.

## 3. Libraries

| Need | Library | Notes |
|------|---------|-------|
| Auth, sessions, password rules, rate limiting | PocketBase built-ins | Confirm rate-limit settings in pinned version |
| Database, schema migrations, admin dashboard | PocketBase (SQLite, Go migrations) | |
| File uploads | PocketBase | |
| Scheduled jobs | PocketBase cron | Confirm time-zone handling; must run in Asia/Kolkata |
| Backups + S3-compatible off-site upload | PocketBase built-ins | Target: Cloudflare R2 |
| Request status machine | `github.com/qmuntal/stateless` | External state storage + guards |
| Working days / holidays | `github.com/rickar/cal/v2` | Custom workday func for 2nd/4th Saturday |
| CSRF protection | Go stdlib `net/http` CrossOriginProtection (Go ≥ 1.25) | Confirm against installed Go |
| Templates | Go stdlib `html/template` | |
| Tests | Go stdlib `testing` + PocketBase `tests` package | |

**Before implementation:** verify maintenance status of `stateless` and `cal` (last release/commit dates),
and pin exact versions of PocketBase, Go and all libraries. Replace any library that appears abandoned.
PocketBase is pre-v1.0 (backward compatibility not guaranteed), so upgrades are deliberate and rehearsed
(§9).

## 4. Architecture

```
Browser (phone / laptop)
      │ HTTPS
Cloudflare Tunnel ──► Ubuntu office machine (127.0.0.1)
                       └─ officeapp (single Go binary, systemd service)
                           ├─ web/        templates + static (htmx, Alpine, Pico)
                           ├─ core/       users, requests, approvals, notifications, layout, auth guards
                           ├─ forms/leave/ leave types, rules, ledger, calendar, policy checks, jobs
                           ├─ migrations/ numbered Go migrations, auto-applied at startup
                           └─ PocketBase  auth, SQLite, files, cron, backups, admin dashboard
```

**Module boundaries**

- `core/` knows nothing about leave. It exposes an approval service (submit, forward, approveFinal,
  reject, cancel, reassign) and a hook interface a form module implements:
  - `Validate(request) error` — called on submit and again on final approval.
  - `OnFinalApproved(tx, request) error` — e.g., leave writes ledger debits.
  - `OnCancelledAfterApproval(tx, request) error` — e.g., leave writes reversal entries.
- `forms/leave/` implements those hooks plus its own screens, rules, ledger and jobs.
- A new form = a new `forms/<name>/` folder + migrations; no changes to other forms.

## 5. Data model

All collections are PocketBase collections. "Append-only" = API rules forbid update/delete, and app
code never updates or deletes them.

### 5.1 Core

| Collection | Fields | Rules |
|------------|--------|-------|
| `users` (auth) | name, email, employee_code (unique), department, designation, date_of_joining, probation_end, can_approve, is_hr_admin, active | Never deleted; deactivated users cannot log in |
| `requests` | form_type, requester, status (`pending` / `approved` / `rejected` / `cancelled` / `cancel_requested`), current_approver, final_approver, submitted_at, form_version | One row per request of any form |
| `approval_steps` | request, seq, actor, action (`submitted` / `forwarded` / `approved_final` / `rejected` / `cancelled` / `cancel_requested` / `cancel_approved` / `reassigned`), comment, to_user, created | **Append-only** |
| `notifications` | recipient, request, message, read_at | |

### 5.2 Leave

| Collection | Fields | Rules |
|------------|--------|-------|
| `leave_types` | code (CL, SL, EL, ML, PL, BL, MRL, CO, LOP), name, active | |
| `leave_rules` | leave_type, effective_from, days_per_year (nullable = no quota), credit_method (`yearly` / `monthly` / `per_event`), monthly_credit, max_consecutive_days, yearly_cap, min_notice_days, half_day_allowed, attachment_after_days, carry_forward_cap, encashable, allowed_in_probation | **Versioned:** a change adds a row with a new `effective_from`; rows already in effect are read-only |
| `leave_requests` | request, leave_type, from_date, from_half, to_date, to_half, working_days, reason, attachment, rule (relation to the `leave_rules` row applied) | |
| `leave_ledger` | employee, leave_type, leave_year, entry_type (`opening` / `credit` / `debit` / `reversal` / `carry_forward` / `lapse` / `adjustment`), days (signed, multiples of 0.5), request (optional), period_key (optional, unique with employee+leave_type for credits), note, created_by | **Append-only** |
| `holidays` | date (unique), name | |
| `settings` | key, value (JSON) — e.g., working-week pattern | |

### 5.3 Balance rules

- **Balance** = sum of ledger entries for (employee, leave type, leave year).
- **Available** = balance − working days of that type in `pending` requests.
- **Rule applied** = the `leave_rules` row whose `effective_from` is the latest on or before the
  submission date; stored on the leave request and reused on final approval.
- **LOP** has no quota; its cap is checked against the sum of LOP debits in the leave year.
- **Opening balances** are loaded once by HR (CSV upload) as `opening` entries.
- **New joiners:** when HR creates an employee, the app writes a pro-rata `credit` for yearly types
  (SL: annual quota × months remaining in the leave year including the joining month ÷ 12, rounded down
  to 0.5). Monthly types start with the next 1st-of-month run.

### 5.4 Scheduled jobs (Asia/Kolkata)

| When | Job |
|------|-----|
| 1st of each month, 00:30 | Credit monthly leave types (CL 1, EL 1.5) to employees active on that date with `date_of_joining` on or before it |
| 1 April, 00:15 | Credit yearly types (SL); carry forward EL up to cap; lapse CL/SL/over-cap EL |

Every credit carries a `period_key` (e.g., `2026-10`), and (employee, leave_type, period_key) is unique,
so re-running a job never double-credits.

## 6. Screens

**Employee:** Home (balance cards, recent requests, pending-approval count, upcoming holidays) ·
Apply for leave (live working-day count, balance after, live policy messages, first approver) ·
My requests (list + detail with history timeline, cancel / request cancellation) ·
Policy & holidays · Notifications (bell with unread count, polled every 30–60 s) · Profile (change password).

**Approver:** Pending my approval inbox · Request detail with requester's balance and history, actions
**Approve & forward** (next approver + comment), **Approve as final**, **Reject** (comment required).

**HR admin:** Employees (create, edit, deactivate, flags, reset password) · Leave rules (view versions,
add new version with effective date) · Holidays · Balances (ledger view, adjustment with reason, CSV
opening-balance upload) · All requests (filters, CSV export, reassign approver).

**PocketBase dashboard (`/_/`):** technical owner only, reachable over Tailscale only (§8).

## 7. Request flow

1. **Submit.** Server computes working days (rickar/cal: excludes Sundays, 2nd/4th Saturdays, holidays;
   half days = 0.5), picks the rule version, runs all checks, then in one transaction: creates
   `requests` (pending, current_approver = chosen) + `leave_requests` + `approval_steps(submitted)` +
   notification.
2. **Approve & forward.** Only `current_approver`; target must have `can_approve`, be active, and not be
   the actor or the requester. Step recorded; notification to new approver.
3. **Approve as final.** Re-runs `Validate` (balance may have changed). On success: status `approved`,
   `final_approver` set, `OnFinalApproved` writes ledger `debit`, step recorded, requester notified.
   On failure: approver sees the reason and may reject.
4. **Reject.** Comment required; status `rejected`; requester notified.
5. **Cancel (pending).** Requester only; status `cancelled`.
6. **Cancel (approved, not started).** Requester requests → status `cancel_requested`, routed to
   `final_approver` → on confirm, status `cancelled` and ledger `reversal`; on decline, back to `approved`.
7. **Cancel (started or past).** HR admin only; ledger `reversal` for the full request.
8. **Reassign.** HR admin changes `current_approver` (e.g., approver absent or deactivated); step recorded.
9. **No editing** of submitted requests — cancel and reapply.

Status transitions are defined once with `stateless`; guards enforce who may act.

### Error handling and consistency
- All rules are enforced server-side; live UI checks are previews that call the same Go functions.
- Each action (status + step + ledger + notification) runs in a single database transaction.
- The transition checks the current status inside the transaction, so double clicks or concurrent
  approvers cannot double-deduct.
- Users see plain-language messages; internal errors are logged, not shown.

## 8. Security

- HTTPS via Cloudflare; app listens on `127.0.0.1` only.
- PocketBase login rate limiting; minimum password length 10.
- Requests to `/_/` (PocketBase dashboard) and superuser APIs are refused when they arrive through the
  tunnel (identified by Cloudflare headers); allowed only via Tailscale.
- CSRF protection via Go `net/http` CrossOriginProtection.
- Authorisation checked server-side on every page/action: employees see own requests; approvers see
  requests assigned to or acted on by them; HR admin sees all.
- Session lifetime 7 days (configurable); deactivated users are logged out.

## 9. Change safety (non-negotiable rules)

1. Approval steps and ledger entries are never edited or deleted.
2. Schema changes are additive; fields holding data are never dropped or renamed — they are retired.
3. Every schema change is a numbered Go migration applied automatically at startup.
4. Rule changes are new versions with an effective date; in-flight and past requests keep their version.
5. Each request stores its `form_version`; workflow changes apply only to new requests.
6. Every release is rehearsed against a copy of the latest live backup before deploy.

## 10. Testing

- **Unit (table-driven):** working-day calculation (weekly offs, holidays, half days, month/year spans),
  every policy check, balance/available maths, credit/carry-forward/lapse, rule-version selection.
- **Integration (PocketBase test app):** submit → forward → final approve deducts; reject; cancel
  pending; cancel approved restores; double approve deducts once; rule change mid-flight leaves
  existing request untouched; monthly job run twice credits once; permission checks per role.
- **Upgrade rehearsal:** run the new binary against a copy of the latest live backup.
- **Manual:** click-through checklist on laptop and phone before each release.
- Test-first development for business rules.

## 11. Deployment and operations

- Build: `GOOS=linux` cross-compiled single binary (architecture of the office machine to be confirmed).
- `deploy` script over Tailscale: backup live DB → copy binary (keep previous) → restart service →
  health check → automatic rollback of binary + DB on failure.
- systemd units for `officeapp` (dedicated user, `Restart=always`, `TZ=Asia/Kolkata`) and `cloudflared`.
- Backups: nightly PocketBase backup, 14-day retention, uploaded to Cloudflare R2; backup before each
  deploy; quarterly restore drill on the laptop.
- Logs: PocketBase built-in request/app logs.
- Code: git repository; private GitHub remote recommended as off-site copy (owner to confirm).

## 12. Open items (owner to provide / confirm before or during implementation)

| Item | Needed by |
|------|-----------|
| Public domain / Cloudflare Tunnel hostname | Deployment |
| Office machine CPU architecture (amd64/arm64) and Ubuntu version | Deployment |
| Cloudflare R2 bucket (or other S3-compatible storage) for off-site backups | Deployment |
| Final leave numbers and policy text (sample in `docs/hr-policy/leave-policy.md`) | Before go-live (editable in app) |
| Holiday list for the current leave year | Before go-live |
| Employee list + opening balances (CSV) | Before go-live |
| Private GitHub repo: yes/no | Setup |
