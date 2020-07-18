package builders

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"github.com/hellgate75/go-network/model"
	"golang.org/x/crypto/acme/autocert"
	"io/ioutil"
)

// Helper for building a model.ClientConfig instance
type ClientConfigBuilder interface {
	// Associate an host and a port to the builder workflow
	WithHost(protocol, address string, port int) ClientConfigBuilder
	// Associate certificate and key files full path to the builder workflow
	WithTLSCerts(certificate string, key string) ClientConfigBuilder
	// Add some more certificate files to the certificate list to the builder workflow
	// If no certificate is settled up first call with associate the main TLS certificate files
	MoreTLSCerts(certificate string, key string) ClientConfigBuilder
	// Add one root CA certificate files to the certificate list to the builder workflow
	WithRootCaCert(certificate string) ClientConfigBuilder
	// Add one client CA certificate files to the certificate list to the builder workflow
	WithClientCaCert(certificate string) ClientConfigBuilder
	// Set up the certificate manager for the auto-scan of certificates for a folder
	WithCertificateManager(dir string) ClientConfigBuilder
	// Add more root CA certificate files to the certificate list to the builder workflow
	MoreClientCaCerts(certificate string) ClientConfigBuilder
	// Add more client CA certificate files to the certificate list to the builder workflow
	MoreRootCaCerts(certificate string) ClientConfigBuilder
	// Set min version different from tls.VersionTLS12
	WithMinVersion(min uint16) ClientConfigBuilder
	// Set the insecure skip verify flag, by default it's false
	WithInsecureSkipVerify(insecure bool) ClientConfigBuilder
	// Set up renegotiation, by default it's sett up to: tls.RenegotiateNever
	WithRenegotiationSupport(renegotiation tls.RenegotiationSupport) ClientConfigBuilder
	// Set up the Client Session Cache manager (suggested: tls.NewLRUClientSessionCache(1024) or more ...)
	WithClientSessionCache(cache tls.ClientSessionCache) ClientConfigBuilder
	// Add more Cipher suites to the preset values : tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	// tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA, tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
	// tls.TLS_RSA_WITH_AES_256_CBC_SHA
	MoreCipherSuites(cipherSuite uint16) ClientConfigBuilder
	// Add more Curve Ids to the current TLS Curve Preferences, adding to the preset values : tls.CurveP521, tls.CurveP384,
	// tls.CurveP256
	MoreCurvePreferences(curve tls.CurveID) ClientConfigBuilder
	// Set preference for Server Size Cipher Suite
	WithPreferServerCipherSuites(preferServerCipherSuites bool)  ClientConfigBuilder
	// Build the model.ClientConfig and report any error occurred during the build process
	Build() (model.ClientConfig, error)
}

type clientConfigBuilder struct{
	protocol	 				string
	address      				string
	port         				int
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

func (b *clientConfigBuilder) WithHost(protocol, address string, port int) ClientConfigBuilder {
	b.protocol = protocol
	b.address = address
	b.port = port
	return b
}

func (b *clientConfigBuilder) WithTLSCerts(certificate string, key string) ClientConfigBuilder {
	cert, err := tls.LoadX509KeyPair(certificate, key)
	if err == nil {
		b.certificates = append(b.certificates, cert)
	}
	return b
}

func (b *clientConfigBuilder) MoreTLSCerts(certificate string, key string) ClientConfigBuilder {
	cert, err := tls.LoadX509KeyPair(certificate, key)
	if err == nil {
		b.certificates = append(b.certificates, cert)
	}
	return b
}

func (b *clientConfigBuilder) WithRootCaCert(certificate string) ClientConfigBuilder {
	caCert, err := ioutil.ReadFile(certificate)
	if err != nil {
		if b.rootCaPool == nil {
			b.rootCaPool = x509.NewCertPool()
		}
		_ = b.rootCaPool.AppendCertsFromPEM(caCert)
	}
	return b
}

func (b *clientConfigBuilder) WithClientCaCert(certificate string) ClientConfigBuilder {
	caCert, err := ioutil.ReadFile(certificate)
	if err != nil {
		if b.caPool == nil {
			b.caPool = x509.NewCertPool()
		}
		_ = b.caPool.AppendCertsFromPEM(caCert)
	}
	return b
}

func (b *clientConfigBuilder) MoreClientCaCerts(certificate string) ClientConfigBuilder {
	caCert, err := ioutil.ReadFile(certificate)
	if err != nil {
		if b.caPool == nil {
			b.caPool = x509.NewCertPool()
		}
		_ = b.caPool.AppendCertsFromPEM(caCert)
	}
	return b
}

func (b *clientConfigBuilder) MoreRootCaCerts(certificate string) ClientConfigBuilder {
	caCert, err := ioutil.ReadFile(certificate)
	if err != nil {
		if b.rootCaPool == nil {
			b.rootCaPool = x509.NewCertPool()
		}
		_ = b.rootCaPool.AppendCertsFromPEM(caCert)
	}
	return b
}

func (b *clientConfigBuilder) WithCertificateManager(dir string) ClientConfigBuilder {
	b.certManager = &autocert.Manager{
		Prompt: autocert.AcceptTOS,
		Cache: autocert.DirCache(dir),
	}
	return b
}

func (b *clientConfigBuilder) WithMinVersion(min uint16) ClientConfigBuilder {
	b.minVersion = min
	return b
}

func (b *clientConfigBuilder) WithInsecureSkipVerify(insecure bool) ClientConfigBuilder {
	b.insecure = insecure
	return b
}

func (b *clientConfigBuilder) WithRenegotiationSupport(renegotiation tls.RenegotiationSupport) ClientConfigBuilder {
	b.renegotiation = renegotiation
	return b
}

func (b *clientConfigBuilder) WithClientSessionCache(cache tls.ClientSessionCache) ClientConfigBuilder {
	b.cache = cache
	return b
}

func (b *clientConfigBuilder) MoreCipherSuites(cipherSuite uint16) ClientConfigBuilder {
	b.cipherSuits = append(b.cipherSuits, cipherSuite)
	return b
}

func (b *clientConfigBuilder) WithPreferServerCipherSuites(preferServerCipherSuites bool)  ClientConfigBuilder {
	b.preferServerCipherSuites = preferServerCipherSuites
	return b
}

func (b *clientConfigBuilder) MoreCurvePreferences(curve tls.CurveID) ClientConfigBuilder {
	b.curvePref = append(b.curvePref, curve)
	return b
}

func (b *clientConfigBuilder) Build() (model.ClientConfig, error) {
	var err error
	var getCert func(info *tls.ClientHelloInfo) (*tls.Certificate, error)
	if b.certManager != nil {
		getCert = b.certManager.GetCertificate
	}
	return model.ClientConfig{
		Host: b.address,
		Port: b.port,
		Protocol: b.protocol,
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

func NewClientConfigBuilder() ClientConfigBuilder{
	return &clientConfigBuilder{
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