---
name: design-system
description: "Agree the interface direction with the human, then turn the brand into buildable design tokens and a component set for the product's platforms, contrast-checked, looked at in a browser, blind-reviewed, with a living style page. Use after /brand, when the user says \"design system\", \"components\", \"tokens\", or invokes /design-system."
---

# /design-system — Brand & Design, station 3

This is what stops every feature from inventing its own button. It must exist as code the line can import, not as a document.

## Steps

1. Read `factory/BRAND.md`, `factory/design/DIRECTION.md`, `factory/SPEC.md` (platforms), `factory/DIALS.md`.
2. **Interface check-in.** Before writing tokens, tell the human in a few lines how you plan to turn the brand into an interface, and ask where they would push differently:
   - density (compact or roomy)
   - corners (sharp, soft, round)
   - depth (borders or shadows)
   - how much the accent colour is used (primary actions only, or also headers and highlights)
   - how dark mode should feel
   - icon style (outline or filled, stroke weight)
   - where the platform's own look wins over the brand
   Record the agreed answers in `factory/design/DIRECTION.md` under `## Interface`. Earlier remarks can shape your questions but do not count as answers. If the human is not available, stop and wait; do not write tokens while you wait.
3. **Tokens.** Write `factory/design/tokens.json`: colour (brand + semantic, light and dark), spacing scale (4-based), radius, type scale, shadows, motion durations. Then generate the platform form of the tokens (CSS variables for web, a Swift/Kotlin constants file or theme object for mobile, the desktop framework's theme file) into `factory/design/tokens/`. Keep one source, generate the rest.
4. **Platform conventions.** For each platform in the spec, write `factory/design/conventions-<platform>.md`: what the platform's own guidelines dictate (navigation patterns, system fonts, touch targets, window chrome, keyboard shortcuts) and where the brand yields to them. Native must feel native.
5. **Components.** List the components the first ten user actions need; on `design-system: quick`, build only those (usually 6–8); on `full`, add the usual set. For each component: states (default, hover, focus, active, disabled, loading, error), sizes, and the copy rules from `voice.md`. Build them in the chosen UI technology only if `/blueprint` has run; otherwise write them as platform-neutral specs in `factory/design/components/<name>.md` and let `/blueprint` implement. (If the human wants to run `/blueprint` first and come back, that is fine; record the order in `STATE.md`.)
6. **Accessibility rules.** Write `factory/design/a11y.md`: minimum contrast, focus visibility, touch target size, motion-reduction, text scaling. These become guardrails in `/tool-up`.
7. **Living style page.** `factory/design/styleguide.html`: every token and component rendered in light and dark, with the exact usage rule under each. It links the real token CSS, component CSS and icon files, so it shows exactly what the line will import. This is what the design reviewer compares against.
8. **Look.** `scripts/look.py factory/design/styleguide.html --out factory/design/shots`. Fix every failure it prints. Then open every screenshot (phone and desktop, light and dark) and look for missing icons, clipped or overlapping text, states that look identical, colours that stray from the brand, and anything that looks unfinished. Fix and re-run until clean.
9. **Blind review.** `dir=$(scripts/blind-copy.sh factory/design styleguide.html <each CSS, token, icon, font and script file the page loads>)`, then `scripts/look.py "$dir/styleguide.html" --root "$dir" --out "$(mktemp -d)"`: a failed asset means the copy is missing a file, so add it and copy again. Launch a fresh `reviewer-blind` with exactly three things: `$dir`, the entry `styleguide.html`, and the absolute path of `scripts/look.py`. On a style page its tasks are telling what is clickable, what each control does, and which state each variant shows. Fix every `STUCK` and `CONFUSING` finding and put `NOTE`s on a list for the human. Re-run steps 8–9 with a fresh copy and a fresh reviewer until a round has no `STUCK` or `CONFUSING` finding.
10. **Show the human** the style page, the screenshots, and a short summary of the blind review and what changed after it. Fold in their reactions before marking done.
11. Update `factory/STATE.md`: `design-system: done`.

## Gate

The human answered the interface check-in at this station, before tokens were written. Every colour pairing used for text passes 4.5:1 (run a contrast check script; do not eyeball). Every component renders in both themes on the style page. `look.py` exits 0 and the screenshots were opened. The last blind review round has no `STUCK` or `CONFUSING` finding.

## Rules

- One source of truth for tokens; everything else is generated.
- A passing script is not a look. Open the screenshots before calling anything done.
- Do not design screens here. That is `/wireframe`.
