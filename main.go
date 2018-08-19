package main

import (
	"fmt"
	"os"

	"github.com/hioki-daichi/imgconv/cmd"
	"github.com/hioki-daichi/imgconv/opt"
)

func main() {
	err := execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(0)
}

func execute() error {
	dirname, options, err := opt.Parse()
	if err != nil {
		return err
	}

	runner := &cmd.Runner{
		OutStream: os.Stdout,
		Decoder:   options.Decoder,
		Encoder:   options.Encoder,
		Force:     options.Force,
	}
	err = runner.Run(dirname)
	if err != nil {
		return err
	}

	return nil
}
