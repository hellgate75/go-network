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

// Helper for building a model.TcpServerConfig instance
type TcpServerConfigBuilder interface {
	// Use Tls encryption over standard plain communication protocol
	UseTlsEncryption(use bool) TcpServerConfigBuilder
	// Associate a custom network than the default 'tcp' one
	WithNetwork(network string) TcpServerConfigBuilder
	// Associate a custom encoding than the default jason format -> responding to 'application/json' Mime type
	WithEncoding(enc encoding.Encoding) TcpServerConfigBuilder
	// Associate an host and a port to the builder workflow
	WithHost(address string, port int) TcpServerConfigBuilder
	// Add a certificate files to the certificate list to the builder workflow
	WithTLSCerts(certificate string, key string) TcpServerConfigBuilder
	// Add some more certificate files to the certificate list to the builder workflow
	MoreTLSCerts(certificate string, key string) TcpServerConfigBuilder
	// Add one root CA certificate files to the certificate list to the builder workflow
	WithRootCaCert(certificate string) TcpServerConfigBuilder
	// Add one client CA certificate files to the certificate list to the builder workflow
	WithClientCaCert(certificate string) TcpServerConfigBuilder
	// Set up the certificate manager for the auto-scan of certificates for a folder
	WithCertificateManager(dir string) TcpServerConfigBuilder
	// Add more root CA certificate files to the certificate list to the builder workflow
	MoreClientCaCerts(certificate string) TcpServerConfigBuilder
	// Add more client CA certificate files to the certificate list to the builder workflow
	MoreRootCaCerts(certificate string) TcpServerConfigBuilder
	// Set min version different from tls.VersionTLS12
	WithMinVersion(min uint16) TcpServerConfigBuilder
	// Set the insecure skip verify flag, by default it's false
	WithInsecureSkipVerify(insecure bool) TcpServerConfigBuilder
	// Set up renegotiation, by default it's sett up to: tls.RenegotiateNever
	WithRenegotiationSupport(renegotiation tls.RenegotiationSupport) TcpServerConfigBuilder
	// Set up the Client Session Cache manager (suggested: tls.NewLRUClientSessionCache(1024) or more ...)
	WithClientSessionCache(cache tls.ClientSessionCache) TcpServerConfigBuilder
	// Add more Cipher suites to the preset values : tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	// tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA, tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
	// tls.TLS_RSA_WITH_AES_256_CBC_SHA
	MoreCipherSuites(cipherSuite uint16) TcpServerConfigBuilder
	// Add more Curve Ids to the current TLS Curve Preferences, adding to the preset values : tls.CurveP521, tls.CurveP384,
	// tls.CurveP256
	MoreCurvePreferences(curve tls.CurveID) TcpServerConfigBuilder
	// Set preference for Server Size Cipher Suite
	WithPreferServerCipherSuites(preferServerCipherSuites bool) TcpServerConfigBuilder
	// Build the model.ServerConfig and report any error occurred during the build process
	Build() (model.TcpServerConfig, error)
}

type serverConfigBuilder struct{
	useTls					 	bool
	address      				string
	port         				int
	network  					string
	enc							encoding.Encoding
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

func (b *serverConfigBuilder) UseTlsEncryption(use bool) TcpServerConfigBuilder {
	b.useTls = use
	return b
}

func (b *serverConfigBuilder) WithNetwork(network string) TcpServerConfigBuilder {
	b.network = network
	return b
}

func (b *serverConfigBuilder) WithHost(address string, port int) TcpServerConfigBuilder {
	b.address = address
	b.port = port
	return b
}

func (b *serverConfigBuilder) WithEncoding(enc encoding.Encoding) TcpServerConfigBuilder {
	b.enc=enc
	return b
}

func (b *serverConfigBuilder) WithTLSCerts(certificate string, key string) TcpServerConfigBuilder {
	cert, err := tls.LoadX509KeyPair(certificate, key)
	if err == nil {
		b.certificates = append(b.certificates, cert)
	}
	return b
}

func (b *serverConfigBuilder) MoreTLSCerts(certificate string, key string) TcpServerConfigBuilder {
	cert, err := tls.LoadX509KeyPair(certificate, key)
	if err == nil {
		b.certificates = append(b.certificates, cert)
	}
	return b
}

func (b *serverConfigBuilder) WithRootCaCert(certificate string) TcpServerConfigBuilder {
	caCert, err := ioutil.ReadFile(certificate)
	if err != nil {
		if b.rootCaPool == nil {
			b.rootCaPool = x509.NewCertPool()
		}
		_ = b.rootCaPool.AppendCertsFromPEM(caCert)
	}
	return b
}

func (b *serverConfigBuilder) WithClientCaCert(certificate string) TcpServerConfigBuilder {
	caCert, err := ioutil.ReadFile(certificate)
	if err != nil {
		if b.caPool == nil {
			b.caPool = x509.NewCertPool()
		}
		_ = b.caPool.AppendCertsFromPEM(caCert)
	}
	return b
}

func (b *serverConfigBuilder) MoreClientCaCerts(certificate string) TcpServerConfigBuilder {
	caCert, err := ioutil.ReadFile(certificate)
	if err != nil {
		if b.caPool == nil {
			b.caPool = x509.NewCertPool()
		}
		_ = b.caPool.AppendCertsFromPEM(caCert)
	}
	return b
}

func (b *serverConfigBuilder) MoreRootCaCerts(certificate string) TcpServerConfigBuilder {
	caCert, err := ioutil.ReadFile(certificate)
	if err != nil {
		if b.rootCaPool == nil {
			b.rootCaPool = x509.NewCertPool()
		}
		_ = b.rootCaPool.AppendCertsFromPEM(caCert)
	}
	return b
}

func (b *serverConfigBuilder) WithCertificateManager(dir string) TcpServerConfigBuilder {
	b.certManager = &autocert.Manager{
		Prompt: autocert.AcceptTOS,
		Cache: autocert.DirCache(dir),
	}
	return b
}

func (b *serverConfigBuilder) WithMinVersion(min uint16) TcpServerConfigBuilder {
	b.minVersion = min
	return b
}

func (b *serverConfigBuilder) WithInsecureSkipVerify(insecure bool) TcpServerConfigBuilder {
	b.insecure = insecure
	return b
}

func (b *serverConfigBuilder) WithRenegotiationSupport(renegotiation tls.RenegotiationSupport) TcpServerConfigBuilder {
	b.renegotiation = renegotiation
	return b
}

func (b *serverConfigBuilder) WithClientSessionCache(cache tls.ClientSessionCache) TcpServerConfigBuilder {
	b.cache = cache
	return b
}

func (b *serverConfigBuilder) MoreCipherSuites(cipherSuite uint16) TcpServerConfigBuilder {
	b.cipherSuits = append(b.cipherSuits, cipherSuite)
	return b
}

func (b *serverConfigBuilder) MoreCurvePreferences(curve tls.CurveID) TcpServerConfigBuilder {
	b.curvePref = append(b.curvePref, curve)
	return b
}

func (b *serverConfigBuilder) WithPreferServerCipherSuites(preferServerCipherSuites bool) TcpServerConfigBuilder {
	b.preferServerCipherSuites = preferServerCipherSuites
	return b
}

func (b *serverConfigBuilder) Build() (model.TcpServerConfig, error) {
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
	return model.TcpServerConfig{
		Host: b.address,
		Port: b.port,
		Encoding: b.enc,
		Network: b.network,
		Config: tlsConfig,
	}, err
}

func NewTcpServerConfigBuilder() TcpServerConfigBuilder {
	return &serverConfigBuilder{
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