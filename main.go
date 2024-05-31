package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	// "regexp"
	"strconv"
	"githhub.com/officialkee/goWebServer/chirps"
)

func main() {
	
	// choose port to serve on
	const port = "8080"
	cfg := ApiConfig{
		fileServerHits: 0,
	}
	
	
	// new multiplexer to serve static files
	router := http.NewServeMux()
	// handle all requests to the root directory to serve, has strip the /app prefix to avoid conflicts.... more commonly used for complex cases
	// router.Handle("/app/*", http.StripPrefix("/app",cfg.middleWareMetricsInc(http.FileServer(http.Dir(".")))))
	router.Handle("/app/*", cfg.middleWareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	// readiness endpoint to confirm the server is up... more commonly used for simple cases
	router.HandleFunc("GET /api/healthz", handlerReadiness)
	// metrics endpoint to show the number of hits to the server... more commonly used for complex cases
	router.HandleFunc("GET /admin/metrics", cfg.numberOfRequests)
	// endoint to reset the number of hits to the server
	router.HandleFunc("/api/reset", cfg.reset)
	// serve on port 8080 using the multiplexer
	router.HandleFunc("POST /api/validate_chirp",chirps.unMarshalChirp)
	router.HandleFunc("POST /api/chirps",cfg.postChirpHandler)
	router.HandleFunc("GET /api/chirps",cfg.getChirpsHandler)
	router.HandleFunc("GET /api/chirps/{id}",cfg.getChirpHandler)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}
	// log the port we are serving on
	log.Printf("Serving on port: %s\n", port)
	// // log the error if we can't serve
	log.Fatal(server.ListenAndServe())

}
// unmarshalChirp is a function to unMarshal the json  chirp from the client and validate
// func unMarshalChirp(w http.ResponseWriter, r *http.Request) {
// 	// create a struct to hold the json data from the client
// 	type  parameters struct {
// 		Body string `json:"body"`	
// 	}
// 	// decode the json data from the client into the struct
// 	body := json.NewDecoder(r.Body)
// 	params := parameters{}
// 	err := body.Decode(&params)
// 	if err != nil {
// 		w.Header().Set("Content-Type", "application/json")
//         w.Write([]byte(`{"error": "Something went wrong"}`))
//         return
// 	}
// 	if len(params.Body) > 140   {	
// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusBadRequest)
//         w.Write([]byte(`{"error": "Chirp is too long"}`))
//         return
// 	}
	
// 	re,_ := regexp.Compile(`(?i)\b(kerfuffle|sharbert|fornax)\b`)
// 	fixed := re.ReplaceAllString(params.Body, "****")
// 	type returnVal struct {
// 		UniqueId int `json:"unique_id"`
// 		Cleaned_body string `json:"cleaned_body"`
// 	}
	
// 	respBody := returnVal{
// 		UniqueId: 1,
// 		Cleaned_body: fixed,
// 	}

// 	data,err := json.Marshal(&respBody)
// 	if err!= nil {
//         w.Header().Set("Content-Type", "application/json")
//         w.Write([]byte(`{"error": "Something went wrong"}`))
//         return
//     }

// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write([]byte(data))
	
// }

// function to write the readiness response to the client // or return whatevr is decided from the function
func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
func (cfg *ApiConfig) numberOfRequests(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.fileServerHits)))
}

// middle ware method to increment the number of hits to the server
func (cfg *ApiConfig) middleWareMetricsInc(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits++
		next.ServeHTTP(w, r)

	})
}

// middle ware method to reset the number of hits to the server
func (cfg *ApiConfig) reset(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
	cfg.fileServerHits = 0

}

// stateful struct to hold the number of hits to the server
type ApiConfig struct {
	fileServerHits int
	dataBase  DataBase                                            
}


type chirp struct {
	Id int `json:"id"`
	Body string `json:"body"`
}

type DataBase struct {
	Chirps  map[int]chirp `json:"chirps"`
}

var Db = DataBase{
    Chirps: make(map[int]chirp),
}


 

func (cfg *ApiConfig) postChirpHandler(w http.ResponseWriter, r *http.Request) {
	// create a struct to hold the json data from the client
    // type  parameters struct {
    //     Body string `json:"body"`    
    // }
    // decode the json data from the client into the struct
    body := json.NewDecoder(r.Body)
    params := chirp{}
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
	if len(cfg.dataBase.Chirps) > 0 {
		params.Id = len(cfg.dataBase.Chirps) + 1
		cfg.dataBase.Chirps[params.Id] = params
		fmt.Println(cfg.dataBase)
	}
	if len(cfg.dataBase.Chirps)== 0 {
		params.Id = 1
		cfg.dataBase = Db
		cfg.dataBase.Chirps[1] = params
		
		fmt.Println(cfg.dataBase)
	}
	
    respBody := params

	data,err := json.Marshal(&respBody)

	if err!= nil {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`{"error": "Something went wrong"}`))
        return
    }

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(data))

}

func (cfg *ApiConfig)getChirpsHandler(w http.ResponseWriter, r *http.Request) {

	respBody := cfg.dataBase
	data,err := json.Marshal(&respBody)

	if err!= nil {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`{"error": "Something went wrong"}`))
        return
    }

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(data))
}

func (cfg *ApiConfig)getChirpHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	nint,err := strconv.Atoi(id)
	if err!= nil {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`{"error": "Something went wrong"}`))
        return
    }
	// fmt.Printf(cfg.dataBase.Chirps[nint].Body)
	respBody := cfg.dataBase.Chirps[nint]

	data,err := json.Marshal(&respBody)
	if err!= nil {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`{"error": "Something went wrong"}`))
        return
    }

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(data))
}