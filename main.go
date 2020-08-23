package main

import (
	"GoCHIP-8/chip8"
	"flag"
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/mp3"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"image/color"
	"log"
	"os"
)

var (
	cpu         chip8.CPU
	audioPlayer *audio.Player
	keyMap      map[ebiten.Key]byte
	square      *ebiten.Image
	romPath     string
	pixelColor  string
	showHelp    bool
	mute        bool
)

func init() {
	flag.StringVar(&romPath, "r", "roms/PONG", "The `path` to ROM")
	flag.StringVar(&pixelColor, "c", "white", "Pixel `color`: white, red, green, blue, yellow, pink, cyan")
	flag.BoolVar(&mute, "m", false, "Mute")
	flag.BoolVar(&showHelp, "h", false, "Show help")
	flag.Usage = usage
	flag.Parse()
}

func createSquare() {
	var err error
	square, err = ebiten.NewImage(10, 10, ebiten.FilterNearest)
	if err != nil {
		log.Fatalln("Failed to create ebiten image")
	}
	r, g, b := parsePixelColor()
	err = square.Fill(color.RGBA{
		R: r,
		G: g,
		B: b,
		A: 255,
	})
	if err != nil {
		log.Fatalln("Failed to fill color for square")
	}
}

func parsePixelColor() (r, g, b uint8) {
	switch pixelColor {
	case "white":
		r, g, b = 255, 255, 255
	case "red":
		r, g, b = 255, 0, 0
	case "green":
		r, g, b = 0, 255, 0
	case "blue":
		r, g, b = 0, 0, 255
	case "yellow":
		r, g, b = 255, 255, 0
	case "pink":
		r, g, b = 255, 0, 255
	case "cyan":
		r, g, b = 0, 255, 255
	default:
		r, g, b = 255, 255, 255
	}
	return
}

/*
       Chip-8                       Keyboard
+────+────+────+────+       +────+────+────+────+
| 1  | 2  | 3  | C  |       | 1  | 2  | 3  | 4  |
+────+────+────+────+       +────+────+────+────+
| 4  | 5  | 6  | D  |       | Q  | W  | E  | R  |
+────+────+────+────+  <=>  +────+────+────+────+
| 7  | 8  | 9  | E  |       | A  | S  | D  | F  |
+────+────+────+────+       +────+────+────+────+
| A  | 0  | B  | F  |       | Z  | X  | C  | V  |
+────+────+────+────+       +────+────+────+────+
*/
func setupKeys() {
	keyMap = make(map[ebiten.Key]byte)
	keyMap[ebiten.Key1] = 0x01
	keyMap[ebiten.Key2] = 0x02
	keyMap[ebiten.Key3] = 0x03
	keyMap[ebiten.Key4] = 0x0C

	keyMap[ebiten.KeyQ] = 0x04
	keyMap[ebiten.KeyW] = 0x05
	keyMap[ebiten.KeyE] = 0x06
	keyMap[ebiten.KeyR] = 0x0D

	keyMap[ebiten.KeyA] = 0x07
	keyMap[ebiten.KeyS] = 0x08
	keyMap[ebiten.KeyD] = 0x09
	keyMap[ebiten.KeyF] = 0x0E

	keyMap[ebiten.KeyZ] = 0x0A
	keyMap[ebiten.KeyX] = 0x00
	keyMap[ebiten.KeyC] = 0x0B
	keyMap[ebiten.KeyV] = 0x0F
}

func getPressedKeys() bool {
	for key, value := range keyMap {
		if ebiten.IsKeyPressed(key) {
			cpu.KeyState[value] = 0x01
			return true
		}
	}
	return false
}

func drawGraphics(screen *ebiten.Image) error {
	for i := 0; i < chip8.DisplayHeight; i++ {
		for j := 0; j < chip8.DisplayWidth; j++ {
			if cpu.Display[i][j] == 0x01 {
				opts := &ebiten.DrawImageOptions{}
				opts.GeoM.Translate(float64(j*10), float64(i*10))
				err := screen.DrawImage(square, opts)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func update(screen *ebiten.Image) error {
	ebiten.SetWindowTitle(fmt.Sprintf("GoCHIP-8 | %s (TPS: %f; FPS: %f)", romPath, ebiten.CurrentTPS(), ebiten.CurrentFPS()))
	for i := 0; i < 10; i++ {
		cpu.NeedDraw = false
		cpu.WaitInput = false
		isKeyPressed := true
		cpu.Run()
		//cpu.Debug()

		if ebiten.IsKeyPressed(ebiten.KeyO) {
			cpu.Reset()
			err := cpu.LoadROM(romPath)
			if err != nil {
				return err
			}
		}

		if cpu.WaitInput {
			isKeyPressed = getPressedKeys()
			if !isKeyPressed {
				cpu.Register.PC -= 2
			}
		}

		if cpu.NeedDraw || !isKeyPressed {
			err := drawGraphics(screen)
			if err != nil {
				return err
			}
		}

		for key, value := range keyMap {
			if ebiten.IsKeyPressed(key) {
				cpu.KeyState[value] = 0x01
			} else {
				cpu.KeyState[value] = 0x00
			}
		}

		if !mute && cpu.Register.ST > 0 {
			if audioPlayer != nil {
				err := audioPlayer.Play()
				if err != nil {
					return err
				}
				err = audioPlayer.Rewind()
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func Run() {
	var err error
	cpu = chip8.NewCPU()
	err = cpu.LoadROM(romPath)
	if err != nil {
		log.Fatalln("Failed to load rom")
	}
	setupKeys()
	createSquare()
	if !mute {
		audioContext, err := audio.NewContext(48000)
		if err != nil {
			log.Println("Failed to create audio context")
		} else {
			f, err := ebitenutil.OpenFile("assets/beep.mp3")
			if err != nil {
				log.Println("Failed to load assets/beep.mp3")
			} else {
				d, err := mp3.Decode(audioContext, f)
				if err != nil {
					log.Println("Failed to decode MP3 file")
				} else {
					audioPlayer, err = audio.NewPlayer(audioContext, d)
					if err != nil {
						log.Println("Failed to create audio player")
					}
				}
			}
		}
	}
	if err := ebiten.Run(update, 640, 320, 1, "PONG"); err != nil {
		panic(err)
	}
}

func usage() {
	_, _ = fmt.Fprintf(os.Stderr, `GoCHIP-8
Usage: ./GoCHIP-8 <-r pathToROM> [-m] [-c color]

Options:
`)
	flag.PrintDefaults()
}

func main() {
	if showHelp {
		flag.Usage()
	} else {
		Run()
	}
}
