# KIT-VERSION

kit: SHIPSHOW core
version: 0.3.0
built: 2026-09-16
built-with: Claude Code (skills, subagents, hooks, worktrees as documented at code.claude.com, Sept 2026)
model-at-build: Claude Opus 5.5 and Haiku 4.5 (kit uses Opus and Haiku only since 0.3.0)

## Core files (owned by /shipshow-update; never edited by /learn)
- CLAUDE.md (above the PRODUCT SECTION marker)
- .claude/skills/{shipshow,dials,brainstorm,research,spike,define,name-it,brand,design-system,wireframe,blueprint,tool-up,slice,isolate,build,prove,ship,learn,simplicity,shipshow-update}
- .claude/agents/{researcher,reviewer-code,reviewer-design,reviewer-blind,verifier}
- .claude/hooks/*, .claude/settings.json
- scripts/*, factory/templates/*, factory/MODELS.md

## Product files (owned by the stations and /learn)
- CLAUDE.md product section, .claude/skills/product-*, .claude/agents/product-review-*, factory/*.md, factory/design, factory/stories, factory/evidence, factory/lessons

## Assumptions to re-check on the next /shipshow-update
- 0.3.0 limits (≤ 5 non-test files, ≤ 60-line story body, ≤ 60-line STATE.md, ≤ 2 review rounds): check `scripts/ready.sh --stats` and the lessons log after the first product tier on 0.3.0; loosen or tighten with evidence.
- guard-worktree.sh reads `cwd` and `tool_input.file_path` from the hook's JSON input; verify both field names are still current.
- `design-timing: after-core` assumes one retrofit story can restyle Core screens; check its size on the first product that uses it.
- Skill frontmatter fields used: name, description, argument-hint, disable-model-invocation, user-invocable. Verify they are still current.
- Subagent frontmatter: name, description, tools, model. Verify `model: inherit` is still supported.
- Hook events used: PreToolUse (Edit|Write|MultiEdit, Bash), PostToolUse, Stop. Verify exit-code-2 semantics and the 8-block Stop override.
- Worktrees: `claude -w <name>` and `.claude/worktrees/`. If Claude Code's native worktree/`/batch` covers new-story.sh, replace the script.
- Review loop: if a native evaluator (`/goal`) or agent teams can replace the hand-run reviewer loop in /ship, prefer the native feature.
- Model table in factory/MODELS.md: re-tier when a new mid-tier model matches the previous top tier.
- reviewer-blind stays blind by a copied folder plus an instruction, not a hook. If Claude Code adds per-subagent file-path restrictions, enforce it with those.
- scripts/look.py needs Python Playwright (`pip install playwright && playwright install chromium`). If Claude Code's browser tools can screenshot headless local pages at set widths and colour schemes, compare and prefer the smaller option.

## Changelog
- 0.3.0 (2026-09-23) — core first, small parallel stories, ceremony by risk. New /spike station; tiered spec (Core / Usable / Launch / Later) sliced one tier at a time; design rooms run alongside the Core tier by default (`design-timing`); story lanes (core / full); review gate is no blocking findings in at most two rounds; stories ≤ 5 non-test files and ≤ 60 lines, contract stories first, parallel waves with disjoint file claims (`scripts/check-stories.sh`, `scripts/ready.sh` waves, hotspots, `--stats`); modular monolith with self-registering modules and an import-boundary check; guard-worktree hook; evidence pinned to the branch head with mutation results, raw captures git-ignored; STATE.md one screen plus HISTORY.md; /learn per wave, repeats become hooks; Opus and Haiku only. Source: factory/updates/2026-09-23/proposal.md.
- 0.2.1 (2026-09-16) — talk steps in /brand, /design-system and /wireframe: earlier remarks do not count as answers; if the human is not available, stop and wait instead of drawing. /wireframe draws no other screens until the human has reacted to the first. Source: SahiBid's /wireframe run inferred the layout from earlier remarks and folded the first-screen check into the final presentation because the founder was away.
- 0.2.0 (2026-09-16) — design stations talk first, look, and get a blind review. /brand agrees the direction with the human (DIRECTION.md) before drawing; /design-system adds an interface check-in, a look.py pass and a blind review; /wireframe draws branded high-fidelity clickable wireframes, agrees layout and a first screen, looks, gets a blind review, and locks approved screenshots as the visual contract for /blueprint, /build, /prove and reviewer-design. New: scripts/look.py, scripts/blind-copy.sh, agents/reviewer-blind. Source: factory/updates/2026-09-16/proposal.md.
- 0.1.0 — initial kit.
