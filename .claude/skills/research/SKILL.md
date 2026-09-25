---
name: research
description: "Research the market, competitors, user pains, constraints and feasibility for the product in factory/BRIEF.md, with sources for every claim. Use after /brainstorm, when the user says \"research\", \"competitors\", \"is there a market\", or invokes /research."
---

# /research — Studio, station 2

Every claim gets a source. Opinions are labelled as opinions.

## Steps

1. Read `factory/BRIEF.md` and `factory/DIALS.md`. Stop if the brief is missing: `/brainstorm` first.
2. Plan the research questions from the brief. Always include the riskiest assumption.
3. **Fan out.** On `research: full`, launch parallel `researcher` subagents (see `.claude/agents/researcher.md`), one per lane, each writing its own file in `factory/research/`:
   - `market.md` — who buys, how many, how they buy today, price norms.
   - `competitors.md` — the 5–8 closest alternatives; for each: what it does, price, what users praise, what users complain about. Mine app-store reviews, G2/Capterra, Reddit, Hacker News, support forums. Quote real complaints.
   - `user-pains.md` — pains ranked by how often they appear across sources, with quotes.
   - `constraints.md` — platform rules (app stores, OS), legal and data rules, payment and identity requirements, anything that forces a decision.
   - `feasibility.md` — what is hard, what exists as a library or service, what needs data you may not have.
   On `research: quick`, do one lane yourself: competitors (top 3) plus the riskiest assumption. 30-minute budget.
4. **Human research.** If the riskiest assumption needs real conversations, write `factory/research/interview-guide.md` (8 questions, no leading ones) and stop: the human runs the conversations and drops notes into `factory/research/interviews/`. When notes exist, synthesise them into `user-pains.md`.
5. **Synthesise** into `factory/RESEARCH.md` using `factory/templates/research-brief.md`: what we learned, the competitor table, top 5 pains, the constraints that force decisions, and a go / no-go / pivot recommendation on the riskiest assumption with reasoning.
6. Update `factory/STATE.md`: `research: done`.

## Gate

Every factual claim in `RESEARCH.md` links to a source. The riskiest assumption has a verdict: tested, or consciously accepted by the human with a note. If the verdict is "pivot", go back to `/brainstorm` step 3 with the new evidence; do not proceed to `/define`.

## Rules

- Do not pad. If a lane finds nothing, say so in two lines.
- Prefer primary sources (the competitor's own pricing page, the store's own guidelines) over blog summaries.
- Never invent a statistic. "Unknown" is a valid finding.
