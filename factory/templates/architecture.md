# ARCHITECTURE — <product name>

## Platforms and frameworks
| Target | Framework | Why (and the tradeoff accepted) |
|---|---|---|

## Shape: modular monolith
One repo, one binary, one deploy. Each module owns its routes, logic, queries, migrations and tests.
```
module/  UI  →  services  →  data
```
The layer rule inside a module: … (UI renders and calls services; services own logic; data owns persistence; nothing skips a layer.)
Between modules: only through each module's interface package. Enforced by: … (the import-boundary check in `make check`).
Registration: how a module adds its routes, wiring and migrations without editing a shared file: …
Separate processes (if any) and the workload reason: …

## Data model
```mermaid
erDiagram
```
Which actions touch which entities.

## Buy vs build
| Concern | Choice | Why |
|---|---|---|
| Auth | | |
| Payments | | |
| Email | | |
| Storage | | |
| Analytics | | |

## Hosting and deploy
Where it runs · preview per story (`scripts/preview.sh`) · production deploy · rollback.

## Modules and seams for parallel work
| Module | Folder | Tier | Interface it exposes | Typical story |
|---|---|---|---|---|

Hotspots found by `scripts/ready.sh` and how they were removed:

## Budgets (from the spec)
Load time · offline · privacy · size.

## Anti-patterns for this codebase
Things a fresh agent will be tempted to do and must not.
