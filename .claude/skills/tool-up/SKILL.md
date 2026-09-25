---
name: tool-up
description: "Generate this product's own kit from everything upstairs produced - house rules, architecture skill, design-system skill, brand-voice skill, wireframe index, glossary, proof playbook, specialist reviewers, guardrails - then trial-run each before opening the line. Use after /blueprint, when the user says \"tool up\", \"generate the skills\", \"open the line\", or invokes /tool-up."
---

# /tool-up — Engineering office, station 2 (the tool-maker)

This is what makes SHIPSHOW self-tooling. Read everything upstairs produced; write the product kit as files the line agents will load. Core skills stay untouched; product skills are prefixed `product-`.

## Inputs

`CLAUDE.md`, `factory/SPEC.md`, `factory/BRAND.md`, `factory/ARCHITECTURE.md`, `factory/design/**`, `factory/RESEARCH.md` (constraints), the skeleton code.

## Steps

0. Write the station marker so the guardrails let this station do its job: `echo tool-up > .claude/station`. Remove it when done: `rm -f .claude/station`.
For each item below, write the file, then read it back as if you were a fresh agent with no other context and fix anything that would confuse you.

1. **House rules → `CLAUDE.md` product section.** Append under the PRODUCT SECTION marker, at most 60 lines: the exact commands (`make check`, run dev, run tests, preview), folder map (one line per top-level folder), naming conventions that differ from the framework's defaults, the three gotchas most likely to bite, and the **risk paths**: the file globs that put a story in `lane: full` (auth, payments, personal data, money, destructive migrations). Nothing derivable from the code. For each line ask: "would removing this cause a mistake?" If not, cut it.

2. **Architecture skill → `.claude/skills/product-architecture/SKILL.md`** (`user-invocable: false`, description says "load when writing or reviewing application code"). The layer rule, where each kind of code goes, a worked example of one feature cut across the layers, and an anti-pattern list from `ARCHITECTURE.md`. Under 300 lines; link to `factory/ARCHITECTURE.md` for the rest.

If the design rooms have not run yet (`design-timing: after-core`), skip items 3–5 now; the retrofit stories write them when the wireframes are approved.

3. **Design-system skill → `.claude/skills/product-design-system/SKILL.md`** (`user-invocable: false`, description: "load when building or reviewing any screen or UI"). How to import tokens and components, the "use this, never hand-roll that" table, the a11y rules, and the path to `styleguide.html`.

4. **Brand-voice skill → `.claude/skills/product-voice/SKILL.md`** (`user-invocable: false`, description: "load when writing any user-facing string, error, email or doc"). The voice guide and the five canonical strings as examples.

5. **Wireframe index → `.claude/skills/product-screens/SKILL.md`** (`user-invocable: false`). The action → screen → wireframe → approved screenshot table, so a story can find its target screen and exactly how it must look.

6. **Glossary → `factory/GLOSSARY.md`** and a two-line pointer in the CLAUDE.md product section. Every domain noun in the spec with one meaning and the code identifier for it.

7. **Proof playbook → `.claude/skills/product-proof/SKILL.md`** (`user-invocable: false`, loaded by `/prove`). Exactly how evidence is captured for this product's platforms: the command that starts the app, the command that records/screenshots at the same widths and themes as `factory/design/wireframes/approved/` (before it exists: 390 px and 1280 px, light theme) (Playwright for web, reusing `scripts/look.py` where it fits; simulator recording for mobile, window capture for desktop, test output plus timing for services, transcript for CLI), where files go (`factory/evidence/<story-id>/`), and the before/after comparison format. Include a working script in `scripts/prove.sh` that takes a story id and produces `before/`, `after/` and `report.md` (with `commit: <branch head sha>`, and raw captures under the git-ignored `raw/`), and **one journey runner** that reads the steps from a story's Check section, so stories never add their own journey scripts. Screenshots only for stories with a screen.

7a. **Mutation → `scripts/mutate.sh <story-id>`.** Runs the stack's mutation-testing tool (choose a current, maintained one and check its latest version) on the files the story changed, and prints surviving mutants with file:line. Used by `/prove` and read by `reviewer-code`, so reviewers never improvise mutation by hand.

8. **Specialist reviewers → `.claude/agents/product-review-*.md`.** Generate only what the product needs, using `.claude/agents/reviewer-design.md` as the model: design-fidelity reviewer (anything with screens), security reviewer (accounts, payments, personal data), performance reviewer (data-heavy or latency budget in spec), platform-guidelines reviewer (mobile or desktop), accessibility reviewer (if a11y is a non-negotiable). Each gets the product context it needs inline, tags every finding `BLOCKING` or `NOTE` (as `reviewer-code` does), starts with `BLOCKING: n` and `SCORE: n/5`, and names the file globs that trigger it.

9. **Guardrails → `.claude/settings.json` hooks + `.claude/hooks/`.** Add: a `PostToolUse` formatter on Edit/Write for the product's languages; a `PreToolUse` block on editing protected paths (secrets, migrations, `factory/SPEC.md` outside `/define`, generated token files); a `Stop` hook that refuses to end a `/build` session while `make check` fails. Keep hooks small and deterministic. Extend `.claude/hooks/guard-protected-files.sh` rather than replacing it.

10. **Backlog seed.** Run `/slice` on the **Core tier only**, with contract stories first and claims derived from the `ARCHITECTURE.md` module map. With `design-timing: after-core`, also add the Usable-tier **retrofit stories** with `status: blocked` and `blocked: waiting for wireframe approval`: first a contract story that swaps the stock stylesheet for the design system and writes items 3–5 above, then one story per module (≤ 5 files each) that restyles its Core screens to their approved wireframes. `/wireframe` releases them. `scripts/check-stories.sh` must pass.

11. **Trial run.** Pick the smallest Core story. In a fresh subagent (`context: fork`), run `/isolate` → `/build` → `/prove` on it with the new kit, and watch for every point where the agent hesitated, guessed, or loaded the wrong file. Fix the kit, not the story. Repeat once. Then discard the trial branch.

12. Update `factory/STATE.md`: `tool-up: done`, `line-open: <date>`. Print the list of generated files and tell the human the line is open: `/shipshow` shows ready stories.

## Gate

Every generated skill was loaded by a fresh agent in the trial run without confusion. `make check` passes on `main`. CLAUDE.md total length under 200 lines.

## Rules

- Skills describe what to do, not why the product exists; the why is in `factory/`.
- Descriptions must be a little pushy about when to load; agents under-trigger skills.
- If you write ALWAYS or NEVER in capitals more than twice, you are writing a hook, not a skill.
- Never touch the core skills. If a core skill needs a product-specific twist, the product skill says so and the core skill's steps say "check for `product-*` skills".
