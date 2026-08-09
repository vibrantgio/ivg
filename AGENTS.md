# AGENTS.md — ivg

IconVG rendering for Go: the compact binary format for icons, logos and
glyphs, decoded, encoded, generated and rendered against a `Rasterizer`
interface rather than a fixed bitmap, so one blob of icon data drives an
image or a Gio op list. `decode`, `encode`, `generate` and `render` are
that pipeline; `raster` declares the rasterizer interface, `raster/gio`
implements it over an `op.Ops` and adds the `Widget` and `GioPaint` helpers
the rest of the organization calls, `raster/img` implements it over a
`draw.Image`; `cmd/mdicons` converts a directory of Material SVGs into a Go
package of IVG blobs. It implements FFV0 of the spec, not FFV1.

**Layer.** Outside ADR-001's tier table: a support library, which the rule
binds in one direction only — every tier may import it, and it may import
nothing in the table itself. Its root module imports nothing else in the
organization. Its nested `ivg/raster/gio` module adds `font`, `style` and
`textdraw` — those edges are the nested module's and not the root's.
Imported by `prism`. Outside the tier table, also by the demo modules
`mvu/example` and `prism/gallery` and the workbench applications
`iconbrowser`, `launcher`, `mindchat` and `todos`. Both directions are
measured rather than typed — `scripts/check-layers.sh --edges` reports the
graph and `scripts/sync-agents.sh` renders these sentences from it — so
correcting them here changes nothing.

**Read the canonical guide before you write code against this module.** It is
the organization's one agent guide — the module inventory with current tags,
the application skeleton, the MVU loop and rx semantics, typography, and the
pitfalls that are not guessable. It lives exactly once, in `vibrantgio/.github`,
and this file links it rather than copying it:

    https://raw.githubusercontent.com/vibrantgio/.github/master/llms.txt

**Modules.** `github.com/vibrantgio/ivg` at the repository root, and one
nested module: `raster/gio/` (`github.com/vibrantgio/ivg/raster/gio`).
Nested-module tags carry the directory as a prefix — `raster/gio/v0.1.6`,
not `v0.1.6`.

**Build and test.** From the repository root, and again inside each nested
module directory — `./...` does not cross a module boundary:

    go build ./... && go test ./...
