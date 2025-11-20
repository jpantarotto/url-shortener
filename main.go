package main

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
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

const Base62Alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func encodeToBase62(data []byte) string {
	num := new(big.Int).SetBytes(data)             // convert the byte slice to a Big Int
	result := make([]byte, 0)                      // initialize result byte slice
	base := big.NewInt(int64(len(Base62Alphabet))) // create big int with value of len(alphabet) == 62

	// while num (dividend) is greater than zero
	for num.Cmp(big.NewInt(0)) > 0 {
		remainder := new(big.Int)
		num.QuoRem(num, base, remainder)
		result = append(result, Base62Alphabet[remainder.Int64()])
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

func createSha256Hash(input string) []byte {
	hasher := md5.New()
	hasher.Write([]byte(input))
	hashBytes := hasher.Sum(nil)

	return hashBytes
}

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

	// Create short URL hash
	hash := createSha256Hash(input.URL)
	base62String := encodeToBase62(hash)

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

	port := ":4000"
	http.HandleFunc("/create", create)

	fmt.Printf("Server Running at http://lochalhost%s\n", port)
	http.ListenAndServe(port, nil)
}
