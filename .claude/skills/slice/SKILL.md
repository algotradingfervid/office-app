---
name: slice
description: "Cut one tier of the spec (or a new request) into small stories planned for parallel waves: each one outcome, one check, a disjoint file claim, a lane and prerequisites, with contract stories first. Use when the user says \"add a feature\", \"break this down\", \"backlog\", \"next tier\", or invokes /slice."
argument-hint: "[tier name, or free-text feature request; empty = the next unsliced tier]"
---

# /slice — the backlog

The backlog is a graph planned for parallel work. Each story claims the files it will touch; stories in the same wave claim disjoint files, so agents never build the same thing twice.

## Steps

0. Write the station marker so the guardrails let this station do its job: `echo slice > .claude/station`. Remove it when done: `rm -f .claude/station`.
1. Read `factory/SPEC.md` (tiers), `factory/ARCHITECTURE.md` (modules and seams), `factory/design/wireframes/INDEX.md` if it exists, the risk paths in the CLAUDE.md product section, and existing `factory/stories/`.
2. **Pick the scope.** If `$ARGUMENTS` names a tier, slice that tier. If it is empty, slice the lowest tier not yet sliced, and only if `tier-e2e` in STATE.md says the tier before it passed (say so and stop otherwise). Then set `tier:` to the new tier and `tier-e2e: pending`. A free-text request that is not in the spec: say so and ask whether to add it to `SPEC.md` first (it changes the spec, so the human decides).
3. **Contracts first.** For every interface, type, table or fake that two or more stories will use, write a `kind: contract` story that delivers only that. The stories that use it depend on the contract, not on each other.
4. **Cut the rest** into stories:
   - One outcome a non-programmer could verify, one check.
   - **≤ 5 non-test files and ≤ ~300 changed lines**; one fresh agent builds and proves it in one sitting. Over the limit: split along a file boundary.
   - Claim exact files or narrow globs inside one module. Never claim a whole module or a central file (route table, main wiring, a shared migration sequence); if a story seems to need one, the architecture has a hotspot: note it for `/blueprint` and use the module's own registration.
   - `lane: full` when the claim touches a risk path in the CLAUDE.md product section (auth, payments, personal data, money, destructive migrations); `core` otherwise. A throwaway measurement is not a story; run `/spike`.
   - A story that needs a component the design system lacks is two stories: the component, then the feature.
   - A story that follows another on the same files `needs` it. It never gets a fallback ("if 015 has not merged, you own feed.go").
5. Write each as `factory/stories/NNN-<slug>.md` from `factory/templates/story.md`, under 60 lines: what must be true, not how to code it.
6. Set `status: ready` on every story whose `needs` are all merged (or empty).
7. Run `scripts/check-stories.sh`. Fix every failure by splitting or re-claiming; do not raise the limits.
8. Run `scripts/ready.sh` and print its waves and hotspot report, plus a Mermaid diagram of the new stories and their `needs`.

## Gate

`scripts/check-stories.sh` passes. The first wave has at least two stories unless the tier is genuinely one story. Every story has a tier, a lane, a check, a screen (or "none: no UI"), an exact file claim and prerequisites.

## Rules

- Small beats clever. Three small stories that run side by side finish before one big one.
- The builder designs inside the claim. A story that dictates function signatures has become a spec of the code; cut it back to the outcome.
