package main

import (
  "net/http"
  "encoding/base64"
  "io/ioutil"
  "fmt"
  "log"
  "image"
  _ "image/jpeg"
  "strings"
  prominentcolor "github.com/EdlinOrg/prominentcolor"
)

func getRandomString() string {
  return "teste.jpg"
}

func getImage(data string) (image.Image, error) {
  //var fileName = getRandomString() 

//  decode, err := base64.StdEncoding.DecodeString(content)
//  if err != nil {
//	log.Fatal(err)
//  }

  reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data))
  img, _, err := image.Decode(reader)

  if err != nil {
  	log.Fatal(err)
	return nil, err
  }

  //fmt.Println(img)

  //file, err := os.Create(fileName)
  //if err != nil {
  //  log.Fatal(err)
  //}
  //defer file.Close()

  //_, err = file.Write(decode)
  //if err != nil {
  //  log.Fatal(err)
  //}

  return img, nil 
}

func processBatch(img image.Image) string {
	var buff strings.Builder
	bitarr := []int{
		prominentcolor.ArgumentAverageMean,
		prominentcolor.ArgumentDefault,
	}

	params := 1 


//	prefix := fmt.Sprintf("K=%d, ", params)
	resizeSize := uint(prominentcolor.DefaultSize)
	bgmasks := prominentcolor.GetDefaultMasks()

	for i := 0; i < len(bitarr); i++ {
		res, err := prominentcolor.KmeansWithAll(params, img, bitarr[i], resizeSize, bgmasks)
		if err != nil {
			log.Println(err)
			continue
		}
		//buff.WriteString(outputTitle(prefix + bitInfo(bitarr[i])))
		buff.WriteString(outputColorRange(res))
	}

	return buff.String()
}

func outputColorRange(colorRange []prominentcolor.ColorItem) string {
	var buff strings.Builder
	buff.WriteString("<table><tr>")

	for _, color := range colorRange {
		fmt.Println(color.AsString())
		//fmt.Println(color.AsString(), color.AsString(), color.Cnt)
	}

	buff.WriteString("</tr></table>")
	return buff.String()
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
  teste := processBatch(image)
  fmt.Println(teste)


  w.Write([]byte("teste"))
}

func main() {
  http.HandleFunc("/", getColorsFromImage)

  if err := http.ListenAndServe(":8080", nil); err != nil {
    panic(err)
  }
}
