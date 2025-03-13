package ui

import (
	"bytes"

	"github.com/gdamore/tcell/v2"
	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"github.com/rivo/tview"
)

type UI struct {
	In  *bytes.Buffer
	Out *bytes.Buffer

	Printer   yqlib.Printer
	Evaluator yqlib.StreamEvaluator
	Parser    yqlib.ExpressionParser

	App   *tview.Application
	Input *tview.InputField
	View  *tview.TextView
}

func NewInput() *tview.InputField {
	input := tview.NewInputField()
	input.SetFieldBackgroundColor(tcell.ColorDefault)
	input.SetLabel(">").SetLabelColor(tcell.ColorGreen)
	input.SetFieldTextColor(tcell.ColorYellow)
	return input
}

func NewTextView() *tview.TextView {
	view := tview.NewTextView()
	view.SetDynamicColors(true)
	return view
}

func New(in *bytes.Buffer, out *bytes.Buffer, printer yqlib.Printer, evaluator yqlib.StreamEvaluator, parser yqlib.ExpressionParser) *UI {
	ui := &UI{
		In:        in,
		Out:       out,
		Printer:   printer,
		Evaluator: evaluator,
		Parser:    parser,
		App:       tview.NewApplication(),
		Input:     NewInput(),
		View:      NewTextView(),
	}
	return ui
}

func (ui *UI) Start() error {
	grid := tview.NewGrid().SetRows(1, 0).
		AddItem(ui.Input, 0, 0, 1, 1, 0, 0, true).
		AddItem(ui.View, 1, 0, 1, 1, 0, 0, true)

	ui.App.SetRoot(grid, true)

	ui.View.SetText(ui.Out.String())

	eval := func(expr string) {
		defer func() {
			ui.Out.Reset()
			// do nothing
			_ = recover()
		}()

		node, err := ui.Parser.ParseExpression(expr)
		if err != nil {
			return
		}

		// copy input
		in := *ui.In

		if err := ui.Evaluator.Evaluate("-", &in, node, ui.Printer); err != nil {
			return
		}
		ui.View.SetText(tview.TranslateANSI(ui.Out.String()))
		ui.View.ScrollToBeginning()
	}

	ui.Input.SetChangedFunc(func(text string) {
		go ui.App.QueueUpdateDraw(func() {
			eval(text)
		})
	})

	ui.Input.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			ui.App.SetFocus(ui.View)
		}
		if event.Key() == tcell.KeyCtrlN {
			row, _ := ui.View.GetScrollOffset()
			ui.View.ScrollTo(row+1, 0)
		}
		if event.Key() == tcell.KeyCtrlP {
			row, _ := ui.View.GetScrollOffset()
			if row > 0 {
				ui.View.ScrollTo(row-1, 0)
			}
		}
		return event
	})

	ui.View.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			ui.App.SetFocus(ui.Input)
		}
		return event
	})

	eval("")

	if err := ui.App.Run(); err != nil {
		return err
	}

	return nil
}
