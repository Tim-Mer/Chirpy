package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
)

func handlerHealthz(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("OK"))
}

func validateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	} else {
		if len(params.Body) > 140 {
			respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		} else {
			type returnVals struct {
				Valid        bool   `json:"valid"`
				Body         string `json:"body"`
				Cleaned_body string `json:"cleaned_body"`
			}
			retVal := returnVals{
				Valid:        true,
				Body:         params.Body,
				Cleaned_body: removeProfanity(params.Body),
			}
			respondWithJSON(w, http.StatusOK, retVal)
		}
	}
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type respError struct {
		Error string `json:"error"`
	}
	respBody := respError{
		Error: msg,
	}
	respondWithJSON(w, code, respBody)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func removeProfanity(s string) string {
	profanity := []string{"kerfuffle", "sharbert", "fornax"}
	spl := strings.Split(s, " ")
	for i, t := range spl {
		if slices.Contains(profanity, strings.ToLower(t)) {
			spl[i] = "****"
		}
	}
	return strings.Join(spl, " ")
}
