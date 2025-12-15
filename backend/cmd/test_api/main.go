package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const baseURL = "http://localhost:3000/api/v1"

func main() {
	fmt.Println("Starting API Test...")

	// 1. Test Register
	fmt.Println("\n[1] Testing Register...")
	uniqueName := fmt.Sprintf("testuser_%d", time.Now().Unix())
	registerBody := map[string]string{
		"username": uniqueName,
		"email":    uniqueName + "@example.com",
		"password": "password123",
	}
	code, resp := sendRequest("POST", "/users/register", registerBody, "")
	if code == 201 {
		fmt.Println("✅ Register Success")
	} else {
		fmt.Printf("❌ Register Failed: %d - %s\n", code, resp)
		return
	}

	// 2. Test Login
	fmt.Println("\n[2] Testing Login...")
	loginBody := map[string]string{
		"username": uniqueName,
		"password": "password123",
	}
	code, resp = sendRequest("POST", "/users/login", loginBody, "")
	var token string
	if code == 200 {
		var result map[string]interface{}
		json.Unmarshal([]byte(resp), &result)
		if t, ok := result["token"].(string); ok {
			token = t
			fmt.Println("✅ Login Success, Token received")
		} else {
			fmt.Println("❌ Login Failed: Token not found in response")
			return
		}
	} else {
		fmt.Printf("❌ Login Failed: %d - %s\n", code, resp)
		return
	}

	// 3. Test Protected Route (Get Me)
	fmt.Println("\n[3] Testing Get Me (Protected)...")
	code, resp = sendRequest("GET", "/users/me", nil, token)
	if code == 200 {
		fmt.Println("✅ Get Me Success")
		fmt.Println("Response:", resp)
	} else {
		fmt.Printf("❌ Get Me Failed: %d - %s\n", code, resp)
	}
}

func sendRequest(method, endpoint string, body interface{}, token string) (int, string) {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		bodyReader = bytes.NewBuffer(jsonBody)
	}

	req, _ := http.NewRequest(method, baseURL+endpoint, bodyReader)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Sprintf("Network error: %v", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, string(respBody)
}
