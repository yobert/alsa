[![](https://godoc.org/github.com/yobert/alsa?status.svg)](https://godoc.org/github.com/yobert/alsa)

Synopsis
--------
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

For a simple example of synthesized playback, the beep command will produce
a sine wave for a few seconds on each detected ALSA output:

    go get github.com/yobert/alsa/cmd/beep
    $GOPATH/beep

And for recording from a microphone into a WAV file:

    go get github.com/yobert/alsa/cmd/record
    $GOPATH/record

This example does recording and playback, but it's got a really
buggy ring buffer going on:

    go get github.com/yobert/alsa/cmd/echoback
    $GOPATH/echoback

Disclaimer
----------
This module makes syscalls with pointers to memory buffers that are in garbage collectable memory. I have a feeling this isn't safe, but it hasn't crashed on me yet.

Contributors
------------
Thanks so much for the help! Thanks! See AUTHORS for a list. Pull requests
welcome from anybody, regardless of skill level.

See Also
--------
You may be interested in https://github.com/jfreymuth/pulse which is a lot less likely to crash and will probably work with your sound card.
