package main

import (
	"fmt"
	"log"
	"os"
	"tui/internal"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func main() {
	state, err := internal.InitState()
	if err != nil {
		log.Fatal(err)
	}
	app := tview.NewApplication()
	layout := tview.NewFlex().SetDirection(tview.FlexRow)

	titleField := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false).
		SetText(state.CurrentDir)

	infoField := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false)

	inputField := tview.NewInputField().SetFieldBackgroundColor(tcell.ColorDarkViolet)
	dirNameField := tview.NewInputField().SetFieldBackgroundColor(tcell.ColorDarkViolet)
	fileNameField := tview.NewInputField().SetFieldBackgroundColor(tcell.ColorDarkViolet)

	list := tview.NewList().ShowSecondaryText(false)

	loadItem(list, state)
	list.SetSelectedFunc(handler(state, list, infoField, titleField))
	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		infoField.SetText("")
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'n':
				inputField.SetLabel("[F[]ile / [D[]irectory?")
				layout.AddItem(inputField, 1, 1, false)
				app.SetFocus(inputField)
			case 'x':
				index := list.GetCurrentItem()
				if state.CurrentDir != "/" && index == 0 {
					break
				}
				inputField.SetLabel("Are you sure? [Y[]es / [N[]o")
				layout.AddItem(inputField, 1, 1, false)
				app.SetFocus(inputField)
			}
		}
		return event
	})
	inputField.SetDoneFunc(func(key tcell.Key) {
		inputText := inputField.GetText()
		switch inputText {
		case "d":
			dirNameField.SetLabel("relative path:")
			layout.AddItem(dirNameField, 1, 1, false)
			app.SetFocus(dirNameField)
		case "f":
			fileNameField.SetLabel("relative path:")
			layout.AddItem(fileNameField, 1, 1, false)
			app.SetFocus(fileNameField)
		case "y":
			selectedRow := ""
			index := list.GetCurrentItem()
			if state.CurrentDir != "/" {
				index--
			}
			if index < 0 {
				break
			}
			selectedRow = state.Files[index].Name()
			if err := state.DeleteFileAndDirectory(selectedRow); err != nil {
				infoField.SetText(err.Error())
				break
			}

			if err := state.RefreshFiles(); err != nil {
				infoField.SetText(err.Error())
				break
			}
			list.Clear()
			loadItem(list, state)

			app.SetFocus(list)
		default:
			app.SetFocus(list)
		}
		inputField.SetText("")
		layout.RemoveItem(inputField)
	})

	dirNameField.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			inputText := dirNameField.GetText()

			err := state.CreateDirectory(inputText)
			if err != nil {
				infoField.SetText(err.Error())
			}
			if err := state.RefreshFiles(); err != nil {
				infoField.SetText(err.Error())
				break
			}
			list.Clear()
			loadItem(list, state)
		}
		dirNameField.SetText("")
		layout.RemoveItem(dirNameField)
		app.SetFocus(list)
	})
	fileNameField.SetDoneFunc(func(key tcell.Key) {
		infoField.SetText("")
		switch key {
		case tcell.KeyEnter:
			inputText := fileNameField.GetText()

			if _, err := os.Stat(inputText); err == nil {
				infoField.SetText("file exists")
				break
			}

			if err := state.CreateFile(inputText); err != nil {
				infoField.SetText(err.Error())
				break
			}
			if err := state.RefreshFiles(); err != nil {
				infoField.SetText(err.Error())
				break
			}
			list.Clear()
			loadItem(list, state)
		}
		fileNameField.SetText("")
		layout.RemoveItem(fileNameField)
		app.SetFocus(list)
	})
	layout.
		AddItem(titleField, 1, 1, false).
		AddItem(list, 0, 1, true).
		AddItem(infoField, 1, 1, false)

	if err := app.SetRoot(layout, true).Run(); err != nil {
		panic(err)
	}

}
func handler(state *internal.State, list *tview.List, infoField *tview.TextView, titleField *tview.TextView) func(int, string, string, rune) {
	return func(index int, mainText string, secondaryText string, shortcut rune) {
		titleField.SetText(state.CurrentDir)
		infoField.Clear()
		selectedRow := ""
		if state.CurrentDir != "/" {
			index--
		}
		if index < 0 {
			state.BackToParentDir()
		} else {
			selectedRow = state.Files[index].Name()

			if err := state.ChangeDir(selectedRow); err != nil {
				infoField.SetText(err.Error())
				return
			}
		}
		list.Clear()
		loadItem(list, state)
		// list.SetSelectedFunc(handler(state, list, infoField, titleField))
		titleField.SetText(state.CurrentDir)
	}
}
func visualizeFiles(state internal.State) []string {
	fileNames := []string{}
	if state.CurrentDir != "/" {
		fileNames = append(fileNames, "[green]..")
	}
	for _, f := range state.Files {
		color := "[white]%s"
		if f.IsDir() {
			color = "[green]%s/"
		}
		fileNames = append(fileNames, fmt.Sprintf(color, f.Name()))
	}
	return fileNames
}

func loadItem(list *tview.List, state *internal.State) {
	filesView := visualizeFiles(*state)
	for _, fv := range filesView {
		list.AddItem(fv, "", 0, nil)
	}
}
