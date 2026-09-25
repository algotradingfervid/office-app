# RESEARCH — Office App

Imported on 2026-09-25. Research ran before the kit was installed; its evidence lives in the design spec.

| Lane | Where | Verdict |
|------|-------|---------|
| Platform choice (ERPNext, NocoBase, Budibase, Appsmith, plain Go, PocketBase) | BRIEF.md "Directions considered" | PocketBase as a Go framework |
| Libraries (stateless, rickar/cal, looplab/fsm, Alpine.js, htmx, Pico.css, Tom Select) | design spec §3 | Only PocketBase, htmx 2.0.x, Pico 2.1.x; plain Go for status table and day count |
| PocketBase v0.40.4 mechanics (cron UTC/no catch-up, no cookie sessions, date fields UTC, single write connection, protected files, trusted proxy) | design spec §13 rows 1–10 | Designed around each |
| Cloudflare Tunnel / Tailscale / R2 / healthchecks.io | design spec §8, §11, §13 | Locally-managed tunnel config, path blocking, Host header, private R2 |
| Indian leave law and IT-company practice (Maternity Benefit Act / Code on Social Security 2020, S&E Acts, paternity, EL accrual) | docs/hr-policy/leave-policy.md (sources listed there) | Sample policy, HR-editable settings |
| India DPDP Act 2023 / Rules 2025 | design spec §8 | Security, audit ≥ 1 year, retention, breach runbook |

## Top pains (for the /define cross-check)
1. Leave balances tracked by hand and disputed.
2. Approvals untraceable (who approved, when).
3. Policy limits (consecutive CL, LOP cap) not enforced.
4. New forms need a new app each time.
5. Changes risk losing past data or stopping work.

## Riskiest-assumption verdict
Balance correctness over time — testable with a fixed clock; no spike needed (see BRIEF.md).
