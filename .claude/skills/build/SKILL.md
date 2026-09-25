---
name: build
description: "Build one story in its worktree, tests first, using the product's architecture, design system and voice skills, stating assumptions and touching only the story's footprint. Use inside a story worktree when the user says \"build\", \"implement\", or invokes /build <story-id>."
argument-hint: "<story-id>"
---

# /build — the line, station 2

Success criteria, not instructions: make the acceptance check true with the smallest clean change.

## Steps

1. Confirm you are in the story's worktree (`git branch --show-current` starts with `story/$0`). If not, stop: `/isolate $0` first.
2. Read the story file, its wireframe and approved screenshots in `factory/design/wireframes/approved/` if they exist yet, and load the `product-*` skills listed in the handoff. Read `factory/GLOSSARY.md` for any domain noun in the story.
3. **Plan in five lines or fewer**: which files change, which layer each change lives in, which components are used, and what you are assuming. Print it. If any assumption touches the spec or wireframe, ask before continuing. Otherwise continue; do not wait for approval on the plan.
4. **Check first.** Write the failing test(s) or scripted check that encodes the acceptance check. Run them; confirm they fail for the right reason.
5. **Build.** The smallest change that makes the check pass, following `product-architecture`. Use design-system components; never hand-roll one. The screen matches its approved wireframe: same components, layout, navigation, colours and words. If it cannot, ask; do not restyle. A Core-tier screen built before the wireframes exist uses the skeleton's stock styling and no custom visual work; the retrofit stories restyle it later. Every user-facing string follows `product-voice`.
6. **Run everything**: `make check` (tests, lint, types, contrast). Fix root causes; never silence a check.
7. **Self-review against the claim**: `git diff --stat`. Write only inside this worktree (a hook blocks the rest). Any file outside the story's `files` claim needs a one-line justification in the commit message or should be reverted; another story in the wave may own it.
8. **Simplify once.** Re-read the diff as the reviewer will. Remove anything speculative, any abstraction used once, any comment that narrates the obvious, any dead code. (See the `simplicity` skill.)
9. Commit in small steps with messages that say what and why. Set `status: proving` and `built: <timestamp>` in the story file. Print `Next: /prove $0`.

## Gate

The acceptance check's test passes, `make check` is green, the diff stays inside the footprint (or every exception is justified).

## Rules

- Do not "improve" neighbouring code. If you see a problem outside the story, write it to `factory/lessons/` as a candidate story and move on.
- If a check fails twice for the same reason, write the failure to `factory/lessons/<story-id>.md` before a third attempt.
- If you are unsure whether the spec intends something, ask. A wrong assumption acted on silently is the most expensive mistake on the line.
