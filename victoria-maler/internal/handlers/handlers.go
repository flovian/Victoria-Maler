package handlers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"ecochain-victoria/internal/bitcoin"
	"ecochain-victoria/internal/config"
	"ecochain-victoria/internal/middleware"
	"ecochain-victoria/internal/repositories"
	"ecochain-victoria/internal/services"
	"ecochain-victoria/internal/utils"
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

func funcMap() template.FuncMap {
	return template.FuncMap{
		"progressPercent": func(raised, target float64) int {
			if target <= 0 {
				return 0
			}
			pct := int(raised / target * 100)
			if pct > 100 {
				return 100
			}
			return pct
		},
	}
}

func (r *Renderer) Page(w http.ResponseWriter, req *http.Request, page string, data map[string]interface{}) {
	if data == nil {
		data = map[string]interface{}{}
	}

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

	tpl, err := template.New("layout").Funcs(funcMap()).ParseFiles(files...)
	if err != nil {
		log.Printf("render %s: template parse error: %v", page, err)
		http.Error(w, "Template parse error", http.StatusInternalServerError)
		return
	}

	if err := tpl.ExecuteTemplate(w, "layout", data); err != nil {
		log.Printf("render %s: execute error: %v", page, err)
		http.Error(w, "Template render error", http.StatusInternalServerError)
	}
}

func pageNameFromPath(path string) string {
	base := filepath.Base(path)
	return base[:len(base)-len(filepath.Ext(base))]
}

func (h *HandlerSet) claims(r *http.Request) (*utils.Claims, bool) {
	return middleware.Claims(r)
}

func decodeJSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func queryID(r *http.Request, key string) (int64, error) {
	return strconv.ParseInt(r.URL.Query().Get(key), 10, 64)
}
