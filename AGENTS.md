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

**Layer.** Outside ADR-001's tier table: a support library the design
system consumes and never depends on. `raster/gio` is the package
everything actually calls: prism requires it for its `icon/gallery` and
`gallery` demos, `mvu/example` for four of its examples, and the workbench
applications `todos`, `iconbrowser`, `launcher` and `mindchat` for their
icons — mindchat also encodes IVG at run time through `ivg`, `encode` and
`generate`. Neither of ivg's own modules imports the design system; four of
the eight demo programs under `raster/gio/example/` do, and only the tier-0
leaves `style` and `textdraw`.

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
