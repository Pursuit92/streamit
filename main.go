package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	demo := flag.Bool("test", false, "run in test mode")
	flag.Parse()
	settings, err := LoadSettings()
	if err != nil {
		fmt.Println(err)
		if err == ErrNewSettingsFile {
			os.Exit(0)
		}
		os.Exit(1)
	}
	st, err := NewStreamer(settings, *demo)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	st.run()
}
