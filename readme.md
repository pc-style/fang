# serpentine

An experimental small library to make user friendly [cobra][] commands.

<p><img src="https://github.com/user-attachments/assets/38c4e0ba-cebd-4455-bd14-edc664651123" width="766"></p>

## Features

- Beautiful help pages: styled help and usage pages
- Automatic `--version`: set it to the [build info][info], or a provided version
- Man pages: adds a hidden `man` command to generate _manpages_ using
  [mango][][^1]
- Completions: adds a `completion` command that generate shell completions

[info]: https://pkg.go.dev/runtime/debug#BuildInfo
[cobra]: https://github.com/spf13/cobra
[mango]: https://github.com/muesli/mango

[^1]:
    Default cobra man pages generates one man page for each command. This is
    generally fine for programs with a lot of sub commands, like git, but its an
    overkill for smaller programs.
    Mango also uses _roff_ directly instead of converting from markdown, so it
    should render better looking man pages.
