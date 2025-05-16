package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

const (
	maximumFileNameLength = 63
)

type TestConnectionRequest struct {
	Connection Connection `json:"connection"`
}

type TestConnectionResponse struct {
	Result bool `json:"result"`
}

func handleTestConnection(c echo.Context) error {
	req := ConnectionTestPayload{}
	if err := c.Bind(&req); err != nil {
		zap.L().Error("invalid request, failed to unmarshall json", zap.Error(err))
		return c.String(http.StatusBadRequest, fmt.Sprintf("failed to unmarshall json: %s", err.Error()))
	}

	// Test the connection to the AAP server
	status := testAAPConnection(req)
	if status == false {
		return c.String(http.StatusForbidden, fmt.Sprintf("failed to connect to AAP"))
	}

	res := TestConnectionResponse{
		Result: true,
	}
	return c.JSON(http.StatusOK, res)
}

type InstallCertificateBundleRequest struct {
	Connection           Connection        `json:"connection"`
	CertificateBundle    CertificateBundle `json:"certificateBundle"`
	InstallationKeystore Keystore          `json:"keystore"`
}

type InstallCertificateBundleResponse struct {
	InstallationKeystore Keystore `json:"keystore"`
}

func handleInstallCertificateBundle(c echo.Context) error {
	req := CertificateBundlePayload{}
	if err := c.Bind(&req); err != nil {
		zap.L().Error("invalid request, failed to unmarshall json", zap.Error(err))
		return c.String(http.StatusBadRequest, fmt.Sprintf("failed to unmarshall json: %s", err.Error()))
	}

	// Launch the AAP job
	_, err := launchAAPJob(req)
	if err != nil {
		zap.L().Error("failed to launch AAP job", zap.Error(err))
		return c.String(http.StatusInternalServerError, fmt.Sprintf("failed to launch AAP job: %s", err.Error()))
	}

	res := InstallCertificateBundleResponse{
		InstallationKeystore: req.Keystore,
	}
	return c.JSON(http.StatusOK, &res)
}

type ConfigureInstallationEndpointRequest struct {
	Connection Connection `json:"connection"`
	Keystore   Keystore   `json:"keystore"`
	Binding    Binding    `json:"binding"`
}

func handleConfigureInstallationEndpoint(c echo.Context) error {
	req := ConfigureInstallationEndpointRequest{}
	if err := c.Bind(&req); err != nil {
		zap.L().Error("invalid request, failed to unmarshall json", zap.Error(err))
		return c.String(http.StatusBadRequest, fmt.Sprintf("failed to unmarshall json: %s", err.Error()))
	}
	/*
		TODO: Write the business logic to configure the endpoint with the certificate previously installed.
	*/

	return c.NoContent(http.StatusOK)
}

func launchAAPJob(payload CertificateBundlePayload) (int, error) {
	// Join the certificate chain array into a single string with pipe separators
	// AAP endpoint only supports string type for extra vars, not array
	certificateChainString := strings.Join(payload.CertificateBundle.CertificateChain, "|")

	// Create a map for the extra vars
	// These need to be enabled as survey variables on AAP
	extraVars := map[string]any{
		"certificate_name":  payload.Keystore.CertificateName,
		"certificate":       payload.CertificateBundle.Certificate,
		"certificate_chain": certificateChainString, // Use the string version instead of array
		"private_key":       payload.CertificateBundle.PrivateKey,
	}

	// Create the job launch request
	jobLaunchRequest := JobLaunchRequest{
		ExtraVars: extraVars,
	}

	// Convert the job launch request to JSON
	jobLaunchRequestJSON, err := json.Marshal(jobLaunchRequest)
	if err != nil {
		zap.L().Error("failed to marshal job launch request", zap.Error(err))
		return 0, fmt.Errorf("failed to marshal job launch request: %s", err.Error())
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	// Use the job template ID from the keystore
	jobTemplateID := payload.Keystore.JobId

	// Create the URL for launching the job template
	url := fmt.Sprintf("%s/api/v2/job_templates/%d/launch/", payload.Connection.AapUrl, jobTemplateID)

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jobLaunchRequestJSON))
	if err != nil {
		zap.L().Error("failed to create HTTP request", zap.Error(err))
		return 0, fmt.Errorf("failed to create HTTP request: %s", err.Error())
	}

	// Set the request headers
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(payload.Connection.Username, payload.Connection.Password)

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		zap.L().Error("failed to send HTTP request", zap.Error(err))
		return 0, fmt.Errorf("failed to send HTTP request: %s", err.Error())
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		responseBody, _ := io.ReadAll(resp.Body)
		zap.L().Error("received non-success status code", zap.Int("status_code", resp.StatusCode), zap.String("response_body", string(responseBody)))
		return 0, fmt.Errorf("received non-success status code %d: %s", resp.StatusCode, string(responseBody))
	}

	// Parse the response to get the job ID
	var jobResponse struct {
		ID int `json:"id"`
	}

	err = json.NewDecoder(resp.Body).Decode(&jobResponse)
	if err != nil {
		zap.L().Error("failed to parse response", zap.Error(err))
		return 0, fmt.Errorf("failed to parse response: %s", err.Error())
	}

	log.Printf("Successfully launched AAP job with ID: %d\n", jobResponse.ID)
	return jobResponse.ID, nil
}

func testAAPConnection(connection ConnectionTestPayload) bool {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: 10 * time.Second, // Timeout for the request
	}

	// Create the URL for the API endpoint
	// Use the /api/v2/credentials endpoint which should be available for authenticated users
	url := fmt.Sprintf("%s/api/v2/credentials", connection.Connection.AapUrl)

	// Create the HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		zap.L().Error("failed to create HTTP request", zap.Error(err))
		return false
	}

	// Set the request headers
	req.SetBasicAuth(connection.Connection.Username, connection.Connection.Password)

	resp, err := client.Do(req)
	if err != nil {
		zap.L().Error("failed to connect to AAP", zap.Error(err))
		return false
	}
	defer resp.Body.Close()

	// Check for valid status code
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return true
	} else {
		return false
	}
}
