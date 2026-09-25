# Placet mark (owner-picked variant C): a large solid square and a smaller outlined square on its corner.
# Fill-rule evenodd over three squares: where the outline crosses the solid it cuts through, and inside
# both the shared block stays solid. In colour, that cut (an L shape) is filled with the other blue so the
# ground never shows through (owner, option 1: swapped tones; palette: burnt orange, owner pick). The one-colour
# mark keeps the cut. No masks, so mark-mono.svg also works as a CSS mask image.
# Run: python3 factory/design/brand/make_mark.py
from pathlib import Path

def sq(x, w):
    return f"M{x} {x}h{w}v{w}h-{w}z"

D = sq(5, 40) + sq(29, 30) + sq(35, 18)   # solid 5..45; outline 29..59, 6 wide (inner 35..53)
CUT = sq(29, 16) + sq(35, 10)             # where the outline crosses the solid

def mark(colour, cut=None):
    extra = f'<path fill="{cut}" fill-rule="evenodd" d="{CUT}"/>' if cut else ""
    return ('<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 64 64" role="img" aria-label="Placet">'
            f'<title>Placet</title><path fill="{colour}" fill-rule="evenodd" d="{D}"/>{extra}</svg>\n')

out = Path(__file__).parent / "concord"
(out / "mark.svg").write_text(mark("#B4470E", "#F4925A"))
(out / "mark-dark.svg").write_text(mark("#F4925A", "#B4470E"))
(out / "mark-mono.svg").write_text(mark("#000"))
