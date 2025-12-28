package graph

import (
	"image"
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

var width float32 = 800
var height float32 = 800

// Colors
var (
	blue  = color.NRGBA{R: 0x40, G: 0x40, B: 0xC0, A: 0xFF}
	red   = color.NRGBA{R: 0xC0, G: 0x40, B: 0x40, A: 0xFF}
	white = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	grey  = color.NRGBA{R: 0xC0, G: 0xC0, B: 0xC0, A: 0xFF}
)

// Sizes
var (
	buttonSize = image.Pt(200, 90)
)

// Intializes some default settings for the GUI window
func (g Graph[T]) guiInitialize(window *app.Window) {
	window.Option(app.Size(unit.Dp(width), unit.Dp(height)))
}

func colorBox(gtx layout.Context, size image.Point, color color.NRGBA) layout.Dimensions {
	defer clip.Rect{Max: size}.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: color}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return layout.Dimensions{Size: size}
}

func colorRoundedPill(gtx layout.Context, size image.Point, color color.NRGBA) layout.Dimensions {
	rect := image.Rectangle{Max: size}
	rad := size.Y / 2
	paint.FillShape(gtx.Ops, color, clip.UniformRRect(rect, rad).Op(gtx.Ops))
	return layout.Dimensions{Size: size}
}

func drawListModal[T any](gtx layout.Context, list []T) layout.Dimensions {
	panic("Oops")
}

type button struct {
	pressed      bool
	theme        *material.Theme
	color        color.NRGBA
	pressedColor color.NRGBA
	onPress      func()
}

func (b *button) Layout(gtx layout.Context, label string) layout.Dimensions {
	clip.Rect{Max: buttonSize}.Push(gtx.Ops).Pop()
	// Handle Input
	{
		event.Op(gtx.Ops, b)
		pointer.CursorColResize.Add(gtx.Ops)
		for {
			ev, ok := gtx.Event(pointer.Filter{
				Target: b,
				Kinds:  pointer.Press | pointer.Drag | pointer.Release | pointer.Cancel,
			})
			if !ok {
				break
			}
			e, ok := ev.(pointer.Event)
			if !ok {
				continue
			}
			switch e.Kind {
			case pointer.Press:
				println("I got pressed!")
			}
		}
	}
	if b.pressed {
		b.color = b.pressedColor
	}
	return layout.UniformInset(unit.Dp(15)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Stack{Alignment: layout.Center}.Layout(gtx, layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return colorRoundedPill(gtx, buttonSize, b.color)
		}), layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			label := material.Label(b.theme, unit.Sp(16), label)
			label.Color = white
			return label.Layout(gtx)
		}))
	})
}

// Draws the menu portion of the screen
func (g Graph[T]) drawMenu(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	onPressAllNodes := func() {
		nodes := g.GetNodes()
		drawListModal(gtx, nodes)
	}
	onPressAllEdges := func() {
		edges := g.edges
		drawListModal(gtx, edges)
	}
	return layout.Stack{Alignment: layout.NW}.Layout(gtx, layout.Stacked(func(gtx layout.Context) layout.Dimensions {
		return colorBox(gtx, gtx.Constraints.Max, grey)
	}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{}.Layout(gtx, layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
				button := button{pressed: false, theme: theme, color: blue, onPress: onPressAllNodes, pressedColor: red}
				return button.Layout(gtx, "View All Nodes")
			}),
				layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
					button := button{pressed: false, theme: theme, color: blue, onPress: onPressAllEdges, pressedColor: red}
					return button.Layout(gtx, "View All Edges")
				}),
			)
		}))
}

func (g Graph[T]) drawLayout(gtx layout.Context, theme *material.Theme) {
	leftMenu := layout.Flexed(0.3, func(gtx layout.Context) layout.Dimensions {
		return g.drawMenu(gtx, theme)
	})
	mainView := layout.Flexed(0.7, func(gtx layout.Context) layout.Dimensions {
		return colorBox(gtx, gtx.Constraints.Max, white)
	})
	layout.Flex{Axis: layout.Horizontal}.Layout(gtx, mainView, leftMenu)
}

// Main event loop handling
func (g Graph[T]) handleFrameEvent(event app.FrameEvent, ops op.Ops, theme *material.Theme) {
	gtx := app.NewContext(&ops, event)
	// Initialize the layout
	g.drawLayout(gtx, theme)
	// Finally draw the frame event onto the window.
	event.Frame(gtx.Ops)
}

func (g Graph[T]) handleDestroyEvent(event app.DestroyEvent) error {
	return event.Err
}

// Renders an interactive GUI based off the current graph. Useful for visualization and debugging purposes. Must be run
// in the main thread.
func (g Graph[T]) GUI() {
	var run func(*app.Window) error
	run = func(w *app.Window) error {
		ops := op.Ops{}
		theme := material.NewTheme()
		g.guiInitialize(w)
		for {
			switch e := w.Event().(type) {
			case app.DestroyEvent:
				return g.handleDestroyEvent(e)
			case app.FrameEvent:
				g.handleFrameEvent(e, ops, theme)
			}
		}
	}
	go func() {
		window := new(app.Window)
		err := run(window)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}
