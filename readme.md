# serpentine

An experimental small library to make user friendly [cobra][] commands.

<p><img src="https://github.com/user-attachments/assets/f40ec67f-8437-412a-ae58-0c7810d3ee1a" width="800"></p>

## Features

- Beautiful help pages: styled help and usage pages
- Automatic `--version`: set it to the [build info][info], or a provided version
- Man pages: adds a hidden `man` command to generate _manpages_ using
  [mango][][^1]
- Completions: adds a `completion` command that generate shell completions
- Silence usage (don't show the help after an user error)
- Beautiful error handling

[info]: https://pkg.go.dev/runtime/debug#BuildInfo
[cobra]: https://github.com/spf13/cobra
[mango]: https://github.com/muesli/mango

[^1]:
    Default cobra man pages generates one man page for each command. This is
    generally fine for programs with a lot of sub commands, like git, but its an
    overkill for smaller programs.
    Mango also uses _roff_ directly instead of converting from markdown, so it
    should render better looking man pages.

## Feedback

We’d love to hear your thoughts on this project. Feel free to drop us a note!

- [Twitter](https://twitter.com/charmcli)
- [The Fediverse](https://mastodon.social/@charmcli)
- [Discord](https://charm.sh/chat)

## License

[MIT](https://github.com/charmbracelet/gum/raw/main/LICENSE)

---

Part of [Charm](https://charm.sh).

<a href="https://charm.sh/"><img alt="The Charm logo" src="https://stuff.charm.sh/charm-badge.jpg" width="400" /></a>

Charm热爱开源 • Charm loves open source
