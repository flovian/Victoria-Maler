package handlers

import (
	"net/http"
	"time"

	"ecochain-victoria/internal/services"
	"ecochain-victoria/internal/utils"
)

func (h *HandlerSet) EvidenceTrackerPage(w http.ResponseWriter, r *http.Request) {
	campaigns, err := h.Campaigns.ListCampaigns()
	if err != nil {
		utils.InternalError(w, "could not load campaigns")
		return
	}
	h.Render.Page(w, r, "evidence_tracker", map[string]interface{}{
		"Campaigns": campaigns,
	})
}

type submitCleanupRequest struct {
	CampaignID  int64   `json:"campaign_id"`
	Title       string  `json:"title"`
	Notes       string  `json:"notes"`
	Location    string  `json:"location"`
	CleanupDate string  `json:"cleanup_date"`
	WasteKg     float64 `json:"waste_kg"`
	Volunteers  int     `json:"volunteers"`
}

func (h *HandlerSet) SubmitCleanup(w http.ResponseWriter, r *http.Request) {
	claims, ok := h.claims(r)
	if !ok {
		utils.Unauthorized(w, "authentication required")
		return
	}

	var req submitCleanupRequest
	if err := decodeJSON(r, &req); err != nil {
		utils.BadRequest(w, "invalid request body")
		return
	}

	report, err := h.Cleanups.SubmitReport(services.SubmitCleanupInput{
		CampaignID:  req.CampaignID,
		Title:       req.Title,
		Notes:       req.Notes,
		Location:    req.Location,
		CleanupDate: req.CleanupDate,
		WasteKg:     req.WasteKg,
		Volunteers:  req.Volunteers,
		CreatedBy:   claims.UserID,
	})
	if err != nil {
		utils.BadRequest(w, err.Error())
		return
	}
	utils.Created(w, report)
}

func (h *HandlerSet) CampaignCleanups(w http.ResponseWriter, r *http.Request) {
	campaignID, err := queryID(r, "campaign_id")
	if err != nil {
		utils.BadRequest(w, "campaign_id required")
		return
	}
	reports, err := h.Cleanups.ListByCampaign(campaignID)
	if err != nil {
		utils.InternalError(w, "could not load cleanup reports")
		return
	}
	utils.OK(w, reports)
}

func (h *HandlerSet) VerifyCleanup(w http.ResponseWriter, r *http.Request) {
	reportID, err := queryID(r, "id")
	if err != nil {
		utils.BadRequest(w, "id required")
		return
	}
	report, err := h.Cleanups.GetReport(reportID)
	if err != nil {
		utils.NotFound(w, "report not found")
		return
	}
	utils.OK(w, h.Cleanups.VerifyReport(report))
}

func (h *HandlerSet) SubmitEvidence(w http.ResponseWriter, r *http.Request) {
	if _, ok := h.claims(r); !ok {
		utils.Unauthorized(w, "authentication required")
		return
	}

	if err := r.ParseMultipartForm(12 << 20); err != nil {
		utils.BadRequest(w, "could not parse upload")
		return
	}

	campaignID, _ := queryID(r, "campaign_id")
	reportID, _ := queryID(r, "report_id")
	caption := r.FormValue("caption")

	file, header, err := r.FormFile("file")
	if err != nil {
		utils.BadRequest(w, "file is required")
		return
	}
	defer file.Close()

	fileBytes := make([]byte, 0)
	buf := make([]byte, 4096)
	for {
		n, err := file.Read(buf)
		fileBytes = append(fileBytes, buf[:n]...)
		if err != nil {
			break
		}
		if len(fileBytes) > 10<<20 {
			utils.BadRequest(w, "file exceeds 10 MB limit")
			return
		}
	}

	evidence, err := h.Evidences.SubmitEvidence(services.SubmitEvidenceInput{
		CampaignID: campaignID,
		ReportID:   reportID,
		FilePath:   header.Filename,
		FileType:   header.Header.Get("Content-Type"),
		Caption:    caption,
		FileData:   fileBytes,
		Timestamp:  time.Now().Unix(),
	})
	if err != nil {
		utils.BadRequest(w, err.Error())
		return
	}
	utils.Created(w, evidence)
}

func (h *HandlerSet) CampaignEvidence(w http.ResponseWriter, r *http.Request) {
	campaignID, err := queryID(r, "campaign_id")
	if err != nil {
		utils.BadRequest(w, "campaign_id required")
		return
	}
	evidences, err := h.Evidences.ListByCampaign(campaignID)
	if err != nil {
		utils.InternalError(w, "could not load evidence")
		return
	}
	utils.OK(w, evidences)
}
