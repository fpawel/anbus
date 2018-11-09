package work

import (
	"fmt"
	"github.com/fpawel/anbus/internal/chart"
	"github.com/pkg/errors"
	"io"
	"os"
	"strings"
)

type MainSvc struct {
	w *worker
}

func (x *MainSvc) PerformTextCommand(v [1]string, _ *struct{}) error {
	c, err := parseTxtCmd(v[0])
	if err != nil {
		return errors.Wrap(err, v[0])
	}
	switch strings.ToUpper(c.name()) {
	case "EXIT":
		if !x.w.rpcWnd.CloseWindow() {
			return errors.New("can not close rpc window")
		}
		return nil
	default:
		if r, err := c.parseModbusRequest(); err == nil {
			x.w.chModbusRequest <- r
			return nil
		}
	}
	return errors.Errorf("нет такой команды: %q", c.name())
}

func (x *MainSvc) OpenArchive(v [1]string, _ *struct{}) error {
	if x.w.chartSvc.series != x.w.series {
		if err := x.w.chartSvc.series.Close(); err != nil {
			return err
		}
	}
	if strings.TrimSpace(v[0]) == "" {
		x.w.chartSvc.series = x.w.series
		return nil
	}

	if series, err := chart.OpenFile(v[0]); err != nil {
		return err
	} else {
		x.w.chartSvc.series = series
	}
	return nil
}

func (x *MainSvc) SaveArchive(v [1]string, _ *struct{}) error {
	return CopyFile(x.w.chartSvc.series.FileName(), v[0])
}

// CopyFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherise, attempt to create a hard link
// between the two files. If that fail, copy the file contents from src to dst.
func CopyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	if err = os.Link(src, dst); err == nil {
		return
	}
	err = copyFileContents(src, dst)
	return
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}
