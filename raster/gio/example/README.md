# Examples of Using IVG with Gio

    import "github.com/reactivego/ivg/raster/gio/example"

[![Go Reference](https://pkg.go.dev/badge/github.com/reactivego/ivg/raster/gio/example.svg)](https://pkg.go.dev/github.com/reactivego/ivg/raster/gio/example#section-documentation)

This package contains example programs demonstrating how to use the IVG (Icon Vector Graphics) package with Gio (gioui.org) - a cross-platform GUI framework for Go. The examples show different ways to render vector graphics and icons in your Gio applications.

## Example: Simple Play Button Icon

![PlayArrow Gio](.assets/playbutton-gio.png?raw=true)

This first example demonstrates the most basic usage - loading and rendering a simple play button icon from an .ivg file that's stored as a byte slice in memory. It serves as a good starting point for understanding how to incorporate vector icons into your Gio application.

```go
// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"

	"github.com/reactivego/ivg/raster/gio"
)

func main() {
	go PlayArrow()
	app.Main()
}

func PlayArrow() {
	window := app.NewWindow(
		app.Title("IVG - PlayArrow"),
		app.Size(768, 768),
	)
	widget, err := gio.Icon(AVPlayArrow, 48, 48, gio.WithColors(Amber400))
	if err != nil {
		log.Fatal(err)
	}
	ops := new(op.Ops)
	for event := range window.Events() {
		if event, ok := event.(system.FrameEvent); ok {
			gtx := layout.NewContext(ops, event)
			widget(gtx)
			event.Frame(ops)
		}
	}
	os.Exit(0)
}

var (
	// From "golang.org/x/exp/shiny/materialdesign/icons"
	AVPlayArrow = []byte{
		0x89, 0x49, 0x56, 0x47, 0x02, 0x0a, 0x00, 0x50,
		0x50, 0xb0, 0xb0, 0xc0, 0x70, 0x64, 0xe9, 0xb8,
		0x20, 0xac, 0x64, 0xe1,
	}
	// From "golang.org/x/exp/shiny/materialdesign/colors"
	Amber400 = color.RGBA{0xff, 0xca, 0x28, 0xff} // rgb(255, 202, 40)
)
```

## Example ActionInfo

![ActionInfo Gio](.assets/actioninfo-gio.png?raw=true)

Generating the .ivg bytes for an icon on the fly and then rendering it. Rendering operations are cached in an icon cache. The function `ActionInfoData()` is called once to programatically generate an .ivg byte slice using the following pipeline:
```
Generator -> Encoder
```
The resulting bytes are stored for later rendering during a system.FrameEvent. When the icon needs to be rendered, call the icon.Cache Render method with the .ivg data bytes and additional arguments.
The icon cache uses the following pipeline to render the icon.

```
Decoder -> Renderer -> Rasterizer
```
The cache stores the resulting op.CallOp keyed on the icon data and parameters used for rendering.

```go
package main

import (
    "image/color"
    "log"
    "os"

    "gioui.org/app"
    "gioui.org/f32"
    "gioui.org/io/system"
    "gioui.org/op"
    "gioui.org/unit"

    "github.com/reactivego/ivg"
    "github.com/reactivego/ivg/encode"
    "github.com/reactivego/ivg/generate"
    "github.com/reactivego/ivg/raster/gio"
)

func main() {
    go ActionInfo()
    app.Main()
}

func ActionInfo() {
    window := app.NewWindow(
        app.Title("IVG - ActionInfo"),
        app.Size(unit.Dp(768), unit.Dp(768)),
    )
    // generate ivg data bytes on the fly for the ActionInfo icon.
    enc := &encode.Encoder{}
    gen := &generate.Generator{Destination: enc}
    gen.Reset(ivg.ViewBox{MinX: 0, MinY: 0, MaxX: 48, MaxY: 48}, ivg.DefaultPalette)
    gen.SetPathData("M24 4C12.95 4 4 12.95 4 24s8.95 20 20 20 20-8.95 "+
        "20-20S35.05 4 24 4zm2 30h-4V22h4v12zm0-16h-4v-4h4v4z", 0, false)
    actionInfoData, err := enc.Bytes()
    if err != nil {
        log.Fatal(err)
    }
    actionInfo, err := gio.NewIcon(actionInfoData)
    if err != nil {
        log.Fatal(err)
    }
    cache := gio.NewIconCache()
    ops := new(op.Ops)
    for next := range window.Events() {
        if frame, ok := next.(system.FrameEvent); ok {
            ops.Reset()
            viewRect := actionInfo.AspectMeet(frame.Size, ivg.Mid, ivg.Mid)
            blue := color.RGBA{0x21, 0x96, 0xf3, 0xff}
            if callOp, err := cache.Rasterize(actionInfo, viewRect, gio.WithColors(blue)); err == nil {
                callOp.Add(ops)
            } else {
                log.Fatal(err)
            }
            frame.Frame(ops)
        }
    }
    os.Exit(0)
}
```
# Example Icons

| ![Icons Gio](.assets/icons-gio.png?raw=true) | ![Icons Vector](.assets/icons-vec.png?raw=true) |
|:---:|:---:|
| Gio | Vec |

The Icons example program takes icons from the package `"golang.org/x/exp/shiny/materialdesign/icons"` and renders them. These icons are just a few layers filled with a single color. The example uses function `Render` from package `"github.com/reactivego/ivg/icon"` for rendering. `Render` uses the following pipeline:
```
Decoder -> Renderer -> Rasterizer
```
By clicking the window you can switch between rendering using the Gio or Vec (`"golang.org/x/image/vector"`) rasterizer.

The rendering using Gio is extremely quick as it only generates a few clipping & blending operations that are put in an operator queue. The actual rendering by Gio takes place on the GPU. For the Vec rasterizer all the pixels of the image need to be pre-generated, which takes relatively a long time.

```go
package main

import (
    "fmt"
    "image"
    "log"
    "os"
    "time"

    "golang.org/x/exp/shiny/materialdesign/colornames"

    "gioui.org/app"
    "gioui.org/f32"
    "gioui.org/io/pointer"
    "gioui.org/io/system"
    "gioui.org/op"
    "gioui.org/op/clip"
    "gioui.org/op/paint"
    "gioui.org/unit"

    "github.com/reactivego/ivg/raster/gio"
)

func main() {
    go Icons()
    app.Main()
}

func Icons() {
    window := app.NewWindow(
        app.Title("IVG - Icons"),
        app.Size(unit.Dp(768), unit.Dp(768)),
    )
    type Backend struct {
        Name   string
        Driver gio.Driver
    }
    backend := Backend{"Gio", gio.Gio}
    Grey300 := color.NRGBAModel.Convert(colornames.Grey300).(color.NRGBA)
    Grey800 := color.NRGBAModel.Convert(colornames.Grey800).(color.NRGBA)
    ops := new(op.Ops)
    backdrop := new(int)
    index := 0
    for next := range window.Events() {
        if frame, ok := next.(system.FrameEvent); ok {
            ops.Reset()

            // clicking on backdrop will switch active renderer
            pointer.InputOp{Tag: backdrop, Types: pointer.Release}.Add(ops)
            for _, next := range frame.Queue.Events(backdrop) {
                if event, ok := next.(pointer.Event); ok {
                    if event.Type == pointer.Release {
                        switch backend.Name {
                        case "Gio":
                            backend = Backend{"Vec", gio.Vec}
                        case "Vec":
                            backend = Backend{"Gio", gio.Gio}
                        }
                    }
                }
            }

            // fill the whole backdrop rectangle
            paint.ColorOp{Color: Grey800}.Add(ops)
            paint.PaintOp{}.Add(ops)

            // device independent content rect calculation
            margin := unit.Dp(12)
            minX := unit.Add(frame.Metric, margin, frame.Insets.Left)
            minY := unit.Add(frame.Metric, margin, frame.Insets.Top)
            maxX := unit.Add(frame.Metric, unit.Px(float32(frame.Size.X)), frame.Insets.Right.Scale(-1), margin.Scale(-1))
            maxY := unit.Add(frame.Metric, unit.Px(float32(frame.Size.Y)), frame.Insets.Bottom.Scale(-1), margin.Scale(-1))
            contentRect := f32.Rect(
                float32(frame.Metric.Px(minX)), float32(frame.Metric.Px(minY)),
                float32(frame.Metric.Px(maxX)), float32(frame.Metric.Px(maxY)))
            contentMin := image.Point(int(contentRect.Min.X),int(contentRect.Min.Y))
            contentSize := image.Point(int(contentRect.Dx()),int(contentRect.Dy()))

            // fill content rect
            paint.ColorOp{Color: Grey300}.Add(ops)
            tstack := op.Offset(contentRect.Min).Push(ops)
            cstack := clip.Rect(image.Rectangle{Max: contentSize}).Push(ops)
            paint.PaintOp{}.Add(ops)
            cstack.Pop()
            tstack.Pop()

            // select next icon and paint
            n := uint(len(IconCollection))
            ico := IconCollection[(uint(index)+n)%n]
            index++
            start := time.Now()
            icon, err := gio.NewIcon(ico.data)
            if err != nil {
                log.Fatal(err)
            }
            viewRect := icon.AspectMeet(contentSize, 0.5, 0.5).Add(contentMin)
            if callOp, err := gio.Rasterize(icon, viewRect, gio.WithColors(colornames.LightBlue600), gio.WithDriver(backend.Driver)); err == nil {
                callOp.Add(ops)
            } else {
                log.Fatal(err)
            }
            msg := fmt.Sprintf("%s (%v)", backend.Name, time.Since(start).Round(time.Microsecond))
            H5 := Style(H5, WithMetric(frame.Metric))
            PrintText(msg, contentRect.Min, 0.0, 0.0, contentRect.Dx(), H5, ops)

            at := time.Now().Add(500 * time.Millisecond)
            op.InvalidateOp{At: at}.Add(ops)
            frame.Frame(ops)
        }
    }
    os.Exit(0)
}
```

# Example Favicon

| ![Favicon Gio](.assets/favicon-gio.png?raw=true) | ![Favicon Vector](.assets/favicon-vec.png?raw=true) |
|:---:|:---:|
| Gio | Vec |

Favicon programatically renders a vector image of a Gopher using multiple layers with translucency. It uses the following pipeline:
```
Generator -> Renderer -> Rasterizer
```
This example hooks up the generator directly to the renderer and forgoes the `Encoder -> Decoder` stages. The rendering using Gio is extremely quick as it only generates a few clipping & blending operations that are put in an operator queue. The actual rendering by Gio takes place on the GPU. For the Vec (`"golang.org/x/image/vector"`) rasterizer all the pixels of the image need to be pre-generated, which takes relatively a long time.

The resulting images are a little bit different. Look under the nose of the Gopher. Gio produces lighter results than Vec. The reason for this is that Vec performs blending in sRGB space whereas Gio performs the blending in Linear space.

# Example Cowbell

| ![Cowbell Gio](.assets/cowbell-gio.png?raw=true) | ![Cowbell Vector](.assets/cowbell-vec.png?raw=true) |
|:---:|:---:|
| Gio | Vec |

Cowbell programatically renders a vector image of a Cowbell using multiple layers with gradients and translucency. It uses the following pipeline:
```
Generator -> Renderer -> Rasterizer
```
The rendering takes relatively long because the gradients need to be pre-generated even for the Gio renderer.

# Example Gradients

| ![Gradients Gio](.assets/gradients-gio.png?raw=true) | ![Gradients Vector](.assets/gradients-vec.png?raw=true) |
|:---:|:---:|
| Gio | Vec |

Gradients uses the following pipeline to programatically render a vector image consisting of multiple different gradients:
```
Generator -> Renderer -> Rasterizer
```
The rendering takes relatively long because the gradients need to be pre-generated even for the Gio renderer. But even considering that, Gio is approximately 8 times faster than Vec.

## Acknowledgement

The code in this package is based on [golang.org/x/exp/shiny/iconvg](https://github.com/golang/exp/tree/master/shiny/iconvg).

The specification of the IconVG format has recently been moved to a separate repository [github.com/google/iconvg](https://github.com/google/iconvg).

## License

All the code in this package is either Unlicense OR MIT (whichever you prefer). See file [raster/LICENSE](raster/LICENSE).
