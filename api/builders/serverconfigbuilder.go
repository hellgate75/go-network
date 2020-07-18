package builders

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"github.com/hellgate75/go-network/model"
	"golang.org/x/crypto/acme/autocert"
	"io/ioutil"
)

// Helper for building a model.ServerConfig instance
type ServerConfigBuilder interface {
	// Associate an host and a port to the builder workflow
	WithHost(address string, port int) ServerConfigBuilder
	// Associate certificate and key files full path to the builder workflow
	WithTLSCerts(certificate string, key string) ServerConfigBuilder
	// Add some more certificate files to the certificate list to the builder workflow
	// If no certificate is settled up first call with associate the main TLS certificate files
	MoreTLSCerts(certificate string, key string) ServerConfigBuilder
	// Add one root CA certificate files to the certificate list to the builder workflow
	WithRootCaCert(certificate string) ServerConfigBuilder
	// Add one client CA certificate files to the certificate list to the builder workflow
	WithClientCaCert(certificate string) ServerConfigBuilder
	// Set up the certificate manager for the auto-scan of certificates for a folder
	WithCertificateManager(dir string) ServerConfigBuilder
	// Add more root CA certificate files to the certificate list to the builder workflow
	MoreClientCaCerts(certificate string) ServerConfigBuilder
	// Add more client CA certificate files to the certificate list to the builder workflow
	MoreRootCaCerts(certificate string) ServerConfigBuilder
	// Set min version different from tls.VersionTLS12
	WithMinVersion(min uint16) ServerConfigBuilder
	// Set the insecure skip verify flag, by default it's false
	WithInsecureSkipVerify(insecure bool) ServerConfigBuilder
	// Set up renegotiation, by default it's sett up to: tls.RenegotiateNever
	WithRenegotiationSupport(renegotiation tls.RenegotiationSupport) ServerConfigBuilder
	// Set up the Client Session Cache manager (suggested: tls.NewLRUClientSessionCache(1024) or more ...)
	WithClientSessionCache(cache tls.ClientSessionCache) ServerConfigBuilder
	// Add more Cipher suites to the preset values : tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	// tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA, tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
	// tls.TLS_RSA_WITH_AES_256_CBC_SHA
	MoreCipherSuites(cipherSuite uint16) ServerConfigBuilder
	// Add more Curve Ids to the current TLS Curve Preferences, adding to the preset values : tls.CurveP521, tls.CurveP384,
	// tls.CurveP256
	MoreCurvePreferences(curve tls.CurveID) ServerConfigBuilder
	// Set preference for Server Cipher Suite
	WithPreferServerCipherSuites(preferServerCipherSuites bool)  ServerConfigBuilder
	// Build the model.ServerConfig and report any error occurred during the build process
	Build() (model.ServerConfig, error)
}

type serverConfigBuilder struct{
	address      				string
	port         				int
	certificate  				string
	key          				string
	caPool       				*x509.CertPool
	rootCaPool   				*x509.CertPool
	certificates 				[]tls.Certificate
	cipherSuits  				[]uint16
	insecure     				bool
	curvePref    				[]tls.CurveID
	certManager  				*autocert.Manager
	minVersion 	 				uint16
	renegotiation 				tls.RenegotiationSupport
	cache						tls.ClientSessionCache
	preferServerCipherSuites 	bool
}

func (b *serverConfigBuilder) WithHost(address string, port int) ServerConfigBuilder {
	b.address = address
	b.port = port
	return b
}

func (b *serverConfigBuilder) WithTLSCerts(certificate string, key string) ServerConfigBuilder {
	b.certificate=certificate
	b.key=key
	return b
}

func (b *serverConfigBuilder) MoreTLSCerts(certificate string, key string) ServerConfigBuilder {
	cert, err := tls.LoadX509KeyPair(certificate, key)
	if err == nil {
		b.certificates = append(b.certificates, cert)
		if b.certificate == "" || b.key == "" {
			b.certificate = certificate
			b.key = key
		}
	}
	return b
}

func (b *serverConfigBuilder) WithRootCaCert(certificate string) ServerConfigBuilder {
	caCert, err := ioutil.ReadFile(certificate)
	if err != nil {
		if b.rootCaPool == nil {
			b.rootCaPool = x509.NewCertPool()
		}
		_ = b.rootCaPool.AppendCertsFromPEM(caCert)
	}
	return b
}

func (b *serverConfigBuilder) WithClientCaCert(certificate string) ServerConfigBuilder {
	caCert, err := ioutil.ReadFile(certificate)
	if err != nil {
		if b.caPool == nil {
			b.caPool = x509.NewCertPool()
		}
		_ = b.caPool.AppendCertsFromPEM(caCert)
	}
	return b
}

func (b *serverConfigBuilder) MoreClientCaCerts(certificate string) ServerConfigBuilder {
	caCert, err := ioutil.ReadFile(certificate)
	if err != nil {
		if b.caPool == nil {
			b.caPool = x509.NewCertPool()
		}
		_ = b.caPool.AppendCertsFromPEM(caCert)
	}
	return b
}

func (b *serverConfigBuilder) MoreRootCaCerts(certificate string) ServerConfigBuilder {
	caCert, err := ioutil.ReadFile(certificate)
	if err != nil {
		if b.rootCaPool == nil {
			b.rootCaPool = x509.NewCertPool()
		}
		_ = b.rootCaPool.AppendCertsFromPEM(caCert)
	}
	return b
}

func (b *serverConfigBuilder) WithCertificateManager(dir string) ServerConfigBuilder {
	b.certManager = &autocert.Manager{
		Prompt: autocert.AcceptTOS,
		Cache: autocert.DirCache(dir),
	}
	return b
}

func (b *serverConfigBuilder) WithMinVersion(min uint16) ServerConfigBuilder {
	b.minVersion = min
	return b
}

func (b *serverConfigBuilder) WithInsecureSkipVerify(insecure bool) ServerConfigBuilder {
	b.insecure = insecure
	return b
}

func (b *serverConfigBuilder) WithRenegotiationSupport(renegotiation tls.RenegotiationSupport) ServerConfigBuilder {
	b.renegotiation = renegotiation
	return b
}

func (b *serverConfigBuilder) WithClientSessionCache(cache tls.ClientSessionCache) ServerConfigBuilder {
	b.cache = cache
	return b
}

func (b *serverConfigBuilder) MoreCipherSuites(cipherSuite uint16) ServerConfigBuilder {
	b.cipherSuits = append(b.cipherSuits, cipherSuite)
	return b
}

func (b *serverConfigBuilder) MoreCurvePreferences(curve tls.CurveID) ServerConfigBuilder {
	b.curvePref = append(b.curvePref, curve)
	return b
}

func (b *serverConfigBuilder) WithPreferServerCipherSuites(preferServerCipherSuites bool)  ServerConfigBuilder {
	b.preferServerCipherSuites = preferServerCipherSuites
	return b
}

func (b *serverConfigBuilder) Build() (model.ServerConfig, error) {
	var err error
	var getCert func(info *tls.ClientHelloInfo) (*tls.Certificate, error)
	if b.certManager != nil {
		getCert = b.certManager.GetCertificate
	}
	return model.ServerConfig{
		Host: b.address,
		Port: b.port,
		CertPath: b.certificate,
		KeyPath: b.key,
		Config: &tls.Config{
			ClientCAs: b.caPool,
			Certificates: b.certificates,
			CipherSuites: b.cipherSuits,
			InsecureSkipVerify: b.insecure,
			CurvePreferences: b.curvePref,
			RootCAs: b.rootCaPool,
			GetCertificate: getCert,
			MinVersion: b.minVersion,
			PreferServerCipherSuites: b.preferServerCipherSuites,
			ClientSessionCache: b.cache,
			Rand: rand.Reader,
			Renegotiation: b.renegotiation,
		},
	}, err
}

func NewServerConfigBuilder() ServerConfigBuilder{
	return &serverConfigBuilder{
		certificates: make([]tls.Certificate, 0),
		cipherSuits: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
		curvePref: []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		minVersion: tls.VersionTLS12,
		renegotiation: tls.RenegotiateNever,

	}
}