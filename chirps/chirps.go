package chirps

import (
    "encoding/json"
    "regexp"
    "net/http"
)

func unMarshalChirp(w http.ResponseWriter, r *http.Request) {
	// create a struct to hold the json data from the client
	type  parameters struct {
		Body string `json:"body"`	
	}
	// decode the json data from the client into the struct
	body := json.NewDecoder(r.Body)
	params := parameters{}
	err := body.Decode(&params)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`{"error": "Something went wrong"}`))
        return
	}
	if len(params.Body) > 140   {	
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(`{"error": "Chirp is too long"}`))
        return
	}
	
	re,_ := regexp.Compile(`(?i)\b(kerfuffle|sharbert|fornax)\b`)
	fixed := re.ReplaceAllString(params.Body, "****")
	type returnVal struct {
		UniqueId int `json:"unique_id"`
		Cleaned_body string `json:"cleaned_body"`
	}
	
	respBody := returnVal{
		UniqueId: 1,
		Cleaned_body: fixed,
	}

	data,err := json.Marshal(&respBody)
	if err!= nil {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`{"error": "Something went wrong"}`))
        return
    }

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(data))
	
}