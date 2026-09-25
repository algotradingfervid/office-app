---
name: brainstorm
description: "Start a new product in the SHIPSHOW factory. Widen the idea, find the real problem and users, propose directions, pick one, and name the riskiest assumption. Use when the user has an app idea, says \"I want to build\", \"new product\", or invokes /brainstorm. This is always the first station."
---

# /brainstorm — Studio, station 1

Think as the product owner, not the engineer. No stack, no code, no features yet. The output is a problem worth solving and one direction to pursue.

## Steps

1. **Listen.** Ask the human to describe the idea in their own words. Do not interrupt with structure yet.
2. **Widen.** Ask, in rounds of at most three questions (use the AskUserQuestion tool when it exists):
   - Who has this problem? When exactly does it bite them? What do they do about it today?
   - What would "better" look like to them, in their words, not yours?
   - Why now? What changed that makes this possible or needed?
   - What would make this a failure even if it worked technically?
3. **Diverge.** Propose 3–5 distinct directions. Each is one paragraph: the user, the sharpest version of the problem, the shape of the solution, and what it deliberately does not do. Make them genuinely different (different user, different wedge, different business model), not variations of one idea.
4. **Converge.** Ask the human to pick one or combine two. Push back once if the pick is the vaguest option.
5. **Name the riskiest assumption.** The one belief that, if false, kills the product. Usually "people want this" or "people will pay" or "we can get the data". Say how it could be tested cheaply (a landing page, ten conversations, a throwaway prototype).
6. **Set the dials.** Run the `/dials` logic inline (one question) and write `factory/DIALS.md`.
7. **Write** `factory/BRIEF.md` using `factory/templates/brief.md`.
8. **Write** `factory/STATE.md` (create from `factory/templates/state.md` if missing) and mark `brainstorm: done`.

## Gate

One direction chosen. One riskiest assumption written down with a cheap test for it. Dials set. Do not proceed to `/research` without all three.

## Route the riskiest assumption

- **Demand** ("people want this", "people will pay"): suggest a disposable prototype before the spec: the ugliest thing that shows the core moment, shown to ten people, then deleted.
- **Can we do it** (data, accuracy, feasibility, an AI step working well enough): `/spike` runs after `/research` and before `/define`, and its number decides whether the spec goes ahead.

Write the route in the brief either way.
