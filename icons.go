package main

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"math"
)

const (
	iconSize   = 22
	iconScale  = 4
	iconPad    = 2
	iconBorder = 1
)

var (
	colorWork  = color.RGBA{220, 60, 50, 255}
	colorBreak = color.RGBA{50, 180, 90, 255}
	colorEmpty = color.RGBA{60, 60, 60, 255}
)

func renderIcon(progress float64, isWork bool) []byte {
	fill := colorWork
	if !isWork {
		fill = colorBreak
	}
	return renderPie(progress, fill, colorEmpty, color.RGBA{255, 255, 255, 255})
}

func renderTemplateIcon(progress float64) []byte {
	return renderPie(progress, color.RGBA{0, 0, 0, 255}, color.RGBA{0, 0, 0, 80}, color.RGBA{0, 0, 0, 255})
}

func renderPie(progress float64, fill, empty, border color.RGBA) []byte {
	hi := iconSize * iconScale
	img := image.NewRGBA(image.Rect(0, 0, hi, hi))

	cx := float64(hi) / 2
	cy := float64(hi) / 2
	outerR := float64(hi)/2 - float64(iconPad*iconScale)
	innerR := outerR - float64(iconBorder*iconScale)

	limit := progress * 2 * math.Pi

	for y := range hi {
		for x := range hi {
			dx := float64(x) + 0.5 - cx
			dy := float64(y) + 0.5 - cy
			dist2 := dx*dx + dy*dy
			if dist2 > outerR*outerR {
				continue
			}
			if dist2 > innerR*innerR {
				img.SetRGBA(x, y, border)
				continue
			}
			angle := math.Atan2(dy, dx) + math.Pi/2
			if angle < 0 {
				angle += 2 * math.Pi
			}
			if angle < limit {
				img.SetRGBA(x, y, fill)
			} else {
				img.SetRGBA(x, y, empty)
			}
		}
	}

	out := image.NewRGBA(image.Rect(0, 0, iconSize, iconSize))
	n := iconScale * iconScale
	for y := range iconSize {
		for x := range iconSize {
			var rv, gv, bv, av int
			for dy := range iconScale {
				for dx := range iconScale {
					c := img.RGBAAt(x*iconScale+dx, y*iconScale+dy)
					rv += int(c.R)
					gv += int(c.G)
					bv += int(c.B)
					av += int(c.A)
				}
			}
			out.SetRGBA(x, y, color.RGBA{
				R: uint8(rv / n),
				G: uint8(gv / n),
				B: uint8(bv / n),
				A: uint8(av / n),
			})
		}
	}

	var buf bytes.Buffer
	_ = png.Encode(&buf, out)
	return buf.Bytes()
}
