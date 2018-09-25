package Messenger

import (
	"bytes"
	"image"
	"image/jpeg"

	"github.com/BlueMonday/go-scryfall"
	"github.com/mitsuse/pushbullet-go"
	"github.com/mitsuse/pushbullet-go/requests"
)

type Pushbullet struct {
	ApiToken string `yaml:"api_token"`
	Channel  string `yaml:"channel"`
	client   *pushbullet.Pushbullet
}

func (pb *Pushbullet) Init() {
	pb.client = pushbullet.New(pb.ApiToken)
}

func (pb *Pushbullet) SendCard(card scryfall.Card) error {
	var err error

	request := requests.NewFile()
	request.Title = CreateMessengeTitle(card)
	request.Body = CreateMessengeBody(card)
	request.FileUrl, err = GetImageUrl(card, pb)
	request.ChannelTag = pb.Channel
	request.FileType = "image/jpeg"
	request.FileName = GetImageName(card)
	if err != nil {
		return err
	}

	_, err = pb.client.PostPushesFile(request)
	return err
}

func (pb *Pushbullet) UploadImage(img image.Image, imgName string) (string, error) {
	uploadResonse, err := pb.client.PostUploadRequest(imgName, "image/jpeg")
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer

	err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 80})
	if err != nil {
		return "", err
	}

	return uploadResonse.FileUrl, pushbullet.Upload(pb.client.Client(), uploadResonse, &buf)
}
