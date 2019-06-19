package main

import (
  "net/http"
  "encoding/base64"
  "io/ioutil"
  "log"
  "image"
  _ "image/jpeg"
  "strings"
  "encoding/json"
  prominentcolor "github.com/EdlinOrg/prominentcolor"
)

type Colors []interface{}

func getImage(data string) (image.Image, error) {
  reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data))
  img, _, err := image.Decode(reader)

  if err != nil {
  	log.Fatal(err)
	return nil, err
  }

  return img, nil
}

func process(img image.Image) Colors {
	resizeSize := uint(prominentcolor.DefaultSize)
	bgmasks := prominentcolor.GetDefaultMasks()

	res, err := prominentcolor.KmeansWithAll(3, img, prominentcolor.ArgumentDefault, resizeSize, bgmasks)
	if err != nil {
		log.Println(err)
	}

	colors := Colors{getColors(res)}

	return colors
}

func getColors(colorRange []prominentcolor.ColorItem) []string {
	var colors []string

	for _, color := range colorRange {
		colors = append(colors, color.AsString())
	}

	return colors
}

func getColorsFromImage(w http.ResponseWriter, r *http.Request) {
  body, err := ioutil.ReadAll(r.Body)
  defer r.Body.Close()
  if err != nil {
  	http.Error(w, err.Error(), 500)
  	return
  }

  bodyString := string(body)

  image, err := getImage(bodyString)
  result := process(image)

  colors, err := json.Marshal(result)
  w.Header().Set("Content-Type", "application/json")
  w.Write(colors)
}

func main() {
  http.HandleFunc("/", getColorsFromImage)

  if err := http.ListenAndServe(":8080", nil); err != nil {
    panic(err)
  }
}
