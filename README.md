alsa
----

This is a golang ALSA client implementation, without cgo!

https://godoc.org/github.com/yobert/alsa <- docs

For a simple example of synthesized playback:

    go get github.com/yobert/alsa/cmd/beep
    $GOPATH/beep
    
This example does recording and playback, but it's got a really
buggy ring buffer going on:

    go get github.com/yobert/alsa/cmd/echoback
    $GOPATH/echoback

Both examples default to your first ALSA card and device.

contributors
------------
Thanks to the wonderful contributions of:

- Christopher Harrington (https://github.com/ironiridis)

