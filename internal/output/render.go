package output

import (
	"encoding/json"
	"io"
	"os"

	"golang.org/x/term"
)

type Renderer struct {
	Out        io.Writer
	ForceJSON  bool
	ForceTable bool
	NoColor    bool
}

func New() *Renderer {
	return &Renderer{Out: os.Stdout}
}

func (r *Renderer) Render(v any) error {
	if r.ForceJSON || (!r.ForceTable && !r.isTTY()) {
		return r.renderJSON(v)
	}
	if ok, err := r.renderTable(v); ok {
		return err
	}
	return r.renderJSON(v)
}

func (r *Renderer) renderJSON(v any) error {
	enc := json.NewEncoder(r.Out)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func (r *Renderer) isTTY() bool {
	f, ok := r.Out.(*os.File)
	if !ok {
		return false
	}
	return term.IsTerminal(int(f.Fd()))
}
