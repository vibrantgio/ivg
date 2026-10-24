# Examples of Using IVG with Gio

The example programs of the module `github.com/vibrantgio/ivg/raster/gio`.
Each directory is one `package main`; run one from the module's directory
with `go run ./example/<directory>`.

| Directory   | Shows |
| ---         | --- |
| `arrow`     | Rendering an icon from `.ivg` bytes held in a byte slice with `raster.Widget`. |
| `info`      | Generating the `.ivg` bytes of an icon in the program, then rendering them with `raster.Widget`. |
| `icons`     | The icons of `golang.org/x/exp/shiny/materialdesign/icons`, switching between the Gio and the image backend by clicking the window. |
| `favicon`   | A vector image built with the generator in several translucent layers, rendered with both backends. |
| `cowbell`   | A vector image built with the generator in layers of gradients and translucency, rendered with both backends. |
| `gradients` | A vector image of several different gradients, rendered with both backends. |
| `logo`      | A logo built with the generator and drawn with a chosen palette. |
| `blend`     | A translucent highlight blended over an opaque background, filled once with Gio paint operations and once with `golang.org/x/image/vector`. |

![Arrow Gio](.assets/playbutton-gio.png?raw=true)

| ![Icons Gio](.assets/icons-gio.png?raw=true) | ![Icons Image](.assets/icons-vec.png?raw=true) |
|:---:|:---:|
| Gio | Image |

## The window call

Every program opens its window the same way and runs the same event loop:

```go
window := new(app.Window)
window.Option(
	app.Title("IVG - Arrow"),
	app.Size(768, 768))
ops := new(op.Ops)
for {
	switch e := window.Event().(type) {
	case app.DestroyEvent:
		os.Exit(0)
	case app.FrameEvent:
		gtx := app.NewContext(ops, e)
		w(gtx)
		e.Frame(gtx.Ops)
	}
}
```

`w` is the `layout.Widget` that `raster.Widget(data, width, height,
options...)` returns, with `raster` standing for the import
`github.com/vibrantgio/ivg/raster/gio`.

## License

All the code in this package is either Unlicense OR MIT (whichever you
prefer). See file [LICENSE](../../LICENSE).
