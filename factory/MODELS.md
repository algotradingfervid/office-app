# MODELS — which model does which job

Quality first, then time, then cost. The reasoning: on the line, rework is the expensive thing. A weaker model that produces a 3/5 PR costs two extra prove→build→ship loops, and those loops burn more tokens and more of the human's attention than the stronger model would have cost up front. Haiku pays off only where the work is high-volume, mechanical, and cheaply checked.

## The four criteria

Ask these about any task before assigning a model:

1. **Blast radius** — how many later things copy this? Spec, architecture, design system, skeleton, and the product kit are copied by every story; a mistake there is multiplied. → strongest model, high effort.
2. **Reversibility and ambiguity** — is the task a judgment call with no cheap check (taste, tradeoffs, naming, what to leave out)? → strongest model. Is it mechanical with an automatic check that catches errors (formatting, capturing screenshots, counting statuses)? → smallest model that passes the check.
3. **Verifiability** — if a deterministic check exists (tests, contrast script, a schema), Haiku in a loop may do for mechanical work. If the only check is a human's eye, use Opus.
4. **Volume** — research fan-out, summarising many pages, classifying lessons: high volume, low stakes per unit. → Opus at medium effort, in parallel; Haiku only when an automatic check catches its mistakes.

One extra rule: **the reviewer is never weaker than the writer.** A weaker reviewer rubber-stamps.

## Default assignment

Two tiers: **Opus** for anything with judgement, code or taste, **Haiku** for mechanical steps with an automatic check. No Sonnet, no Fable (the founder's rule, 2026-09-19, made the kit default on 2026-09-23). Omitting `model` inherits the session's Opus.

| Station / agent | Model | Effort | Why |
|---|---|---|---|
| /brainstorm, /spike, /define, /blueprint, /tool-up | Opus | high | judgement, blast radius, human-only checks |
| /research lanes and synthesis | Opus | medium (lanes), high (verdict) | sources are the check; the verdict is judgement |
| /name-it generation, checks and shortlist | Opus | medium | taste; checks are lookups |
| /brand directions, /design-system, /wireframe | Opus | high | taste, blast radius, human-only check |
| Token generation, icon export, contrast script, `scripts/look.py` screenshots | Haiku or a script | low | deterministic |
| /slice | Opus | medium | the wave plan shapes parallelism for the whole tier |
| /isolate, /shipshow status, scripts/ready.sh, check-stories.sh | Haiku or a script | low | mechanical |
| /build | Opus | medium–high | code quality dominates cost |
| /prove capture, mutate.sh | Haiku or a script | low | mechanical; the playbook is deterministic |
| /prove verdict | Opus | medium | visual judgement against a wireframe |
| reviewer-code, reviewer-design, reviewer-blind, product-review-* | Opus | high | never weaker than the writer |
| hook prompts, yes/no evaluators | Haiku | low | single judgements, many times |
| /learn classification, rule writing and trim | Opus | medium–high | every rule is copied into every future story |
| /shipshow-update research and proposal | Opus | high | changes the factory itself |

## How this is wired

- Subagents set `model:` in their frontmatter (`researcher` → opus; reviewers → `inherit`, so they run at the session's Opus).
- Skills run in the session's model. Start every session with Opus. Workflow scripts set `model: 'opus'` or `'haiku'` on each agent, never another model.
- `/shipshow-update` revisits this table when models change.

## Signals to change an assignment

- A mechanical station passes its automatic check first time for ten stories in a row on Opus → try Haiku on it and watch the loop count.
- Loop counts rise after a downgrade → go back up; the tokens saved were spent on rework.
- The human starts finding gaps the reviewers missed → raise the reviewer tier first, not the writer's.
