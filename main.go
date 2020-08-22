package main

import (
	"GoCHIP-8/chip8"
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/audio"
	"github.com/hajimehoshi/ebiten/audio/mp3"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"image/color"
	"log"
)

var (
	cpu         chip8.CPU
	audioPlayer *audio.Player
	keyMap      map[ebiten.Key]byte
	square      *ebiten.Image
)

func init() {
	var err error
	square, err = ebiten.NewImage(10, 10, ebiten.FilterNearest)
	if err != nil {

	}
	err = square.Fill(color.White)
	if err != nil {

	}
}

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
	ebiten.SetWindowTitle(fmt.Sprintf("(TPS: %f; FPS: %f)", ebiten.CurrentTPS(), ebiten.CurrentFPS()))
	for i := 0; i < 10; i++ {
		cpu.NeedDraw = false
		cpu.WaitInput = false
		isKeyPressed := true
		cpu.Run()
		//cpu.Debug()

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

		if cpu.Register.ST > 0 {
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

func main() {
	audioContext, err := audio.NewContext(48000)
	if err != nil {
		log.Println("Failed to create audio context")
	} else {
		f, err := ebitenutil.OpenFile("assets/beep.mp3")
		if err != nil {
			log.Println("Failed to create audio context")
		}
		d, err := mp3.Decode(audioContext, f)
		if err != nil {

		}
		audioPlayer, err = audio.NewPlayer(audioContext, d)
		if err != nil {

		}
	}
	setupKeys()
	cpu = chip8.NewCPU()
	err = cpu.LoadROM("roms/TETRIS")
	if err != nil {
		log.Fatalln("Failed to load rom")
	}
	if err := ebiten.Run(update, 640, 320, 1, "PONG"); err != nil {
		panic(err)
	}
}
