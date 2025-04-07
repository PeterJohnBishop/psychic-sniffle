package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit"
	vegeta "github.com/tsenart/vegeta/lib"
)

func main() {
	StressGetServer()
}

func StressGetServer() {
	gofakeit.Seed(time.Now().UnixNano())
	rate := vegeta.Rate{Freq: 100, Per: time.Second}
	duration := 10 * time.Second
	attacker := vegeta.NewAttacker()

	serviceURL := "http://127.0.0.1:49457/"

	targeter := func(t *vegeta.Target) error {
		t.Method = "GET"
		t.URL = serviceURL
		return nil
	}

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Randomized Test") {
		metrics.Add(res)
	}
	metrics.Close()

	fmt.Println("Requests:", metrics.Requests)
	fmt.Printf("Success Rate: %.2f%%\n", metrics.Success*100)
	fmt.Println("Avg Latency:", metrics.Latencies.Mean)
	fmt.Println("99th Percentile:", metrics.Latencies.P99)
	fmt.Println("Errors:", metrics.Errors)
}

type Payload struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func generateRandomPayload() []byte {
	p := Payload{
		Name:     gofakeit.Username(),
		Email:    gofakeit.Email(),
		Password: gofakeit.Password(true, true, true, false, false, 8),
	}
	jsonPayload, _ := json.Marshal(p)
	return jsonPayload
}

func StressCreateUser() {
	gofakeit.Seed(time.Now().UnixNano())
	rate := vegeta.Rate{Freq: 100, Per: time.Second}
	duration := 10 * time.Second
	attacker := vegeta.NewAttacker()

	targeter := func(t *vegeta.Target) error {
		t.Method = "POST"
		t.URL = "http://localhost:8080/register"
		t.Header = map[string][]string{
			"Content-Type": {"application/json"},
		}
		t.Body = generateRandomPayload()
		return nil
	}

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Randomized Test") {
		metrics.Add(res)
	}
	metrics.Close()

	fmt.Println("Requests:", metrics.Requests)
	fmt.Printf("Success Rate: %.2f%%\n", metrics.Success*100)
	fmt.Println("Avg Latency:", metrics.Latencies.Mean)
	fmt.Println("99th Percentile:", metrics.Latencies.P99)
	fmt.Println("Errors:", metrics.Errors)
}
