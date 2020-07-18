package builders

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"github.com/hellgate75/go-network/model"
	"github.com/hellgate75/go-network/model/encoding"
	"golang.org/x/crypto/acme/autocert"
	"io/ioutil"
)

// Helper for building a model.ClientConfig instance
type TcpClientConfigBuilder interface {
	// Use Tls encryption over standard plain communication protocol
	UseTlsEncryption(use bool) TcpClientConfigBuilder
	// Associate a custom network than the default 'tcp' one
	WithNetwork(network string) TcpClientConfigBuilder
	// Associate a custom encoding than the default jason format -> responding to 'application/json' Mime type
	WithEncoding(enc encoding.Encoding) TcpClientConfigBuilder
	// Associate an host and a port to the builder workflow
	WithHost(address string, port int) TcpClientConfigBuilder
	// Associate certificate and key files full path to the builder workflow
	WithTLSCerts(certificate string, key string) TcpClientConfigBuilder
	// Add some more certificate files to the certificate list to the builder workflow
	// If no certificate is settled up first call with associate the main TLS certificate files
	MoreTLSCerts(certificate string, key string) TcpClientConfigBuilder
	// Add one root CA certificate files to the certificate list to the builder workflow
	WithRootCaCert(certificate string) TcpClientConfigBuilder
	// Add one client CA certificate files to the certificate list to the builder workflow
	WithClientCaCert(certificate string) TcpClientConfigBuilder
	// Set up the certificate manager for the auto-scan of certificates for a folder
	WithCertificateManager(dir string) TcpClientConfigBuilder
	// Add more root CA certificate files to the certificate list to the builder workflow
	MoreClientCaCerts(certificate string) TcpClientConfigBuilder
	// Add more client CA certificate files to the certificate list to the builder workflow
	MoreRootCaCerts(certificate string) TcpClientConfigBuilder
	// Set min version different from tls.VersionTLS12
	WithMinVersion(min uint16) TcpClientConfigBuilder
	// Set the insecure skip verify flag, by default it's false
	WithInsecureSkipVerify(insecure bool) TcpClientConfigBuilder
	// Set up renegotiation, by default it's sett up to: tls.RenegotiateNever
	WithRenegotiationSupport(renegotiation tls.RenegotiationSupport) TcpClientConfigBuilder
	// Set up the Client Session Cache manager (suggested: tls.NewLRUClientSessionCache(1024) or more ...)
	WithClientSessionCache(cache tls.ClientSessionCache) TcpClientConfigBuilder
	// Add more Cipher suites to the preset values : tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	// tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA, tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
	// tls.TLS_RSA_WITH_AES_256_CBC_SHA
	MoreCipherSuites(cipherSuite uint16) TcpClientConfigBuilder
	// Add more Curve Ids to the current TLS Curve Preferences, adding to the preset values : tls.CurveP521, tls.CurveP384,
	// tls.CurveP256
	MoreCurvePreferences(curve tls.CurveID) TcpClientConfigBuilder
	// Set preference for Server Size Cipher Suite
	WithPreferServerCipherSuites(preferServerCipherSuites bool) TcpClientConfigBuilder
	// Build the model.ClientConfig and report any error occurred during the build process
	Build() (model.TcpClientConfig, error)
}

type tcpClientConfigBuilder struct{
	useTls					 bool
	network                  string
	enc                  	 encoding.Encoding
	address                  string
	port                     int
	caPool                   *x509.CertPool
	rootCaPool               *x509.CertPool
	certificates             []tls.Certificate
	cipherSuits              []uint16
	insecure                 bool
	curvePref                []tls.CurveID
	certManager              *autocert.Manager
	minVersion               uint16
	renegotiation            tls.RenegotiationSupport
	cache                    tls.ClientSessionCache
	preferServerCipherSuites bool
}

func (b *tcpClientConfigBuilder) UseTlsEncryption(use bool) TcpClientConfigBuilder {
	b.useTls = use
	return b
}

func (b *tcpClientConfigBuilder) WithNetwork(network string) TcpClientConfigBuilder {
	b.network = network
	return b
}

func (b *tcpClientConfigBuilder) WithHost(address string, port int) TcpClientConfigBuilder {
	b.address = address
	b.port = port
	return b
}

func (b *tcpClientConfigBuilder) WithEncoding(enc encoding.Encoding) TcpClientConfigBuilder {
	b.enc=enc
	return b
}

func (b *tcpClientConfigBuilder) WithTLSCerts(certificate string, key string) TcpClientConfigBuilder {
	cert, err := tls.LoadX509KeyPair(certificate, key)
	if err == nil {
		b.certificates = append(b.certificates, cert)
	}
	return b
}

func (b *tcpClientConfigBuilder) MoreTLSCerts(certificate string, key string) TcpClientConfigBuilder {
	cert, err := tls.LoadX509KeyPair(certificate, key)
	if err == nil {
		b.certificates = append(b.certificates, cert)
	}
	return b
}

func (b *tcpClientConfigBuilder) WithRootCaCert(certificate string) TcpClientConfigBuilder {
	caCert, err := ioutil.ReadFile(certificate)
	if err != nil {
		if b.rootCaPool == nil {
			b.rootCaPool = x509.NewCertPool()
		}
		_ = b.rootCaPool.AppendCertsFromPEM(caCert)
	}
	return b
}

func (b *tcpClientConfigBuilder) WithClientCaCert(certificate string) TcpClientConfigBuilder {
	caCert, err := ioutil.ReadFile(certificate)
	if err != nil {
		if b.caPool == nil {
			b.caPool = x509.NewCertPool()
		}
		_ = b.caPool.AppendCertsFromPEM(caCert)
	}
	return b
}

func (b *tcpClientConfigBuilder) MoreClientCaCerts(certificate string) TcpClientConfigBuilder {
	caCert, err := ioutil.ReadFile(certificate)
	if err != nil {
		if b.caPool == nil {
			b.caPool = x509.NewCertPool()
		}
		_ = b.caPool.AppendCertsFromPEM(caCert)
	}
	return b
}

func (b *tcpClientConfigBuilder) MoreRootCaCerts(certificate string) TcpClientConfigBuilder {
	caCert, err := ioutil.ReadFile(certificate)
	if err != nil {
		if b.rootCaPool == nil {
			b.rootCaPool = x509.NewCertPool()
		}
		_ = b.rootCaPool.AppendCertsFromPEM(caCert)
	}
	return b
}

func (b *tcpClientConfigBuilder) WithCertificateManager(dir string) TcpClientConfigBuilder {
	b.certManager = &autocert.Manager{
		Prompt: autocert.AcceptTOS,
		Cache: autocert.DirCache(dir),
	}
	return b
}

func (b *tcpClientConfigBuilder) WithMinVersion(min uint16) TcpClientConfigBuilder {
	b.minVersion = min
	return b
}

func (b *tcpClientConfigBuilder) WithInsecureSkipVerify(insecure bool) TcpClientConfigBuilder {
	b.insecure = insecure
	return b
}

func (b *tcpClientConfigBuilder) WithRenegotiationSupport(renegotiation tls.RenegotiationSupport) TcpClientConfigBuilder {
	b.renegotiation = renegotiation
	return b
}

func (b *tcpClientConfigBuilder) WithClientSessionCache(cache tls.ClientSessionCache) TcpClientConfigBuilder {
	b.cache = cache
	return b
}

func (b *tcpClientConfigBuilder) MoreCipherSuites(cipherSuite uint16) TcpClientConfigBuilder {
	b.cipherSuits = append(b.cipherSuits, cipherSuite)
	return b
}

func (b *tcpClientConfigBuilder) WithPreferServerCipherSuites(preferServerCipherSuites bool) TcpClientConfigBuilder {
	b.preferServerCipherSuites = preferServerCipherSuites
	return b
}

func (b *tcpClientConfigBuilder) MoreCurvePreferences(curve tls.CurveID) TcpClientConfigBuilder {
	b.curvePref = append(b.curvePref, curve)
	return b
}

func (b *tcpClientConfigBuilder) Build() (model.TcpClientConfig, error) {
	var err error
	var getCert func(info *tls.ClientHelloInfo) (*tls.Certificate, error)
	if b.certManager != nil {
		getCert = b.certManager.GetCertificate
	}
	var tlsConfig *tls.Config
	if b.useTls {
		tlsConfig = &tls.Config{
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
		}
	}
	return model.TcpClientConfig{
		Host: b.address,
		Port: b.port,
		Network: b.network,
		Encoding: b.enc,
		Config: tlsConfig,
	}, err
}

func NewTcpClientConfigBuilder() TcpClientConfigBuilder {
	return &tcpClientConfigBuilder{
		enc: encoding.EncodingJSONFormat,
		network: "tcp",
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