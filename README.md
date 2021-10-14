# chip8

A CHIP-8 emulator written in Go.

![screenshot](https://user-images.githubusercontent.com/21070577/123540302-215be680-d6f3-11eb-84fa-72c4698c1691.png)

## Description

CHIP-8 is an interpreted programming language designed in the late 1970s for writing games. It was designed to run on 8-bit systems such as the [COSMAC VIP](https://en.wikipedia.org/wiki/COSMAC_VIP). The CHIP-8 interpreter and programs are meant to run on a virtual machine, which consists of the following components:

- 4 kilobytes of memory
- 16 8-bit general-purpose registers
- A call stack for subroutines
- 64 x 32 pixel monochrome display
- A hexidecimal keypad
- A monotone beeper
- 2 60Hz countdown timers, one for audio and other for instruction delay

For more information on CHIP-8 and how to write programs on it, see the following link: <http://devernay.free.fr/hacks/chip8/C8TECH10.HTM>.

## Setup

This project requires [Ebiten](https://ebiten.org) for rendering the emulator. Make sure to install all the [system dependencies](https://ebiten.org/documents/install.html) necessary to run Ebiten.

This project uses `Make`. To build the emulator, run the following:

```log
make build
```

This will create the executable file `ch8` in the `bin` directory of this project.

You can also install the emulator on your system using the following command:

```log
make install
```

This will build the emulator and place the built executable within `/usr/local/bin`. You may need to run this command using `sudo` access.

You can always uninstall the emulator by running the following:

```log
make uninstall
```

## Usage

A CLI is used to operate the emulator:

```log
A CHIP-8 emulator written in Go.

Usage:
  chip8 [flags]

Examples:
$ chip8 roms/Logo.ch8

Flags:
  -h, --help                help for chip8
      --hertz-io duration   set the speed of the IO timers. (default 16ms)
      --hertz-vm duration   set the speed of the virtual machine's CPU cycle. (default 2ms)
  -s, --scale int           set the scale factor of the screen (default 10)
  -t, --tps int             set the max ticks-per-second (TPS) of the renderer (default 60)
  -v, --volume float        set the volume of the emulator (default 0.5)
```

### Key Mapping

The following shows the keys that are virtually mapped to the CHIP-8 keypad:

```log
   Key                   Hex
|1|2|3|4|        \    |0|1|2|3|
|Q|W|E|R|   ------\   |4|5|6|7|
|A|S|D|F|   ------/   |8|9|A|B|
|Z|X|C|V|        /    |C|D|E|F|
```

### Emulation

The emulator provides a few basic functions for control:

| Key |      Description |
| :-- | ---------------: |
| `[` | Resume emulation |
| `]` |  Pause emulation |
| `\` |  Reset emulation |

_Note: Pausing emulation will only pause the virtual machine. However, it will not pause the timers or keypad._

## References

- [CHIP-8 - Wikipedia](https://en.wikipedia.org/wiki/CHIP-8)
- [Cowgod's CHIP-8 Technical Reference](http://devernay.free.fr/hacks/chip8/C8TECH10.HTM)
- [Guide to making a CHIP-8 emulator](https://tobiasvl.github.io/blog/write-a-chip-8-emulator/#specifications)
