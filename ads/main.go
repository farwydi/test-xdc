package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/touch", func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		w.Header().Add("Access-Control-Allow-Origin",
			origin)
		w.Header().Add("Access-Control-Allow-Credentials", "true")

		touchCookie, err := r.Cookie("touch")
		if errors.Is(err, http.ErrNoCookie) {
			expire := time.Now().AddDate(0, 0, 1)
			buf := make([]byte, 16)
			_, err := rand.Read(buf)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			touchCookie = &http.Cookie{
				Name:       "touch",
				Value:      hex.EncodeToString(buf),
				Path:       "/",
				Expires:    expire,
				RawExpires: expire.Format(time.UnixDate),
				MaxAge:     86400,
				Secure:     true,
				SameSite:   http.SameSiteNoneMode,
			}
			http.SetCookie(w, touchCookie)
		} else if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		err = json.NewEncoder(w).Encode(touchCookie.Value)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.WriteHeader(200)
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
