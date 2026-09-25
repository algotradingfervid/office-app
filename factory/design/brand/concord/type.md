# Type — Concord

IBM Plex, SIL Open Font License, self-hosted as woff2 (the CSP allows no font origins but our own).
Files: `factory/design/brand/fonts/ibm-plex-*.woff2`.

| Role | Face | Setting | Fallback stack |
|---|---|---|---|
| Headings | **IBM Plex Sans** (variable, wght 100–700) | 600, sentence case, -0.01em | `"IBM Plex Sans", "Segoe UI", "Helvetica Neue", Arial, sans-serif` |
| Body, buttons, labels | **IBM Plex Sans** | 400 body, 500 buttons and labels | same |
| Reference numbers, dates, balances | **IBM Plex Mono** | 400 / 500, tabular | `"IBM Plex Mono", "SFMono-Regular", Menlo, Consolas, monospace` |

Why: Plex Sans is engineered and business-like without being bland, and reads well on phones. Plex Mono
has a faint typewriter character, which is the "typed on a form" analog touch for request codes, dates
and day counts. Avoided: Inter (AI default), Space Grotesk.

## Scale (ratio 1.2, base 16px; nothing below 14px)
| Token | Size / line height | Use |
|---|---|---|
| `--step-4` | 33 / 1.15 | Page title on desktop; a balance figure |
| `--step-3` | 28 / 1.2 | Page title on phone |
| `--step-2` | 23 / 1.25 | Section heading |
| `--step-1` | 19 / 1.35 | Card heading |
| `--step-0` | 16 / 1.5 | Body, buttons, inputs |
| `--step--1` | 14 / 1.45 | Labels, captions, table meta |
