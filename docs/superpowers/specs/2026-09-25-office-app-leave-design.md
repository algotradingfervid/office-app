# Office App — Phase 1 (Platform Core + Leave) Design

**Date:** 2026-09-25 (revised 2026-09-25 after independent review — see §13)
**Status:** Draft for review
**Related:** `docs/hr-policy/leave-policy.md` (sample leave policy — source of the default leave rules)

---

## 1. Goal and success criteria

An in-house web app for ~200–250 employees that replaces paper/e-mail forms and grows **one form at a time**.
Phase 1 delivers the shared platform core plus the **Leave** form with balances. Later forms (Voucher,
Advance, Procurement, Comp-off claim, …) reuse the same core.

Phase 1 is done when:

1. HR admin can create employees (one by one or CSV import), mark approvers and HR admins, and load
   opening leave balances.
2. An employee can apply for leave, see live day count and policy errors/warnings, and pick a first approver.
3. Approvers can approve & forward, approve as final, or reject; final approval deducts the balance.
4. Balances are credited automatically (monthly CL/EL, yearly SL) and closed/carried forward at year end,
   and a missed job run (machine off) is caught up automatically without double credits.
5. Every action is recorded in an uneditable history; every balance is explainable from its ledger;
   every admin change (employees, flags, rules, holidays, adjustments) is in an audit log.
6. HR can change leave rules with an effective date without affecting past or in-flight requests.
7. HR can export leave taken and LOP days per employee per month (CSV) for payroll.
8. The app runs on the office Ubuntu machine, reachable publicly via Cloudflare Tunnel, with off-site
   backups twice a day, backup/uptime alerts, and a one-command deploy with automatic rollback.

## 2. Decisions and constraints

| Topic | Decision |
|-------|----------|
| Cost / licensing | Free and open-source only (free tiers of Cloudflare, healthchecks.io, UptimeRobot are acceptable) |
| Simplicity | Very simple; ERPNext-class systems rejected as too complex |
| Language | **Backend entirely in Go (≥ 1.27, required by PocketBase v0.40).** No non-Go backend services |
| Reuse | Use popular, well-maintained libraries where they remove real work; hand-write only business rules and glue that is shorter than the library integration (§3) |
| Platform | **PocketBase v0.40.x** used as a Go framework (embedded SQLite, auth, files, cron, backups, admin UI) |
| Front end | Go `html/template` + **htmx 2.x** + **Pico.css 2.x**, vendored into the binary with `embed`. No Alpine.js, no JS build step (§3) |
| Login | Server-rendered login form; email **or employee code** + password; accounts created by HR admin (no self sign-up); HttpOnly cookie holding a PocketBase auth token (§8) |
| Approvals | Employee picks the first approver from users with *Can approve*; each approver forwards, approves as final, or rejects |
| Notifications | In-app only (Phase 1). Microsoft 365 e-mail later, from the same events |
| Hosting | Office Ubuntu machine; admin access via Tailscale; public access via Cloudflare Tunnel with a **locally-managed** `config.yml` (domain TBD by owner) |
| Locale | India; time zone Asia/Kolkata (set explicitly in code, not relied on from the OS); leave year April–March |
| Working week | Mon–Sat; Sundays and 2nd & 4th Saturdays are off (configurable); one pattern and one holiday list for all employees |
| Dates | Leave dates are calendar dates stored as text `YYYY-MM-DD`, never as datetimes (§5) |

### Out of scope for Phase 1
E-mail notifications, other forms (Voucher/Advance/Procurement), **Compensatory Off (earning and use —
becomes the "Comp-off claim" form in a later phase)**, payroll integration (a CSV export is in scope),
leave encashment payout calculation, restricted/optional holidays, per-location holiday lists or weekly
offs, team leave calendar, partial cancellation / early return (cancel and reapply), mobile app,
SSO with Microsoft 365, self sign-up, self-service password reset.

### Assumptions
- Three roles: **Employee** (everyone), **Approver** (flag), **HR admin** (flag). A user can hold several.
- All employees work the same weekly-off pattern and holiday list (single location).
- Eligibility that needs personal facts the app does not store (maternity: number of children, 80 days
  worked; paternity: date of birth of child; bereavement: relationship) is checked by the approver, not
  the app.
- Not every employee may have a company e-mail; login by employee code covers them (see §12).

## 3. Libraries

| Need | Library | Notes |
|------|---------|-------|
| Auth, password hashing, auth tokens, rate limiting | PocketBase built-ins | Cookie session is thin glue over `record.NewAuthToken()` / `app.FindAuthRecordByToken()` (§8) |
| Database, schema migrations, admin dashboard | PocketBase (SQLite, Go migrations) | Registered migrations run automatically on `serve`; `Automigrate` **off** in production |
| File uploads | PocketBase file field (`protected`) | Served only through the app's own authorised handler (§8) |
| Scheduled jobs | PocketBase cron | Defaults to **UTC** and does not catch up missed runs → `app.Cron().SetTimezone(Asia/Kolkata)` at startup, and jobs are idempotent catch-up jobs (§5.4) |
| Backups + S3-compatible off-site upload | PocketBase built-ins | Target: Cloudflare R2 (endpoint `<acct>.r2.cloudflarestorage.com`, region `auto`, path-style). Uses `VACUUM INTO`, short write lock only |
| Request status machine | **Plain Go transition table** (no library) | See decision below |
| Working days / holidays | **Plain Go function** over the holiday table (no library) | See decision below |
| CSRF protection | Go stdlib `net/http` `CrossOriginProtection` | Uses `Sec-Fetch-Site`, falls back to `Origin` vs `Host` (needs correct Host — §8) |
| Templates | Go stdlib `html/template` | |
| CSV import/export | Go stdlib `encoding/csv` | One fixed format per file; explicit column mapping gives per-row error messages |
| Front end | htmx 2.0.x, Pico.css 2.1.x | Vendored files, pinned versions, checksums recorded |
| Approver picker | Native `<select>` (approvers are a small subset of users); Tom Select 2.x only if the list exceeds ~50 | |
| Tests | Go stdlib `testing` + PocketBase `tests` package (`tests.NewTestApp`, `tests.ApiScenario`) | |

**Library decisions from the review (evidence checked 2026-09-25):**

- **`qmuntal/stateless` dropped.** It is maintained (v1.8.0, Feb 2026) but for five statuses a
  `map[status][]action` table plus a per-action permission check is shorter and easier to test, and the
  library does not help with the actual hard part (checking status and writing inside one DB
  transaction). `looplab/fsm` was considered for the same reason and rejected.
- **`rickar/cal/v2` dropped.** It is well maintained (v2.1.31, Sep 2026) but has no India package; every
  holiday from our table would have to be converted to a `cal.Holiday` with `StartYear/EndYear` and
  `Func: CalcDayOfMonth` (omitting `Func` silently disables the holiday), and half-days are not modelled.
  A ~20-line function over a `map[date]bool` of holidays plus the weekly-off pattern is clearer and
  fully unit-tested.
- **Alpine.js dropped.** Every interactive need (live day count, balance-after, policy messages) is a
  server round-trip via htmx calling the same Go validation functions. Removing Alpine lets us run a
  strict CSP (`script-src 'self'`); Alpine's standard build needs `'unsafe-eval'` and its CSP build
  forbids common syntax.
- **htmx 2.x, not 4.x.** 2.0.11 (Sep 2026) is the maintained stable line; 4.0 is new. Configure
  `allowEval: false`, `includeIndicatorStyles: false` (ship indicator CSS in our stylesheet) so CSP needs
  no `unsafe-*`. State-changing actions are never GET.
- **Pico.css kept.** Last release 2.1.1 (Mar 2025) — dormant but stable, and it is CSS only; vendored,
  so abandonment has no runtime risk.

**Before implementation:** pin exact versions (Go 1.27.x, PocketBase v0.40.x, htmx 2.0.x, Pico 2.1.1) in
`go.mod` / a `web/static/VERSIONS` file. PocketBase is pre-v1.0 (breaking changes happen, e.g., v0.40
moved to `encoding/json/v2`), so upgrades are deliberate and rehearsed (§9).

## 4. Architecture

```
Browser (phone / laptop)
      │ HTTPS
Cloudflare edge ── Tunnel (cloudflared, local config.yml) ──► 127.0.0.1:8090
                    │  blocks /_/* and /api/* (http_status:404)
                    │  sets httpHostHeader = public hostname
Tailscale (admins) ── tailscale serve ──► 127.0.0.1:8090   (dashboard /_/ reachable only this way)

officeapp (single Go binary, systemd service)
 ├─ web/         templates + static (htmx, Pico) via embed.FS
 ├─ core/        users, sessions, requests, approvals, notifications, audit log, layout, auth guards,
 │               clock, working-day calendar
 ├─ forms/leave/ leave types, rules, ledger, policy checks, jobs, reports
 ├─ migrations/  numbered Go migrations, applied automatically on start
 └─ PocketBase   auth, SQLite, files, cron, backups, admin dashboard
```

**Module boundaries**

- `core/` knows nothing about leave. It exposes an approval service (submit, forward, approveFinal,
  reject, cancel, requestCancel, decideCancel, reassign, recordOnBehalf) and a hook interface a form
  module implements. All hooks receive the transaction app (`txApp`) — never the outer `app`, which
  would read outside the transaction or deadlock on PocketBase's single write connection:
  - `Validate(txApp, request, now) (warnings []string, err error)` — called on submit (and preview),
    and again on final approval. Errors block; warnings are shown to the approver.
  - `OnFinalApproved(txApp, request) error` — e.g., leave writes ledger debits.
  - `OnCancelledAfterApproval(txApp, request) error` — e.g., leave writes reversal entries.
  - `Summary(request) string` — one line for inboxes and notifications.
- `core/` owns a `Clock` interface (`Now()` in Asia/Kolkata); all date logic takes it, so tests can
  fix "today".
- `forms/leave/` implements those hooks plus its own screens, rules, ledger and jobs.
- A new form = a new `forms/<name>/` folder + migrations; no changes to other forms.

## 5. Data model

All collections are PocketBase collections. **All collection API rules are `nil` (superuser-only)**: the
app is server-rendered, so the PocketBase REST API is not used by browsers, and all access goes through
the app's own handlers with their own authorisation (§8).

"Append-only" means: (1) API rules `nil`; (2) `OnRecordUpdate` / `OnRecordDelete` hooks on the collection
return an error — this also blocks edits from the superuser dashboard; (3) app code never updates or
deletes them. Corrections are new entries.

Calendar dates (leave dates, holiday dates, date_of_joining, probation_end, effective_from) are **text
fields `YYYY-MM-DD`** validated by pattern. PocketBase `date` fields store UTC datetimes, so an IST
midnight would be saved as the previous day. Timestamps (created, submitted_at, read_at) stay as
PocketBase datetimes in UTC and are displayed in Asia/Kolkata.

Days are PocketBase numbers constrained to multiples of 0.5 (halves are exact in binary floating point,
so sums never drift).

### 5.1 Core

| Collection | Fields | Rules |
|------------|--------|-------|
| `users` (auth) | name, email (optional), employee_code (unique, login identity), department, designation, date_of_joining, probation_end, can_approve, is_hr_admin, active, must_change_password | Never deleted. Identity fields = email, employee_code |
| `requests` | form_type, requester, status (`pending` / `approved` / `rejected` / `cancelled` / `cancel_requested`), current_approver, final_approver, submitted_at, form_version, recorded_by_hr (bool) | One row per request of any form |
| `approval_steps` | request, seq, actor, action (`submitted` / `forwarded` / `approved_final` / `rejected` / `cancelled` / `cancel_requested` / `cancel_approved` / `cancel_declined` / `reassigned` / `recorded_by_hr`), comment, to_user, warnings (JSON), created | **Append-only**; unique (request, seq) |
| `notifications` | recipient, request, message, read_at | Deleted after 180 days by a cleanup job |
| `audit_log` | actor, action, target_collection, target_id, before (JSON), after (JSON), created | **Append-only**; written by the app for every HR/admin change (users, flags, rules, holidays, settings, adjustments, imports); kept ≥ 1 year (DPDP Rules, §8) |

### 5.2 Leave

| Collection | Fields | Rules |
|------------|--------|-------|
| `leave_types` | code (CL, SL, EL, ML, PL, BL, MRL, CO, LOP), name, active | CO seeded **inactive** in Phase 1 |
| `leave_rules` | leave_type, effective_from, days_per_year (nullable = no quota), credit_method (`yearly` / `monthly` / `per_event`), monthly_credit, count_mode (`working_days` / `calendar_days`), max_consecutive_days, yearly_cap, max_times_per_employment (nullable), min_notice_days, max_backdate_days, half_day_allowed, attachment_after_days, carry_forward_cap, encashable, allowed_in_probation | **Versioned:** a change adds a row with a new `effective_from`. Rows in effect (effective_from ≤ today) or referenced by any request are read-only; future-dated, unreferenced rows may be edited/deleted |
| `leave_requests` | request, leave_type, from_date, from_session (`full` / `second_half`), to_date, to_session (`full` / `first_half`), days, leave_year, reason, attachment (protected file, PDF/JPG/PNG, ≤ 5 MB), rule (relation to the `leave_rules` row applied) | Single-day request: from_date = to_date and one session value `full` / `first_half` / `second_half` |
| `leave_ledger` | employee, leave_type, leave_year (e.g. `2026-27`), entry_type (`opening` / `credit` / `debit` / `reversal` / `carry_forward` / `lapse` / `adjustment`), days (signed), request (optional), period_key (optional), note, created_by | **Append-only**. Partial unique index on (employee, leave_type, entry_type, period_key) `WHERE period_key != ''` |
| `holidays` | date (unique), name | Adding/removing a holiday in a past or current period is audited |
| `settings` | key, value (JSON) — weekly-off pattern, `accrual_start_period`, session defaults | Audited |

Default rule values (from the sample policy):
CL: monthly 1, working_days, max 2 consecutive, notice 1, backdate 0, half-day yes, CF 0.
SL: yearly 8, working_days, attachment after 2, notice 0, backdate 7, half-day yes, CF 0.
EL: monthly 1.5, working_days, notice 7, backdate 0, half-day yes, CF cap 30, not in probation.
ML: per_event, **calendar_days**, max 182. PL: per_event, max 5. BL: per_event, max 3.
MRL: per_event, max 3, max_times_per_employment 1. LOP: no quota, yearly_cap 15, max 5 consecutive.

### 5.3 Day counting and policy checks

**Day count.** `count_mode = working_days`: every date in range that is not a weekly off or holiday counts
1, except that `second_half` on from_date and `first_half` on to_date count 0.5. `calendar_days`
(maternity): every date counts. A request whose first or last date is a non-working day is rejected
(start and end on working days), and a half-day session on a non-working day is rejected.

**Checks on submit** (errors unless marked *warning*):

1. Leave type active and allowed in probation (`today < probation_end` → probation).
2. Dates valid; **request must not cross the leave-year boundary** (31 Mar / 1 Apr) — the employee submits
   two requests. This keeps each request charged to exactly one year's balance.
3. **No overlap** with the employee's own `pending`, `approved` or `cancel_requested` requests (any type,
   half-day aware: a first-half and a second-half on the same date do not overlap).
4. `from_date ≥ today − max_backdate_days`.
5. `from_date ≥ today + min_notice_days` — *warning* only (policy says "where possible"); the approver
   sees "short notice".
6. `days ≤ max_consecutive_days`, where days of **adjacent requests of the same type** (separated only by
   non-working days) are added together, so the limit cannot be bypassed by splitting.
7. Quota types: `days ≤ available` (below). No negative balances.
8. `yearly_cap`: approved + pending days of that type in the leave year + this request ≤ cap (this is how
   the LOP cap works).
9. `max_times_per_employment`: count of approved/pending requests of that type ever < limit.
10. Attachment present when `days > attachment_after_days`.
11. LOP while paid balance (CL or EL available > 0) exists — *warning* to the approver.

**Final approval** re-runs all checks with today's data, **recomputes `days`** (a holiday may have been
added), stores the recomputed value if it changed (shown to the approver before confirming), and debits
that value.

### 5.4 Balance rules

- **Balance** = sum of ledger entries for (employee, leave type, leave year).
- **Available** = balance − days of that type in `pending` requests for that leave year.
- **Rule applied** = the `leave_rules` row whose `effective_from` is the latest on or before the
  submission date; stored on the leave request and reused on final approval and cancellation.
  Credit jobs use the rule in effect on the credit date.
- **Opening balances** are loaded by HR (CSV, dry-run preview with per-row errors, then commit) as
  `opening` entries **as at go-live, including the go-live month's credit**. The setting
  `accrual_start_period` (e.g. `2026-11`) is the first month the monthly job credits; no automatic credit
  is ever written for an earlier period.
- **New joiners:** when HR creates an employee whose date_of_joining is inside the current leave year and
  on/after go-live, the app writes a pro-rata `credit` for yearly types
  (SL: annual quota × months remaining in the leave year including the joining month ÷ 12, rounded down
  to 0.5; period_key `<leave_year>-prorata`, e.g. `2026-27-prorata`). Monthly types credit **in advance on the 1st** for employees with
  date_of_joining on or before that 1st, so a mid-month joiner's first credit is the next 1st.
- **Exit:** HR deactivates the user (§6). Balances remain in the ledger for final settlement; the balance
  report on the exit date is the input for encashment/recovery (calculation out of scope).

### 5.5 Scheduled jobs (Asia/Kolkata)

PocketBase cron runs in UTC by default and never replays missed runs, and the office machine may be off
(power cuts). So every job is a **reconcile job**: it computes what entries *should* exist up to today and
writes only the missing ones. It runs **at startup and every hour** (cron set to Asia/Kolkata).
Idempotency comes from the partial unique index on `period_key`.

| Job | What it ensures |
|-----|-----------------|
| Monthly credit | For each period P from `accrual_start_period` to the current month: every employee active on the 1st of P with date_of_joining ≤ 1st of P has a `credit` with period_key `P` for each monthly type (CL 1, EL 1.5) |
| Year close | Once today ≥ 1 April of year Y+1, for each employee and type of leave year Y: carry forward (EL up to cap) and lapse the rest (CL, SL, over-cap EL). Old year gets `carry_forward` (−x) and `lapse` (−y), so its remaining balance equals only the days held by its pending requests; new year gets `carry_forward` (+x). period_key `Y-close` |
| Year credit | Every active employee has the yearly `credit` (SL) for the current leave year (period_key = leave year) |
| Notification cleanup | Deletes notifications older than 180 days |

**Ordering:** within one run, year close → year credit → monthly credit, so April's CL/EL go into the new year.

**Late changes to a closed year.** Days held by `pending` requests of year Y at close time are excluded
from carry-forward/lapse. When such a request is later approved, rejected or cancelled, or HR reverses a
past approved request of year Y, the app **recomputes year Y's close** (what carry_forward/lapse should
be given the final balance) and writes the difference as `adjustment` entries in years Y and Y+1 linked
to that request. Year Y therefore always closes to exactly 0.

Mid-year the EL balance may exceed the carry-forward cap; the cap applies only at year close.

## 6. Screens

**Everyone:** Login (email or employee code) · Forced password change when `must_change_password` ·
Profile (change password; changing it logs out other sessions).

**Employee:** Home (balance cards, recent requests, pending-approval count, upcoming holidays) ·
Apply for leave (live day count, balance after, live policy errors and warnings via htmx, first approver) ·
My requests (list + detail with history timeline, cancel / request cancellation) ·
Policy & holidays · Notifications (bell with unread count, htmx poll every 60 s).

**Approver:** Pending my approval inbox · Request detail with requester's balance, other leave in the
same period, warnings, and history; actions **Approve & forward** (next approver + comment),
**Approve as final**, **Reject** (comment required).

**HR admin:**
- Employees: create, edit, CSV import (dry run first), flags, reset password (sets a temporary
  password and `must_change_password`), deactivate.
  **Deactivate** is blocked while the user is the `current_approver` of any request (HR reassigns first);
  it cancels the user's own `pending` requests, refreshes their token key (all sessions end), and sets
  `active = false`.
- Leave rules (view versions, add new version with effective date).
- Holidays (per leave year; warning if the list for next year is empty in March).
- Balances (ledger view, adjustment with reason, CSV opening-balance upload).
- All requests (filters incl. "pending > 3 days", CSV export, reassign approver).
- **Record leave on behalf** (e.g., unauthorised absence as LOP, LOP beyond the cap, employee
  without access): creates the request directly `approved` with step `recorded_by_hr`; HR may override
  cap/notice/backdate checks with a mandatory reason; balance and overlap checks still apply.
- Reports: leave taken and LOP days per employee per calendar month (split by month), balances as at a
  date — CSV.
- Audit log view.

**PocketBase dashboard (`/_/`):** technical owner only, reachable over Tailscale only (§8). Collections
are never edited there in production — schema changes go through migrations.

## 7. Request flow

1. **Submit.** Server computes days, picks the rule version, runs all checks, then in one transaction:
   creates `requests` (pending, current_approver = chosen) + `leave_requests` + `approval_steps(submitted)`
   + notification. First approver: `can_approve`, active, not the requester.
2. **Approve & forward.** Only `current_approver`; target must have `can_approve`, be active, and not be
   the actor or the requester. A request can be forwarded at most 5 times. Step recorded; notification
   to new approver.
3. **Approve as final.** Re-runs checks (§5.3). On success: status `approved`, `final_approver` set,
   `current_approver` cleared, `OnFinalApproved` writes ledger `debit`, step recorded, requester
   notified. On failure: approver sees the reason and may reject.
4. **Reject.** Comment required; status `rejected`; requester notified.
5. **Cancel (pending).** Requester only; status `cancelled`.
6. **Cancel (approved, not started — from_date > today).** Requester requests → status
   `cancel_requested`, routed to `final_approver` (or to HR admins if that user is inactive) → on confirm,
   status `cancelled` and ledger `reversal`; on decline, back to `approved`.
7. **Cancel (started or past).** HR admin only; ledger `reversal` for the full request (and year-close
   recompute if the year is closed). If the employee returned early, HR cancels and records the actual
   days on the employee's behalf.
8. **Reassign.** HR admin changes `current_approver` of a `pending` or `cancel_requested` request;
   step recorded.
9. **No editing** of submitted requests — cancel and reapply.

**Transition table** (in `core/`):

| From | Action | To | Who |
|------|--------|----|-----|
| — | submit | pending | requester |
| — | recordOnBehalf | approved | HR admin |
| pending | forward | pending | current_approver |
| pending | approveFinal | approved | current_approver |
| pending | reject | rejected | current_approver |
| pending | cancel | cancelled | requester |
| pending / cancel_requested | reassign | (unchanged) | HR admin |
| approved | requestCancel | cancel_requested | requester, only if not started |
| approved | hrCancel | cancelled | HR admin |
| cancel_requested | approveCancel | cancelled | current_approver |
| cancel_requested | declineCancel | approved | current_approver |

### Error handling and consistency
- All rules are enforced server-side; live UI checks are previews that call the same Go functions.
- Each action (status + step + ledger + notification + audit) runs in one `app.RunInTransaction`, and
  **every read inside it uses `txApp`**. PocketBase runs write transactions on a single connection, so
  transactions are serialised: two submits cannot both pass the balance check, and double clicks or
  concurrent approvers cannot double-deduct (the status is re-read inside the transaction).
- Users see plain-language messages; internal errors are logged with a request ID shown to the user.

## 8. Security

- **Network.** App listens on `127.0.0.1:8090` only. Public traffic arrives via cloudflared with a
  locally-managed `config.yml` whose ingress rules, before the catch-all, return `http_status:404` for
  `path: ^/_/` and `path: ^/api/` (the browser never needs the PocketBase REST API or dashboard).
  Admins reach the app, including `/_/`, via `tailscale serve` (tailnet only).
- **Defence in depth in the app.** A middleware refuses `/_/` and `/api/` (except `/api/health`) for any
  request carrying `CF-Connecting-IP` / `Cf-Ray` (these are always added by Cloudflare and cannot be
  removed by a client going through the tunnel).
- **Host header.** cloudflared rewrites Host to the origin's hostname by default; set
  `originRequest.httpHostHeader` to the public hostname so `CrossOriginProtection`'s fallback and any
  absolute URLs work.
- **Real client IP.** Set PocketBase `TrustedProxy.Headers = ["CF-Connecting-IP"]`; without it every
  public client appears as `127.0.0.1` and shares a single rate-limit bucket.
- **Sessions.** On login the app verifies the password with PocketBase, issues `record.NewAuthToken()`,
  and stores it in a cookie `HttpOnly; Secure; SameSite=Lax; Path=/`. A middleware loads the user with
  `app.FindAuthRecordByToken(token, TokenTypeAuth)` and rejects inactive users on every request.
  Token lifetime 7 days (`AuthToken.Duration`). Logout clears the cookie; password change,
  password reset and deactivation call `RefreshTokenKey()`, which invalidates all of that user's tokens.
- **Rate limiting** (PocketBase rate limiter, which also matches custom routes): `POST /login` 5 per
  minute per IP for guests; general limit for authenticated users. Optionally the one free Cloudflare WAF
  rate-limit rule on `/login`.
- **Passwords.** Minimum length 10 (`PasswordField.Min`), no composition rules; temporary passwords
  force a change at next login.
- **CSRF.** `http.CrossOriginProtection` checked in a global router middleware; all state changes are
  POST; cookies `SameSite=Lax`.
- **Security headers** set by middleware: CSP `default-src 'self'; script-src 'self'; style-src 'self';
  img-src 'self' data:; frame-ancestors 'none'`, `X-Content-Type-Options: nosniff`,
  `Referrer-Policy: same-origin`. HSTS enabled at Cloudflare.
- **Authorisation** checked server-side on every page/action: employees see own requests; approvers see
  requests assigned to or acted on by them; HR admin sees all. Attachments (medical certificates) are
  PocketBase *protected* files, streamed only by the app's handler (`app.NewFilesystem()` → `Serve`)
  after the same check, with `Content-Disposition: attachment` for non-images and `Cache-Control: private,
  no-store`. Upload types restricted by content sniffing, not extension.
- **Personal data (India DPDP Act 2023 / Rules 2025).** Core duties apply from 13 May 2027 (possibly
  brought forward). Leave and medical records are processed for employment (Sec. 7, no consent needed),
  but security, retention and breach duties apply: access control and audit log as above (kept ≥ 1 year),
  full-disk encryption (LUKS) on the office machine, private R2 bucket, least-privilege R2 token,
  a one-page breach runbook (who notifies the Board and employees), and a retention setting for
  attachments of ex-employees (owner to decide the period; purge by a job, recorded in the audit log).

## 9. Change safety (non-negotiable rules)

1. Approval steps, ledger entries and the audit log are never edited or deleted (hooks enforce it).
2. Schema changes are additive; fields holding data are never dropped or renamed — they are retired.
3. Every schema change is a numbered Go migration applied automatically at startup; `Automigrate` is off
   in production and collections are not edited in the dashboard.
4. Rule changes are new versions with an effective date; in-flight and past requests keep their version.
5. Each request stores its `form_version`; workflow changes apply only to new requests.
6. Every release is rehearsed against a copy of the latest live backup before deploy.
7. PocketBase upgrades are separate releases, read against the CHANGELOG, and rehearsed like any other.

## 10. Testing

- **Unit (table-driven, fixed `Clock`):** day counting (weekly offs, 2nd/4th Saturday, holidays, half-day
  sessions, calendar-day mode, month/year spans, start/end on non-working day); every policy check
  including overlap, adjacency, cross-year block, backdate, notice warning, yearly cap, once-per-employment;
  balance/available maths; rule-version selection; transition table (every allowed and forbidden action
  per role).
- **Integration (PocketBase `tests.NewTestApp`):** submit → forward → final approve deducts; reject;
  cancel pending; cancel approved restores; double approve deducts once; two concurrent submits exceeding
  balance → one fails; rule change mid-flight leaves existing request untouched; monthly job run twice
  credits once; **jobs catch up after a simulated two-month outage**; year close with pending request,
  then approve/reject after close → year closes to 0; `accrual_start_period` respected; deactivation
  rules; append-only hooks reject update/delete; permission checks per role on every route and on
  attachment download; `/_/` and `/api/*` refused with Cloudflare headers.
- **Upgrade rehearsal:** run the new binary against a copy of the latest live backup.
- **Manual:** click-through checklist on laptop and phone before each release.
- Test-first development for business rules.

## 11. Deployment and operations

- Build: Go 1.27.x, `GOOS=linux` cross-compiled single binary (architecture of the office machine to be
  confirmed), static assets embedded.
- `deploy` script over Tailscale: copy new binary → stop service → snapshot `pb_data` (tar) → swap binary
  (keep previous) → start → health check `GET /api/health` plus a login-page check → on failure: stop,
  restore snapshot and previous binary, start. Deploy outside office hours (downtime of seconds).
- systemd units for `officeapp` and `cloudflared`: dedicated non-root user, `Restart=always`,
  `StateDirectory=officeapp`, `NoNewPrivileges=yes`, `ProtectSystem=strict`, `ProtectHome=yes`,
  `PrivateTmp=yes`, `PrivateDevices=yes`, `UMask=0077`, `CapabilityBoundingSet=`,
  `SystemCallFilter=@system-service`; checked with `systemd-analyze security`. Time zone is set in code
  (§3), not relied on from `TZ`.
- Backups: PocketBase backups at **02:00 and 14:00 IST** (cron time zone set in code), uploaded to a
  private Cloudflare R2 bucket (free tier 10 GB is ample), 14-day retention by R2 lifecycle rule; local
  keep 3; backup before each deploy; the backup includes uploaded attachments stored in `pb_data`
  (attachments stay on local storage, not S3). After each successful backup the app pings a
  **healthchecks.io** check; a missed ping alerts the owner. Quarterly restore drill on the laptop.
  Recovery point ≤ 12 h is accepted; Litestream (Go, maintained, supports R2) is the upgrade path if that
  becomes too much.
- Monitoring: UptimeRobot (free) on the public URL; machine on a UPS; NTP time sync enabled.
- Logs: PocketBase request/app logs, retention 30 days; audit log in the database (§5.1).
- Code: git repository; private GitHub remote recommended as off-site copy (owner to confirm).

## 12. Open items (owner to provide / confirm before or during implementation)

| Item | Needed by |
|------|-----------|
| Public domain / Cloudflare Tunnel hostname | Deployment |
| Office machine CPU architecture (amd64/arm64), Ubuntu version, disk encryption, UPS | Deployment |
| Cloudflare R2 bucket (or other S3-compatible storage) for off-site backups | Deployment |
| healthchecks.io / UptimeRobot accounts and alert recipient | Deployment |
| Do all employees have company e-mail? (If not, employee code is their login.) | Before build of login |
| Go-live date and `accrual_start_period` | Before go-live |
| Final leave numbers and policy text (sample in `docs/hr-policy/leave-policy.md`, aligned with this spec on 2026-09-25). Owner to confirm: CO not offered until the Comp-off claim form; SL backdating 7 days; forward limit 5 | Before go-live (editable in app) |
| Retention period for ex-employee attachments; named person for breach notification | Before go-live |
| Holiday list for the current leave year | Before go-live |
| Employee list + opening balances (CSV) | Before go-live |
| Private GitHub repo: yes/no | Setup |

## 13. Review log (2026-09-25)

Independent review of the first draft. Facts were checked against PocketBase v0.40.4 source, Go 1.25.5
`net/http` source, library repositories, and Cloudflare/Tailscale/healthchecks/R2 docs.

| # | Hole found in draft | Resolution |
|---|---------------------|------------|
| 1 | PocketBase v0.40 requires Go 1.27; spec assumed Go ≥ 1.25 (1.25.5 installed) | Pin Go 1.27.x (§2, §3, §11) |
| 2 | PocketBase cron runs in **UTC** by default and never replays missed runs; office machine may be off | `SetTimezone(Asia/Kolkata)`; jobs are hourly + startup reconcile jobs (§5.5) |
| 3 | PocketBase has no cookie sessions (auth only via `Authorization` header) — server-rendered login undefined | Cookie holding a PocketBase auth token, validated per request; `RefreshTokenKey()` for logout-everywhere (§8) |
| 4 | Behind the tunnel every client is `127.0.0.1` → shared rate-limit bucket | `TrustedProxy.Headers = CF-Connecting-IP` (§8) |
| 5 | cloudflared rewrites Host by default → CSRF fallback and URLs break | `httpHostHeader` set (§8) |
| 6 | PocketBase REST API (`/api/*`) and dashboard exposed publicly; per-collection API rules undefined | All API rules `nil`; `/_/` and `/api/` blocked in cloudflared ingress and in app middleware (§5, §8) |
| 7 | "Append-only" enforced only by API rules — dashboard edits bypass it | `OnRecordUpdate/Delete` hooks reject changes (§5) |
| 8 | PocketBase `date` fields are UTC datetimes → off-by-one for IST dates | Calendar dates as `YYYY-MM-DD` text (§5) |
| 9 | Medical certificates served by public file URL | Protected files, streamed by an authorised handler (§8) |
| 10 | Concurrency claim unproven | Verified: `RunInTransaction` uses the single write connection; rule that all reads use `txApp` (§4, §7) |
| 11 | No overlap check between requests | Added (§5.3) |
| 12 | Leave spanning 31 Mar/1 Apr unhandled | Blocked; submit two requests (§5.3) |
| 13 | Pending requests at year end + later cancellations of closed-year leave broke carry-forward/lapse | Pending holds excluded at close; recompute-and-adjust on later change (§5.5) |
| 14 | Carry-forward ledger sign convention unclear; period_key uniqueness clashed across entry types | Old year −, new year +; unique (employee, type, entry_type, period_key) (§5.2, §5.5) |
| 15 | Go-live mid-year: catch-up would double-credit the opening month | `accrual_start_period` setting (§5.4) |
| 16 | Maternity leave counted as working days (law: 26 weeks = calendar days) | `count_mode` per rule (§5.2, §5.3) |
| 17 | "Min notice" is "where possible" in policy but modelled as a hard block; no limit on backdating | Notice = warning to approver; `max_backdate_days` (§5.3) |
| 18 | Max-consecutive rule bypassable by adjacent requests; marriage leave "once" not modelled | Adjacent same-type requests summed; `max_times_per_employment` (§5.3) |
| 19 | CO in leave types but no way to earn it; expiry needs per-credit tracking | CO deferred to a later "Comp-off claim" form; type inactive (§2) |
| 20 | LOP beyond cap "HR exception", unauthorised absence, early return — no mechanism | HR *Record leave on behalf* (§6) |
| 21 | Half-day semantics (`from_half`/`to_half`) ambiguous; holiday added after submission | Explicit sessions; final approval recomputes days (§5.2, §5.3) |
| 22 | Deactivating an approver strands inbox; deactivated user sessions stay valid | Deactivation blocked until reassigned; token key refreshed (§6) |
| 23 | Password reset without e-mail undefined; staff without e-mail | Temporary password + forced change; login by employee code (§5.1, §6) |
| 24 | No audit of admin changes (flags, rules, holidays) — also DPDP Rule 6 logging | `audit_log`, append-only, ≥ 1 year (§5.1) |
| 25 | Payroll needs LOP per month; 250 employees created by hand | LOP/leave monthly CSV; employee CSV import (§1, §6) |
| 26 | DPDP Act/Rules not considered for medical data | Safeguards, retention, breach runbook (§8) |
| 27 | Nightly-only backup (24 h loss), no alert if backup/app silently fails | Twice-daily backups, healthchecks.io ping, UptimeRobot, UPS (§11) |
| 28 | Deploy "backup live DB" while running + DB rollback ill-defined | Stop → snapshot → swap → health check → restore on failure (§11) |
| 29 | `stateless`, `rickar/cal`, Alpine.js add more integration than they save | Replaced by a transition table, a tested day-count function, htmx-only UI (§3) |
| 30 | No way to fix "today" in tests | `Clock` interface (§4, §10) |
