---
name: define
description: "Interview the product owner until the product is concrete, then write the spec they approve (value proposition, users, first ten user actions grouped into build tiers, out-of-scope, platforms, metrics, acceptance checks). Use after /research (and /spike when one ran), or when the user says \"write the spec\", \"PRD\", \"define the product\", or invokes /define."
---

# /define — Studio, station 4 (human gate)

The spec is written for a stranger: a fresh agent with no memory of this conversation must be able to build from it.

## Steps

0. Write the station marker so the guardrails let this station do its job: `echo define > .claude/station`. Remove it when done: `rm -f .claude/station`.
1. Read `factory/BRIEF.md`, `factory/RESEARCH.md`, `factory/SPIKE.md` (if a spike ran: its verdict and numbers are facts, not options), `factory/DIALS.md`.
2. **Interview** in rounds of at most three questions, until nothing important is vague. Cover, in this order:
   - The value proposition in one sentence a user would say to a friend.
   - The primary user and the one secondary user, if any.
   - The first ten things a user does, in order, from opening the product to getting value. This becomes the backlog.
   - What version one deliberately does not do. Be ruthless; write the "not yet" list.
   - Platforms: web, iOS, Android, macOS, Windows, Linux, API, CLI. Which first.
   - Success metrics: two or three numbers that mean it is working, with a target and a date.
   - Non-negotiables: offline, privacy, accessibility, languages, budget.
3. For each of the ten actions write an **acceptance check**: a sentence a non-programmer could verify by looking ("a new user can sign up with email and see an empty dashboard within 60 seconds").
3a. **Tier the actions**, and agree the tiers with the human:
   - **Core** — the smallest loop that proves the product's value (for a matching product: fetch → read → match → ranked list, measured on the spike's data). Usually two or three actions. Its end-to-end check is a number where the spike produced one.
   - **Usable** — what one real user needs to use Core unaided (sign-in, one profile, the main screens).
   - **Launch** — what paying strangers need (billing, alerts, legal pages, onboarding polish).
   - **Later** — everything else.
   Each tier gets one end-to-end check. The line builds one tier at a time.
4. **Write** `factory/SPEC.md` using `factory/templates/spec.md`.
5. **Cross-check** the spec against the research: every top-5 pain is either addressed by an action or explicitly deferred with a reason. List any that fell through.
6. Present the spec and ask for approval. Apply edits. When approved, set `spec-approved: <date>` in `factory/STATE.md` and `define: done`.

## Gate

The human has approved the spec, including its tiers, in writing. Do not proceed to the next station on a draft (`/blueprint` when `design-timing: after-core`, `/name-it` when `early`).

## Rules

- Ten actions, not thirty. If the human wants thirty, the extra twenty go in "later".
- Core is small on purpose. If an action is not needed to prove the value, it is not Core.
- Every acceptance check is observable. "Works well" is not a check; "loads in under two seconds on a phone" is.
- Do not choose a stack here. That is `/blueprint`.
