package helpers

import (
	"log"
	"os"
)

type WorkDir struct {
	dir string
}

// TODO: make configuration dir absolute on start.

func (w *WorkDir) Dir() string {
	return w.dir
}

func (w *WorkDir) Create(dir string) error {

	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				log.Panicf("Could not create working directory at %s, exiting...", dir)
				return err
			} else {
				w.dir = dir
				return nil
			}
		} else {
			return err
		}
	}
	w.dir = dir
	return nil
}
