# Kit update 2026-09-16 → v0.2.0: design stations talk first, look, and get a blind review

Source: the first product run on the kit (SahiBid, `TenderFinder/shipshow/factory/lessons/LOG.md`, four rows dated 2026-09-16), plus the human's written agreement in that session to change the /brand and /wireframe gates.

## What went wrong on v0.1.0

1. `/brand` generated directions from the spec and the agent's own reading of it. The colour scheme took two rework rounds ("needs to be different", then "too many colours").
2. Nothing in the kit looks at rendered design output. The gates are scripts (contrast, spec↔CSS). The agent improvised Playwright screenshots, and the human still had to screenshot pages to point at problems.
3. The design-system critic was briefed with the product and spec, so it judged against intent, not against what a newcomer sees.
4. Grey-box wireframes gave the human no real picture of the output, and the line had no visual target beyond structure.

## Changes

| File | Change | Reason | Risk |
|---|---|---|---|
| `scripts/look.py` (new) | Serve pages over local HTTP, screenshot every page at phone and desktop widths in light and dark, and fail on console errors, missing assets, horizontal scroll, broken local links and pages the entry can't reach | Lesson 2; a deterministic helper instead of improvised Playwright each time | low |
| `scripts/blind-copy.sh` (new) | Copy only renderable files (html, css, js, images, fonts) into a temp folder outside the repo | Lesson 3; the blind reviewer must not be able to read the documents | low |
| `.claude/agents/reviewer-blind.md` (new) | First-time-user reviewer: sees only rendered pages, maps the navigation, walks the obvious tasks, reports STUCK / CONFUSING / NOTE, scores SIMPLICITY 0–5 | Lesson 3 | medium: blindness rests on the copied folder plus an instruction, not a hook |
| `/brand` | Direction talk (feel, references, colour appetite, logo approach, type), played back as `factory/design/DIRECTION.md` and OK'd before any direction is drawn; look.py on previews; talk through palette and mark when presenting | Lesson 1 | **high: changes a human gate** (agreed by the human) |
| `/design-system` | Interface check-in before tokens (density, corners, depth, accent use, icons, platform yield) recorded in DIRECTION.md; look.py on the style page; blind review; show the human | Lessons 1–3 | medium |
| `/wireframe` | Branded high-fidelity, clickable wireframes linked to the real design-system CSS and icons; layout talk first; one screen agreed before the rest; look.py with link crawl; blind review; approved screenshots locked in `wireframes/approved/` as the visual contract | Lessons 2–4 | **high: changes a human gate** (agreed by the human) |
| `/dials` | The wireframe row says hi-fi and clickable on both settings | Consistency with /wireframe | low |
| `/blueprint`, `/build`, `/prove`, `/ship`, `/tool-up`, `reviewer-design` | Build and judge against the approved hi-fi screenshots. Styling differences from them are GAPs; differences from real data are NOTES | Lesson 4: "once approved, the same design language is incorporated" | medium |
| `CLAUDE.md` | Two house rules: talk through design direction before drawing; look at visual work with look.py before showing it | Lessons 1–2 apply beyond one station | low |
| `README.md`, `factory/MODELS.md`, `factory/KIT-VERSION.md`, `factory/updates/CHANGELOG.md` | Station table, agent list, model row for reviewer-blind, version 0.2.0 | Bookkeeping | low |

## Not changed

- The five human gates stay five: spec, name, brand, wireframes, merge. The design-system check-in is a conversation, not a gate.
- No new hook. Blindness is enforced by the copied folder, not by a guard.

## Verification

- `look.py` and `blind-copy.sh` run against SahiBid's real style page and wireframes (read-only, output to a scratch folder).
- Fresh subagents given the new /brand and /wireframe skills are put under pressure ("skip the chat, you know the spec"; "brief the reviewer") and checked for compliance.
- Grep for leftover "low fidelity" or "grey box" wording; frontmatter parses; no skill over 500 lines.
