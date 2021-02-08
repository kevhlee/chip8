# chip8

A CHIP-8 emulator written in Go.

```log
 ██████╗██╗  ██╗██╗██████╗        █████╗
██╔════╝██║  ██║██║██╔══██╗      ██╔══██╗
██║     ███████║██║██████╔╝█████╗╚█████╔╝
██║     ██╔══██║██║██╔═══╝ ╚════╝██╔══██╗
╚██████╗██║  ██║██║██║           ╚█████╔╝
 ╚═════╝╚═╝  ╚═╝╚═╝╚═╝            ╚════╝
```

## Setup

This project requires [Ebiten](https://ebiten.org) for rendering the emulator. Make sure to install all the [system dependencies](https://ebiten.org/documents/install.html) necessary to run Ebiten.

**Building** the emulator:

```bash
make
```

**Installing** the emulator into your machine:

```bash
make install
```

**Uninstalling** the emulator from your machine:

```bash
make uninstall
```

## Usage

```log
Usage:
  ch8 [command]

Examples:
$ ch8 run roms/Logo.ch8

Available Commands:
  help        Help about any command
  run         Run a CHIP-8 ROM file

Flags:
  -h, --help   help for ch8

Use "ch8 [command] --help" for more information about a command.
```

## References

- [Cowgod's CHIP-8 Technical Reference](http://devernay.free.fr/hacks/chip8/C8TECH10.HTM)