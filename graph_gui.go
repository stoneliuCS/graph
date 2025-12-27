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

var GUI_SP_WIDTH float32 = 800
var GUI_SP_HEIGHT float32 = 800

// Intializes some default settings for the GUI window
func (g Graph[T]) guiInitialize(window *app.Window, theme *material.Theme) {
	window.Option(app.Size(unit.Dp(GUI_SP_WIDTH), unit.Dp(GUI_SP_HEIGHT)))
}

func (g Graph[T]) drawLayout(gtx layout.Context) {
	layout.Flex{
		Axis:    layout.Vertical,
		Spacing: layout.SpaceStart,
	}.Layout(gtx)
}

// Main event loop handling
func (g Graph[T]) handleFrameEvent(event app.FrameEvent, ops op.Ops) {
	gtx := app.NewContext(&ops, event)
	// Initialize the layout
	g.drawLayout(gtx)
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
		theme := material.NewTheme()
		ops := op.Ops{}
		g.guiInitialize(w, theme)
		for {
			switch e := w.Event().(type) {
			case app.DestroyEvent:
				return g.handleDestroyEvent(e)
			case app.FrameEvent:
				g.handleFrameEvent(e, ops)
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
