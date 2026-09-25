# Palette — Stamp-pad

One loud accent, ink and paper. Everything is a flat fill: no gradients, no tints-of-tints.
Contrast ratios are WCAG 2.x, computed 2026-09-25.

| Role | Name | Hex | Why |
|---|---|---|---|
| Accent | Stamp-pad blue | `#2340C8` | The blue of an office rubber stamp; the one loud colour. Text on paper 6.87:1 |
| Accent on dark | Stamp-pad blue, lit | `#93A8FF` | Same ink read on carbon; text on carbon 7.77:1 |
| Neutral, light ground | Paper | `#F4EEDF` | Warm form paper, not white; the analog ground |
| Neutral, dark ground | Carbon | `#1B1914` | Carbon-copy black-brown, the dark-mode ground |
| Ink | Ink | `#17150F` | Outlines, hard shadows and body text on paper (15.77:1) |
| Support | Manila | `#E7C35A` | Folder manila for highlights, tags and "pending" chips; a fill only, ink text on it 10.73:1 |
| Muted text | Pencil | `#5B5547` / `#B3AB98` on carbon | Secondary text; 6.40:1 and 7.69:1 |
| Success | Ledger green | `#1E6B35` / `#7FCB8F` on carbon | Approved, credited; 5.65:1 and 9.07:1 |
| Warning | Ochre | `#8A4B00` / `#F0B24A` on carbon | Near a limit, pending too long; 5.88:1 and 9.34:1 |
| Danger | Red-stamp | `#B8231B` / `#FF8A7A` on carbon | Rejected, destructive actions; 5.51:1 and 7.67:1 |

## Use
- A typical screen shows **paper, ink and one blue**. Manila and the semantic colours appear only on
  status chips and messages. Never more than one blue block per view.
- Blue is for the primary action, links, and the "approved" stamp. Paper text on a blue button: 6.87:1.
- Hard shadow: `4px 4px 0 ink` on paper, `4px 4px 0 paper` on carbon. No blur, ever.
- Avoided: purple/violet (Qubitech's colour, and an AI default), Fervid orange (Placet stands apart).
