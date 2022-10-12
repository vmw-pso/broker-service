package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
}

type AuthPayload struct {
	Username string `json:"username"`
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
			s.handleSignin(requestPayload.Auth)
		case "log":
		default:
			s.tools.ErrorJSON(w, errors.New("unknown action"))
		}
	}
}

func (s *server) handleSignin(payload AuthPayload) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			s.tools.ErrorJSON(w, err)
			return
		}

		request, err := http.NewRequest("POST", "http://authentication-service/signin", bytes.NewBuffer(jsonData))
		if err != nil {
			s.tools.ErrorJSON(w, err)
			return
		}

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			s.tools.ErrorJSON(w, err)
			return
		}
		defer response.Body.Close()

		if response.StatusCode == http.StatusUnauthorized {
			s.tools.ErrorJSON(w, errors.New("invalid username or password"))
			return
		} else if response.StatusCode != http.StatusAccepted {
			s.tools.ErrorJSON(w, errors.New("error calling auth service"))
			return
		}

		// return token
	}
}
