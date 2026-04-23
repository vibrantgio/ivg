// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"fmt"
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"

	"github.com/vibrantgio/ivg"
	"github.com/vibrantgio/ivg/decode"
	"github.com/vibrantgio/ivg/encode"
	"github.com/vibrantgio/ivg/generate"
	raster "github.com/vibrantgio/ivg/raster/gio"

	"golang.org/x/exp/shiny/materialdesign/colornames"
)

const shouldDisassemble = false

func main() {
	go Logo()
	app.Main()
}

func Logo() {
	window := app.NewWindow(
		app.Title("IVG - Logo"),
		app.Size(768, 768),
	)

	data, err := LogoIVG()
	if err != nil {
		log.Fatal(err)
	}

	if shouldDisassemble {
		dis, err := decode.Disassemble(data)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(dis))
	}

	var widget layout.Widget
	if true {
		widget, _ = raster.Widget(data, 40, 40, raster.WithColors(colornames.LightBlue600, colornames.Yellow100))
	} else {
		widget, _ = raster.Widget(data, 40, 40)
	}

	ops := new(op.Ops)
	for next := range window.Events() {
		if event, ok := next.(system.FrameEvent); ok {
			gtx := layout.NewContext(ops, event)
			layout.UniformInset(24).Layout(gtx, widget)
			event.Frame(ops)
		}
	}
	os.Exit(0)
}

func LogoIVG() ([]byte, error) {
	// generate ivg data bytes on the fly for the logo.
	enc := &encode.Encoder{}
	// dlog := &ivg.DestinationLogger{Destination: enc}
	dlog := enc
	gen := &generate.Generator{Destination: dlog}

	// Palette that can be referenced from CReg array, gets overidden with colors from by externally set palette.
	pal := ivg.DefaultPalette
	pal[0] = color.RGBA{0x74, 0xAA, 0x9C, 0xff}
	pal[1] = color.RGBAModel.Convert(color.White).(color.RGBA)
	pal[2] = colornames.DeepOrange500
	pal[3] = colornames.Amber500

	gen.Reset(ivg.ViewBox{MinX: 0, MinY: 0, MaxX: 2406, MaxY: 2406}, pal)

	gen.SetCReg(0, true, ivg.PaletteIndexColor(1)) // CReg[0] => palette[1] Yellow100
	gen.SetCReg(0, true, ivg.PaletteIndexColor(0)) // CReg[1] => palette[0] LightBlue600
	gen.SetCReg(0, true, ivg.PaletteIndexColor(3)) // CReg[2] => palette[3] Amber500
	gen.SetCReg(0, true, ivg.PaletteIndexColor(2)) // CReg[3] => palette[2] DeepOrange500

	gen.SetCSel(1)
	// gen.SetCSel(3)

	// Shield for CSel == 1 uses adj == 0: color CReg[CSel-adj == 1] => palette[0] LightBlue600
	// Shield for CSel == 3 uses adj == 0: color CReg[CSel-adj == 3] => palette[2] DeepOrange500
	gen.SetPathData("M1 578.4C1 259.5 259.5 1 578.4 1h1249.1c319 0 577.5 258.5 577.5 577.4V2406H578.4C259.5 2406 1 2147.5 1 1828.6V578.4z", 0)

	// Logo for CSel == 1 uses adj == 1: color CReg[CSel-adj == 0] => palette[1] Yellow100
	// Logo for CSel == 3 uses adj == 1: color CReg[CSel-adj == 2] => palette[3] Amber500
	gen.SetPathData("M1107.3 299.1c-198 0-373.9 127.3-435.2 315.3C544.8 640.6 434.9 720.2 370.5 833c-99.3 171.4-76.6 386.9 56.4 533.8-41.1 123.1-27 257.7 38.6 369.2 98.7 172 297.3 260.2 491.6 219.2 86.1 97 209.8 152.3 339.6 151.8 198 0 373.9-127.3 435.3-315.3 127.5-26.3 237.2-105.9 301-218.5 99.9-171.4 77.2-386.9-55.8-533.9v-.6c41.1-123.1 27-257.8-38.6-369.8-98.7-171.4-297.3-259.6-491-218.6-86.6-96.8-210.5-151.8-340.3-151.2zm0 117.5-.6.6c79.7 0 156.3 27.5 217.6 78.4-2.5 1.2-7.4 4.3-11 6.1L952.8 709.3c-18.4 10.4-29.4 30-29.4 51.4V1248l-155.1-89.4V755.8c-.1-187.1 151.6-338.9 339-339.2zm434.2 141.9c121.6-.2 234 64.5 294.7 169.8 39.2 68.6 53.9 148.8 40.4 226.5-2.5-1.8-7.3-4.3-10.4-6.1l-360.4-208.2c-18.4-10.4-41-10.4-59.4 0L1024 984.2V805.4L1372.7 604c51.3-29.7 109.5-45.4 168.8-45.5zM650 743.5v427.9c0 21.4 11 40.4 29.4 51.4l421.7 243-155.7 90L597.2 1355c-162-93.8-217.4-300.9-123.8-462.8C513.1 823.6 575.5 771 650 743.5zm807.9 106 348.8 200.8c162.5 93.7 217.6 300.6 123.8 462.8l.6.6c-39.8 68.6-102.4 121.2-176.5 148.2v-428c0-21.4-11-41-29.4-51.4l-422.3-243.7 155-89.3zM1201.7 997l177.8 102.8v205.1l-177.8 102.8-177.8-102.8v-205.1L1201.7 997zm279.5 161.6 155.1 89.4v402.2c0 187.3-152 339.2-339 339.2v-.6c-79.1 0-156.3-27.6-217-78.4 2.5-1.2 8-4.3 11-6.1l360.4-207.5c18.4-10.4 30-30 29.4-51.4l.1-486.8zM1380 1421.9v178.8l-348.8 200.8c-162.5 93.1-369.6 38-463.4-123.7h.6c-39.8-68-54-148.8-40.5-226.5 2.5 1.8 7.4 4.3 10.4 6.1l360.4 208.2c18.4 10.4 41 10.4 59.4 0l421.9-243.7z", 1)

	return enc.Bytes()
}
