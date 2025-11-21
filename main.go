package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/jpantarotto/url-shortener/config"
	"github.com/jpantarotto/url-shortener/db"
)

type InputUrl struct {
	URL string `json:"url"`
}

type TinyUrl struct {
	Original string `json:"orignal"`
	Tiny     string `json:"tiny"`
}

var urls = make(map[string]string)
var counter int64

const Base62Alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func encodeToBase62(num int64) string {
	result := make([]byte, 0)          // initialize result byte slice
	base := int64(len(Base62Alphabet)) // create int64 with value of len(alphabet) == 62

	// while num (dividend) is greater than zero
	for num > 0 {
		num = num / base
		remainder := num % base
		result = append(result, Base62Alphabet[remainder])
	}
	// reverse the result []byte with no extra space
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	if len(result) == 0 {
		return string(Base62Alphabet[0])
	}

	return string(result)
}

// func get(w http.ResponseWriter, req *http.Request) {
// 	tinyUrl := req.PathValue("tinyUrl")

// }

func create(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading request body: %v", err)
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	// Print the received body (for demonstration)
	fmt.Printf("Received POST request with body: %s\n", body)

	// Parse JSON input
	jsonData := []byte(body)
	var input InputUrl
	err = json.Unmarshal(jsonData, &input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error unmarshalling JSON: %v", err)
		w.WriteHeader(http.StatusBadRequest)
	}

	fmt.Printf("Counter: %v", counter)
	// Create short URL hash
	// hash := createSha256Hash(input.URL)
	base62String := encodeToBase62(counter)
	counter++

	urls[base62String] = input.URL

	responseData := TinyUrl{
		Original: input.URL,
		Tiny:     base62String,
	}

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(responseData); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %v", err)
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusMovedPermanently)

	// Log the response
	fmt.Printf("Original: %s, Tiny: %s", input.URL, base62String)
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to load config: %v\n", err)
		os.Exit(1)
	}

	conn, err := db.Connect(cfg.DB)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	counter = cfg.CounterStart
	port := ":" + cfg.Port
	http.HandleFunc("/create", create)
	// http.HandleFunc("GET /{tinyUrl}", get)

	fmt.Printf("Server Running at http://lochalhost%s\n", port)
	http.ListenAndServe(port, nil)
}
