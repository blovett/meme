package image

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/nomad-software/meme/cli"
	"github.com/nomad-software/meme/image/stream"
	"github.com/nomad-software/meme/output"
)

// Save the passed image to disk.
func Save(opt cli.Options, st stream.Stream) (string, error) {
	var name string

	if opt.OutName != "" {
		name = opt.OutName
	} else {
		name = tempName(st.FileExt())
	}

	file, err := os.Create(name)
	if err != nil {
		return "", output.WrapError(err, "could not create image file")
	}
	defer file.Close()

	if _, err := io.Copy(file, &st); err != nil {
		return "", output.WrapError(err, "could not save image stream to file")
	}

	return name, nil
}

// Generate a temporary file name.
func tempName(ext string) string {
	dir := os.TempDir()
	return filepath.Join(dir, fmt.Sprintf("meme.%s", ext))
}
