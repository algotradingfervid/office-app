# Palette — Concord (burnt orange & graphite)

Owner picked option 1 of `../moodboard-orange.html` on 2026-09-25. One accent, warm paper neutrals,
graphite structure, status colours kept away from the accent. Flat fills only: **no gradients anywhere**.
Contrast is WCAG 2.x, computed 2026-09-25; every text pairing used in `preview.html` is checked by script.

| Role | Name | Light | Dark | Why |
|---|---|---|---|---|
| Accent | Burnt orange | `#B4470E` | `#F4925A` | Primary action, links, the mark, the one "look here" colour. White on it 5.46:1; as text 4.89:1 on paper, 5.46:1 on white; dark 7.45:1 on dark card |
| Mark second tone | Ember / burnt | `#F4925A` | `#B4470E` | Where the mark's outline crosses the solid (owner: swapped tones). Graphic only, never text |
| Neutral ground | Paper / Char | `#F6F2EA` | `#141312` | Warm off-white page (analog touch); warm near-black in dark |
| Card | White / Char card | `#FFFFFF` | `#1D1B19` | Content surfaces |
| Support | Graphite | `#2F3237` | `#ECE8E1` | Structure: the heavy rule under table heads and section labels, secondary emphasis |
| Text | Ink | `#1E2024` | `#ECE8E1` | 14.7:1 on paper; 13.5:1 on dark card |
| Muted text | Pencil | `#5A5D63` | `#A8A49D` | Labels, meta; ≥ 5.9:1 light, ≥ 7:1 dark |
| Hairline | Rule | `#E0D9CC` | `#34302C` | The register lines between rows |
| Highlight | Ember tint | `#FBE9DE` | `#3A2518` | Selected row, current item; a fill, text on it stays ink |
| Approved | Green word | `#2E6B34` | `#7FCB86` | Status word + ✓ only, never a fill or brand colour |
| Pending | Slate | `#4F5B73` | `#A9B4CC` | Status word + ◯; slate so it never looks like the orange accent |
| Rejected / danger | Berry | `#A3213F` | `#F28CA0` | Status word + ✕, error rules, the destructive button; cool red, far from the orange |

## Use (mobile first)
- A typical phone screen: paper, white cards, ink text, hairlines, and burnt orange on one primary button
  plus links. Never two orange buttons on one screen.
- Status = word + symbol + colour; colour is never the only signal.
- Focus ring: 3px accent outline, 2px offset (≥ 3:1 against both grounds).
- Corners 4px; no shadows, no blur, no gradients.
- Fervid uses orange too (`#DA4F27`); Placet's is deeper and browner, a relative rather than a copy.
