package main

import (
    "os"
    "io/fs"
    "log"
    "fmt"
    "strings"
    "encoding/json"

    tea "github.com/charmbracelet/bubbletea"
)

type DirState struct {
    allDir []string
    toRenderDir []string
    currentDir string
}

type model struct {
    path string
    typed bool
    cursor int
    selected string
}

type DirFrequence struct {
    Name string `json:"name"`
    Frequence int `json:"frequence"`
}

const basePath string = "/usr/local/directory-picker-go"
const dirFrequencePath string = "/usr/local/directory-picker-go/frequence.json"

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
                if pathLen == 0 {
                    break
                }

                m.path = m.path[:pathLen - 1]

                dirState.GetSimilarDir(m.path)
            case "enter":
                dirState.currentDir = dirState.toRenderDir[m.cursor]
                dirState.SaveDirFrequence()
                createShellScript(dirState.currentDir)
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
    }

    /* else {
        if m.checkDirExists() {
            s += "Exists\n\n"
        } else {
            s += "Dont exist\n\n"
        }
    } */

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

        if !d.IsDir() {
            return nil
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
    d.toRenderDir = nil

    for _, dir := range d.allDir {
        if len(d.toRenderDir) == 10 {
            break
        }

        if (strings.Contains(dir, path)) {
            d.AddDir(dir)
        }
    }
}

func (d *DirState) SaveDirFrequence() {
    data, err := os.ReadFile(dirFrequencePath)

    var dirFrequences []DirFrequence

    if err != nil {
        if !os.IsNotExist(err)  {
            log.Println("Unknown error")
            log.Println(err)
            os.Exit(1)
        }
    }
    
    /* 
        Since we check for errors different than IsNotExist before, if an error
        occurs, surely is the IsNotExist error, so we must initialize the array data with
        the first dir frequence and dont need to load the previous JSON data, since it
        dont exist yet.

        If we dont get the IsNotExist error, the file exists and some data
        already is saved, so we need to check if the dir frequence that we are going
        to save already exists (just update) or we must create it.
    */

    if err != nil {
        frequence := DirFrequence{
            Name: d.currentDir,
            Frequence: 1,
        }

        dirFrequences = append(dirFrequences, frequence)
    } else {
        err = json.Unmarshal(data, &dirFrequences)

        if err != nil {
            log.Println("Unmarshall existing frequences")
            log.Println(err)
            os.Exit(1)
        }

        frequenceExists := false
        for _, dir := range dirFrequences {
            if dir.Name == d.currentDir {
                frequenceExists = true
                dir.Frequence++
            }
        }

        if !frequenceExists {
            frequence := DirFrequence{
                Name: d.currentDir,
                Frequence: 1,
            }

            dirFrequences = append(dirFrequences, frequence)
        }
    }

    updatedDirFrequences, err := json.Marshal(dirFrequences)

    if err != nil {
        log.Println("Marshal updatedFrequences")
        log.Println(err)

        os.Exit(1)
    }

    err = os.WriteFile(dirFrequencePath, []byte(updatedDirFrequences), 0666)

    if err != nil {
        log.Println("Error writing to frequences file")
        log.Println(err)

        os.Exit(1)
    }
} 

func (d *DirState) AddDir(path string) {
    d.toRenderDir = append(d.toRenderDir, path)
}
