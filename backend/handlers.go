package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"gorm.io/gorm"
)

type Server struct {
	db *gorm.DB

	unlockMu      sync.Mutex
	pendingUnlock bool
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

// POST /api/v1/access — dipanggil firmware ESP8266 setiap kartu ditap.
func (s *Server) handleAccess(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UID string `json:"uid"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.UID == "" {
		writeError(w, http.StatusBadRequest, "uid is required")
		return
	}

	var card Card
	granted := s.db.Where("uid = ?", req.UID).First(&card).Error == nil

	log := AccessLog{UID: req.UID, Label: card.Label, Granted: granted}
	s.db.Create(&log)

	writeJSON(w, http.StatusOK, map[string]any{
		"granted": granted,
		"label":   card.Label,
	})
}

// GET /api/v1/cards
func (s *Server) handleListCards(w http.ResponseWriter, r *http.Request) {
	var cards []Card
	s.db.Order("created_at desc").Find(&cards)
	writeJSON(w, http.StatusOK, cards)
}

// POST /api/v1/cards
func (s *Server) handleCreateCard(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UID   string `json:"uid"`
		Label string `json:"label"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.UID == "" || req.Label == "" {
		writeError(w, http.StatusBadRequest, "uid and label are required")
		return
	}

	card := Card{UID: req.UID, Label: req.Label}
	if err := s.db.Create(&card).Error; err != nil {
		writeError(w, http.StatusConflict, "uid already registered")
		return
	}

	writeJSON(w, http.StatusCreated, card)
}

// DELETE /api/v1/cards/{id}
func (s *Server) handleDeleteCard(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid card id")
		return
	}

	if err := s.db.Delete(&Card{}, id).Error; err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete card")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GET /api/v1/logs
func (s *Server) handleListLogs(w http.ResponseWriter, r *http.Request) {
	var logs []AccessLog
	s.db.Order("created_at desc").Limit(100).Find(&logs)
	writeJSON(w, http.StatusOK, logs)
}

// POST /api/v1/unlock — dipanggil dashboard (butuh login) untuk buka pintu dari jarak jauh.
func (s *Server) handleUnlock(w http.ResponseWriter, r *http.Request) {
	s.unlockMu.Lock()
	s.pendingUnlock = true
	s.unlockMu.Unlock()

	s.db.Create(&AccessLog{UID: "-", Label: "Remote Unlock", Granted: true})

	writeJSON(w, http.StatusOK, map[string]string{"status": "unlock requested"})
}

// GET /api/v1/unlock/poll — dipanggil firmware ESP8266 secara berkala (tanpa login,
// dipercaya karena hanya bisa diakses dari jaringan lokal yang sama).
func (s *Server) handlePollUnlock(w http.ResponseWriter, r *http.Request) {
	s.unlockMu.Lock()
	pending := s.pendingUnlock
	s.pendingUnlock = false
	s.unlockMu.Unlock()

	writeJSON(w, http.StatusOK, map[string]bool{"unlock": pending})
}

// GET /healthz
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if sqlDB, err := s.db.DB(); err != nil || sqlDB.Ping() != nil {
		writeError(w, http.StatusServiceUnavailable, "database unreachable")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
