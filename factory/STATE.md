# STATE — the floor board
<!-- One screen: at most 60 lines. Dated notes and history go to factory/HISTORY.md. -->

product-name: unnamed (working: Office App)
kit-version: see factory/KIT-VERSION.md

## Floor 1: prove it
brainstorm: done           # imported from pre-kit session 2026-09-25
research: done             # imported; evidence in design spec §3, §13 and the leave policy
spike: skipped: riskiest assumption (balance correctness) is tested in the Core/Usable e2e checks; owner agreed 2026-09-25
define: done               # spec-approved: 2026-09-25
blueprint: done            # skeleton: login + home; make check passes; scripts/preview.sh
tool-up: done              # line-open: 2026-09-25 (two trial runs on story 001; kit fixes in lessons LOG)

## Floor 2: make it a product (design-timing: see DIALS.md)
name-it: pending
brand: pending
design-system: pending
wireframe: pending

## Line
tier: core
tier-e2e: pending
agents-recommended: 4      # owner wants small, fast parallel stories
stories: scripts/ready.sh prints counts, waves and hotspots

## Waiting on the human
- (nothing)

## Open decisions
- Repo is PUBLIC (github.com/algotradingfervid/office-app): never commit secrets, real employee data or .env
- Launch tier must ensure no demo users (public password) exist in production; `officeapp seed` is for dev/previews only
- Design spec §12 open items (domain, machine arch, R2, e-mail coverage, go-live date)
