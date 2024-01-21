# chip8

A CHIP-8 emulator written in Go.

![screenshot](https://user-images.githubusercontent.com/21070577/123540302-215be680-d6f3-11eb-84fa-72c4698c1691.png)

## Description

CHIP8 is an interpreted programming language designed in the 1970s for writing games. The CHIP-8 programs are run on a virtual machine, which consists of the following components:

- 4 kilobytes of memory
- 16 8-bit general-purpose registers
- A call stack for subroutines
- 64 x 32 pixel monochrome display
- A hexidecimal keypad
- A monotone beeper
- 2 60Hz countdown timers, one for audio and other for instruction delay

For more information on CHIP-8 and how to write programs on it, see the following link: <http://devernay.free.fr/hacks/chip8/C8TECH10.HTM>.

## Setup

This project requires Go 1.18+ and [Ebiten](https://ebiten.org). Make sure to install all the [system dependencies](https://ebitengine.org/en/documents/install.html) necessary to run Ebiten.

To build the emulator, run the following:

```log
make build
```

This will create the executable file `chip8` in the root directory of the project.

## Usage

The executable takes in the path to a CHIP-8 ROM file.

```shell
$ ./chip8 <path to ROM>
```

For example,

```shell
$ ./chip8 roms/Logo.ch8
```

The executable also takes in a `-scale` flag for changing the window size of the emulator.

```shell
$ ./chip8 -scale=10 roms/Logo.ch8   // Scales the window size by a factor of 10
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

| Key   |         Description |
|:------|--------------------:|
| `[`   |    Resume emulation |
| `]`   |     Pause emulation |
| `\`   |     Reset emulation |
| `Esc` | Terminate emulation |

_Note: Pausing emulation will only pause the virtual machine. However, it will not pause the timers or keypad._

## References

- [CHIP-8 - Wikipedia](https://en.wikipedia.org/wiki/CHIP-8)
- [Cowgod's CHIP-8 Technical Reference](http://devernay.free.fr/hacks/chip8/C8TECH10.HTM)
- [Guide to making a CHIP-8 emulator](https://tobiasvl.github.io/blog/write-a-chip-8-emulator/#specifications)
