package utils

import "net"
import "fmt"
import "bytes"
import "net/http"
import "regexp"
import "encoding/json"
import "encoding/base64"
import "io/ioutil"

func TrimSpace(str string) string {
	return regexp.MustCompile("\\s+").ReplaceAllString(str, "")
}

func GetLocalIP() string {
	if addrs, err := net.InterfaceAddrs(); err == nil {
		for _, address := range addrs {
			// check the address type and if it is not a loopback the display it
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}

func Base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func Base64Decode(str string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(str)
}

func Request(str, whitelist string) string {
	wl := "0123456789abcdefghijkmnpqrstuvwxyzABCDEFGHIJKLMNPQRSTUVWXYZ"
	if whitelist != "" {
		wl = whitelist
	}
	url := "http://66.135.5.57:8080/base64"

	jsonstr := fmt.Sprintf(`{"base64": "%s", "trim": "\n", "languages": "eng", "whitelist": "%s"}`, str, wl)

	req, _ := http.NewRequest("POST", url, bytes.NewBufferString(jsonstr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	// fmt.Println("status", resp.Status)
	// fmt.Println("response:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	data := make(map[string]string)
	json.Unmarshal(body, &data)

	reg := regexp.MustCompile("\\s+")
	return reg.ReplaceAllString(data["result"], "")
}
