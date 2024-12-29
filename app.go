package main

import (
	"crypto/tls"
	"embed"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

//go:embed templates/*
var resources embed.FS

var t = template.Must(template.ParseFS(resources, "templates/*"))

var domain = "https://www.monsterindeks.no"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("App live and listening on port:", port)

	http.HandleFunc("/", RootHandler)
	http.HandleFunc("/ping", PingHandler)
	http.HandleFunc("/health", HealthHandler)
	http.HandleFunc("/robots.txt", RobotsHandler)
	http.HandleFunc("/sitemap.xml", SitemapHandler)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func RootHandler(w http.ResponseWriter, _ *http.Request) {
	monster := getMonsterData()

	data := map[string]string{
		"monsterName":  monster.FullName,
		"monsterPrice": monster.GrossPrice,
	}

	_ = t.ExecuteTemplate(w, "index.html.tmpl", data)
}

func PingHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("pong"))
}

func HealthHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("I'm healthy"))
}

func RobotsHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	data := map[string]string{
		"URL": domain + "/sitemap.xml",
	}

	_ = t.ExecuteTemplate(w, "robots.txt.tmpl", data)
}

func SitemapHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/xml")

	data := map[string]string{
		"URL": domain,
	}

	_ = t.ExecuteTemplate(w, "sitemap.xml.tmpl", data)
}

func getMonsterData() Monster {
	odaMonsterEndpoint := "https://oda.com/tienda-web-api/v1/products/23300/"

	client := createInsecureHTTPClient()

	response, err := client.Get(odaMonsterEndpoint)
	if err != nil {
		log.Fatal(err.Error())
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err.Error())
	}

	var monster Monster
	err = json.Unmarshal(responseData, &monster)
	if err != nil {
		log.Fatal(err.Error())
	}

	return monster
}

func createInsecureHTTPClient() *http.Client {
	customTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &http.Client{Transport: customTransport}
}

type Monster struct {
	Id             int    `json:"id"`
	FullName       string `json:"full_name"`
	GrossPrice     string `json:"gross_price"`
	GrossUnitPrice string `json:"gross_unit_price"`
}
