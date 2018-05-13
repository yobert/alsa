[![](https://godoc.org/github.com/yobert/alsa?status.svg)](https://godoc.org/github.com/yobert/alsa)

alsa
----
This is a golang ALSA client implementation, without cgo! Unfortunately,
doing it without cgo means throwing away many years of compatibility work
that has been put into libalsa. So be warned, this library is not likely
to work with a lot of the more colorful audio cards out there, and is not
likely to work on platforms other than x86_64. (Though, someone has nicely
done some work on ARM. Thanks!)

But fear not! Go is fun, and I tried to keep the library on the simple
side, so adding in support for what your audio card needs might actually
be just a nice afternoon of programming. The hardest part for me was just
trying to understand all of the alsa terminology.

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
Thanks so much for the help! Thanks! See AUTHORS for a list. Pull requests
welcome from anybody, regardless of skill level.

