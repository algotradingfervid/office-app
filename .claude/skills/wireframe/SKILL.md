---
name: wireframe
description: "Plan the UX with the human, then draw branded high-fidelity, clickable wireframes for the key screens (real design system, colours, navigation and copy), look at them in a browser, get a blind review, and lock the approved screens as the visual contract every feature is built and proved against. Use after /design-system, when the user says \"wireframes\", \"screens\", \"mockups\", \"user flow\", \"UX\", or invokes /wireframe."
---

# /wireframe — Brand & Design, station 4 (human gate)

High fidelity on purpose: the real brand, tokens, components, navigation and words. The human approves what the product will actually look like, and the line gets an exact target to build to.

## Steps

1. Read `factory/SPEC.md` (the ten actions), `factory/BRAND.md` (voice), `factory/design/DIRECTION.md`, `factory/design/` (tokens CSS, components, icons, `styleguide.html`, platform conventions), `factory/DIALS.md`.
2. **Flows.** For each of the ten actions draw the flow as a Mermaid diagram in `factory/design/flows.md`: entry point → screens → decision points → success state → the two most likely failure states. Note where the user can leave and come back.
3. **Screen map.** `factory/design/screen-map.md`: every screen the flows touch, grouped by area, with a one-line purpose and which actions use it. Merge screens that do the same job.
4. **Layout talk.** Before drawing, agree the layout with the human in a short conversation:
   - the navigation model per platform (bottom tabs, sidebar, top bar) and its destinations
   - what the home screen puts first
   - how phone and desktop differ
   - density
   - one or two apps whose layout they like
   Record the answers in `factory/design/DIRECTION.md` under `## Layout`. Earlier remarks can shape your questions but do not count as answers. If the human is not available, stop and wait; do not draw screens while you wait.
5. **One screen first.** Build the most-used screen as in step 6, run step 7 on it, and show it to the human. Do not draw the other screens until the human has reacted to this one.
6. **Wireframes.** On `wireframe: full`, every screen in the map; on `quick`, the three screens that matter most. Each is `factory/design/wireframes/<screen-id>.html`, a branded, high-fidelity, responsive page:
   - **Linked, not copied.** Link the design system's real token CSS, component CSS and icon files by relative path, so a design-system fix shows on every screen.
   - **System parts only.** Use design-system components and tokens only. If the same arrangement (page header, list row, tab bar) appears on two or more screens, make it a design-system component and use that.
   - **Clickable.** The navigation and every primary action link to the right screen file, so the whole set clicks through as a prototype from `factory/design/wireframes/index.html`. The first screen a new user sees also clicks through to the rest: stand in for emails, codes and sign-in with a plain link.
   - **One page per screen, responsive.** It reads right at phone and desktop width (where the spec has both) and in light and dark. Never draw phone and desktop frames side by side. For native mobile or desktop, draw at device size following `conventions-<platform>.md`.
   - **Real words and realistic data.** Plausible names, amounts and dates from the product's domain; never lorem ipsum or "Item 1".
   - **States.** Empty, loading and error for every screen that has them, reachable from the page (a small state switcher) without editing code.
   - **Annotations off by default.** What each region is and what happens on tap or click, behind a toggle that starts off, so the screen reads like the product.
   - **Prototype controls stay apart.** The state switcher and annotations toggle are prototype controls, not product UI: keep them in one small bar set apart from the screen. The design-system rule does not apply to them.
7. **Look.** `scripts/look.py factory/design/wireframes --out factory/design/wireframes/shots --crawl factory/design/wireframes/index.html`. It fails on console errors, missing assets, sideways scroll, broken links and screens the index cannot reach; fix every one. Then open every screenshot and fix what looks broken, cramped, off-brand or unclear. Re-run until clean.
8. **Blind review.** `dir=$(scripts/blind-copy.sh factory/design wireframes <each CSS, token, icon, font and script file the screens load>)`, then `scripts/look.py --root "$dir" --crawl "$dir/wireframes/index.html" --out "$(mktemp -d)"`: a failed asset means the copy is missing a file, so add it and copy again. Launch a fresh `reviewer-blind` with exactly three things: `$dir`, the entry screen (the first screen a new user lands on, not `index.html`), and the absolute path of `scripts/look.py`. Save its report as `factory/design/wireframes/blind-review-<round>.md`.
   Fix every `STUCK` and `CONFUSING` finding. When a finding is about a term this product's users use daily in their work (for example "GST" for an Indian seller), keep the term and list it for the human with that reason. Put `NOTE`s on a list for the human. Re-run steps 7–8, with a fresh copy and a fresh reviewer each round, until a round has no `STUCK` or `CONFUSING` finding apart from kept terms.
9. **Index.** `factory/design/wireframes/INDEX.md`: a table of action → screens → wireframe files → screenshots. `/slice` copies this into each story.
10. **Present** to the human, screen by screen: the phone and desktop screenshots, the clickable prototype (`index.html`), the blind reviewer's "what I think this is", what you changed because of the review, and any terms you kept. Revise; after each round re-run steps 7–8 on the changed screens.
11. **Approve and lock.** When the human approves, copy the final screenshots to `factory/design/wireframes/approved/`, then set `wireframes-approved: <date>` and `wireframe: done` in `factory/STATE.md`. From here the approved screens and screenshots are the visual contract: `/blueprint`, `/build`, `/prove` and `reviewer-design` build and judge against them, and any visual change comes back through this station, including a design-system change that alters an approved screen (re-run step 7 and refresh `approved/` with the human's OK). With `design-timing: after-core`, approval also releases the retrofit stories that `/tool-up` put in the Usable tier: set the ones whose `needs` are merged from `blocked` to `ready` (the rest to `backlog`) and remove their `blocked:` line.

## Gate

The human agreed the layout (step 4) and the first screen (step 5) before the rest was drawn, then approved the full set. Every one of the ten actions maps to at least one wireframe. `look.py` exits 0 on the final set. The last blind review round has no `STUCK` or `CONFUSING` finding apart from kept terms. `approved/` holds the screenshots of the approved set.

## Rules

- Real words on every button, heading and empty state. UX writing is part of this station.
- High fidelity is not a license to invent style. No colour, size, radius or shadow outside the tokens. If a screen needs a component or token the design system lacks, add it to `/design-system`'s files (spec, CSS, style page) now; do not improvise it here or on the line later.
- Do not add screens the spec does not need. If you think one is missing, ask; the spec changes first.
