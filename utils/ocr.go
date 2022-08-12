package utils

import (
	"io/ioutil"
	"os"
	"strings"
	//"github.com/otiai10/gosseract"
	"github.com/otiai10/gosseract/v2"
)

func Ocr(str, languages, whitelist string) (string, error) {
	trim := "\n"

	client := gosseract.NewClient()
	defer client.Close()

	tempfile, err := ioutil.TempFile("", "ocrserver"+"-")
	if err != nil {
		return "", err
	}
	defer func() {
		tempfile.Close()
		os.Remove(tempfile.Name())
	}()

	b, err := Base64Decode(str)
	if err != nil {
		return "", err
	}

	tempfile.Write(b)

	client.SetImage(tempfile.Name())

	if languages != "" {
		client.Languages = strings.Split(languages, ",")
	}

	if whitelist != "" {
		client.SetWhitelist(whitelist)
	}

	text, err := client.Text()
	if err != nil {
		return "", err
	}

	return strings.Trim(text, trim), nil

}

func CaptchaOcr(imgBytes []byte, languages, whitelist string) (string, error) {
	text, err := Ocr(ImageProcess(imgBytes), languages, whitelist)
	if err != nil {
		return "", err
	}

	return TrimSpace(text), nil
}
