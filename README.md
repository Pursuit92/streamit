# Streamit

Simple Twitch.tv streaming app for Linux. Pretty much just a simple wrapper
around ffmpeg.

## Installing

Assuming you have Go installed and your $GOPATH set up, simply do:

```bash
go get github.com/Pursuit92/streamit
```

Binaries to come!

## Usage

On first run, streamit will create a settings.json file at
```$HOME/.config/streamit/settings.json```. This should be sufficient for most
people starting out, you just need to fill in your Twitch stream key.

Log files will default to being created at ```$HOME/.config/streamit/log/```.
These are recreated at every run and shouldn't take up much space.

Streaming will start a few seconds after you run ```streamit``` and will end as
soon as it exits. Currently, the best way to see that you're *actually*
streaming is to check your Twitch dashboard since ffmpeg doesn't really provide
any "streaming-started" notification.
