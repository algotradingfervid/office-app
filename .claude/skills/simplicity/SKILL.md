---
name: simplicity
description: "Coding standards for every line of application code written or reviewed in this factory - load whenever writing, refactoring, simplifying or reviewing code. Based on Andrej Karpathy's observations about how agents go wrong and Anthropic's anti-over-engineering guidance."
user-invocable: false
---

# Simplicity standards

Agents tend to over-complicate code and APIs, bloat abstractions, add defensive code for cases that cannot happen, copy-paste instead of reuse, change comments and code they did not understand as side effects, and leave dead code behind. Karpathy's summary: they will "implement a bloated construction over 1000 lines when 100 would do", and asked to simplify, they "just can't". So the standard is written down and checked, not hoped for.

## The rules

1. **Smallest change that makes the check pass.** Then stop. Features, refactors and "improvements" beyond the story are not in scope.
2. **100 lines beats 1,000.** If a solution feels large, the design is probably wrong; ask what would make it small.
3. **No speculative abstractions.** An abstraction is earned by the third concrete use, not predicted by the first.
4. **No defensive code for impossible cases.** No try/except around things that cannot fail, no validation of inputs the type system already guarantees, no fallbacks for scenarios the spec excludes.
5. **Explicit over clever.** A reader with no context should follow the code top to bottom. Prefer a plain function over a pattern.
6. **Match the existing style.** The codebase's conventions win over your preferences, even when yours are better. Consistency is what lets the next agent copy the right pattern.
7. **Surgical diffs.** Touch only what the story needs. Do not reformat, rename, reorder or "tidy" untouched code. Do not change or delete a comment you do not fully understand.
8. **Delete dead code.** If your change makes something unused, remove it in the same change.
9. **Comments say why, not what.** A comment that narrates the line below it is noise; a comment that explains a non-obvious decision is gold.
10. **Tests verify, they do not define.** Never special-case to make a test pass. If a test is wrong, say so and fix the test.
11. **Minimal dependencies.** Adding a package is a decision with a cost; justify it in the commit message or do without.
12. **Surface assumptions.** Before acting on a guess about the spec, the data or another module, write the assumption down (in the plan or a comment) or ask.

## The self-review question

Before committing, read the diff as a sceptical senior engineer and ask of every hunk: "If I removed this, would the acceptance check still pass?" If yes, remove it.
