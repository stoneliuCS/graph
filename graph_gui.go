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
	buttonSize          = image.Pt(200, 90)
	buttonSizeClickable = buttonSize.Mul(2)
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
	return colorBox(gtx, gtx.Constraints.Max, blue)
}

// Hold all available buttons to press in the gui.
type widgets struct {
	nodeButton *button
	edgeButton *button
}

type button struct {
	pressed bool
	theme   *material.Theme
	onPress func(gtx layout.Context)
	color   color.NRGBA
}

func (b *button) Layout(gtx layout.Context, label string) layout.Dimensions {
	defer clip.Rect{Max: buttonSizeClickable}.Push(gtx.Ops).Pop()
	// Handle Input
	{
		event.Op(gtx.Ops, b)
		for {
			ev, ok := gtx.Event(pointer.Filter{
				Target: b,
				Kinds:  pointer.Press | pointer.Cancel | pointer.Release,
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
				b.pressed = true
			case pointer.Release:
				b.pressed = false
			}
		}
	}
	if b.pressed {
		b.color = red
		b.onPress(gtx)
	} else {
		b.color = blue
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
func (g Graph[T]) drawMenu(gtx layout.Context, widgets *widgets) layout.Dimensions {
	return layout.Stack{Alignment: layout.NW}.Layout(gtx, layout.Stacked(func(gtx layout.Context) layout.Dimensions {
		return colorBox(gtx, gtx.Constraints.Max, grey)
	}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{}.Layout(gtx, layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
				return widgets.nodeButton.Layout(gtx, "View All Nodes")
			}),
				layout.Flexed(0.5, func(gtx layout.Context) layout.Dimensions {
					return widgets.edgeButton.Layout(gtx, "View All Edges")
				}),
			)
		}))
}

func (g Graph[T]) drawLayout(gtx layout.Context, widgets *widgets) {
	leftMenu := layout.Flexed(0.3, func(gtx layout.Context) layout.Dimensions {
		return g.drawMenu(gtx, widgets)
	})
	mainView := layout.Flexed(0.7, func(gtx layout.Context) layout.Dimensions {
		return colorBox(gtx, gtx.Constraints.Max, white)
	})
	layout.Flex{Axis: layout.Horizontal}.Layout(gtx, mainView, leftMenu)
}

// Main event loop handling
func (g Graph[T]) handleFrameEvent(event app.FrameEvent, ops op.Ops, widgets *widgets) {
	gtx := app.NewContext(&ops, event)
	// Initialize the layout
	g.drawLayout(gtx, widgets)
	// Finally draw the frame event onto the window.
	event.Frame(gtx.Ops)
}

func (g Graph[T]) handleDestroyEvent(event app.DestroyEvent) error {
	return event.Err
}

// Renders an interactive GUI based off the current graph. Useful for visualization and debugging purposes. Must be run
// in the main thread.
func (g Graph[T]) GUI() {
	// Initialize all widgets
	var initializeWidgets func(theme *material.Theme) widgets
	initializeWidgets = func(theme *material.Theme) widgets {
		onPressAllNodes := func(gtx layout.Context) {
			nodes := g.GetNodes()
			drawListModal(gtx, nodes)
		}
		onPressAllEdges := func(gtx layout.Context) {
			edges := g.edges
			drawListModal(gtx, edges)
		}
		var widgets widgets
		widgets.nodeButton = &button{
			pressed: false,
			theme:   theme,
			onPress: onPressAllNodes,
			color:   blue,
		}
		widgets.edgeButton = &button{
			pressed: false,
			theme:   theme,
			onPress: onPressAllEdges,
			color:   blue,
		}
		return widgets
	}
	// The main run method.
	var run func(*app.Window) error
	run = func(w *app.Window) error {
		ops := op.Ops{}
		theme := material.NewTheme()
		g.guiInitialize(w)
		widgets := initializeWidgets(theme)
		for {
			switch e := w.Event().(type) {
			case app.DestroyEvent:
				return g.handleDestroyEvent(e)
			case app.FrameEvent:
				g.handleFrameEvent(e, ops, &widgets)
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
