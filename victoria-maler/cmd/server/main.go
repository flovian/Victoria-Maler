package main

import (
    "html/template"
    "log"
    "os"
    "net/http"
    "path/filepath"
)

func main() {
    // Static assets
    fs := http.FileServer(http.Dir("web/static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    // Routes
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        // decide which page to load
        path := r.URL.Path
        if path == "/" || path == "" {
            path = "/home.html"
        }
        page := filepath.Base(path) // e.g. home.html
        name := page[:len(page)-len(filepath.Ext(page))]

        // build template set for this request: layouts + the requested page
        layouts := []string{
            "web/templates/layouts/base.html",
            "web/templates/layouts/header.html",
            "web/templates/layouts/footer.html",
        }
        pagePath := "web/templates/pages/" + name + ".html"
        files := append(layouts, pagePath)

        tpl, err := template.ParseFiles(files...)
        if err != nil {
            http.Error(w, "Template parse error", http.StatusInternalServerError)
            log.Printf("template parse error for %s: %v", name, err)
            return
        }

        if err := tpl.ExecuteTemplate(w, "base.html", nil); err != nil {
            http.Error(w, "Template render error", http.StatusInternalServerError)
            log.Printf("template execute error for %s: %v", name, err)
        }
    })

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    addr := ":" + port
    log.Printf("Starting server on %s — open http://localhost:%s/", addr, port)
    if err := http.ListenAndServe(addr, nil); err != nil {
        log.Fatalf("server error: %v", err)
    }
}

