package main

import (
  "net/http"
  "encoding/base64"
  "io/ioutil"
  "log"
  "image"
  _ "image/jpeg"
  _ "image/gif"
  _ "image/png"
  "html/template"
  "strings"
  "encoding/json"
  prominentcolor "github.com/EdlinOrg/prominentcolor"
)

// An Response represents data result
type Response struct {
	Color []string `json:"colors"`
}

// An Request represents an image encoded by base64
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

func getColorsFromImage(w http.ResponseWriter, r *http.Request) {
  body, err := ioutil.ReadAll(r.Body)
  defer r.Body.Close()
  if err != nil {
  	http.Error(w, err.Error(), 500)
  	return
  }

  var request Request
  json.Unmarshal([]byte(body), &request)

  image, err := getImage(request.Image)
  result := process(image)

  colors, err := json.Marshal(result)
  w.Header().Set("Content-Type", "application/json")
  w.Write(colors)
}

func index(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}

func main() {
  http.HandleFunc("/upload", getColorsFromImage)
  http.HandleFunc("/", index)

  if err := http.ListenAndServe(":8080", nil); err != nil {
    panic(err)
  }
}
