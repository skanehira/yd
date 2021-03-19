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
	//input.SetFieldBackgroundColor(tcell.ColorDefault)
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

	ui.Input.SetChangedFunc(func(text string) {
		go ui.App.QueueUpdateDraw(func() {
			defer func() {
				ui.Out.Reset()
				// do nothing
				recover()
			}()
			expr := text
			node, err := ui.Parser.ParseExpression(expr)
			if err != nil {
				return
			}

			// copy input
			in := *ui.In

			if err := ui.Evaluator.Evaluate("-", &in, node, ui.Printer); err != nil {
				return
			}

			ui.View.SetText(tview.TranslateANSI(ui.Out.String())).ScrollToBeginning()
		})
	})

	ui.Input.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			ui.App.SetFocus(ui.View)
		}
		return event
	})

	ui.View.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			ui.App.SetFocus(ui.Input)
		}
		return event
	})

	if err := ui.App.Run(); err != nil {
		return err
	}

	return nil
}
