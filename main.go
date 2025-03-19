package main

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/square/go-jose.v2"
)

func rfc3339TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(time.RFC3339))
}

func main() {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	loggerConfig.EncoderConfig.TimeKey = "time"
	loggerConfig.EncoderConfig.EncodeTime = rfc3339TimeEncoder
	l, err := loggerConfig.Build()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(l)

	//Replace log statement with connector getting built.
	l.Info("Sample Connector Starting")

	e := echo.New()

	e.GET("/healthz", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	g := e.Group("/v1")
	addPayloadEncryptionMiddleware(g)
	/*
		The connector code needs to provide handlers (business logic) for the 3 APIs listed below.
		The request body of each API can be defined by the developer.
		The response of the APIs can be customized to some extent, but base booleans for success or failures are mandatory.
	*/

	/*
		The testconnection API takes the connection string needed and responds back with a success or failure connection to the system we are building a connection for.
	*/
	g.POST("/testconnection", handleTestConnection)

	/*
		The installcertificatebundle API takes the certificate bundle, connections details, and target information (key store) on where the certificate bundle needs to be installed.
		A handler must be defined to provide business logic on how to install a certificate bundle to the external system.
	*/
	g.POST("/installcertificatebundle", handleInstallCertificateBundle)

	/*
		 			The configureinstallationendpoint API takes the certificate which was installed previously
				and configures it to a specific endpoint provided in previous request.
	*/
	g.POST("/configureinstallationendpoint", handleConfigureInstallationEndpoint)

	/*
		The discovercertificates API takes the remote host connection information and discovery preferences to
		retrieve all relevant certificates from the host.
	*/
	g.POST("/discovercertificates", handleDiscoverCertificates)

	//Do not change the port for the service. It has to run on port 8080.
	if err := e.Start(":8080"); err != nil {
		zap.L().Error("Error starting Echo", zap.Error(err))
	}
}

func addPayloadEncryptionMiddleware(g *echo.Group) {
	privateKeyPemData, err := os.ReadFile("/keys/payload-encryption-key.pem")
	if err != nil {
		zap.L().Error("payload encryption key not found or readable", zap.Error(err))
		return
	}
	p, _ := pem.Decode(privateKeyPemData)
	if p == nil {
		zap.L().Error("payload encryption key not in PEM format")
		return
	}
	pk, err := x509.ParsePKCS1PrivateKey(p.Bytes)
	if err != nil {
		zap.L().Error("payload encryption key not properly encoded", zap.Error(err))
		return
	}
	zap.L().Info("adding payload encryption middleware")
	g.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				return err
			}
			object, err := jose.ParseEncrypted(string(body))
			if err != nil {
				return err
			}
			decrypted, err := object.Decrypt(pk)
			if err != nil {
				return err
			}
			req.Body = io.NopCloser(bytes.NewReader(decrypted))
			return next(c)
		}
	})
}
