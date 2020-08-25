package main

import (
	"GoCHIP-8/chip8"
	"flag"
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/mp3"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"image/color"
	"log"
	"math/rand"
	"os"
	"time"
)

var (
	cpu         chip8.CPU
	audioPlayer *audio.Player
	keyMap      map[ebiten.Key]byte
	view        *ebiten.Image
	romPath     string
	pixelColor  string
	fullScreen  bool
	showHelp    bool
	mute        bool
	counter     float64
	clockSpeed  int
	r, g, b     uint8 // Pixel color RGB
	paused      bool
	debug       bool
)

func init() {
	rand.Seed(time.Now().UnixNano())
	flag.StringVar(&romPath, "rom", "roms/PONG", "The `path` to ROM")
	flag.StringVar(&pixelColor, "color", "white", "Pixel `color`: white, red, green, blue, yellow, pink, cyan")
	flag.IntVar(&clockSpeed, "clock", 400, "CPU `clock speed` in Hz")
	flag.BoolVar(&mute, "mute", false, "Mute")
	flag.BoolVar(&debug, "debug", false, "Debug mode")
	flag.BoolVar(&fullScreen, "full", false, "Full screen")
	flag.BoolVar(&showHelp, "h", false, "Show help")
	flag.Usage = usage
	flag.Parse()
	view, _ = ebiten.NewImage(chip8.DisplayWidth*10, chip8.DisplayHeight*10, ebiten.FilterDefault)
	r, g, b = parsePixelColor()
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

type Game struct{}

func (game *Game) Layout(int, int) (screenWidth, screenHeight int) {
	screenWidth, screenHeight = chip8.DisplayWidth*10, chip8.DisplayHeight*10
	return
}

func (game *Game) Draw(screen *ebiten.Image) {
	if cpu.NeedDraw {
		_ = view.Fill(color.Black)
		for i := 0; i < chip8.DisplayHeight; i++ {
			for j := 0; j < chip8.DisplayWidth; j++ {
				if cpu.Display[i][j] == 0x01 {
					ebitenutil.DrawRect(view, float64(j*10), float64(i*10), 10, 10, color.RGBA{
						R: r,
						G: g,
						B: b,
						A: 255,
					})
				}
			}
		}
		cpu.NeedDraw = false
	}
	opts := &ebiten.DrawImageOptions{}
	_ = screen.DrawImage(view, opts)
}

func (game *Game) Update(*ebiten.Image) error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		paused = !paused
	}

	if paused && inpututil.IsKeyJustPressed(ebiten.KeyN) {
		step()
	}

	if !paused {
		for counter > 0 {
			step()
			counter -= float64(ebiten.MaxTPS())
		}
		counter += float64(clockSpeed)
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

	if ebiten.IsKeyPressed(ebiten.KeyI) {
		cpu.Reset()
		err := cpu.LoadROM(romPath)
		if err != nil {
			return err
		}
		paused = false
	}

	return nil
}

func step() {
	cpu.WaitInput = false
	if debug {
		cpu.Debug()
	}
	cpu.Run()
	if cpu.WaitInput {
		if !getPressedKeys() {
			cpu.Register.PC -= 2
		}
	}

	for key, value := range keyMap {
		if ebiten.IsKeyPressed(key) {
			cpu.KeyState[value] = 0x01
		} else {
			cpu.KeyState[value] = 0x00
		}
	}
}

func Run() {
	var err error
	cpu = chip8.NewCPU()
	err = cpu.LoadROM(romPath)
	if err != nil {
		log.Fatalln("Failed to load rom")
	}
	setupKeys()
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
	ebiten.SetFullscreen(fullScreen)
	ebiten.SetWindowSize(chip8.DisplayWidth*10, chip8.DisplayHeight*10)
	ebiten.SetWindowTitle(fmt.Sprintf("GoCHIP-8 | %s", romPath))
	if err := ebiten.RunGame(&Game{}); err != nil {
		panic(err)
	}
}

func usage() {
	_, _ = fmt.Fprintf(os.Stderr, `GoCHIP-8
Usage: ./GoCHIP-8 <-path pathToROM> [-clock clock_speed] [-color color] [-mute] [-full]

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
