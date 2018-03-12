package alsa

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/yobert/alsa/alsatype"
)

type Card struct {
	Path   string
	Title  string
	Number int

	fh       *os.File
	pversion alsatype.PVersion
	cardinfo cardInfo
}

func (card Card) String() string {
	return card.Title
}

func OpenCards() ([]*Card, error) {
	dir := "/dev/snd/"
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var ret []*Card
	for _, fi := range fis {
		var n int
		if _, err := fmt.Sscanf(fi.Name(), "controlC%d", &n); err != nil {
			continue
		}

		cardPath := path.Join(dir, fi.Name())
		fh, err := os.Open(cardPath)
		if err != nil {
			return ret, err
		}

		card := Card{
			Path:   cardPath,
			Number: n,
			fh:     fh,
		}

		err = ioctl(fh.Fd(), ioctl_encode_ptr(cmdRead, &card.pversion, cmdControlVersion), &card.pversion)
		if err != nil {
			return ret, err
		}

		err = ioctl(fh.Fd(), ioctl_encode_ptr(cmdRead, &card.cardinfo, cmdControlCardInfo), &card.cardinfo)
		if err != nil {
			return ret, err
		}

		card.Title = gstr(card.cardinfo.Name[:])
		ret = append(ret, &card)
	}

	return ret, nil
}

func CloseCards(cards []*Card) {
	for _, card := range cards {
		card.fh.Close()
	}
}
