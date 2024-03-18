package main

import (
    "os"
    "io/fs"
    "log"
    "fmt"
    "strings"

    tea "github.com/charmbracelet/bubbletea"
)

type DirState struct {
    allDir []string
    toRenderDir []string
}

type model struct {
    path string
    typed bool
    cursor int
    selected string
}

// TODO: This should be dynamic where your go binary is located
const basePath string = "/home/william/Projects/directory-finder-go"

var dirIgnore = []string {
    "afs",
    "var",
    "boot",
    "dev",
    "lost+found",
    "proc",
    "root",
    "run",
}

var dirState DirState

func main() {
    app := tea.NewProgram(initialModel())

    if _, err := app.Run(); err != nil {
        log.Println(err)
        os.Exit(1)
    }
}

func initialModel() model {
    return model{
        path: "",
        cursor: 0,
    }
}

func (m model) Init() tea.Cmd {
    m.getAllDir()
    return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
        case tea.KeyMsg:
            switch msg.String() {
            case "ctrl+c", "q":
                return m, tea.Quit
            case "up":
                m.cursor--
            case "down":
                m.cursor++
            case "backspace":
                pathLen := len(m.path) 
                m.path = m.path[:pathLen - 1]

                dirState.GetSimilarDir(m.path)
            case "enter":
                createShellScript(dirState.toRenderDir[m.cursor])
                return m, tea.Quit
            default:
                if !m.typed {
                    m.typed = true
                }

                m.path += msg.String()

                dirState.GetSimilarDir(m.path)
            }
    }

    return m, nil
}

func (m model) View() string {
    s := ""
    if !m.typed {
        s += "Type an directory path\n\n"
    } else {
        if m.checkDirExists() {
            s += "Exists\n\n"
        } else {
            s += "Dont exist\n\n"
        }
    }

    for i, dir := range dirState.toRenderDir {
        if m.cursor == i {
            s += "> "
        }

        s += dir
        s += "\n"
    }

    s += "\n"
    s += m.path
    s += "\n"

    return s
}

func (m model) getAllDir() {
    root := "/"
    fileSystem := os.DirFS(root)
    fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            log.Println(err)
            os.Exit(1)
        }

        skip := false
        for i := 0; i < len(dirIgnore); i++ {
            if (dirIgnore[i] == path) {
                skip = true
                break
            }
        }

        if skip {
            return fs.SkipDir
        }

        dirState.allDir = append(dirState.allDir, path)

        return nil
    })
} 

func (m model) checkDirExists() bool {
    _, err := os.Stat(m.path)

    if err != nil {
        if os.IsNotExist(err) {
            return false
        } else {
            return true
        }
    }

    return true
}

func createShellScript(path string) {
    f := fmt.Sprintf("%s/change-directory.sh", basePath)
    err := os.Remove(f)

    if err != nil {
        if !os.IsNotExist(err) {
            os.Exit(1)
        }
    }

    s := fmt.Sprintf("#!/bin/bash\ncd %s", fmt.Sprintf("/%s", path))
    err = os.WriteFile(f, []byte(s), 0666)

    if err != nil {
        log.Println(err)
        os.Exit(1)
    }
}

func (d *DirState) GetSimilarDir(path string) {
    dirState.toRenderDir = nil

    for _, dir := range dirState.allDir {
        if len(dirState.toRenderDir) == 10 {
            break
        }

        if (strings.Contains(dir, path)) {
            dirState.AddDir(dir)
        }
    }
}

func (d *DirState) AddDir(path string) {
    d.toRenderDir = append(d.toRenderDir, path)
}
