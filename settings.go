package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	settingsPaths = [...]string{
		"$HOME/.config/streamit/settings.json",
		"$HOME/.streamit/settings.json",
	}
	ErrNewSettingsFile = errors.New("new settings file created")
)

// INRES="1920x1080" # input resolution
// OUTRES="1920x1080" # output resolution
// FPS="15" # target FPS
// GOP="30" # i-frame interval, should be double of FPS,
// GOPMIN="15" # min i-frame interval, should be equal to fps,
// THREADS="2" # max 6
// CBR="1000k" # constant bitrate (should be between 1000k - 3000k)
// QUALITY="ultrafast"  # one of the many FFMPEG preset
// AUDIO_RATE="44100"
// STREAM_KEY="live_89091880_bDWwThFObk1HQy6O7DXm9K52pfqJDb"
// SERVER="live-lax" # twitch server in LA, see http://bashtech.net/twitch/ingest.php for list

type Settings struct {
	InRes     string `json:"in_res"`
	OutRes    string `json:"out_res"`
	FPS       int    `json:"fps"`
	Threads   int    `json:"threads"`
	CBR       int    `json:"cbr"`
	Quality   string `json:"quality"`
	AudioRate int    `json:"audio_rate"`
	StreamKey string `json:"stream_key"`
	Server    string `json:"server"`
	LogPath   string `json:"log_path"`
	LogLevel  string `json:"log_level"`
}

func (s *Settings) validate() error {
	if s.Threads <= 0 || s.Threads > 6 {
		return fmt.Errorf("invalid number of threads: %d", s.Threads)
	}

	if s.CBR < 1000 || s.CBR > 3000 {
		return fmt.Errorf("cbr should be between 1000k - 3000k: %dk", s.CBR)
	}

	if s.StreamKey == "" {
		return errors.New("missing stream_key")
	}

	return nil
}

func DefaultSettings() *Settings {
	return &Settings{
		InRes:     "1920x1080",
		OutRes:    "1920x1080",
		FPS:       15,
		Threads:   2,
		CBR:       1000,
		Quality:   "ultrafast",
		AudioRate: 44100,
		Server:    "live-lax",
		LogLevel:  "error",
	}
}

func LoadSettingsPath(path string) (*Settings, error) {
	file, err := os.Open(os.ExpandEnv(path))
	if err != nil {
		return nil, err
	}

	settings := &Settings{}
	err = json.NewDecoder(file).Decode(settings)
	if err != nil {
		return nil, err
	}

	return settings, nil
}

func LoadSettings() (*Settings, error) {
loop:
	for _, v := range settingsPaths {
		settings, err := LoadSettingsPath(v)
		if err != nil {
			perr := err.(*os.PathError)
			if strings.TrimSpace(perr.Err.Error()) == "no such file or directory" {
				continue loop
			}
			return nil, err
		}

		return settings, nil
	}

	defaultPath := settingsPaths[0]
	dir := filepath.Dir(defaultPath)
	err := os.MkdirAll(os.ExpandEnv(dir), 0755)
	if err != nil {
		return nil, err
	}

	file, err := os.Create(os.ExpandEnv(defaultPath))
	if err != nil {
		return nil, err
	}

	settings := DefaultSettings()
	settings.LogPath = filepath.Join(dir, "log")

	js, err := json.MarshalIndent(settings, "", "    ")
	if err != nil {
		return nil, err
	}

	js = append(js, '\n')

	_, err = file.Write(js)
	if err != nil {
		return nil, err
	}

	return nil, ErrNewSettingsFile
}
