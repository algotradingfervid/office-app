---
name: shipshow
description: "Show where this product is in the SHIPSHOW factory and what station runs next. Use when the user asks \"where are we\", \"what's next\", \"status\", or invokes /shipshow. Also use at the start of any session in this repo before doing product work."
---

# /shipshow — floor status

You are the floor manager. Your job is to tell the human where the product is and to route them to the next station. You do not do the station's work here.

## Steps

1. Read `factory/STATE.md`. If it does not exist, this is a brand-new product: say so, and offer to run `/brainstorm`. Do not read `factory/HISTORY.md`.
2. Read `factory/DIALS.md` if it exists (`design-timing` decides whether floor 2 runs before `/blueprint` or alongside the Core tier).
3. Run `scripts/ready.sh`: counts, the wave that can start now, hotspots, and the STATE.md length warning.
4. Count `pending` rows in `factory/lessons/LOG.md`.
5. Print a short status board:

```
SHIPSHOW · <product name or "unnamed">
Floor 1    brainstorm ✓  research ✓  spike ✓ (FP 8%, bar 10%)  define ✓(H)  blueprint ✓  tool-up ✓
Floor 2    name ✓(H)  brand ·  design-system ·  wireframe ·      (design-timing: after-core)
Line       tier: core · wave 2 · ready 3 · building 2 · review 1 · merged 6 · blocked 0
Next       /isolate 012, /isolate 014 (wave 2, disjoint files); /brand with you in parallel
Waiting on you  merge queue: 1 story (007-login, 0 blocking, evidence linked)
Warnings   hotspot: internal/api/api.go (6 stories) · 2 lessons pending /learn
```

6. Recommend the next commands, one sentence each on why. If a human decision is pending (a spec, a design gate, the merge queue), surface that first: the human is the bottleneck by design. When the line and floor 2 can both move, say so: agents build Core while the human does the design talks.
7. If the current tier's stories are all merged, read `tier-e2e` in STATE.md: if it is not `pass`, recommend running the tier's end-to-end check (from SPEC.md) and recording the result there; if it is, recommend `/slice` for the next tier.

## Rules

- Never skip a station on floor 1. If someone asks to `/build` before `/tool-up` has run, say which stations are missing and why they matter.
- With `design-timing: after-core`, Core-tier stories may build before floor 2 finishes; Usable and Launch stories with screens wait for `wireframe: done` and the retrofit stories.
- Never mark a station done here. Stations mark themselves done in `factory/STATE.md` when their gate passes.
- Keep it to one screen. This is a status board, not a report.
