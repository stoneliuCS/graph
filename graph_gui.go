package graph

import (
	"fmt"
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
	blue     = color.NRGBA{R: 0x40, G: 0x40, B: 0xC0, A: 0xFF}
	red      = color.NRGBA{R: 0xC0, G: 0x40, B: 0x40, A: 0xFF}
	white    = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	grey     = color.NRGBA{R: 0xC0, G: 0xC0, B: 0xC0, A: 0xFF}
	greyDark = color.NRGBA{R: 0xB0, G: 0xB0, B: 0xB0, A: 0xFF}
	black    = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
)

// Sizes
var (
	buttonSize          = image.Pt(200, 90)
	buttonSizeClickable = buttonSize.Mul(2)
	textButtonSize      = 12
	modalTitleSize      = unit.Sp(20)
	defaultInset        = layout.UniformInset(0)
)

// Color Constants
var (
	buttonUnpressedColor = blue
	buttonPressedColor   = red
	modalContainerColor  = greyDark
	modalTitleColor      = black
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

// Hold all available buttons to press in the gui.
type widgets struct {
	nodeButton      *button
	edgeButton      *button
	modalExitButton *button
	modalNodeList   *layout.List
	modalEdgeList   *layout.List
}

type button struct {
	pressed   bool
	theme     *material.Theme
	color     color.NRGBA
	onPressed func()
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
			}
		}
	}
	if b.pressed {
		b.color = buttonPressedColor
		b.onPressed()
	} else {
		b.color = buttonUnpressedColor
	}
	return layout.Stack{Alignment: layout.Center}.Layout(gtx, layout.Stacked(func(gtx layout.Context) layout.Dimensions {
		return colorRoundedPill(gtx, buttonSize, b.color)
	}), layout.Stacked(func(gtx layout.Context) layout.Dimensions {
		label := material.Label(b.theme, unit.Sp(textButtonSize), label)
		label.Color = white
		return label.Layout(gtx)
	}))

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

func drawModalFromButtonListFormat[T any](
	gtx layout.Context,
	modalTitle string,
	b *button,
	widgets *widgets,
	listToDraw []T,
	modalList *layout.List,
) layout.Dimensions {
	size := gtx.Constraints.Max.Div(2)
	// Restrain the constraints to be this exact size.
	gtx.Constraints = layout.Exact(size)
	modalContainer := layout.Stacked(func(gtx layout.Context) layout.Dimensions {
		return defaultInset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return colorBox(gtx, size, modalContainerColor)
		})
	})
	widgets.modalExitButton.onPressed = func() {
		b.pressed = false
		widgets.modalExitButton.pressed = false
	}
	var titleAndExitButton layout.Widget = func(gtx layout.Context) layout.Dimensions {
		title := layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			modalTitle := material.Label(b.theme, modalTitleSize, modalTitle)
			modalTitle.Color = modalTitleColor
			return modalTitle.Layout(gtx)
		})
		exitButton := layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return widgets.modalExitButton.Layout(gtx, "Exit")
		})
		return layout.Flex{
			Axis:    layout.Horizontal,
			Spacing: layout.SpaceBetween,
		}.Layout(gtx, title, exitButton)
	}
	var mainContent layout.Widget = func(gtx layout.Context) layout.Dimensions {
		itemsToDraw := listToDraw
		return modalList.Layout(gtx, len(itemsToDraw), func(gtx layout.Context, index int) layout.Dimensions {
			itemToDraw := itemsToDraw[index]
			s := fmt.Sprintf("%+v", itemToDraw)
			lbl := material.Label(b.theme, unit.Sp(14), s)
			return lbl.Layout(gtx)
		})
	}
	content := layout.Stacked(func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceBetween}.Layout(gtx,
			layout.Rigid(titleAndExitButton),
			layout.Flexed(1, mainContent),
		)
	})
	return layout.Stack{Alignment: layout.Center}.Layout(gtx, modalContainer, content)
}

func (g Graph[T]) drawLayout(gtx layout.Context, widgets *widgets) {
	leftMenu := layout.Flexed(0.3, func(gtx layout.Context) layout.Dimensions {
		return g.drawMenu(gtx, widgets)
	})
	mainView := layout.Flexed(0.7, func(gtx layout.Context) layout.Dimensions {
		return colorBox(gtx, gtx.Constraints.Max, white)
	})
	main := layout.Flex{Axis: layout.Horizontal}.Layout(gtx, mainView, leftMenu)
	// Handle main UI stack handling
	layout.Stack{Alignment: layout.Center}.Layout(gtx, layout.Stacked(func(gtx layout.Context) layout.Dimensions {
		return main
	}), layout.Stacked(func(gtx layout.Context) layout.Dimensions {
		// This callback only renders if a modal button has been pressed
		// Handle modal drawings
		if !widgets.nodeButton.pressed && !widgets.edgeButton.pressed {
			return layout.Dimensions{}
		}
		if widgets.nodeButton.pressed {
			return drawModalFromButtonListFormat(
				gtx,
				"View All Nodes",
				widgets.nodeButton,
				widgets,
				g.GetNodes(),
				widgets.modalNodeList,
			)
		} else {
			return drawModalFromButtonListFormat(
				gtx,
				"View all Edges",
				widgets.edgeButton,
				widgets,
				g.edges,
				widgets.modalEdgeList,
			)
		}
	}))
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
		var widgets widgets
		widgets.nodeButton = &button{
			pressed:   false,
			theme:     theme,
			color:     blue,
			onPressed: func() {},
		}
		widgets.edgeButton = &button{
			pressed:   false,
			theme:     theme,
			color:     blue,
			onPressed: func() {},
		}
		widgets.modalExitButton = &button{
			pressed:   false,
			theme:     theme,
			color:     blue,
			onPressed: func() {},
		}
		widgets.modalNodeList = &layout.List{Axis: layout.Vertical}
		widgets.modalEdgeList = &layout.List{Axis: layout.Vertical}
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
