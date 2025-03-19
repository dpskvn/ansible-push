package main

import (
	"fmt"
	"net/http"

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
	req := TestConnectionRequest{}
	if err := c.Bind(&req); err != nil {
		zap.L().Error("invalid request, failed to unmarshall json", zap.Error(err))
		return c.String(http.StatusBadRequest, fmt.Sprintf("failed to unmarshall json: %s", err.Error()))
	}

	/*
		TODO: Write the logic on how do we connect to the server. A separate connect method can be written and called.
			Also, a separate connect method can be used for other functions as well.
	*/

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
	req := InstallCertificateBundleRequest{}
	if err := c.Bind(&req); err != nil {
		zap.L().Error("invalid request, failed to unmarshall json", zap.Error(err))
		return c.String(http.StatusBadRequest, fmt.Sprintf("failed to unmarshall json: %s", err.Error()))
	}

	/*
		TODO: Write your business logic to install the certificate bundle in this area.
	*/

	res := InstallCertificateBundleResponse{
		InstallationKeystore: req.InstallationKeystore,
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

func handleDiscoverCertificates(c echo.Context) error {
	req := DiscoverCertificatesRequest{}
	if err := c.Bind(&req); err != nil {
		zap.L().Error("invalid request, failed to unmarshall json", zap.Error(err))
		return c.String(http.StatusBadRequest, fmt.Sprintf("failed to unmarshall json: %s", err.Error()))
	}

	res := DiscoverCertificatesResponse{
		Messages: []*DiscoveredCertificate{},
	}

	return c.JSON(http.StatusOK, &res)
}
