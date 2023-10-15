package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/go-mp3"
	"os"
	"time"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

var player *oto.Player
var otoCtx *oto.Context

func setUpContext() {

	op := &oto.NewContextOptions{}

	// Usually 44100 or 48000. Other values might cause distortions in Oto
	op.SampleRate = 44100

	// Number of channels (aka locations) to play sounds from. Either 1 or 2.
	// 1 is mono sound, and 2 is stereo (most speakers are stereo).
	op.ChannelCount = 2

	// Format of the source. go-mp3's format is signed 16bit integers.
	op.Format = oto.FormatSignedInt16LE

	// Remember that you should **not** create more than one context
	otoCtxTemp, readyChan, err := oto.NewContext(op)
	otoCtx = otoCtxTemp
	if err != nil {
		panic("oto.NewContext failed: " + err.Error())
	}
	// It might take a bit for the hardware audio devices to be ready, so we wait on the channel.
	<-readyChan
}

func playFile(filename string) {

	stopAllFiles()
	// Read the mp3 file into memory
	fileBytes, err := os.ReadFile(filename)
	if err != nil {
		panic("reading " + filename + " failed: " + err.Error())
	}

	// Convert the pure bytes into a reader object that can be used with the mp3 decoder
	fileBytesReader := bytes.NewReader(fileBytes)

	// Decode file
	decodedMp3, err := mp3.NewDecoder(fileBytesReader)
	if err != nil {
		panic("mp3.NewDecoder failed: " + err.Error())
	}

	// Prepare an Oto context (this will use your default audio device) that will
	// play all our sounds. Its configuration can't be changed later.
	if otoCtx == nil {
		setUpContext()
	}

	// Create a new 'player' that will handle our sound. Paused by default.
	player = otoCtx.NewPlayer(decodedMp3)

	// Play starts playing the sound and returns without waiting for it (Play() is async).
	player.Play()

	// We can wait for the sound to finish playing using something like this
	for player.IsPlaying() {
		time.Sleep(time.Millisecond)
	}

	// Now that the sound finished playing, we can restart from the beginning (or go to any location in the sound) using seek
	// newPos, err := player.(io.Seeker).Seek(0, io.SeekStart)
	// if err != nil{
	//     panic("player.Seek failed: " + err.Error())
	// }
	// println("Player is now at position:", newPos)
	// player.Play()

	// If you don't want the player/sound anymore simply close
	err = player.Close()
	if err != nil {
		panic("player.Close failed: " + err.Error())
	}
	player = nil

}

func stopAllFiles() {
	if player != nil {
		err := player.Close()
		if err != nil {
			panic("player.Close failed: " + err.Error())
		}
	}
}

var m map[uint8]string

func setUpMap() {
	if m == nil {
		m = make(map[uint8]string)
	}
	m[0] = "mp3/System_Of_A_Down_-_Chop_Suey_Official_HD_Video_from40.980125998571154s_to209s.mp3"
	m[1] = "mp3/Smash_Mouth_-_All_Star_Audio.mp3"
	m[2] = "mp3/Rick_Astley_-_Never_Gonna_Give_You_Up_Rick_roll_song.mp3"
	m[3] = "mp3/Nickelback_-_Photograph_OFFICIAL_VIDEO.mp3"
	m[4] = "mp3/He-Man_What_s_Going_On_Song_High_Quality.mp3"
	m[5] = "mp3/Bag_Raiders_-_Shooting_Stars_Official_Video.mp3"
}

// Greet returns a greeting for the given name
func (a *App) PlayAudio(song uint8) {
	setUpMap()
	if song > 5 {
		return
	}
	playFile(m[song])
}
func (a *App) StopAudio() {
	stopAllFiles()
}
