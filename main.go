package main

import (
	"fmt"
	"log"
	"tui/internal"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

// func main() {
// 	s, _ := internal.BackToParentDir()
// 	fmt.Println(s)
// 	os.Chdir("/")
// 	s1, _ := internal.BackToParentDir()
// 	fmt.Println(s1)
// }
func main() {
	var err error
	state, err := internal.InitState()
	if err != nil {
		log.Fatal(err)
	}
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	x, y := ui.TerminalDimensions()
	p := widgets.NewParagraph()
	p.SetRect(0, y-1, x, y)
	p.Border = false

	l := widgets.NewList()
	l.Rows = visualizeFiles(*state)
	l.Title = state.CurrentDir
	l.WrapText = false
	l.Border = false
	l.SetRect(0, 0, x, y-1)

	ui.Render(l, p)

	previousKey := ""
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "j", "<Down>":
			l.ScrollDown()
		case "k", "<Up>":
			l.ScrollUp()
		case "<Enter>":
			p.Text = ""

			index := l.SelectedRow
			selectedRow := ""
			if state.CurrentDir != "/" {
				index--
			}
			if index < 0 {
				state.BackToParentDir()
			} else {
				selectedRow = state.Files[index].Name()

				if err := state.ChangeDir(selectedRow); err != nil {
					p.Text = err.Error()
				}
			}
			l.Title = state.CurrentDir
			l.Rows = visualizeFiles(*state)
		case "<Resize>":
			rs, ok := e.Payload.(ui.Resize)
			if !ok {
				p.Text = "Internal error"
			}
			fmt.Println(rs.Width, rs.Height)
			l.SetRect(0, 0, rs.Width, rs.Height-1)
			p.SetRect(0, rs.Height-1, rs.Width, rs.Height)
			ui.Clear()
		}
		if previousKey == "g" {
			previousKey = ""
		} else {
			previousKey = e.ID
		}
		ui.Render(l, p)
	}
}

func visualizeFiles(state internal.State) []string {
	fileNames := []string{}
	if state.CurrentDir != "/" {
		fileNames = append(fileNames, "[..](fg:green)")
	}
	for _, f := range state.Files {
		color := "[%s](fg:red)"
		if f.IsDir() {
			color = "[%s](fg:green)"
		}
		fileNames = append(fileNames, fmt.Sprintf(color, f.Name()))
	}
	return fileNames
}
