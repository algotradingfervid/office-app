---
id: 018
title: Restyle sign-in and home to the wireframes
tier: usable
lane: core
kind: feature
status: blocked
needs: ["017", "007"]
files: ["internal/core/auth/templates/login.html", "internal/core/home/templates/home.html", "internal/forms/leave/templates/balance_card.html"]
screen: auth/home screens
check: "The auth/home screens match their approved wireframes at 390 and 1280 px"
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

# Restyle sign-in and home to the wireframes

## What
Restyle the Core-tier auth/home templates to their approved wireframes using the design-system components from 017. Markup and classes only; no behaviour changes.

## Check
The screens' existing journeys (stories 007, 011, 013, 014, 015) re-run by `scripts/prove.sh 018`, screenshots beside the approved wireframes.

## Out of scope
Handlers, services, new features.

## Constraints
Released by /wireframe after 017 merges.
