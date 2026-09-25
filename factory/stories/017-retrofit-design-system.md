---
id: 017
title: Swap stock styling for the design system
tier: usable
lane: core
kind: contract
status: blocked
needs: []
files: ["internal/core/web/static/app.css", "internal/core/web/templates/base.html", ".claude/skills/product-design-system/SKILL.md", ".claude/skills/product-voice/SKILL.md", ".claude/skills/product-screens/SKILL.md"]
screen: all screens
check: "Every page loads the design-system stylesheet and layout, and the product skills for design system, voice and screens exist"
blocked: waiting for wireframe approval
agent: 
started: 
built: 
proved: 
reviewed: 
merged: 
review-rounds: 0
pr: 
---

# Swap stock styling for the design system

## What
Replace the stock stylesheet and layout with the approved design system (tokens, components, app shell with navigation), and write the `product-design-system`, `product-voice` and `product-screens` skills (/tool-up items 3–5).

## Check
Screenshots of login and home at 390/1280 beside the approved wireframes.

## Out of scope
Restyling module screens (018–020).

## Constraints
Released by /wireframe. Keep CSP: no inline styles or scripts.
