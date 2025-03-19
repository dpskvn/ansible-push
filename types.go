package main

/*
All structures here matches the manifest template.
They are kept here to provide example and contract between manifest and backend.
Please update / add / remove as needed.
*/

type DiscoveryType string

const (
	AllDiscoveryTypes                     DiscoveryType = "all"
	MonitorDiscoveryType                  DiscoveryType = "monitor"
	VirtualServerDiscoveryType            DiscoveryType = "virtualServer"
	inactiveClientSslProfileDiscoveryType DiscoveryType = "clientSsl"
	inactiveServerSslProfileDiscoveryType DiscoveryType = "serverSsl"
)

type PartitionNames []string

type Connection struct {
	HostnameOrAddress string `json:"hostnameOrAddress"`
	Username          string `json:"username"`
	Password          string `json:"password"`
	Port              int    `json:"port"`
}

type sslProfileType string

const (
	ClientSSLProfile = sslProfileType("client")
	ServerSSLProfile = sslProfileType("server")
)

type SSLProfile struct {
	Name      string         `json:"name"`
	Type      sslProfileType `json:"type"`
	Partition string         `json:"partition"`
}

type CertificateBundle struct {
	Certificate      []byte   `json:"certificate"`
	PrivateKey       []byte   `json:"privateKey"`
	CertificateChain [][]byte `json:"certificateChain"`
}

type Keystore struct {
	Partition       string `json:"partition"`
	CertificateName string `json:"certificateName"`
	ChainName       string `json:"chainName"`
}

type Binding struct {
	SSLProfile     string         `json:"sslProfile"`
	ParentProfile  string         `json:"parentProfile"`
	SSLProfileType sslProfileType `json:"sslProfileType"`
	ServerName     string         `json:"serverName"`
}

type DiscoverCertificatesRequest struct {
	Configuration DiscoverCertificatesConfiguration `json:"discovery"`
	Connection    Connection                        `json:"connection"`
	Control       DiscoveryControl                  `json:"discoveryControl"`
	Page          DiscoveryPage                     `json:"discoveryPage"`
}

type DiscoveryControl struct {
	MaxResults int `json:"maxResults"`
}

type DiscoverCertificatesConfiguration struct {
	ExcludeExpiredCertificates bool          `json:"excludeExpiredCertificates"`
	ExcludeInactiveProfiles    bool          `json:"excludeInactiveProfiles"`
	Partition                  string        `json:"partition"`
	TimeStamp                  string        `json:"timeStamp"`
	Type                       DiscoveryType `json:"discoveryType"`

	partitions PartitionNames
}

type DiscoveryPage struct {
	Type      DiscoveryType `json:"discoveryType"`
	Paginator string        `json:"paginator"`
}

type DiscoverCertificatesResponse struct {
	Page     *DiscoveryPage           `json:"discoveryPage"`
	Messages []*DiscoveredCertificate `json:"messages"`
}

type DiscoveredCertificate struct {
	Certificate       string                     `json:"certificate"`
	CertificateChain  []string                   `json:"certificateChain"`
	Installations     []*CertificateInstallation `json:"installations"`
	MachineIdentities []*MachineIdentity         `json:"machineIdentities"`
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
