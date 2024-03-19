# Demonstration
https://github.com/WilliamKSilva/directory-picker-go/assets/75429175/db55cac3-7315-4ff8-81c1-56e26e955313

# Install
- Check the release page to download the .tar.gz for linux and follow the instructions: https://github.com/WilliamKSilva/directory-picker-go/releases

# Next steps
- [X] Record basic video showing the current state of the TUI
- [X] Add better docs showing how to install and run the TUI without
the need of building from source
- Most visited directories show up at the top of the list
- Some basic customization (dont want the TUI to be so bloated)

# How to build
- You need Go toolchain installed
- Clone this repo
- `go build`
- Add an alias on your .bashrc or .zshrc: 
    `alias dp='sudo PATH-TO-REPO/directory-picker-go/directory-picker-go && source PATH-TO-REPO/directory-picker-go/change-directory.sh'`
