// e2e/e2e_suite_test.go
package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stannisl/pull-request-service/internal/app"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type E2ETestSuite struct {
	suite.Suite
	app               *app.App
	postgresContainer testcontainers.Container
	baseURL           string
	ctx               context.Context
	httpClient        *http.Client
	dbConnStr         string
}

func (suite *E2ETestSuite) SetupSuite() {
	suite.ctx = context.Background()
	suite.httpClient = &http.Client{
		Timeout: 30 * time.Second,
	}
	suite.baseURL = "http://localhost:8080"

	// Запускаем PostgreSQL контейнер
	err := suite.startPostgreSQL()
	suite.Require().NoError(err)

	// Ждем пока база данных станет полностью доступна
	suite.waitForDBReady()

	// Запускаем приложение
	err = suite.startApplication()
	suite.Require().NoError(err)

	// Ждем пока приложение станет доступно
	suite.waitForAppReady()
}

func (suite *E2ETestSuite) startPostgreSQL() error {
	var err error

	// Запускаем контейнер PostgreSQL
	suite.postgresContainer, err = postgres.Run(suite.ctx,
		"postgres:15-alpine",
		postgres.WithDatabase("test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second)), // Увеличим таймаут
	)
	if err != nil {
		return fmt.Errorf("failed to start postgres container: %w", err)
	}

	// Получаем реальный порт
	endpoint, err := suite.postgresContainer.Endpoint(suite.ctx, "")
	if err != nil {
		return fmt.Errorf("failed to get postgres endpoint: %w", err)
	}

	// Формируем строку подключения
	suite.dbConnStr = fmt.Sprintf("postgres://test:test@%s/test?sslmode=disable", endpoint)

	log.Printf("PostgreSQL container started at: %s", endpoint)
	log.Printf("Database connection string: %s", suite.dbConnStr)

	return nil
}

func (suite *E2ETestSuite) waitForDBReady() {
	os.Setenv("APP_DATABASE_CONN_URL", suite.dbConnStr)
	os.Setenv("APP_DATABASE_RETRIES", "10")
	os.Setenv("APP_DATABASE_RETRIES_SECONDS_DELAY", "3")
	os.Setenv("APP_DATABASE_DRIVER_NAME", "postgres")

	log.Println("Environment variables set for database connection")

	time.Sleep(3 * time.Second)
}

func (suite *E2ETestSuite) startApplication() error {
	suite.app = &app.App{}

	ctx, cancel := context.WithTimeout(suite.ctx, 60*time.Second)
	defer cancel()

	err := suite.app.Setup(ctx)
	if err != nil {
		return fmt.Errorf("failed to setup app: %w", err)
	}

	// Запускаем приложение в горутине
	go func() {
		log.Println("Starting HTTP server...")
		if err := suite.app.StartAndServeHTTP(suite.ctx); err != nil {
			log.Printf("Application stopped with error: %v", err)
		}
	}()

	return nil
}

func (suite *E2ETestSuite) waitForAppReady() {
	log.Println("Waiting for application to become ready...")

	suite.Require().Eventually(func() bool {
		resp, err := suite.httpClient.Get(suite.baseURL + "/health")
		if err != nil {
			log.Printf("Health check failed: %v", err)
			return false
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			log.Println("Application is ready!")
			return true
		}

		return false
	}, 30*time.Second, 1*time.Second, "Application failed to become ready")
}

func (suite *E2ETestSuite) TearDownSuite() {
	log.Println("Tearing down test suite...")

	if suite.app != nil {
	}

	if suite.postgresContainer != nil {
		if err := suite.postgresContainer.Terminate(suite.ctx); err != nil {
			log.Printf("Failed to terminate postgres container: %v", err)
		} else {
			log.Println("PostgreSQL container terminated")
		}
	}
}

// Вспомогательные методы для работы с HTTP
func (suite *E2ETestSuite) makeRequest(method, path string, body interface{}) (*http.Response, error) {
	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, suite.baseURL+path, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	return suite.httpClient.Do(req)
}

func (suite *E2ETestSuite) parseResponse(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(target)
}

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, new(E2ETestSuite))
}
