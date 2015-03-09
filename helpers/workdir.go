package helpers

import (
	"log"
	"os"
	"runtime"
)

type WorkDir struct {
	dir string
}

func (w *WorkDir) userHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

func (w *WorkDir) Dir() string {
	return w.dir
}

func (w *WorkDir) Create(dir string) error {

	rel_path := "/.vamp_router"
	socket_dir := "/sockets"

	// allow the setting of a custom dir
	if len(dir) > 0 {
		if err := os.MkdirAll(dir+socket_dir, 0755); err != nil {
			return err
		} else {
			w.dir = dir
			return nil
		}
	}

	// if no custom dir is provided, just use the standard user $HOME + the relative path
	if err := os.MkdirAll(w.userHomeDir()+rel_path+socket_dir, 0755); err != nil {
		log.Panicf("Could not create working directorty at %s, exiting...", w.userHomeDir()+rel_path)
		return err
	}
	w.dir = w.userHomeDir() + rel_path
	return nil
}
