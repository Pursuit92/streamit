package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type Streamer struct {
	Settings       *Settings
	Cmd            *exec.Cmd
	LogOut, LogErr io.WriteCloser
	Out, Err       *log.Logger
	Notify         bool
	demo           bool
}

func NewStreamer(settings *Settings, demo bool) (*Streamer, error) {
	if err := settings.validate(); err != nil {
		if !(err == ErrNoStreamKey && demo) {
			return nil, err
		}
	}

	st := &Streamer{Settings: settings, demo: demo}

	err := st.getLogs()
	if err != nil {
		return nil, err
	}

	st.buildCmd()
	st.checkNotify()

	return st, err
}

func (st *Streamer) buildCmd() {
	opts := []string{
		"-f", "x11grab",
		"-s", st.Settings.Res,
		"-r", fmt.Sprint(st.Settings.FPS),
		"-i", ":0.0",
		"-f", "alsa",
		"-i", "pulse",
		"-f", "flv",
		"-ac", "2",
		"-ar", fmt.Sprint(st.Settings.AudioRate),
		"-vcodec", "libx264",
		"-g", fmt.Sprint(st.Settings.FPS * 2),
		"-keyint_min", fmt.Sprint(st.Settings.FPS),
		"-b:v", fmt.Sprintf("%dk", st.Settings.CBR),
		"-minrate", fmt.Sprintf("%dk", st.Settings.CBR),
		"-maxrate", fmt.Sprintf("%dk", st.Settings.CBR),
		"-pix_fmt", "yuv420p",
		"-s", st.Settings.Res,
		"-preset", st.Settings.Quality,
		"-tune", "film",
		"-acodec", "libmp3lame",
		"-threads", fmt.Sprint(st.Settings.Threads),
		"-strict", "normal",
		"-bufsize", fmt.Sprintf("%dk", st.Settings.CBR),
		"-loglevel", st.Settings.LogLevel,
	}

	if st.demo {
		opts = append(opts, "-")
	} else {
		opts = append(opts, fmt.Sprintf("rtmp://%s.twitch.tv/app/%s", st.Settings.Server, st.Settings.StreamKey))
	}

	cmd := exec.Command("ffmpeg",
		opts...,
	)

	// fmt.Println(strings.Join(cmd.Args, " "))
	// os.Exit(1)

	st.Cmd = cmd
}

func (st *Streamer) run() error {
	if st.demo {
		mplayer := exec.Command("mplayer", "-")
		mpIn, err := mplayer.StdinPipe()
		if err != nil {
			return err
		}
		st.Cmd.Stdout = mpIn
		mplayer.Stdout = st.LogOut
		mplayer.Stderr = st.LogErr
		mplayer.Start()
	} else {
		st.Cmd.Stdout = st.LogOut
	}

	st.Cmd.Stderr = st.LogErr

	defer func() {
		st.LogOut.Close()
		st.LogErr.Close()
	}()

	st.notify("Starting streamit!")
	err := st.Cmd.Start()
	if err != nil {
		return err
	}

	err = st.Cmd.Wait()
	st.notify("Streamit exited!")

	return err
}

func (st *Streamer) getLogs() error {
	logDir := os.ExpandEnv(st.Settings.LogPath)
	info, err := os.Stat(logDir)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		err := os.MkdirAll(logDir, 0755)
		if err != nil {
			return err
		}
	} else {
		if !info.IsDir() {
			return fmt.Errorf("log path is not a directory: %s", logDir)
		}
	}

	logOut, err := os.Create(filepath.Join(logDir, "streamit.out.log"))
	if err != nil {
		return err
	}

	logErr, err := os.Create(filepath.Join(logDir, "streamit.err.log"))
	if err != nil {
		return err
	}

	st.LogOut = logOut
	st.LogErr = logErr
	st.Out = log.New(logOut, "[streamit] ", log.LstdFlags)
	st.Err = log.New(logErr, "[streamit] ", log.LstdFlags)

	return nil
}

func (st *Streamer) checkNotify() {
	_, err := exec.LookPath("notify-send")
	if err == nil {
		st.Notify = true
		return
	}
	st.Err.Println("notify-send not found - notifications disabled")
}

func (st *Streamer) notify(f string, rest ...interface{}) {
	st.Out.Printf(f, rest...)
	if st.Notify {
		exec.Command("notify-send", fmt.Sprintf(f, rest...)).Run()
	}
}
