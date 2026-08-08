package handlers

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"ecochain-victoria/internal/bitcoin"
	"ecochain-victoria/internal/config"
	"ecochain-victoria/internal/middleware"
	"ecochain-victoria/internal/repositories"
	"ecochain-victoria/internal/services"
)

// HandlerSet bundles every dependency handlers need.
type HandlerSet struct {
	Config         *config.Config
	Render         *Renderer
	Auth           *services.AuthService
	Campaigns      *services.CampaignService
	Donations      *services.DonationService
	Cleanups       *services.CleanupService
	Evidences      *services.EvidenceService
	BitcoinRecords repositories.BitcoinRecordRepository
	Wallet         *bitcoin.Wallet
}

// Renderer renders pages from the layouts + page template set.
type Renderer struct{}

func NewRenderer() *Renderer { return &Renderer{} }

func (r *Renderer) Page(w http.ResponseWriter, req *http.Request, page string, data map[string]interface{}) {
	if data == nil {
		data = map[string]interface{}{}
	}
	data["Content"] = page

	if claims, ok := middleware.Claims(req); ok {
		data["IsLoggedIn"] = true
		data["UserName"] = claims.Email
		data["UserRole"] = claims.Role
	} else {
		data["IsLoggedIn"] = false
		data["UserName"] = ""
		data["UserRole"] = ""
	}

	layouts := []string{
		"web/templates/layouts/base.html",
		"web/templates/layouts/header.html",
		"web/templates/layouts/footer.html",
	}
	pagePath := "web/templates/pages/" + page + ".html"
	files := append(layouts, pagePath)

	tpl, err := template.ParseFiles(files...)
	if err != nil {
		log.Printf("render %s: template parse error: %v", page, err)
		http.Error(w, "Template parse error", http.StatusInternalServerError)
		return
	}

	if err := tpl.ExecuteTemplate(w, "base.html", data); err != nil {
		log.Printf("render %s: execute error: %v", page, err)
		http.Error(w, "Template render error", http.StatusInternalServerError)
	}
}

func pageNameFromPath(path string) string {
	base := filepath.Base(path)
	return base[:len(base)-len(filepath.Ext(base))]
}
