# Demonstration
https://github.com/WilliamKSilva/directory-picker-go/assets/75429175/db55cac3-7315-4ff8-81c1-56e26e955313

# Installation
- Check the releases page to download the .tar.gz for linux: https://github.com/WilliamKSilva/directory-picker-go/releases
- To install the application run:
    `rm -rf /usr/local/directory-picker-go && tar -C /usr/local -xzf directory-picker-go.tar.gz` (You may need to run as sudo)
- Will also need to add an alias on your shell config file (.bashrc or .zshrc):
    `alias dp='sudo /usr/local/directory-picker-go/bin/directory-picker-go && source /usr/local/directory-picker-go/change-directory.sh'`
    This is needed so the shell can invoke the cd command based on the path saved on an .sh file stored on the root of the application at /usr/local

# Next steps
- [X] Record basic video showing the current state of the TUI
- [X] Add better docs showing how to install and run the TUI without
the need of building from source
- [X] Most visited directories show up at the top of the list
- [ ] Add initial most common directories for user to choose before some frequence is already established
Example: If you type "local" before having already visited /usr/local/, /usr/local will not appear in the list, will you have to type "usr/local" to go to this directory for the first time.
- [ ] Some basic customization (dont want the TUI to be so bloated)

# How to build
- You need Go toolchain installed
- Clone this repo
- `go build`
- Add an alias on your .bashrc or .zshrc: 
    `alias dp='sudo PATH-TO-REPO/directory-picker-go/directory-picker-go && source PATH-TO-REPO/directory-picker-go/change-directory.sh'`
