# Bash prompt

This is an experimental bash prompt, with a few features:
- date/time, load average, and last command exit status info
- highlights when you are root, or on a remote machine
- git integration
- highlights files and commands which print no end-of-line nicely

It depends on having a few symbols from https://github.com/ryanoasis/nerd-fonts
installed. It has also been built using the Solarized Light palette from
http://ethanschoonover.com/solarized .

## Helper file notes

- `nonl` — file with no newline, `cat` it to see behaviour
- `colourtable.sh` — displays combinations of fg/bg colours
- `source_me.sh` — shows the form of `PS1` needed to activate the prompt
