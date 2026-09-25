# Type — Stamp-pad

Two faces, both Google Fonts under the SIL Open Font License, self-hosted as woff2 (the CSP allows no
font origins other than our own). Files: `factory/design/brand/fonts/`.

| Role | Face | Setting | Fallback stack |
|---|---|---|---|
| Display (headings, the wordmark) | **Archivo** (variable, wght 100–900, wdth 62–125) | weight 800, width 112%, tight tracking, often uppercase | `"Archivo", "Arial Black", "Helvetica Neue", Arial, sans-serif` |
| Body (sentences, buttons, form labels) | **Archivo** | weight 400/600, width 100% | `"Archivo", "Helvetica Neue", Arial, sans-serif` |
| Typed (numbers, dates, request codes, small labels) | **Courier Prime** 400/700 | tabular by nature; labels uppercase with 0.06em tracking | `"Courier Prime", "Courier New", Courier, monospace` |

Why: Archivo's heavy, wide cuts give the chunky, stamped headings of retrobrutalism, and its normal width
reads cleanly at phone sizes. Courier Prime is a real typewriter face, so balances and dates look typed
onto a form. Avoided: Space Grotesk, Space Mono and Inter (AI defaults) and pixel monos (digital, not analog).

## Scale (ratio 1.25, base 16px; nothing below 14px)
| Token | Size / line height | Use |
|---|---|---|
| `--step-5` | 49 / 1.0 | Page hero numbers (a balance) |
| `--step-4` | 39 / 1.05 | Page title on desktop |
| `--step-3` | 31 / 1.1 | Page title on phone |
| `--step-2` | 25 / 1.2 | Section heading |
| `--step-1` | 20 / 1.3 | Card heading |
| `--step-0` | 16 / 1.5 | Body, buttons, inputs |
| `--step--1` | 14 / 1.4 | Typed labels (uppercase Courier Prime), captions |
