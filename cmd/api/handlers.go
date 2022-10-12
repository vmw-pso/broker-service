package main

import (
	"errors"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (s *server) handleRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestPayload RequestPayload

		err := s.tools.ReadJSON(w, r, &requestPayload)
		if err != nil {
			s.tools.ErrorJSON(w, err)
			return
		}

		switch requestPayload.Action {
		case "signin":
		case "log":
		default:
			s.tools.ErrorJSON(w, errors.New("unknown action"))
		}
	}
}
