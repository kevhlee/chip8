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

## References

- [Cowgod's CHIP-8 Technical Reference](http://devernay.free.fr/hacks/chip8/C8TECH10.HTM)