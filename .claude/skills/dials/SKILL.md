---
name: dials
description: "Set or change how much ceremony this product gets (quick vs full) for research, naming, brand, design system and wireframes, and whether design runs before or alongside the Core tier. Use when the user says \"this is a weekend project\", \"go quick\", \"do this properly\", or invokes /dials."
argument-hint: "[room=quick|full ...] e.g. brand=quick wireframes=full"
---

# /dials — size the ceremony to the product

Every upstairs room always runs, but each has a `quick` and a `full` setting. Quick still produces the artifact; it produces a smaller one in less time and burns fewer tokens.

| Room | quick | full |
|---|---|---|
| research | one researcher, top 3 competitors, 30-minute budget | parallel researchers: market, competitors, user pains, constraints, feasibility |
| name-it | working name in one round, domain check only | long list, five styles, domain + app store + trademark-style + language checks |
| brand | one direction: one mark, one palette, one type pairing | three directions, refined winner, brand one-pager, full logo set |
| design-system | tokens + the 6 components the wireframes need | tokens, full component set, accessibility rules, living style page |
| wireframe | the 3 screens that matter, branded hi-fi and clickable | all key screens, branded hi-fi and clickable, flows, screen map, UX copy |

One more dial says **when** the design rooms run:

| Dial | after-core (default) | early |
|---|---|---|
| design-timing | `/name-it`, `/brand`, `/design-system`, `/wireframe` run while the line builds the Core tier on stock styling; retrofit stories (a stylesheet swap, then one per module) restyle Core screens once the wireframes are approved | the design rooms run before `/blueprint`, as in kit 0.2; choose this when the look itself is the riskiest assumption |

Ceremony on the line is set per story, not here: each story's `lane` (core or full) is chosen by `/slice`.

## Steps

1. If `$ARGUMENTS` is empty, ask the human one question: "Is this a throwaway experiment, a personal tool, or something you will put in front of other people?" Map: throwaway → all quick; personal tool → research quick, brand quick, others full; for other people → all full. `design-timing` is `after-core` unless the human says the look is what could make or break the product. Confirm the mapping, let them override per room.
2. Write `factory/DIALS.md`:

```markdown
# Dials
product-kind: personal-tool
research: quick
name-it: quick
brand: quick
design-system: full
wireframe: full
design-timing: after-core
set-on: 2026-09-16
reason: "one-line reason the human gave"
```

3. Say which rooms are quick and remind the human they can turn any dial up later without redoing the others.
