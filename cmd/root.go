package cmd

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"github.com/skanehira/yid/ui"
	"github.com/spf13/cobra"
	"golang.org/x/term"
	logging "gopkg.in/op/go-logging.v1"
)

var rootCmd = &cobra.Command{
	Use: "yid",
}

func exitError(msg interface{}) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

func Execute() {
	file := rootCmd.PersistentFlags().StringP("file", "f", "", "yaml file")
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		var f io.Reader
		if term.IsTerminal(int(os.Stdin.Fd())) {
			if *file == "" {
				_ = rootCmd.Help()
				return
			}
			var err error
			f, err = os.Open(*file)
			if err != nil {
				exitError(err)
				return
			}
		} else {
			f = os.Stdin
		}

		b, err := ioutil.ReadAll(f)
		if err != nil {
			exitError(err)
		}
		in := bytes.NewBuffer(b)
		logging.SetLevel(logging.ERROR, "")

		out := &bytes.Buffer{}
		printer := yqlib.NewPrinter(out, false, true, true, 2, true)
		eval := yqlib.NewStreamEvaluator()
		parser := yqlib.NewExpressionParser()

		// workaround
		// for fix "inappropriate ioctl for device"
		// this error is cause in tcell/v2 use stdin/stdout when initialize
		if runtime.GOOS != "windows" {
			stdin, e := os.OpenFile("/dev/tty", os.O_RDONLY, 0)
			if e != nil {
				exitError(err)
			}
			os.Stdin = stdin
		}

		if err := ui.New(in, out, printer, eval, parser).Start(); err != nil {
			exitError(err)
		}
	}

	if err := rootCmd.Execute(); err != nil {
		exitError(err)
	}
}
