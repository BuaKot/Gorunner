package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

var urlStore = make(map[string]string)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func main() {
	rand.Seed(time.Now().UnixNano())

	// Route สำหรับหน้าเว็บ
	http.HandleFunc("/", handleIndex)
	// Route สำหรับการย่อ URL (รับค่าจาก Form)
	http.HandleFunc("/shorten", handleShorten)

	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	// ถ้า path ไม่ใช่ / (เช่น /abc) ให้ไปที่ Redirect logic
	if r.URL.Path != "/" {
		handleRedirect(w, r)
		return
	}
	// ส่งไฟล์ index.html (เราจะสร้างไฟล์นี้ในขั้นตอนถัดไป)
	http.ServeFile(w, r, "index.html")
}

func handleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	longURL := r.FormValue("url")
	if longURL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	code := generateCode(6)
	urlStore[code] = longURL

	shortURL := fmt.Sprintf("http://localhost:8080/%s", code)

	// ส่งผลลัพธ์กลับเป็นข้อความแบบง่ายๆ หรือจะทำหน้า Result ก็ได้
	fmt.Fprintf(w, "Success! Your short URL is: %s", shortURL)
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Path[1:]
	longURL, exists := urlStore[code]
	if !exists {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, longURL, http.StatusFound)
}

func generateCode(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
