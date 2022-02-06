// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"fmt"
	"image"
	"image/color"
	"sort"

	"gioui.org/f32" // f32 is used for shape calculations.

	// system is used for system events (e.g. closing the window).
	"gioui.org/layout" // layout is used for layouting widgets.
	// op is used for recording different operations.
	"gioui.org/op/clip"  // clip is used to draw the cell shape.
	"gioui.org/op/paint" // paint is used to paint the cells.
)

// BoardStyle draws Board with rectangles.
type BoardStyle struct {
	CellSizePx int
}

func (board BoardStyle) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Flex{}.Layout(gtx,
		// layout.Flexed(.5, func(gtx layout.Context) layout.Dimensions {
		// 	return layout.Stack{}.Layout(gtx,
		// 		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
		// 			return board.layoutFilled(gtx, false)
		// 		}),
		// 		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
		// 			return board.layoutStroke(gtx, false)
		// 		}),
		// 	)
		// }),
		layout.Flexed(.5, func(gtx layout.Context) layout.Dimensions {
			return layout.Stack{}.Layout(gtx,
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return board.layoutFilled(gtx, true)
				}),
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return board.layoutStroke(gtx, true)
				}),
			)
		}),
	)
}

func (board BoardStyle) layoutStroke(gtx layout.Context, isSell bool) layout.Dimensions {
	sizeX := gtx.Constraints.Max.X
	if isSell {
		sizeX = 0
	}
	size := image.Point{X: sizeX, Y: 1000}
	var stroke clip.Path
	// var length float32 = 100
	// var depth float32 = 20

	var x float32 = float32(size.X)
	var y float32 = float32(size.Y)
	// fmt.Println(size.X, size.Y)
	stroke.Begin(gtx.Ops)
	// 9650.00000003 800
	fmt.Println("...............................", x, y)
	sort.Slice(rawBids, func(i, j int) bool {
		return rawBids[i][0] > rawBids[j][0]
	})

	for i := 0; i < len(rawBids); i++ {
		// fmt.Println(bids[i].price, bids[i].amount)

		var length float32 = 0
		var depth float32 = rawBids[i][1] / 3
		if i < len(rawBids)-1 {
			length = (rawBids[i][0] - rawBids[i+1][0]) * 80000 // price
		}
		// fmt.Println(length, depth, rawBids[i][0], rawBids[i][1])

		stroke.MoveTo(f32.Pt(x, y))
		y -= depth
		stroke.LineTo(f32.Pt(x, y))
		stroke.Close()
		stroke.MoveTo(f32.Pt(x, y))
		if isSell {
			x += length
		} else {
			x -= length
		}
		stroke.LineTo(f32.Pt(x, y))
		stroke.Close()
	}
	if isSell {
		stroke.MoveTo(f32.Pt(x, y))
		stroke.LineTo(f32.Pt(float32(gtx.Constraints.Max.X), y))
	} else {
		stroke.MoveTo(f32.Pt(x, y))
		stroke.LineTo(f32.Pt(0, y))
	}
	stroke.Close()

	// for i := 0; i < 10; i++ {
	// 	stroke.MoveTo(f32.Pt(x, y))
	// 	y -= depth
	// 	stroke.LineTo(f32.Pt(x, y))
	// 	stroke.Close()
	// 	stroke.MoveTo(f32.Pt(x, y))
	// 	if isSell {
	// 		x += length
	// 	} else {
	// 		x -= length
	// 	}
	// 	stroke.LineTo(f32.Pt(x, y))
	// 	stroke.Close()
	// }

	clip.Stroke{
		Path:  stroke.End(),
		Width: 1,
	}.Op().Push(gtx.Ops)

	var color color.NRGBA
	if isSell {
		sizeX = 0
		color = rgb(0xed6d47)
	} else {
		sizeX = gtx.Constraints.Max.X
		color = rgb(0x41bf53)
	}

	paint.Fill(gtx.Ops, color)

	return layout.Dimensions{Size: size}
}

func (board BoardStyle) layoutFilled(gtx layout.Context, isSell bool) layout.Dimensions {
	var sizeX int
	var color color.NRGBA
	if isSell {
		sizeX = 0
		color = rgb(0xed6d47)
	} else {
		sizeX = gtx.Constraints.Max.X
		color = rgb(0xE1F8EF)
	}

	size := image.Point{X: sizeX, Y: 1000}
	var filled clip.Path

	// var length float32 = 100
	// var depth float32 = 20
	var x float32 = float32(size.X)
	var y float32 = float32(size.Y)
	filled.Begin(gtx.Ops)

	sort.Slice(rawBids, func(i, j int) bool {
		return rawBids[i][0] > rawBids[j][0]
	})
	var nextX float32
	for i := 0; i < len(rawBids); i++ {
		var length float32 = 0
		var depth float32 = rawBids[i][1] / 3

		if i < len(rawBids)-1 {
			length = (rawBids[i][0] - rawBids[i+1][0]) * 80000 // price
		}
		fmt.Println(length, depth, rawBids[i][0], rawBids[i][1])

		if isSell {
			nextX = x + length
		} else {
			nextX = x - length
		}
		filled.MoveTo(f32.Pt(x, float32(size.Y)))
		filled.LineTo(f32.Pt(nextX, float32(size.Y)))
		y -= depth
		filled.LineTo(f32.Pt(nextX, y))
		filled.LineTo(f32.Pt(x, y))
		filled.Close()
		if isSell {
			x += length
		} else {
			x -= length
		}
	}

	// Fill the rest of the chart
	if isSell {
		filled.MoveTo(f32.Pt(x, float32(size.Y)))
		filled.LineTo(f32.Pt(float32(gtx.Constraints.Max.X), float32(size.Y)))
		filled.LineTo(f32.Pt(float32(gtx.Constraints.Max.X), y))
		filled.LineTo(f32.Pt(x, y))
	} else {
		filled.MoveTo(f32.Pt(x, float32(size.Y)))
		filled.LineTo(f32.Pt(0, float32(size.Y)))
		filled.LineTo(f32.Pt(0, y))
		filled.LineTo(f32.Pt(x, y))
	}
	filled.Close()

	// for i := 0; i < 10; i++ {
	// 	var nextX float32
	// 	if isSell {
	// 		nextX = x + length
	// 	} else {
	// 		nextX = x - length
	// 	}
	// 	filled.MoveTo(f32.Pt(x, float32(size.Y)))
	// 	filled.LineTo(f32.Pt(nextX, float32(size.Y)))
	// 	y -= depth
	// 	filled.LineTo(f32.Pt(nextX, y))
	// 	filled.LineTo(f32.Pt(x, y))
	// 	filled.Close()
	// 	if isSell {
	// 		x += length
	// 	} else {
	// 		x -= length
	// 	}
	// }

	clip.Outline{Path: filled.End()}.Op().Push(gtx.Ops)
	paint.ColorOp{Color: color}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	return layout.Dimensions{Size: size}
}

func rgb(c uint32) color.NRGBA {
	return argb(0xff000000 | c)
}

func argb(c uint32) color.NRGBA {
	return color.NRGBA{A: uint8(c >> 24), R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c)}
}
