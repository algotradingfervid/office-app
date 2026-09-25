---
name: isolate
description: "Start one story in its own git worktree with a clean context, after checking it is ready and does not collide with stories already in flight. Use when starting any story, when the user says \"start story\", \"pick up\", or invokes /isolate <story-id>."
argument-hint: "<story-id>"
disable-model-invocation: true
---

# /isolate — the line, station 1

One story, one worktree, one agent. Nothing else.

## Steps

0. Write the station marker so the guardrails let this station do its job: `echo isolate > .claude/station`. Remove it when done: `rm -f .claude/station`.
1. Read `factory/stories/$0.md`. Refuse if `status` is not `ready`: say which `needs` are unmerged.
2. **Collision check.** Run `scripts/ready.sh`: this story must be in the wave that can start now. If its claim overlaps a story in flight, refuse and name that story. Wait for it to merge; never give this story a fallback ownership of the shared file.
3. Set `status: building`, `agent: <session id or name>`, `started: <YYYY-MM-DD HH:MM>` in the story file on `main`, commit with message `story $0: start`, and push if there is a remote. This comes **before** the worktree, so the branch starts with it; from here on only the story branch edits the story file, until it merges (editing it on both sides makes every merge conflict).
4. Run `scripts/new-story.sh $0`. It creates the worktree at `.claude/worktrees/story-$0` on branch `story/$0-<slug>` from `origin/main` (or `main` with no remote), installs dependencies, copies env files listed in `.worktreeinclude`, and prints the path.
5. Print the handoff block the building agent needs, and nothing else:
   ```
   STORY 007 · Email sign-up
   Worktree: .claude/worktrees/story-007
   Check:    A new user can sign up with email and see an empty dashboard within 60 seconds
   Screen:   factory/design/wireframes/signup.html
   Files:    src/auth/*, src/screens/signup/*
   Load:     the product-* skills that exist (product-architecture, product-proof; product-design-system, product-voice, product-screens once floor 2 is done)
   Next:     /build 007  (in a fresh session opened in the worktree: `claude -w story-007`)
   ```

## Gate

Story is ready, no footprint collision, worktree exists on a fresh branch from `origin/main`.

## Rules

- If this session already carries another story's history, stop and tell the human to open a fresh session in the worktree. Clean context is the point.
- Never start two stories in one worktree.
