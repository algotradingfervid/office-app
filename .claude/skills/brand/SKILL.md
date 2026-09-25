---
name: brand
description: "Agree the design direction with the human (feel, colour, logo) and only then create brand directions (logo mark, palette, type pairing, voice), let the human pick one, and write the brand one-pager and logo set. Use after /name-it, when the user says \"logo\", \"brand\", \"colours\", \"look and feel\", or invokes /brand."
---

# /brand — Brand & Design, station 2 (human gate)

You produce real files, not descriptions: SVG logos, a palette with hex values, named fonts, and a voice guide with example sentences. The look comes out of a conversation with the human. The spec says what the product does, not how it should look and feel.

## Steps

0. Write the station marker so the guardrails let this station do its job: `echo brand > .claude/station`. Remove it when done: `rm -f .claude/station`.
1. Read `factory/SPEC.md`, `factory/NAME.md`, `factory/RESEARCH.md` (what competitors look like: do not look like them), `factory/DIALS.md`.
2. **Direction talk.** Before drawing anything, talk with the human: a few questions at a time, a conversation rather than a questionnaire. Cover:
   - **Feel**: three words for how the product should feel, and one thing it must not feel like.
   - **References**: brands, apps, places or objects whose look they like, and some they dislike.
   - **Colour**: how many colours (one accent, or a small family), how loud (muted, balanced, bold), colours that mean something to these users, colours to avoid.
   - **Logo**: symbol, lettermark or wordmark; motifs from the product's world; what to stay away from.
   - **Type**: modern or traditional; scripts it must render (for example Devanagari).
   You may offer your own reading of the spec as a starting point, labelled as your suggestion. When words are not enough, make `factory/design/brand/moodboard.html` with 4–6 palette strips and 3–4 mark styles as rough sketches to react to (sketches, not directions).
   Write what you heard to `factory/design/DIRECTION.md` under `## Brand` (feel, references, colour appetite, logo approach, type, must-avoid, open questions). Show it to the human and get an explicit OK before step 3. Earlier remarks can shape your questions but do not count as answers. If the human is not available, stop and wait; do not draw anything while you wait.
   If the human asks to skip the talk, ask only the colour and logo questions. Their answers count as the OK: record them, mark every other field `open: directions may vary this`, and continue.
3. **Directions.** On `brand: full`, produce three directions; on `quick`, one. All of them stay inside the agreed `## Brand` section and differ from each other in real ways. Each direction is a folder `factory/design/brand/<direction-slug>/` containing:
   - `mark.svg`: a logo mark that works at 16px and 512px, single colour and full colour, light and dark background versions (`mark-dark.svg`). Simple geometry; no gradients as a crutch; no clip art; no text effects. If the mark needs the name, make a `wordmark.svg` too.
   - `palette.md`: 1 accent, 2 neutrals (light ground, dark ground), 1 support colour, semantic colours (success, warning, danger), each with hex and a one-line reason. Keep to the agreed colour count and loudness. Verify the accent passes 4.5:1 contrast for text on both grounds, or specify a text-safe variant.
   - `type.md`: a display face and a body face (Google Fonts or system fonts, with fallback stacks) and a type scale.
   - `voice.md`: how the product speaks: three adjectives, three "we say / we don't say" pairs, and five real strings: a welcome, an empty state, an error, a success toast, a destructive-action confirm.
   - `preview.html`: one self-contained page showing the mark (including at 16px), wordmark, palette swatches, type samples and the five strings, in light and dark. Its first line says which parts of `DIRECTION.md` this direction answers and how it differs from the others. This is what the human judges.
4. **Look before showing.** `scripts/look.py factory/design/brand/*/preview.html --out factory/design/brand/shots`. Fix every failure it prints. Then open every screenshot and fix what is clipped, illegible, broken, or outside the agreed direction. Re-run until clean.
5. **Present and talk it through.** Open the previews and screenshots. For each direction, walk the human through the palette (what the accent is used for, how many colours a typical screen shows) and the mark (what it depicts, how it reads at 16px), and ask what to keep, change or combine. Do not pick for them. Revise, re-run step 4, and repeat until they pick.
6. **Refine** the winner: tighten the mark, export `icon-512.png`, `icon-192.png`, `favicon.svg`, and platform icons if the spec has mobile or desktop (use a script; keep it in `factory/design/brand/export.sh`).
7. **Write** `factory/BRAND.md` (the one-pager): name, tagline, mark usage (clear space, minimum size, don'ts), palette, type, voice, and file paths. One page. This is what `/design-system` and `/tool-up` read.
8. Update `factory/STATE.md`: `brand: done`.

## Gate

The human OK'd the `## Brand` section of `DIRECTION.md` at this station, before any direction was drawn. `look.py` exits 0 on the previews. The human picked a direction. `BRAND.md` exists and every file path in it resolves.

## Rules

- Avoid the defaults every AI reaches for: purple-to-blue gradients, a lone neon accent on black, Inter or Space Grotesk as the "safe" face, rounded-everything. Pick things this product's users would recognise as theirs.
- Never copy a known logo or a competitor's shape. Original marks only.
- Taste is the human's. Ask before you assume, offer real differences, and let them choose.
- When new feedback contradicts an earlier answer ("colourful", then "too many colours"), ask which one wins and update `DIRECTION.md` before redrawing.
