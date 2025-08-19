package server

import (
	"encoding/json"
	"net/http"
)

func (s *Server) handleGetCheckins(res http.ResponseWriter, req *http.Request) {

	checkins, err := s.recurseClient.CurrentCheckins()
	if err != nil {
		s.logger.Error("could not get checkins", "error", err)
		http.Error(res, "could not get checkins", http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(checkins)
	if err != nil {
		s.logger.Error("json marshal failed", "error", err)
		http.Error(res, "json marshal failed", http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(b)
}
