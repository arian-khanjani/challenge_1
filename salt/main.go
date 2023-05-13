package main

import (
	"crypto/rand"
	"math/big"
	"net/http"
)

func main() {
	http.HandleFunc("/generate-salt", generateSalt)
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}

func generateSalt(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		//w.Write([]byte("welcome"))
		salt, err := generateRandomString()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Write([]byte(salt))
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func generateRandomString() (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, 12)
	for i := 0; i < 12; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}
