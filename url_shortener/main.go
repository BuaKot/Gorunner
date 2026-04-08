package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ctx         = context.Background()
	redisClient *redis.Client
	charset     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func main() {
	// 1. เชื่อมต่อ Redis
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// ตรวจสอบการเชื่อมต่อ
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("เชื่อมต่อ Redis ไม่ได้: %v", err))
	}

	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/shorten", handleShorten)

	fmt.Println("Server running with Redis at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		handleRedirect(w, r)
		return
	}
	http.ServeFile(w, r, "index.html")
}

func handleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	longURL := r.FormValue("url")
	code := generateCode(6)

	// 2. บันทึกลง Redis (ตั้งให้หมดอายุใน 24 ชั่วโมง)
	err := redisClient.Set(ctx, code, longURL, 24*time.Hour).Err()
	if err != nil {
		http.Error(w, "Error saving to Redis", http.StatusInternalServerError)
		return
	}

	shortURL := fmt.Sprintf("http://localhost:8080/%s", code)
	fmt.Fprintf(w, "Success! Your short URL is: %s", shortURL)
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Path[1:]

	// 3. ดึงข้อมูลจาก Redis
	longURL, err := redisClient.Get(ctx, code).Result()
	if err == redis.Nil {
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
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
