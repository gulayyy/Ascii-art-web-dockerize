package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/template"
)

func main() {
	// HTML şablonunu dosyadan yükle
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	var err error
	tmpl, err := template.ParseFiles("templates/index.html") // index.html kodumuzu çekmemizi sağladı
	if err != nil {
		fmt.Println("HTML şablonu yüklenirken bir hata oluştu", err)
		return
	}
	// pointer olarak kullanmamızın sebebi formdaki verilen kopyalanmamsını ve veri boyutunun aşılmamasını engelliyor
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { // handelfunc url yoluna http isteklerini işler. Burdaki '\'  http isteklerini ele alır ve ilk olarak metodu kontrol eder.
		if r.Method != http.MethodPost { // w http.ResponseWriter: http yanıtını oluşturmak için  kullanılır. r *http.Request: İstemciden gelen bilgileri alır.
			// HTML şablonunu kullanarak sayfayı gönder
			err := tmpl.Execute(w, nil) // Bu metod, belirtilen HTML şablonunu verilen veri setiyle doldurur ve çıktıyı w adlı http.ResponseWriter nesnesine yazar.
			if err != nil {
				fmt.Println("Şablonu sayfaya yazarken bir hata oluştu:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		} // post işlemlerimizi yapmamızı sağlarken get fonksiyonu saadece ilk açılan sayfamızın açılmasını sağlıyor

		ad := r.FormValue("Enter Your Text")
		secim := r.FormValue("Select style")
		for _, char := range ad {
			if strings.Contains("öüğışçİÖÜĞŞÇ", string(char)) {
				http.Error(w, "Bad request(lütfen geçerli istekler girin)", http.StatusBadRequest)
				return
			}
		}

		asciiArt := AScii(ad, secim) // geçtiğimiz diğer sayfaya ait hmtl kodu
		fmt.Fprintf(w, `    
			<!DOCTYPE html>
			<html lang="en">
			<head>
				<meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<title>ASCII Art Result</title>
			</head>
			<body>
				<h1></h1>
				<pre>%s</pre>
				</body>
				</html>
			`, asciiArt)
	})

	// Web sunucusunu başlat
	fmt.Println("Port Çalişiyor...")
	http.ListenAndServe(":8080", nil)
}

func AScii(ad string, secim string) string {
	word := ad
	words := secim
	var asciiArt strings.Builder

	word = strings.ReplaceAll(word, "\\n", "\n")

	var content string
	if words == "standard" {
		dosya, err := os.ReadFile("standard.txt")
		if err != nil {
			asciiArt.WriteString("Not Found, if nothing is found, for example templates or banners.")
		}
		content = string(dosya)
	} else if words == "shadow" {
		dosya, err := os.ReadFile("shadow.txt")
		if err != nil {
			asciiArt.WriteString("Not Found, if nothing is found, for example templates or banners.")
		}
		content = string(dosya)
	} else if words == "thinkertoy" {
		dosya, err := os.ReadFile("thinkertoy.txt")
		if err != nil {
			asciiArt.WriteString("Not Found, if nothing is found, for example templates or banners.")
		}
		content = string(dosya)
	}

	s := strings.Split(content, "\n")

	for i, satir := range strings.Split(word, "\n") {
		if satir == "" {
			if i != 0 {
				asciiArt.WriteString("\n")
			}
			continue
		}

		for h := 1; h < 9; h++ {
			for _, char := range satir {
				printAsciiArtForCharacter(&asciiArt, char, h, s)
			}
			asciiArt.WriteString("\n")
		}
	}
	return asciiArt.String()
}

func printAsciiArtForCharacter(asciiArt *strings.Builder, char rune, h int, satir []string) {
	index := (int(char) - 32) * 9

	if index >= 0 && index+8 <= len(satir) {
		asciiArt.WriteString(satir[index+h])
	}
}
