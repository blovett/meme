package font

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/nomad-software/meme/cli"
	"github.com/nomad-software/meme/data"
)

var (
	// Path is the location of the font file.
	Path string
)

// SetPath overrides the font path at runtime, based on opt.Font. It returns
// an error if opt.Font is set but does not resolve to a usable font, since
// that's user input and shouldn't take the process down.
func SetPath(opt cli.Options) error {
	if opt.Font == "" {
		return nil
	}

	// direct file
	if _, err := os.Stat(opt.Font); err == nil {
		Path = opt.Font
		return nil
	}

	// try fc-match (Linux/macOS with fontconfig)
	if fc, err := exec.LookPath("fc-match"); err == nil {
		out, err := exec.Command(fc, "-f", "%{file}\\n", opt.Font).Output()
		if err == nil {
			f := strings.TrimSpace(string(out))
			if f != "" {
				if _, err := os.Stat(f); err == nil {
					Path = f
					return nil
				}
			}
		}
	}

	return fmt.Errorf("invalid font: %s", opt.Font)
}

// Write the embedded font to the temporary directory. The embedded font is
// baked into the binary at build time, so a failure here indicates a broken
// build/environment rather than something a caller can recover from.
func init() {
	if Path != "" {
		return
	}

	Path = filepath.Join(os.TempDir(), filepath.Base(data.Font))

	if _, err := os.Stat(Path); os.IsNotExist(err) {
		file, err := os.Create(Path)
		if err != nil {
			panic(fmt.Errorf("could not create font file: %w", err))
		}
		defer file.Close()

		stream, err := data.Files.ReadFile(data.Font)
		if err != nil {
			panic(fmt.Errorf("could not read embedded font: %w", err))
		}

		if _, err := file.Write(stream); err != nil {
			panic(fmt.Errorf("could not write font file: %w", err))
		}
	}
}
