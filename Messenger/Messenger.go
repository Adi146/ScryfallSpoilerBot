package Messenger

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"net/http"
	"strings"

	"github.com/BlueMonday/go-scryfall"
)

type Messenger interface {
	SendCard(card scryfall.Card) error
	UploadImage(img image.Image, imgName string) (string, error)
}

func CreateMessengeTitle(card scryfall.Card) string {
	return card.Name
}

func CreateMessengeBody(card scryfall.Card) string {
	if card.CardFaces == nil {
		return fmt.Sprintf("%s\n\n%s", card.TypeLine, card.OracleText)
	} else {
		typeLines := make([]string, len(card.CardFaces))
		oracleTexts := make([]string, len(card.CardFaces))
		for index, face := range card.CardFaces {
			typeLines[index] = face.TypeLine
			oracleTexts[index] = *face.OracleText
		}
		return fmt.Sprintf("%s\n\n%s", strings.Join(typeLines, " // "), strings.Join(oracleTexts, "\n//\n"))
	}
}

func GetImageUrl(card scryfall.Card, messenger Messenger) (string, error) {
	if card.ImageURIs != nil {
		return card.ImageURIs.Normal, nil
	} else {
		images := make([]image.Image, 0)
		for _, face := range card.CardFaces {
			img, err := downloadImage(face.ImageURIs.Normal)
			if err != nil {
				return "", err
			}
			images = append(images, img)
		}
		return messenger.UploadImage(joinImages(images), GetImageName(card))
	}
}

func GetImageName(card scryfall.Card) string {
	return card.Name + ".jpg"
}

func downloadImage(url string) (image.Image, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return jpeg.Decode(resp.Body)
}

func joinImages(images []image.Image) image.Image {
	sumWidth := 0
	maxHeight := 0
	for _, img := range images {
		sumWidth += img.Bounds().Size().X
		if img.Bounds().Size().Y > maxHeight {
			maxHeight = img.Bounds().Size().Y
		}
	}

	dstImg := image.NewRGBA(image.Rect(0, 0, sumWidth, maxHeight))

	minPoint := image.Point{0, 0}
	for _, img := range images {
		location := image.Rectangle{minPoint, image.Point{minPoint.X + img.Bounds().Dx(), img.Bounds().Dy()}}
		draw.Draw(dstImg, location, img, image.Point{0, 0}, draw.Src)

		minPoint = image.Point{location.Max.X, 0}
	}

	return dstImg
}
