package rigabe

import (
	"net/http"
	"encoding/base64"
	"io/ioutil"
	"log"
	"image"
	_ "image/jpeg"
	_ "image/gif"
	_ "image/png"
	"strings"
	"encoding/json"
	prominentcolor "github.com/EdlinOrg/prominentcolor"
)

type Response struct {
	Color []string `json:"colors"`
}

type Request struct {
	Image string `json:"image"`
}

func getImage(data string) (image.Image, error) {
  reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data))
  img, _, err := image.Decode(reader)

  if err != nil {
	return nil, err
  }

  return img, nil
}

func process(img image.Image) Response {
	resizeSize := uint(prominentcolor.DefaultSize)
	bgmasks := prominentcolor.GetDefaultMasks()

	res, err := prominentcolor.KmeansWithAll(3, img, prominentcolor.ArgumentDefault, resizeSize, bgmasks)
	if err != nil {
		log.Println(err)
	}

	colors := Response{getColors(res)}

	return colors
}

func getColors(colorRange []prominentcolor.ColorItem) []string {
	var colors []string

	for _, color := range colorRange {
		colors = append(colors, color.AsString())
	}

	return colors
}

func Rigabe(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var request Request
	json.Unmarshal(body, &request)

	image, err := getImage(request.Image)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}

	result := process(image)

	colors, err := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.Write(colors)
}
