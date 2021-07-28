# :video_game: GoCHIP-8

[![License: GPL v3](https://img.shields.io/github/license/lgiki/GoCHIP-8?style=flat-square)](https://www.gnu.org/licenses/gpl-3.0) [![GitHub release (latest by date)](https://img.shields.io/github/v/release/LGiki/GoCHIP-8?style=flat-square)](https://github.com/LGiki/GoCHIP-8/releases)

A [CHIP-8](https://en.wikipedia.org/wiki/CHIP-8) emulator written in [Go](https://golang.org/). 

# Screenshot

![Screenshot](https://github.com/LGiki/GoCHIP-8/raw/master/images/screenshot.png)

# Build

1. Install all [dependencies](https://ebiten.org/documents/install.html) needed by [ebiten](https://github.com/hajimehoshi/ebiten).

2. ```bash
   git clone https://github.com/LGiki/GoCHIP-8.git
   cd GoCHIP-8
   go mod download
   go build
   ```

3. Then you can find the output program named `GoCHIP-8` under GoCHIP-8 folder.

# Command-line Flags

You can view all command-line flags via `./GoCHIP-8 -h` or `./GoCHIP-8 --help`.

## ROM

You can specify the ROM path by `-rom` parameter, for example: `-rom roms/PONG`.

Default ROM path is roms/PONG.

## Clock Speed

You can specify the CPU clock speed by `-clock` parameter, for example: `-clock 400`.

Default clock speed is 400 Hz.

## Pixel Color

You can specify the pixel color using `-color` parameter, you can choose from the following colors: white, red, green, blue, yellow, pink, cyan. For example: `-color cyan`.

Default pixel color is white.

## Full Screen

If you pass `-full` parameter on command line, the program will run in full screen mode.

## Mute Mode

If you pass `-mute` parameter on command line, the program will run in mute mode.

## Debug Mode

If you pass `-debug` parameter on command line, the program will run in debug mode.

You can view current opcode and all registers value in debug mode.

# Keyboard Configuration

## Key Mapping

```
       Chip-8                     Keyboard
+────+────+────+────+       +────+────+────+────+
| 1  | 2  | 3  | C  |       | 1  | 2  | 3  | 4  |
+────+────+────+────+       +────+────+────+────+
| 4  | 5  | 6  | D  |       | Q  | W  | E  | R  |
+────+────+────+────+  <=>  +────+────+────+────+
| 7  | 8  | 9  | E  |       | A  | S  | D  | F  |
+────+────+────+────+       +────+────+────+────+
| A  | 0  | B  | F  |       | Z  | X  | C  | V  |
+────+────+────+────+       +────+────+────+────+
```

## Additional Keys

- `Escape`:  Exit
- `P`: Pause or unpause emulation loop
- `N`: Step through while paused
- `I`: Initialize(Reset) the CPU

# References

- [https://en.wikipedia.org/wiki/CHIP-8](https://en.wikipedia.org/wiki/CHIP-8)
- [http://devernay.free.fr/hacks/chip8/C8TECH10.HTM](http://devernay.free.fr/hacks/chip8/C8TECH10.HTM)
- http://mattmik.com/files/chip8/mastering/chip8.html