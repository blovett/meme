package main

import (
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/fatih/color"
	"github.com/nomad-software/meme/cli"
	"github.com/nomad-software/meme/font"
	"github.com/nomad-software/meme/image"
	"github.com/nomad-software/meme/output"
)

func main() {
	opt := cli.ParseOptions()

	if opt.Help {
		opt.PrintUsage()
		return
	}

	if opt.ListTemplates {
		for _, id := range cli.ImageIds {
			fmt.Fprintln(output.Stdout, color.CyanString("%s", id))
		}
		return
	}

	if err := run(opt); err != nil {
		fmt.Fprintln(output.Stderr, color.RedString(err.Error()))
		os.Exit(1)
	}
}

// run performs the actual meme generation. It is the only place in this
// command that's allowed to exit the process on failure; every package it
// calls into returns errors (or panics for invariants this program itself
// establishes, like its own embedded assets) rather than exiting directly,
// so callers embedding this library get a chance to handle failures
// themselves.
func run(opt cli.Options) error {
	if err := opt.Valid(); err != nil {
		return err
	}

	if err := font.SetPath(opt); err != nil {
		return err
	}

	st, err := image.Load(opt)
	if err != nil {
		return err
	}

	st, err = image.RenderImage(opt, st)
	if err != nil {
		return err
	}

	if opt.ClientID != "" {
		url, err := image.Upload(opt, st)
		if err != nil {
			return err
		}
		output.Info(url)
		return nil
	}

	file, err := image.Save(opt, st)
	if err != nil {
		return err
	}
	output.Info(file)
	return nil
}
