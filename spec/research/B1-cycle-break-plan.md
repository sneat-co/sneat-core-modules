# B1 cycle-break — execution plan (sneat-core-modules)

**Goal:** remove every `spaceus|linkage|userus|invitus → contactus` import so that, after
contactus is later extracted to its own module, no `core-modules → contactus` module edge exists.
B1 is verifiable in place: `go build ./...` + `go test ./...` green, and
`grep -rl sneat-core-modules/contactus {spaceus,linkage,userus,invitus}` returns nothing.

**Branch:** `refactor/b1-contactus-cycle-break` · single branch, one big change.

## Categories of back-edge (from code audit 2026-06-20)

1. **Foundational data + constants** — `briefs4contactus`, `const4contactus`.
   Embedded structs / constants; not invertible. **Fix: relocate** out of `contactus/` to a new
   low package `contactusmodels/` (stays in core-modules; both contactus and siblings import it).
2. **Facade behavior** — `facade4contactus.RefuseToJoinSpace` (spaceus). **Fix: invert** behind a
   consumer-owned interface; contactus registers an impl.
3. **DAL entry types + transaction workers** — `dal4contactus.{ContactEntry, ContactusSpaceEntry,
   New*, RunContactusSpaceWorker/…Readonly/…NoUpdate, ContactusSpaceWorkerParams, RunContactWorker,
   ContactWorkerParams}`, `dto4contactus.ContactRequest` (spaceus, userus, invitus).
   **Fix: wrap-then-invert** — coarse `facade4contactus` method per usage, then invert the facade.

## Phases

### Phase 1 — Relocate cat-1 (kills linkage entirely + all cat-1 edges)
- `git mv contactus/briefs4contactus contactusmodels/briefs4contactus`
- `git mv contactus/const4contactus contactusmodels/const4contactus`
- Rewrite import path `…/contactus/briefs4contactus` → `…/contactusmodels/briefs4contactus`
  and `…/contactus/const4contactus` → `…/contactusmodels/const4contactus` across all `*.go`
  (package identifiers unchanged → no other edits).
- Verify: `go build ./...` + `go test ./...` green.
- Verify: linkage no longer imports contactus.

### Phase 2 — Wrap-then-invert cat-2/3, per consumer
For each of spaceus, userus, invitus: introduce consumer-owned interface(s) for the contactus
behavior they need; contactus provides + registers the implementation at init; replace direct
`dal4contactus`/`facade4contactus`/`dto4contactus` use with the interface. Verify green after each.

### Phase 3 — Final verification
- `go build ./...` + `go test ./...` green.
- `grep -rl sneat-core-modules/contactus spaceus linkage userus invitus` → empty.
- contactus still builds (it keeps its forward deps on spaceus/linkage/userus/invitus — allowed).

## Status (2026-06-20)

**DONE & VERIFIED (this branch, full module `go build ./...` + `go test ./...` green):**
- Phase 1 relocation → `contactusmodels/{briefs4contactus,const4contactus}`.
- **spaceus** — 0 back-edges. `http_join_refuse` made self-contained; `create_space` routes
  contactus record creation through `facade4spaceus.ContactusSpaceContributor` (registered by
  `contactus.Extension()`); `member_helpers.go` logic moved into `contactus/spaceus_contributor.go`.
- **userus** — 0 back-edges. `set_user_country` routes through
  `facade4userus.ContactusCountryUpdater` (impl `contactus/userus_contributor.go`);
  `GetUserSpaceContactID` takes a local `contactusSpaceContactsReader` interface.
- **linkage** — 0 back-edges (prod + test; test uses local fakes).
- **auth** test regression fixed (registers a fake contributor; stays contactus-free).

**DONE (follow-up branch `refactor/b1-invitus-cycle-break`) — invitus, accessor inversion:**
invitus now imports zero contactus-module packages (incl. tests). All 4 modules at 0 back-edges →
**module cycle fully broken.** Approach: `facade4invitus` owns `ContactusAccess` (+ `SpaceContactsSession`,
`ContactSession`, `MemberContact`, local `ContactRequest`) registered by `contactus.Extension()`;
impl `contactus/invitus_contributor.go` adapts the contactus workers/DBOs. invitus keeps owning each
transaction (one worker callback per op) so claim/join/create atomicity is preserved. Write-only
response fields `Contact`/`ContactusSpace` dropped (no readers). `go build`/`go test ./...` green.

**(historical) DEFERRED — invitus (7 files), separate follow-up:**
`claim_personal_invite.go`, `join_space.go`, `create_invite_to_contact.go`, `get_personal_invite.go`,
`join_info.go`, `create_invite_response.go`, `claim_personal_invite_test.go`.

Why deferred: invitus's claim/join flows mutate contactus contacts + invite status + user records in
ONE atomic `RunContactusSpaceWorker`/`RunContactWorker` transaction. "Coarse contactus methods" would
split that across transactions and break atomicity. Correct fix = **option A (accessor inversion)**:
contactus exposes a generic transaction runner that hands invitus an interface over the contactus
space/contact records (`GetContactBriefByUserID`, `Contacts`, `AddContact`, `AddUserID`, member
load/update); invitus keeps owning its transaction and its invite/user logic. Also: `join_info.go`'s
single read-only member load can use a simple `GetContactBrief(ctx, getter, spaceID, contactID)` method.
Note: `contactus → invitus/dbo4invitus` (contact_invites.go) is a FORWARD edge — allowed, leave it.
