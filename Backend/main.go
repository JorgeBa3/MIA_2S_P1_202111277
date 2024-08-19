package main

import (
	"encoding/json"
	"main/analyzer"
	"net/http"
)

// Middleware para manejar CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Función para manejar la ruta /greet
func greetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var requestBody map[string]string
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Solicitud no válida", http.StatusBadRequest)
		return
	}

	// Analiza el comando o comentario proporcionado
	result, err := analyzer.Analyzer(requestBody["name"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Responde con el resultado del análisis
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"messages": result})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/greet", greetHandler)
	http.ListenAndServe(":8080", corsMiddleware(mux))
}
