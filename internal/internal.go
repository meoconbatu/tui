package internal

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
)

// State type
type State struct {
	CurrentDir string
	Files      []os.FileInfo
}

// InitState func
func InitState() (*State, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	files, err := getFiles(currentDir)
	if err != nil {
		return nil, err
	}
	return &State{CurrentDir: currentDir, Files: files}, nil
}

func getFiles(dir string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(dir)
}

// ChangeDir func
func (st *State) ChangeDir(dir string) error {
	f, err := os.Stat(dir)
	if err != nil {
		log.Fatal(err)
	}
	if !f.IsDir() {
		return errors.New("Dir is not a forder")
	}
	err = os.Chdir(dir)
	if err != nil {
		return err
	}
	st.CurrentDir, _ = os.Getwd()
	st.Files, _ = getFiles(st.CurrentDir)
	return nil
}

// BackToParentDir func
func (st *State) BackToParentDir() error {
	err := os.Chdir("..")
	if err != nil {
		return err
	}
	st.CurrentDir, _ = os.Getwd()
	st.Files, _ = getFiles(st.CurrentDir)
	return nil
}
