package main

/*
All structures here matches the manifest template.
They are kept here to provide example and contract between manifest and backend.
Please update / add / remove as needed.
*/

type Connection struct {
	HostnameOrAddress string `json:"hostnameOrAddress"`
	Username          string `json:"username"`
	Password          string `json:"password"`
	Port              int    `json:"port"`
}

type CertificateBundle struct {
	Certificate      []byte   `json:"certificate"`
	PrivateKey       []byte   `json:"privateKey"`
	CertificateChain [][]byte `json:"certificateChain"`
}

type Keystore struct {
	CertificateName string `json:"certificateName"`
	JobId           int    `json:"jobId"`
}

type Binding struct {
	SSLProfile    string `json:"sslProfile"`
	ParentProfile string `json:"parentProfile"`
	ServerName    string `json:"serverName"`
}

type CertificateInstallation struct {
	Hostname  string `json:"hostname"`
	IpAddress string `json:"ipAddress"`
	Port      int    `json:"port"`
}

type MachineIdentity struct {
	Keystore Keystore `json:"keystore"`
	Binding  Binding  `json:"binding"`
}

type CertificateBundlePayload struct {
	CertificateBundle struct {
		Certificate      string   `json:"certificate"`
		CertificateChain []string `json:"certificateChain"`
		PrivateKey       string   `json:"privateKey"`
	} `json:"certificateBundle"`
	Connection struct {
		AapUrl         string `json:"aapUrl"`
		ConnectionType string `json:"connectionType"`
		Username       string `json:"username"`
		Password       string `json:"password"`
	} `json:"connection"`
	Keystore struct {
		CertificateName string `json:"certificateName"`
		JobId           int    `json:"jobId"`
	} `json:"keystore"`
}

// AAP job launch request
type JobLaunchRequest struct {
	ExtraVars map[string]interface{} `json:"extra_vars"`
}

type ConnectionTestPayload struct {
	Connection struct {
		AapUrl         string `json:"aapUrl"`
		ConnectionType string `json:"connectionType"`
		Username       string `json:"username"`
		Password       string `json:"password"`
	} `json:"connection"`
}
