---
name: shipshow-update
description: "Research what has changed in Claude models, Claude Code capabilities, and agentic-engineering best practice since the factory was last updated, then propose and apply upgrades to the core kit (skills, agents, hooks, CLAUDE.md) with a changelog. Use when a new model or Claude Code feature ships, every few weeks, or when the user invokes /shipshow-update."
argument-hint: "[--apply | --dry-run (default)]"
disable-model-invocation: true
---

# /shipshow-update — the factory upgrades itself

Models get better; scaffolding that was load-bearing last quarter becomes overhead this quarter, and new capabilities make new stations possible. This skill keeps the core kit current. It never touches product files or `product-*` skills; those belong to `/learn`.

## Steps

0. Write the station marker so the guardrails let this station do its job: `echo shipshow-update > .claude/station`. Remove it when done: `rm -f .claude/station`.
1. Read `factory/KIT-VERSION.md` (last update date, model at the time, the list of core files and their purposes, and the open questions from the last run).

2. **Research**, with sources for every claim, in parallel `researcher` subagents:
   - **Models.** What Claude models are current, what changed (context size, tool use, thinking/effort controls, known behaviour changes such as over-engineering tendencies or verification habits). Sources: Anthropic model docs and release notes, the Claude 4.x/5.x prompting best-practices page.
   - **Claude Code.** New or changed features: skills frontmatter fields, subagent options, hooks events, worktree and parallel-agent features, agent teams, `/goal`, `/loop`, workflows, headless flags, computer use, browser tools. Sources: code.claude.com docs, the changelog, Anthropic engineering blog.
   - **Practice.** What Anthropic's engineering posts and the well-reviewed community frameworks (BMAD, Spec Kit, Superpowers, GSD, OpenSpec, Compound Engineering, gstack, Ras Mic's skills) changed or dropped in the period, and what practitioners now warn against. Include Karpathy's latest notes on working with agents.
   Write findings to `factory/updates/<date>/research.md`.

3. **Diff against the kit.** For each core file, ask three questions and write the answer:
   - *Obsolete?* Does the current model do this correctly without being told? (Evidence: our own `factory/lessons/LOG.md` shows the mistake stopped happening, or Anthropic says so.) → propose removing the instruction.
   - *Wrong?* Does a frontmatter field, hook event, flag or command in the file no longer exist or behave differently? → propose the fix.
   - *Missing?* Does a new capability let a station do its job better (e.g. a native evaluator replacing a hand-written review loop, a native worktree feature replacing `scripts/new-story.sh`, a new hook event making a guardrail deterministic)? → propose the change, smallest version first.
   Also check the assumptions in the "Rules of the floor" (fresh context per story, writer/reviewer split, evidence over claims, hooks over instructions, isolation over swarms). If new evidence contradicts one, say so explicitly rather than silently keeping it.

4. **Propose.** Write `factory/updates/<date>/proposal.md`: a table of file → change → reason → source → risk (low/medium/high). Keep the kit smaller where possible; every addition must name what it replaces or why nothing does.

5. **Apply** (only with `--apply`, or after the human says yes to the proposal):
   - Make the edits on a branch `kit/update-<date>`.
   - Run the tool-up trial (step 11 of `/tool-up`) on the smallest existing story using the updated kit, in a fresh subagent, and record whether anything regressed.
   - Update `factory/KIT-VERSION.md`: date, model, summary of changes, open questions to check next time.
   - Append to `factory/updates/CHANGELOG.md`.
   - Open a PR for the human to merge. The factory does not merge its own upgrades.

6. Print the proposal summary and the PR link.

## Gate

Every proposed change cites a source or a lessons-log entry. The trial run passes on the updated kit. High-risk changes (removing a guardrail, changing the review loop) are flagged for explicit human approval.

## Rules

- Prefer removing over adding. The best update is one that deletes scaffolding the model no longer needs.
- Never change the five human gates without the human's written agreement.
- If a source is a community blog rather than Anthropic or a maintainer, label it as such and weight it lower.
- Schedule: suggest running this after every major model release and otherwise monthly. The human can wire it to a scheduled task or `/loop`; this skill does not schedule itself.
