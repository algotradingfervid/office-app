---
name: blueprint
description: "Choose the stack and a modular-monolith architecture for the product and build the walking skeleton: the thinnest deployable, end-to-end version, branded if the design rooms have run. Use after /define (or after /wireframe when design-timing is early), when the user says \"architecture\", \"stack\", \"skeleton\", \"set up the project\", or invokes /blueprint."
---

# /blueprint — Engineering office, station 1

Plan thoroughly the things that are expensive to change (data model, architecture shape, platform, hosting). Everything else is a story.

## Steps

0. Write the station marker so the guardrails let this station do its job: `echo blueprint > .claude/station`. Remove it when done: `rm -f .claude/station`.
1. Read `factory/SPEC.md` (tiers), `factory/RESEARCH.md` (constraints, feasibility), `factory/SPIKE.md` (the approach that passed, and its measurement), `factory/design/` if the design rooms have run, `factory/DIALS.md`.
2. **Decide**, and write the reasoning into `factory/ARCHITECTURE.md` using `factory/templates/architecture.md`:
   - Platform and framework per target (prefer boring, well-documented, agent-friendly choices with large ecosystems; note the tradeoff you are making).
   - Architecture shape: a **modular monolith** by default: one repo, one binary, one deploy. Each feature area is a module (a folder) that owns its routes, handlers, logic, queries, migrations and tests, with the layer rule inside it (for most apps: UI → services → data). Modules talk only through small interfaces. A separate process only for a workload reason (for example a background fetch or read worker), from the same codebase. Microservices only with a named reason (separate teams, separate scaling), written down. This becomes the architecture skill in `/tool-up`.
   - **No central files a story must edit**: each module registers its own routes and wiring; migrations are named by timestamp, not by sequence number; configuration is read per module. A file every feature has to touch is a hotspot that forces stories to run one after another.
   - Data model: the core entities, their relationships, and which ones the ten actions touch. Draw it (Mermaid).
   - Buy vs build: auth, payments, email, storage, analytics. Name the service and why.
   - Hosting and deploy: where it runs, how a preview gets built per story, how production deploys.
   - Seams for parallel work: the module map (module → folder → tier → typical story), so `/slice` can give stories in one wave disjoint files. When `scripts/ready.sh` reports a hotspot later, fix it here.
   - Non-functional budgets from the spec: load time, offline, privacy.
   If any decision is genuinely a coin flip with real cost, ask the human. Otherwise decide and record.
3. **Walking skeleton.** Build the thinnest end-to-end slice: the app starts; with `design-timing: early` it shows the real name and logo on the real design tokens and one real screen built to match its approved wireframe (screenshot it with `scripts/look.py` at the same widths and themes as `factory/design/wireframes/approved/` and compare side by side); with `after-core` it uses a plain stock stylesheet in one file that the retrofit stories later replace, and a working name. Then: one real service call, one real data read/write, deploys to a preview URL (or builds an installable / runs in the simulator), and has one passing test at each layer. Nothing more.
4. **Wire the checks**: test runner, linter/formatter, type check, contrast check (once tokens exist), an **import-boundary check** (a module may import another module's interface package only, never its internals; use the stack's own mechanism, such as Go `internal/` packages, plus a lint rule), the spike's measurement as a runnable eval if one exists, a `make check` (or equivalent single command) that runs them all, and CI that runs `make check` on every push. Write `scripts/preview.sh` that produces a preview URL or a runnable build for a branch.
5. Implement the design-system components in the chosen UI technology if `/design-system` has run and left them as specs. With `after-core`, this happens in the retrofit story when the wireframes are approved.
6. Commit to `main` (this is the one time the office commits to main; after `/tool-up`, nobody does). Update `factory/STATE.md`: `blueprint: done`.

## Gate

`make check` passes. The skeleton runs on every target platform in the spec. A preview or build can be produced from a clean clone with one command.

## Rules

- Prefer one language across layers where the platform allows; agents handle fewer context switches better.
- Minimal dependencies. Every dependency is a thing an agent will misuse when its API changes.
- The skeleton is clean, not quick. Every later feature copies its patterns, including how a module registers itself.
- Parallelism comes from module boundaries and file ownership, not from separate services.
