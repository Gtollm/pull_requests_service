package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"os"
	"pull-request-review/config"
	"pull-request-review/internal/app"
	"pull-request-review/internal/infrastructure/adapters/logger"
	"pull-request-review/internal/infrastructure/database"
	"testing"
)

var (
	baseURL string
	db      *database.Database
)

func TestMain(m *testing.M) {
	log := logger.NewZerologLogger()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error(err, "Failed to load configuration")
		os.Exit(1)
	}

	cfg.Server.Port = "8081"

	ctx := context.Background()
	db = database.NewDatabase(cfg.Database, log)
	if err := db.Connect(ctx); err != nil {
		log.Error(err, "Failed to connect to database")
		os.Exit(1)
	}

	if err := cleanupTestData(ctx); err != nil {
		log.Error(err, "Failed to cleanup test data")
	}

	go app.Run(cfg, db, log)

	baseURL = fmt.Sprintf("http://localhost:%s", cfg.Server.Port)

	code := m.Run()

	db.Close()
	os.Exit(code)
}

func cleanupTestData(ctx context.Context) error {
	_, err := db.GetPool().Exec(
		ctx, `TRUNCATE TABLE review_assignments, pull_requests, users, teams RESTART IDENTITY CASCADE;`,
	)
	if err != nil {
		return err
	}
	return nil
}

func TestHealthCheck(t *testing.T) {
	resp, err := http.Get(baseURL + "/health")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if result["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got '%v'", result["status"])
	}
}

func TestCreateTeam(t *testing.T) {
	teamData := map[string]interface{}{
		"team_name": "test-team-e2e",
		"members": []map[string]string{
			{
				"user_id":  uuid.New().String(),
				"username": "user1",
			},
			{
				"user_id":  uuid.New().String(),
				"username": "user2",
			},
		},
	}

	jsonData, err := json.Marshal(teamData)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	resp, err := http.Post(
		baseURL+"/team/add",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		var errorResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errorResp)
		t.Fatalf("Expected status 201, got %d. Response: %v", resp.StatusCode, errorResp)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	team, ok := result["team"].(map[string]interface{})
	if !ok {
		t.Fatal("Response doesn't contain team object")
	}

	if team["team_name"] != "test-team-e2e" {
		t.Errorf("Expected team_name 'test-team-e2e', got '%v'", team["team_name"])
	}
}

func TestGetTeam(t *testing.T) {
	resp, err := http.Get(baseURL + "/team/get?team_name=test-team-e2e")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	team, ok := result["team"].(map[string]interface{})
	if !ok {
		t.Fatal("Response doesn't contain team object")
	}

	if team["team_name"] != "test-team-e2e" {
		t.Errorf("Expected team_name 'test-team-e2e', got '%v'", team["team_name"])
	}

	members, ok := team["members"].([]interface{})
	if !ok || len(members) != 2 {
		t.Errorf("Expected 2 members, got %v", members)
	}
}