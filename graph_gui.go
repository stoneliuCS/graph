package graph

import (
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

type gui struct {
	window     *app.Window
	operations op.Ops
	theme      *material.Theme
}

var GUI_SP_WIDTH float32 = 800
var GUI_SP_HEIGHT float32 = 800

// Intializes some default settings for the GUI window
func (g gui) initialize() {
	g.window.Option(app.Size(unit.Dp(GUI_SP_WIDTH), unit.Dp(GUI_SP_HEIGHT)))
}

func (g gui) layout(gtx layout.Context) {
	layout.Flex{
		Axis:    layout.Vertical,
		Spacing: layout.SpaceStart,
	}.Layout(gtx)
}

// Main event loop handling
func (g gui) handleFrameEvent(event app.FrameEvent) {
	gtx := app.NewContext(&g.operations, event)
	// Initialize the layout
	g.layout(gtx)
	// Finally draw the frame event onto the window.
	event.Frame(gtx.Ops)
}

func (g gui) handleDestroyEvent(event app.DestroyEvent) error {
	return event.Err
}

// Renders an interactive GUI based off the current graph. Useful for visualization and debugging purposes. Must be run
// in the main thread.
func (g Graph[T]) GUI() {
	var run func(*app.Window) error
	run = func(w *app.Window) error {
		theme := material.NewTheme()
		gui := gui{
			window:     w,
			theme:      theme,
			operations: op.Ops{},
		}
		gui.initialize()
		for {
			switch e := w.Event().(type) {
			case app.DestroyEvent:
				return gui.handleDestroyEvent(e)
			case app.FrameEvent:
				gui.handleFrameEvent(e)
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
