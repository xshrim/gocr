package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"gocr/utils"
	"strconv"
	"strings"
	"time"
)

func Test(w http.ResponseWriter, r *http.Request) {
	am := "GET"
	if r.Method != am {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Only %s Method Allowed!\n", am)
		return
	}

	delay := strings.TrimPrefix(r.URL.Path, "/test/")
	if r.URL.Path == "/test" {
		delay = ""
	}

	if delay != "" {
		var err error
		sec, err := strconv.Atoi(delay)
		dur := time.Duration(sec) * time.Second
		if err != nil {
			dur, err = time.ParseDuration(delay)
			if err != nil {
				return
			}
		}
		time.Sleep(dur)
		fmt.Fprintf(w, "(%s later) ", dur)
	}

	fmt.Fprintf(w, "[Headers]:\n")

	for name, headers := range r.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func File(w http.ResponseWriter, r *http.Request) {
	am := "POST"
	if r.Method != am {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Only %s Method Allowed!\n", am)
		return
	}

	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Retrieving File Failed: %v\n", err)
		return
	}
	defer file.Close()

	body := r.FormValue("json")

	bd := make(map[string]string)
	_ = json.Unmarshal([]byte(body), &bd)
	languages, ok := bd["languages"]
	if !ok {
		languages = "eng"
	}
	whitelist, ok := bd["whitelist"]
	if !ok {
		whitelist = "0123456789abcdefghijkmnpqrstuvwxyzABCDEFGHIJKLMNPQRSTUVWXYZ"
	}

	imgBytes, err := ioutil.ReadAll(file)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Reading File Failed: %v\n", err)
		return
	}

	log.Printf("Data: %s %d | Args: %s %s\n", handler.Filename, handler.Size, languages, whitelist)

	text, err := utils.CaptchaOcr(imgBytes, languages, whitelist)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Image Recognize Failed: %v\n", err)
		return
	}

	fmt.Fprintf(w, text)

}

func Base64(w http.ResponseWriter, r *http.Request) {
	am := "POST"
	if r.Method != am {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Only %s Method Allowed!\n", am)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Retrieving Body Failed: %v\n", err)
		return
	}

	bd := make(map[string]string)
	if err := json.Unmarshal(body, &bd); err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Reading Body Failed: %v\n", err)
		return
	}
	languages, ok := bd["languages"]
	if !ok {
		languages = "eng"
	}
	whitelist, ok := bd["whitelist"]
	if !ok {
		whitelist = "0123456789abcdefghijkmnpqrstuvwxyzABCDEFGHIJKLMNPQRSTUVWXYZ"
	}

	imgBytes, err := utils.Base64Decode(bd["data"])
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Decoding Body Failed: %v\n", err)
		return
	}

	log.Printf("Data: %s %d | Args: %s %s\n", "base64", len(imgBytes), languages, whitelist)

	text, err := utils.CaptchaOcr(imgBytes, languages, whitelist)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Image Recognize Failed: %v\n", err)
		return
	}

	fmt.Fprintf(w, text)

}
