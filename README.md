# chip8

A CHIP-8 emulator written in Go.

## Setup

To **build** the emulator, run the following command:

```bash
make
```

To **install** the emulator into your machine, run the following command:

```bash
make install
```

To **uninstall** from your machine, run the following command:

```bash
make uninstall
```

## Usage

```log
 ██████╗██╗  ██╗██╗██████╗        █████╗
██╔════╝██║  ██║██║██╔══██╗      ██╔══██╗
██║     ███████║██║██████╔╝█████╗╚█████╔╝
██║     ██╔══██║██║██╔═══╝ ╚════╝██╔══██╗
╚██████╗██║  ██║██║██║           ╚█████╔╝
 ╚═════╝╚═╝  ╚═╝╚═╝╚═╝            ╚════╝

A CHIP-8 emulator written in Go.

Usage:
  ch8 [command]

Available Commands:
  help        Help about any command
  run         Run a CHIP-8 ROM

Flags:
  -h, --help   help for ch8

Use "ch8 [command] --help" for more information about a command.
```

## References

- [Cowgod's CHIP-8 Technical Reference](http://devernay.free.fr/hacks/chip8/C8TECH10.HTM)