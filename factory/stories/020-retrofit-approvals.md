---
id: 020
title: Restyle inbox, my requests and request pages to their wireframes
tier: usable
lane: core
kind: feature
status: blocked
needs: ["017", "011", "013", "015"]
files: ["internal/core/approvals/templates/inbox.html", "internal/core/approvals/templates/my_requests.html", "internal/core/approvals/templates/request.html"]
screen: approvals screens
check: "The approvals screens match their approved wireframes at 390 and 1280 px"
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

# Restyle inbox, my requests and request pages to their wireframes

## What
Restyle the Core-tier approvals templates to their approved wireframes using the design-system components from 017. Markup and classes only; no behaviour changes.

## Check
The screens' existing journeys (stories 007, 011, 013, 014, 015) re-run by `scripts/prove.sh 020`, screenshots beside the approved wireframes.

## Out of scope
Handlers, services, new features.

## Constraints
Released by /wireframe after 017 merges.
