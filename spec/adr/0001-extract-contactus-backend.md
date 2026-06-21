# ADR 0001 — Extract `contactus` backend into its own module

**Status:** Accepted & implemented (2026-06-21)
**Scope:** Workstream B (backend) of decoupling `contactus` into the dedicated `github.com/sneat-co/contactus` repo. (Frontend = Workstream A, tracked separately.)
**Related:** `spec/research/contactus-repo-extraction-plan.md`, `spec/research/B1-cycle-break-plan.md`.

## Context

`contactus` lived inside the single module `github.com/sneat-co/sneat-core-modules`. It had a **bidirectional** relationship with four sibling core modules — `spaceus`, `linkage`, `userus`, `invitus` — plus the root composition aggregator. Legal within one module (the package graph was acyclic), but extracting `contactus` as-is would have created a **module-level cycle** `contactus ↔ sneat-core-modules`, which makes releases chicken-and-egg and is a long-term maintenance trap.

## Decision

Extract `contactus` into a standalone module `github.com/sneat-co/contactus/backend`, breaking the cycle first so the dependency is strictly one-directional: **`contactus/backend → sneat-core-modules`** (never the reverse).

The back-edges were broken by **category**, because one tool did not fit all:

1. **Foundational data & constants are not invertible.** `briefs4contactus` (e.g. `ContactBrief`, `ContactBase` — the latter is *embedded* in `userus.UserDbo`) and `const4contactus` are shared vocabulary. They were **relocated** out of `contactus/` into a new low package **`sneat-core-modules/contactusmodels/`** that both core-modules and the extracted module import. They did **not** move into `contactus/backend`.

2. **Behaviour was inverted via registered contributors** (mirroring the pre-existing `facade4linkage.RegisterDboFactory` pattern). Each consumer module owns an interface + registration hook; `contactus` provides and registers the implementation at startup:
   - `facade4spaceus.ContactusSpaceContributor` — builds the contactus records created during space creation.
   - `facade4userus.ContactusCountryUpdater` — updates a contact's country within a space.
   - `facade4invitus.ContactusAccess` (+ `SpaceContactsSession` / `ContactSession` / `MemberContact`) — an **accessor inversion**, chosen because invitus's claim/join flows mutate contactus contacts + invite + user records in **one atomic transaction**; coarse contactus methods would have split that across transactions and broken atomicity. invitus keeps owning its transaction and calls contactus through the accessor.

3. **Composition root moved to the application.** `sneat_core_modules.CoreExtensions()` no longer returns `contactus.Extension()`. Keeping it would force core-modules to import the extracted module and re-create the cycle. Apps now register `contactusext.Extension()` themselves (e.g. `sneat-go-backend` `standardExtensions()`).

The extracted module mirrors the `assetus`/`listus` dedicated-repo layout: packages flattened under `backend/` (`dal4contactus`, `dbo4contactus`, `dto4contactus`, `facade4contactus`, `api4contactus`, …) and the root package as **`contactusext`** (`Extension()` + the contributor registrations).

## Resulting architecture

```
github.com/sneat-co/contactus/backend
        │  requires
        ▼
github.com/sneat-co/sneat-core-modules
        ├── contactusmodels/{briefs4contactus, const4contactus}   ← shared low package
        ├── spaceus, linkage, userus, invitus                      ← own the inversion interfaces
        └── (no dependency on contactus/backend — cycle-free)
```

Consumer apps (sneat-go-backend, sneat-go, …) depend on **both** modules and wire `contactusext.Extension()` at their composition root.

## Release map

| Tag | Repo | What |
|---|---|---|
| `v0.38.59` | sneat-core-modules | B1 cycle-break (relocation + inversions); contactus still present |
| `contactus/backend v0.1.0` | contactus | B2 — extracted module |
| `v0.38.60` | sneat-core-modules | B3 — contactus dropped from `CoreExtensions()` |
| `calendarius backend/v0.2.1`, `logistus backend/v0.1.0`, `sneat-bots v0.1.3`, `sneat-go-backend v0.59.2`, `debtus backend/v0.1.1`, sneat-go | (consumers) | B3 — repointed to `contactus/backend` |
| `v0.38.61` | sneat-core-modules | B4 — orphaned `contactus/` packages removed |

Consumer release order was dependency-driven: **sneat-bots → debtus → sneat-go** (and the independents calendarius/logistus/sneat-go-backend), because each repo could only drop its local `replace` once its siblings' repointed versions were published.

## Import migration (for any remaining/out-of-tree consumer)

- `sneat-core-modules/contactus/briefs4contactus` → `sneat-core-modules/contactusmodels/briefs4contactus`
- `sneat-core-modules/contactus/const4contactus` → `sneat-core-modules/contactusmodels/const4contactus`
- `sneat-core-modules/contactus/<other pkg>` → `contactus/backend/<other pkg>`
- `sneat-core-modules/contactus` (root) → `contactus/backend/contactusext`
- Add `contactusext.Extension()` at the app composition root (core-modules no longer provides it).
- Require `sneat-core-modules v0.38.60+` and `contactus/backend v0.1.0+`.

## Consequences

**Positive:** no module cycle; contactus versions and ships independently; core-modules is smaller; the inversion interfaces make the contactus coupling explicit and testable (consumer tests use fakes, staying contactus-free).

**Costs / watch-items:**
- Apps **must** register `contactusext.Extension()` — forgetting it silently drops contactus wiring.
- The `contactusmodels` split is a permanent seam; new shared contact vocabulary belongs there, not in `contactus/backend`.
- Invitus's accessor interface re-exposes a broad slice of the contactus DAL; keep it as narrow as the flows require.
- B1 changed some invitus response shapes (`ClaimPersonalInviteResponse.ContactusSpace` removed; `ContactRequest` became a local mirror) — out-of-tree consumers may hit the same drift fixed in debtus/sneat-bots.
