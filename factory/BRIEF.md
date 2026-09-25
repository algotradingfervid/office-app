# BRIEF — Office App (working name)

Imported on 2026-09-25 from the pre-kit brainstorming session (Claude Code, same day). The full design
it led to is `docs/superpowers/specs/2026-09-25-office-app-leave-design.md`.

## The problem (in the user's words)
The owner runs a company of ~200–250 employees (India). Day-to-day forms — leave, leave approval,
vouchers and voucher approval, procurement requests, advance requests — are handled by hand. The owner
wants one app every employee uses regularly, where new forms can be added "as and when required",
without all requirements known up front, and where changing a workflow never breaks existing data or
stops people working.

## Why now
The company has grown past the point where paper and e-mail approvals are traceable; the owner wants
leave balances and approvals enforced by rules (e.g., CL credited monthly, max 2 consecutive days,
LOP capped per year) instead of by memory.

## Directions considered
1. **Configure an existing ERP (ERPNext/Frappe)** — every form exists already · rejected: too complex.
2. **No-code platform (NocoBase, Budibase, Appsmith)** — forms without code · rejected: not Go, heavier
   to run and understand.
3. **Plain Go from scratch** — full control · rejected: rebuilds auth, admin UI and files.
4. **PocketBase as a Go framework + server-rendered pages, one form module at a time** — chosen.

## Chosen direction
A single Go binary built on PocketBase (SQLite, auth, files, cron, backups), server-rendered with Go
templates + htmx + Pico.css. A shared core (users, requests, an approver-chosen forwarding workflow,
append-only history, notifications, audit log) plus one self-contained module per form. Leave (with
balances) is the first form. Free and open source, very simple, backend only in Go, popular libraries
over hand-written code. Hosted on an office Ubuntu machine, published via Cloudflare Tunnel, managed over
Tailscale.

## Riskiest assumption
**Correctness of leave balances over time** (monthly credits, year close, carry-forward, missed job runs,
late cancellations): if balances drift, HR stops trusting the app and employees go back to paper.
Demand is not the risk — the owner is the customer and mandates use.
Cheapest test: the Core tier's end-to-end check plus table-driven tests on a fixed clock, including a
simulated two-month outage and a year close with pending requests. Route: no separate `/spike`
(the PocketBase mechanics were already checked against v0.40.4 source in the design review, §13);
the risk is tested inside the Core and Usable tiers.

## Dials
See factory/DIALS.md.
