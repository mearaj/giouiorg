// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"fmt"
	"os"
	"strings"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/event"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget/material"
)

func main() {
	type command struct {
		name string
		run  func() error
	}

	var commands = []*command{
		// drawing section
		{name: "draw-operations", run: drawLoop(addColorOperation)},
		{name: "draw-paint", run: drawLoop(drawRedRect)},
		{name: "draw-transformation", run: drawLoop(drawRedRect10PixelsRight)},
		{name: "draw-clip", run: drawLoop(redButtonBackground)},
		{name: "draw-stack", run: drawLoop(redButtonBackgroundStack)},
		{name: "draw-draworder", run: drawLoop(drawOverlappingRectangles)},
		{name: "draw-macro", run: drawLoop(drawFiveRectangles)},
		{name: "draw-animation", run: drawLoop(drawProgressBarInternal)},
		{name: "draw-cache", run: drawLoop(drawWithCache)},
		{name: "draw-image", run: drawLoop(drawImageInternal)},

		{name: "button-low", run: drawQueueLoop(doButton)},
		{name: "external-changes", run: externalChanges},
		{name: "button-visual", run: contextLoop(handleButtonVisual)},
		{name: "button", run: contextLoop(handleButton)},

		{name: "layout-inset", run: contextLoop(inset)},
		{name: "layout-stack", run: contextLoop(stacked)},
		{name: "layout-list", run: contextLoop(listing)},
		{name: "layout-flex", run: contextLoop(flexed)},

		{name: "theme", run: themeLoop(themedApplication)},

		{name: "split-visual", run: themeLoop(exampleSplitVisual)},
		{name: "split-ratio", run: themeLoop(exampleSplitRatio)},
		{name: "split-interactive", run: themeLoop(exampleSplit)},
	}

	var cmdname string
	if len(os.Args) >= 2 {
		cmdname = os.Args[1]
	}

	var cmd *command
	for _, c := range commands {
		if strings.EqualFold(c.name, cmdname) {
			cmd = c
		}
	}

	if len(os.Args) <= 1 || cmd == nil {
		if cmdname != "" {
			fmt.Fprintf(os.Stderr, "unknown command %q\n", cmdname)
		}
		fmt.Fprint(os.Stderr, "basics [command]:\n")
		for _, cmd := range commands {
			fmt.Fprintf(os.Stderr, "\t%q\n", cmd.name)
		}
		os.Exit(1)
	}

	go func() {
		err := cmd.run()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()
	app.Main()
}

func drawLoop(draw func(*op.Ops)) func() error {
	return func() error {
		// START DRAWLOOP OMIT
		window := app.NewWindow()
		for {
			select {
			case e := <-window.Events():
				switch e := e.(type) {
				case system.DestroyEvent:
					// The window was closed.
					return e.Err
				case system.FrameEvent:
					// A request to draw the window state.
					ops := new(op.Ops)
					// Draw the state into ops.
					draw(ops)
					// Update the display.
					e.Frame(ops)
				}
			}
		}
		// END DRAWLOOP OMIT
	}
}

func drawQueueLoop(draw func(*op.Ops, event.Queue)) func() error {
	return func() error {
		// START DRAWQUEUELOOP OMIT
		window := app.NewWindow()
		for {
			select {
			case e := <-window.Events():
				switch e := e.(type) {
				case system.DestroyEvent:
					// The window was closed.
					return e.Err
				case system.FrameEvent:
					// A request to draw the window state.
					ops := new(op.Ops)
					// Draw the state into ops based on events in e.Queue.
					draw(ops, e.Queue)
					// Update the display.
					e.Frame(ops)
				}
			}
		}
		// END DRAWQUEUELOOP OMIT
	}
}

func contextLoop(draw func(*layout.Context)) func() error {
	return func() error {
		// START CONTEXTLOOP OMIT
		gtx := new(layout.Context)
		window := app.NewWindow()
		for {
			select {
			case e := <-window.Events():
				switch e := e.(type) {
				case system.DestroyEvent:
					// The window was closed.
					return e.Err
				case system.FrameEvent:
					// Reset the layout.Context for a new frame.
					gtx.Reset(e.Queue, e.Config, e.Size)

					// Draw the state into ops based on events in e.Queue.
					draw(gtx)

					// Update the display.
					e.Frame(gtx.Ops)
				}
			}
		}
		// END CONTEXTLOOP OMIT
	}
}

func themeLoop(draw func(*layout.Context, *material.Theme)) func() error {
	return func() error {
		// START THEMELOOP OMIT
		gofont.Register()
		th := material.NewTheme()

		gtx := new(layout.Context)
		window := app.NewWindow()
		for {
			select {
			case e := <-window.Events():
				switch e := e.(type) {
				case system.DestroyEvent:
					// The window was closed.
					return e.Err
				case system.FrameEvent:
					// Reset the layout.Context for a new frame.
					gtx.Reset(e.Queue, e.Config, e.Size)

					// Draw the state into ops based on events in e.Queue.
					draw(gtx, th)

					// Update the display.
					e.Frame(gtx.Ops)
				}
			}
		}
		// END THEMELOOP OMIT
	}
}
