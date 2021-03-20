package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"runtime"

	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"github.com/skanehira/yd/ui"
	"github.com/spf13/cobra"
	"golang.org/x/term"
	logging "gopkg.in/op/go-logging.v1"
)

var rURL = regexp.MustCompile(`https?://[\w!\?/\+\-_~=;\.,\*&@#\$%\(\)'\[\]]+`)

var rootCmd = &cobra.Command{
	Args: rootArgs,
	Use:  "yd",
}

func rootArgs(cmd *cobra.Command, args []string) error {
	if len(args) > 0 && args[0] == "version" {
		return errors.New("invalid arguments")
	}
	return nil
}

func exitError(msg interface{}) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

func Execute() {
	file := rootCmd.PersistentFlags().StringP("file", "f", "", "yaml file")
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		var (
			b   []byte
			err error
		)
		if term.IsTerminal(int(os.Stdin.Fd())) {
			if len(args) == 0 && *file == "" {
				_ = rootCmd.Help()
				return
			}

			var f string
			if *file != "" {
				f = *file
			} else {
				f = args[0]
			}

			if rURL.MatchString(f) {
				resp, err := http.Get(f)
				if err != nil {
					exitError(err)
				}

				b, err = ioutil.ReadAll(resp.Body)
				if err != nil {
					exitError(err)
				}
			} else {
				b, err = ioutil.ReadFile(f)
				if err != nil {
					exitError(err)
				}
			}
		} else {
			b, err = ioutil.ReadAll(os.Stdin)
			if err != nil {
				exitError(err)
			}

			// workaround
			// for fix "inappropriate ioctl for device"
			// this error is cause in tcell/v2 use stdin/stdout when initialize
			if runtime.GOOS != "windows" {
				stdin, err := os.OpenFile("/dev/tty", os.O_RDONLY, 0)
				if err != nil {
					exitError(err)
				}
				os.Stdin = stdin
			}
		}

		in := bytes.NewBuffer(b)
		logging.SetLevel(logging.ERROR, "")

		out := &bytes.Buffer{}
		printer := yqlib.NewPrinter(out, false, true, true, 2, true)
		eval := yqlib.NewStreamEvaluator()
		parser := yqlib.NewExpressionParser()

		if err := ui.New(in, out, printer, eval, parser).Start(); err != nil {
			exitError(err)
		}
	}

	if err := rootCmd.Execute(); err != nil {
		exitError(err)
	}
}
