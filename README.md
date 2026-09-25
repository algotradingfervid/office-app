# SHIPSHOW

**Turn your shitshow into a ship show.**

A self-tooling app factory for Claude Code. Give it an idea; it brainstorms, researches, proves the riskiest part with a throwaway spike and a number, writes a tiered spec you approve, lays down a walking skeleton and its own product-specific tooling, and builds the Core tier first in small stories that run in parallel waves. While Core is built, it talks through the look with you, names and brands the product, and draws the design system and branded hi-fi wireframes you approve; then it builds the rest tier by tier. You decide five times: spec, name, brand, wireframes, merge.

Works for web, mobile, desktop, services and CLI tools. The stations never change; `/tool-up` cuts the tools to fit the product.

## Install

```bash
# 1. Copy this folder to become your new product's repo
cp -r shipshow my-product && cd my-product && git init && git add . && git commit -m "SHIPSHOW kit"
git branch -M main   # the kit assumes `main`

# 2. Open Claude Code here and say hello to the floor manager
claude
> /shipshow
```

`/shipshow` reads `factory/STATE.md` and tells you the next station. That is the only command you need to remember.

## The stations

| Room | Command | Produces | Gate |
|---|---|---|---|
| Studio | `/brainstorm` | `factory/BRIEF.md`, `factory/DIALS.md` | one direction, one riskiest assumption |
| Studio | `/research` | `factory/RESEARCH.md` + `factory/research/*` | every claim sourced; assumption verdict |
| Studio | `/spike` | `factory/SPIKE.md`, `spikes/*` (throwaway) | a number against an agreed bar; go / change / stop |
| Studio | `/define` | `factory/SPEC.md` with Core / Usable / Launch / Later tiers | **you approve** |
| Brand & Design | `/name-it` | `factory/NAME.md` | **you pick** |
| Brand & Design | `/brand` | `factory/design/DIRECTION.md`, `factory/BRAND.md`, `factory/design/brand/*` | you OK the direction, then **you pick** |
| Brand & Design | `/design-system` | `factory/design/tokens.json`, components, `styleguide.html` | contrast passes, both themes, looked at, blind review |
| Brand & Design | `/wireframe` | `factory/design/flows.md`, branded hi-fi clickable `wireframes/*`, `wireframes/approved/` | looked at, blind review, **you approve** |
| Engineering | `/blueprint` | `factory/ARCHITECTURE.md`, the walking skeleton, `make check` | skeleton runs and deploys |
| Engineering | `/tool-up` | `.claude/skills/product-*`, `product-review-*` agents, hooks, backlog | trial run passes |
| Line | `/slice` → `/isolate N` → `/build N` → `/prove N` → `/ship N` | small stories in parallel waves, merged with evidence | no blocking findings (≤ 2 rounds), **you merge per wave** |
| Return | `/learn` (per wave) | updated kit | rules stay short; repeats become hooks |
| Maintenance | `/shipshow-update` | proposal + PR upgrading the core kit | trial passes, you merge |
| Anytime | `/dials` | `factory/DIALS.md` | — |

Order: floor 1 (brainstorm → research → spike → define → blueprint → tool-up → Core tier), then floor 2 (name-it → brand → design-system → wireframe) runs alongside the Core tier by default (`design-timing: after-core`), and the line builds Usable and Launch against the approved design. Set `design-timing: early` when the look is the risk.

## Running features in parallel

```bash
scripts/ready.sh                 # waves that can start now, hotspots, STATE length
scripts/check-stories.sh         # the /slice gate: size, exact claims, no collisions
scripts/ready.sh --stats         # median hours per station, review rounds per lane
claude -w story-007              # open a fresh session in a story's worktree
> /build 007
```

Each story is small (one outcome, ≤ 5 non-test files) and claims its files; stories in a wave claim disjoint files, and contract stories come first so the stories that use them run side by side. The architecture is a modular monolith whose modules register themselves, so no story has to edit a shared route table or migration sequence. Each story gets its own git worktree and a clean context, and a hook keeps it there. Ceremony follows the story's lane: `core` gets one code reviewer, `full` (auth, money, personal data) gets every specialist. Your merge-review rate is the factory's true capacity, so review the evidence once per wave, not the diff.

## Self-improvement

- `/learn` runs after every wave: twice wrong becomes a rule, wrong again after the rule becomes a hook or script, thrice repeated becomes a skill, improvised UI goes back into the design system, and rules that stop earning their place are cut.
- `/shipshow-update` researches new models, Claude Code features and current practice, proposes changes to the core kit with sources, trial-runs them, and opens a PR. Run it after major model releases or monthly (wire it to a scheduled task or `/loop` if you like).
- `factory/MODELS.md` says which model tier does which job and when to change it.

## Layout

```
CLAUDE.md                core house rules (+ product section appended by /tool-up)
.claude/skills/          core stations; product-* skills are generated
.claude/agents/          researcher, reviewer-code, reviewer-design, reviewer-blind, verifier (+ product-review-*)
.claude/hooks/           guardrails (protected files, main branch, story worktree, formatting, stop-check)
.claude/settings.json    permissions + hook wiring
factory/                 the documents: STATE (one screen), HISTORY, DIALS, BRIEF, RESEARCH, SPIKE, SPEC, NAME, BRAND, ARCHITECTURE, GLOSSARY, MODELS, KIT-VERSION
factory/design/          tokens, components, flows, wireframes, styleguide
factory/stories/         one file per story (status, needs, files, screen, check)
factory/evidence/        before/after proof per story
factory/lessons/         what went wrong, feeding /learn
scripts/                 new-story, finish-story, ready, check-stories, stories.py, look.py, blind-copy (+ preview, prove and mutate, written by /blueprint and /tool-up)
```

## Provenance

Built from Ras Mic's software factory (isolate / build / prove / ship, Greptile 5/5 loop), BMAD's analyst→PM→UX→architect ordering, Superpowers' tests-first and fresh-subagent-per-task, GSD's parallel research and dependency-aware backlog, Anthropic's Claude Code guidance (short CLAUDE.md, skills with progressive disclosure, hooks as enforcement, writer/reviewer split, evidence over claims), and Karpathy's agentic-engineering notes (small reviewable chunks, success criteria over instructions, surface assumptions, simplicity is the human's job).

No warranty, no unions, no ellipses.
