package main

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

type config struct {
	ext string
	// size int64
	list     bool
	del      bool
	wLogFile io.Writer
}

func main() {

	f, err := os.OpenFile("deletedLogFile", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer f.Close()

	root := "."
	conf := config{
		ext:      ".ll",
		list:     true,
		del:      true,
		wLogFile: f,
	}
	if err := run(root, os.Stdout, conf); err != nil {
		fmt.Fprintln(os.Stdout, err)

	}

}

func run(root string, out io.Writer, cfg config) error {

	logger := slog.New(slog.NewJSONHandler(cfg.wLogFile, nil))
	foundFiles := false
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filterOut(path, cfg.ext, info) {
			return nil
		}
		foundFiles = true
		if cfg.list {
			listFile(path, out)
		}
		if cfg.del {
			return deleteFile(path, logger)
		}
		return nil
	})
	if err != nil {
		return err
	}

	if !foundFiles {
		fmt.Fprintln(out, "No files with the specified extension found.")
	}

	return nil
}

func filterOut(path string, ext string, info os.FileInfo) bool {
	if info.IsDir() {
		return true
	}
	if ext != "" && !strings.EqualFold(filepath.Ext(path), ext) {
		return true
	}
	return false
}

func listFile(path string, out io.Writer) error {
	_, err := fmt.Fprintln(out, path)
	return err
}

func deleteFile(path string, delLogger *slog.Logger) error {
	if err := os.Remove(path); err != nil {
		return err
	}
	delLogger.Info("Deleted File", "path", path)
	return nil
}
