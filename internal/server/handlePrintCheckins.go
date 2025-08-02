package server

import (
	"fmt"
	"net/http"
)

func (s *Server) handlePrintCheckins(res http.ResponseWriter, req *http.Request) {

	checkins, err := s.recurseClient.CurrentCheckins()
	if err != nil {
		s.logger.Error("could not get checkins", "error", err)
		http.Error(res, "could not get checkins", http.StatusInternalServerError)
		return
	}

	checkinStrings := "Current Checkins: \n"
	for _, checkin := range checkins {
		profile, err := s.recurseClient.Profile(checkin.Person.ID)
		if err != nil {
			s.logger.Error("could not get profile", "error", err)
			continue
		}
		checkinStrings += fmt.Sprintf("- %s, %s, %s\n", checkin.Person.Name, checkin.CreatedAt.Format("15:04PM"), profile.Stints[0].Batch.ShortName)
	}

	err = s.printer.Text(checkinStrings)
	if err != nil {
		s.logger.Error("could not print checkins", "error", err)
		http.Error(res, "could not print checkins", http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
}
