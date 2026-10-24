# ivg

    import "github.com/vibrantgio/ivg"

[![Go Reference](https://pkg.go.dev/badge/github.com/vibrantgio/ivg.svg)](https://pkg.go.dev/github.com/vibrantgio/ivg#section-documentation)

Package `ivg` renders [IconVG](https://github.com/google/iconvg) icons through a vector `Rasterizer` interface. IconVG is a compact binary format for vector graphic icons.

The code is based on [golang.org/x/exp/shiny/iconvg](https://github.com/golang/exp/tree/master/shiny/iconvg), which renders to a bitmap. Here an icon is drawn through a `Rasterizer`, so the output target is chosen by the rasterizer: an image through `raster/img`, or Gio operations through `raster/gio`.

The package is named `ivg` rather than `iconvg` to keep it clear of the original implementation's namespace.

## File Format Versions

The IconVG format was simplified to version 1 (FFV1) so that a later version can support animation; the first format is called FFV0. FFV1 represents static vector graphics icons, and a future FFV2 is to represent animated ones.

The rationale is discussed in a GitHub proposal: [File Format Versions 1, 2 and Beyond](https://github.com/google/iconvg/issues/4#issue-905297018)

The spec of each File Format Version:
- [IconVG FFV0](spec/iconvg-spec-v0.md)
- [IconVG FFV1](https://github.com/google/iconvg/blob/97b0c08e6e298f5f3606f79f3fb38cc0d64d3198/spec/iconvg-spec.md)

> NOTE: This package implements the [FFV0](spec/iconvg-spec-v0.md) version of the IconVG format.

## Packages

| Package | What it has |
|---|---|
| `ivg` | The `Destination` interface, `ViewBox`, `Color` and the palette and gradient encoding shared by the other packages. |
| `decode` | `Decode`, which reads an IconVG blob and calls the methods of a `Destination`; `DecodeViewBox`; `Disassemble`; the options `WithPalette` and `WithColorAt`. |
| `encode` | `Encoder`, a `Destination` that writes an IconVG blob; `Bytes` returns it. |
| `generate` | `Generator`, which takes SVG-style path data (`SetPathData`), gradients (`SetLinearGradient`, `SetCircularGradient`, `SetEllipticalGradient`, `SetGradient`) and transforms (`SetTransform`) and calls the methods of the `Destination` it is given. |
| `render` | `Renderer`, a `Destination` that draws onto a `raster.Rasterizer`; `Gradient`, the image a gradient is drawn from, which implements `raster.GradientConfig`. |
| `raster` | The `Rasterizer` and `GradientConfig` interfaces. |
| `raster/img` | `Rasterizer`, which draws into a `draw.Image` with `golang.org/x/image/vector`; `NewRasterizer`. |
| `raster/gio` | A nested module, `github.com/vibrantgio/ivg/raster/gio`: `Widget` returns a Gio `layout.Widget` for an IconVG blob, scaled to fit and centred; `GioPaint` paints with Gio clip paths, `ImagePaint` through an image (chosen with `WithImageBackend`); `WithColors` sets palette colours; `NewRasterizer` returns a `Rasterizer` that records Gio operations. |
| `mdicons` | A parser for the Material Design icon sets as Google publishes them: `Parse` reads a Material icons folder structure of SVG files and writes Go source holding the icons as IconVG data. |
| `cmd/disivg` | A command that disassembles an IconVG file into its opcodes. |
| `cmd/mdicons` | A command that converts a Material icons folder structure of SVG files into a Go package of IconVG data, through `mdicons.Parse`. |

## Code Organization

The packages combine into different IconVG render pipelines. The `Destination` interface is implemented both by the `Encoder` in package `encode` and by the `Renderer` in package `render`. The `Generator` in package `generate` only uses a `Destination`, so the same calls either produce a data blob or draw directly to a `Rasterizer` through the `Renderer`.

```go
// Destination handles the actions decoded from an IconVG graphic's opcodes.
//
// When passed to Decode, the first method called (if any) will be Reset. No
// methods will be called at all if an error is encountered in the encoded form
// before the metadata is fully decoded.
type Destination interface {
	Reset(viewbox ViewBox, palette [64]color.RGBA)
	CSel() uint8
	SetCSel(cSel uint8)
	NSel() uint8
	SetNSel(nSel uint8)
	SetCReg(adj uint8, incr bool, c Color)
	SetNReg(adj uint8, incr bool, f float32)
	SetLOD(lod0, lod1 float32)

	StartPath(adj uint8, x, y float32)
	ClosePathEndPath()
	ClosePathAbsMoveTo(x, y float32)
	ClosePathRelMoveTo(x, y float32)

	AbsHLineTo(x float32)
	RelHLineTo(x float32)
	AbsVLineTo(y float32)
	RelVLineTo(y float32)
	AbsLineTo(x, y float32)
	RelLineTo(x, y float32)
	AbsSmoothQuadTo(x, y float32)
	RelSmoothQuadTo(x, y float32)
	AbsQuadTo(x1, y1, x, y float32)
	RelQuadTo(x1, y1, x, y float32)
	AbsSmoothCubeTo(x2, y2, x, y float32)
	RelSmoothCubeTo(x2, y2, x, y float32)
	AbsCubeTo(x1, y1, x2, y2, x, y float32)
	RelCubeTo(x1, y1, x2, y2, x, y float32)
	AbsArcTo(rx, ry, xAxisRotation float32, largeArc, sweep bool, x, y float32)
	RelArcTo(rx, ry, xAxisRotation float32, largeArc, sweep bool, x, y float32)
}
```

A parser of SVG files reads the SVG and calls methods on a `Destination` to produce a binary data blob. For the Material Design icons:

```
mdicons.ParsePathData -> [Destination] encode.Encoder -> []byte
```

For SVGs with gradients and transforms, a program calls the `Generator`, which is hooked up to a `Destination` to produce the binary data blob:

```
generate.Generator -> [Destination] encode.Encoder -> []byte
```

To render the icon, the binary data blob is passed to `decode.Decode`, which calls methods on a `Renderer` hooked up to a `Rasterizer`.

```go
// Rasterizer is a 2-D vector graphics rasterizer.
type Rasterizer interface {
	// Reset resets a Rasterizer as if it was just returned by NewRasterizer.
	// This includes setting z.DrawOp to draw.Over.
	Reset(w, h int)
	// Size returns the width and height passed to NewRasterizer or Reset.
	Size() image.Point
	// Bounds returns the rectangle from (0, 0) to the width and height passed to
	// Reset.
	Bounds() image.Rectangle
	// Pen returns the location of the path-drawing pen: the last argument to the
	// most recent XxxTo call.
	Pen() (x, y float32)
	// MoveTo starts a new path and moves the pen to (ax, ay). The coordinates
	// are allowed to be out of the Rasterizer's bounds.
	MoveTo(ax, ay float32)
	// LineTo adds a line segment, from the pen to (bx, by), and moves the pen to
	// (bx, by). The coordinates are allowed to be out of the Rasterizer's
	// bounds.
	LineTo(bx, by float32)
	// QuadTo adds a quadratic Bézier segment, from the pen via (bx, by) to (cx,
	// cy), and moves the pen to (cx, cy). The coordinates are allowed to be out
	// of the Rasterizer's bounds.
	QuadTo(bx, by, cx, cy float32)
	// CubeTo adds a cubic Bézier segment, from the pen via (bx, by) and (cx, cy)
	// to (dx, dy), and moves the pen to (dx, dy). The coordinates are allowed to
	// be out of the Rasterizer's bounds.
	CubeTo(bx, by, cx, cy, dx, dy float32)
	// ClosePath closes the current path.
	ClosePath()
	// Draw aligns r.Min in z with sp in src and then replaces the rectangle r in
	// z with the result of a Porter-Duff composition. The vector paths added by
	// the XxxTo calls become the mask for drawing src onto z.
	Draw(r image.Rectangle, src image.Image, sp image.Point)
}
```

Decoding a blob and rendering it to a `Rasterizer`:

```
[]byte -> decode.Decode -> [Destination] render.Renderer -> [Rasterizer] img.Rasterizer
```

The `Generator` can also be hooked up to the `Renderer` directly, which skips encoding and decoding:

```
generate.Generator -> [Destination] render.Renderer -> [Rasterizer] img.Rasterizer
```

## Examples

The directory `raster/gio/example` has one program per directory:

- `arrow` renders an icon held as IconVG bytes.
- `info` generates an icon with the `Generator`, encodes it and renders it.
- `blend` paints a translucent colour over an opaque one twice, once with Gio paint operations and once with `golang.org/x/image/vector`, to show that both blend alike.
- `logo` generates a logo with the `Generator` and renders it with two palette colours.
- The following examples switch between the Gio and the image backend when the image is clicked, to show the rendering and speed differences:
    - `icons` renders the icons of `golang.org/x/exp/shiny/materialdesign/icons`.
    - `favicon` a vector image with several blended layers.
    - `cowbell` a vector image with several blended layers including gradients.
    - `gradients` a vector image with many gradients.

## For coding assistants

Read the org guide before you write code against this module: the plan
root's [`AGENTS.md`](https://github.com/vibrantgio/organisation/blob/master/AGENTS.md).
The build and test commands are reached from that file.

## Acknowledgement

The code in this package is based on [golang.org/x/exp/shiny/iconvg](https://github.com/golang/exp/tree/master/shiny/iconvg).

The specification of the IconVG format is in the repository [github.com/google/iconvg](https://github.com/google/iconvg).

## License

Everything under the raster folder is Unlicense OR MIT (whichever you prefer). See file [raster/LICENSE](raster/LICENSE).

All the other code is governed by a BSD-style license that can be found in the [LICENSE](LICENSE) file.
