---
name: learn
description: "After a wave of stories merges, turn what went wrong into rules, repeated procedures into skills, and improvised UI into design-system components, then trim anything that stopped earning its place. Use when the line drains or five stories have merged since the last /learn, when the user says \"retro\", \"what did we learn\", or invokes /learn [story-ids]."
argument-hint: "[story-ids of the wave]"
---

# /learn — the return path

The factory compounds only if lessons flow back upstairs. Twice wrong becomes a rule. Wrong again after the rule becomes a script or a hook. Thrice repeated becomes a skill. A rule that no longer prevents a mistake gets cut. Run once per wave: when the line drains (nothing building, proving or in review) or five stories have merged since the last run, over every story merged since the last LOG.md entry.

## Steps

0. Write the station marker so the guardrails let this station do its job: `echo learn > .claude/station`. Remove it when done: `rm -f .claude/station`.
1. Gather, for every story merged since the last LOG.md entry: `factory/lessons/<id>.md` (if any), the PR review threads (`gh pr view <pr> --comments`), the prove→build loop count and `review-rounds`, and the diff footprint versus the story's claimed `files`. Also read every `pending` row in `factory/lessons/LOG.md`: this run clears them.
2. Classify each lesson:
   - **Rule** — the agent did something wrong that a one-line instruction would prevent, and it has now happened in two or more stories (check earlier lessons). → Add one verifiable line to the CLAUDE.md product section or the relevant `product-*` skill. Then delete a line that has stopped earning its place, to keep the file the same length.
   - **Skill** — the same multi-step procedure was worked out by hand in three or more stories. → Write `.claude/skills/product-<name>/SKILL.md` following the tool-up conventions, with the procedure as steps and the earlier stories as examples.
   - **Guardrail** — something happened that must never happen (edited a protected file, skipped a check), **or a lesson came back after its rule was written**. → A hook (`.claude/hooks/`, `.claude/settings.json`) or a check in a script (`prove.sh`, the merge step, `check-stories.sh`). Instructions are requests; hooks and scripts are enforcement. Never write a second rule for the same mistake.
   - **Component** — the story hand-rolled UI the design system lacked. → Add the component spec and implementation to the design system and the style page; open a story if it needs its own work.
   - **Footprint** — the story touched files outside its claim, or a hotspot shows in `scripts/ready.sh`. → Fix the module map or registration in `ARCHITECTURE.md` (and open a story to remove the hotspot), so future waves stay disjoint.
   - **Core kit** — the lesson is about how a station works, not about this product. → Write it to LOG.md as `core kit → /shipshow-update`; do not edit core files here.
   - **Noise** — a one-off. → Leave it in the lessons file. Do not add a rule for a thing that happened once.
3. Write `factory/lessons/LOG.md` entry: date, story, lessons, what changed in the kit.
4. Once every ten merged stories, do a **trim pass**: read the CLAUDE.md product section and every `product-*` skill and remove anything the current model no longer needs (things it does right without being told). Note what was removed in the log.
5. Update `factory/STATE.md` counts (keep it one screen; dated notes go to `factory/HISTORY.md`). Print `scripts/ready.sh --stats` so the human sees time per station and review rounds per lane.

## Rules

- A rule must be verifiable ("use `Button` from `@/ui`, never a raw `<button>`"), not aspirational ("write clean UI").
- Never let CLAUDE.md pass 200 lines or a skill pass 500. If it would, split reference material into a linked file.
