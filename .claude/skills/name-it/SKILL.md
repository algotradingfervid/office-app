---
name: name-it
description: "Generate product name candidates, run availability and risk checks, and give the human a shortlist to pick from. Use after /define (design-timing early) or alongside the Core tier after /tool-up (after-core), when the user says \"name the product\", \"what should we call it\", or invokes /name-it."
---

# /name-it — Brand & Design, station 1 (human gate)

## Steps

1. Read `factory/SPEC.md`, `factory/RESEARCH.md` (competitor names), `factory/DIALS.md`.
2. **Generate.** On `name-it: full`, produce 40 candidates across five styles, eight each: descriptive (says what it does), invented (coined word), metaphor (borrowed object or idea), compound (two real words), personal (a name or nickname). On `quick`, 10 candidates, any style.
3. **Cull** to 12 by: easy to say aloud, easy to spell after hearing it once, not confusable with a competitor, no unfortunate meaning in English, Spanish, Hindi, Mandarin, Arabic, French, German, Portuguese (check with web search when unsure).
4. **Check** each survivor and record results in a table:
   - Web domain: `.com` and one sensible alternative (search "<name>.com" and try fetching; unavailable is a fact, available is a guess until registered).
   - App stores (if the spec has mobile): search the App Store and Google Play for the exact name.
   - Trademark-style: web-search "<name> trademark" and "<name> software"; note obvious clashes in the same category. This is a smell test, not legal advice; say so.
   - Social handles: note whether the obvious handle looks taken on two platforms.
   On `quick`, do the domain check only.
5. **Shortlist** five (three on quick). For each: the name, one tagline, one sentence on why it fits the spec, and the check results.
6. Present the shortlist and ask the human to pick or send you back with a direction. Do not pick for them.
7. **Write** `factory/NAME.md`: chosen name, tagline, the shortlist and checks for the record. Update `factory/STATE.md`: `name-it: done`, `product-name: <name>`.

## Gate

The human picked. Record their pick verbatim.

## Rules

- Never claim a domain or trademark is definitely free. Say "appeared available on <date>".
- Do not write the logo here. That is `/brand`.
