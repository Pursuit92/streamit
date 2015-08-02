package main

import (
	"fmt"
	"os"
)

func main() {
	settings, err := LoadSettings()
	if err != nil {
		fmt.Println(err)
		if err == ErrNewSettingsFile {
			os.Exit(0)
		}
		os.Exit(1)
	}
	st, err := NewStreamer(settings)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	st.run()
}
