package main

import (
	"crypto/tls"
	"embed"
	"html/template"
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
	_ = t.ExecuteTemplate(w, "index.html.tmpl", nil)
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

func createInsecureHTTPClient() *http.Client {
	customTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &http.Client{Transport: customTransport}
}
