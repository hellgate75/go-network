package builders

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"github.com/hellgate75/go-network/model"
	"golang.org/x/crypto/acme/autocert"
	"io/ioutil"
	"regexp"
	"strings"
)

// Helper for building a model.PipeNodeConfig instance
type PipeNodeConfigBuilder interface {
	// Use Tls encryption over standard plain communication protocol
	UseTlsEncryption(use bool) PipeNodeConfigBuilder
	// Associate a custom network than the default 'tcp' one
	WithNetwork(network string) PipeNodeConfigBuilder
	// Associate an inHost and a inPort to the builder workflow, setting-up or tear-sown input pipe node mode
	WithInHost(address string, port int) PipeNodeConfigBuilder
	// Associate an outHost and a outPort to the builder workflow, setting-up or tear-sown output pipe node mode
	WithOutHost(address string, port int) PipeNodeConfigBuilder
	// Add a certificate files to the certificate list to the builder workflow
	WithTLSCerts(certificate string, key string) PipeNodeConfigBuilder
	// Add some more certificate files to the certificate list to the builder workflow
	MoreTLSCerts(certificate string, key string) PipeNodeConfigBuilder
	// Add one root CA certificate files to the certificate list to the builder workflow
	WithRootCaCert(certificate string) PipeNodeConfigBuilder
	// Add one client CA certificate files to the certificate list to the builder workflow
	WithClientCaCert(certificate string) PipeNodeConfigBuilder
	// Set up the certificate manager for the auto-scan of certificates for a folder
	WithCertificateManager(dir string) PipeNodeConfigBuilder
	// Add more root CA certificate files to the certificate list to the builder workflow
	MoreClientCaCerts(certificate string) PipeNodeConfigBuilder
	// Add more client CA certificate files to the certificate list to the builder workflow
	MoreRootCaCerts(certificate string) PipeNodeConfigBuilder
	// Set min version different from tls.VersionTLS12
	WithMinVersion(min uint16) PipeNodeConfigBuilder
	// Set the insecure skip verify flag, by default it's false
	WithInsecureSkipVerify(insecure bool) PipeNodeConfigBuilder
	// Set up renegotiation, by default it's sett up to: tls.RenegotiateNever
	WithRenegotiationSupport(renegotiation tls.RenegotiationSupport) PipeNodeConfigBuilder
	// Set up the Client Session Cache manager (suggested: tls.NewLRUClientSessionCache(1024) or more ...)
	WithClientSessionCache(cache tls.ClientSessionCache) PipeNodeConfigBuilder
	// Add more Cipher suites to the preset values : tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	// tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA, tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
	// tls.TLS_RSA_WITH_AES_256_CBC_SHA
	MoreCipherSuites(cipherSuite uint16) PipeNodeConfigBuilder
	// Add more Curve Ids to the current TLS Curve Preferences, adding to the preset values : tls.CurveP521, tls.CurveP384,
	// tls.CurveP256
	MoreCurvePreferences(curve tls.CurveID) PipeNodeConfigBuilder
	// Set preference for Server Size Cipher Suite
	WithPreferServerCipherSuites(preferServerCipherSuites bool) PipeNodeConfigBuilder
	// Build the model.ServerConfig and report any error occurred during the build process
	Build() (model.PipeNodeConfig, error)
}

type pipeNodeConfigBuilder struct{
	useTls                   bool
	network					 string
	inAddress                string
	inPort                   int
	outAddress               string
	outPort                  int
	pipeType				 model.PipeType
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

func (b *pipeNodeConfigBuilder) UseTlsEncryption(use bool) PipeNodeConfigBuilder {
	b.useTls = use
	return b
}
func (b *pipeNodeConfigBuilder) WithNetwork(network string) PipeNodeConfigBuilder{
	b.network = network
	return b
}


func containsAlpha(s string) bool {
	pattern := `^w+$`
	matched, err := regexp.Match(pattern, []byte(s))
	if err != nil {
		return false
	}
	return matched
}

func isValidAddress(addr string) bool {
	return len(strings.TrimSpace(addr)) == 0 ||
			strings.Count(addr, ".") >= 4 ||
			strings.Count(addr, ":") >= 4 ||
			containsAlpha(addr)
}

func isValidPort(port int, netwotk string) bool {
	return port > 0 ||
		(port == 0 && netwotk != "tcp" && netwotk != "udp")
}

func calculateInModeType(builder *pipeNodeConfigBuilder) {
	if isValidAddress(builder.inAddress)  && isValidPort(builder.inPort, builder.network) {
		if builder.pipeType == model.OutputPipe {
			builder.pipeType = model.InputOutputPipe
		} else if builder.pipeType == model.NoTypeSelected {
			builder.pipeType = model.InputPipe
		}
	} else {
		if builder.pipeType == model.InputPipe {
			builder.pipeType = model.NoTypeSelected
		} else if builder.pipeType == model.InputOutputPipe {
			builder.pipeType = model.OutputPipe
		}
	}
}

func calculateOutModeType(builder *pipeNodeConfigBuilder) {
	if isValidAddress(builder.outAddress)  && isValidPort(builder.outPort, builder.network) {
		if builder.pipeType == model.InputPipe {
			builder.pipeType = model.InputOutputPipe
		} else if builder.pipeType == model.NoTypeSelected {
			builder.pipeType = model.OutputPipe
		}
	} else {
		if builder.pipeType == model.OutputPipe {
			builder.pipeType = model.NoTypeSelected
		} else if builder.pipeType == model.InputOutputPipe {
			builder.pipeType = model.InputPipe
		}
	}
}

func (b *pipeNodeConfigBuilder) WithInHost(address string, port int) PipeNodeConfigBuilder {
	b.inAddress = address
	b.inPort = port
	calculateInModeType(b)
	return b
}

func (b *pipeNodeConfigBuilder) WithOutHost(address string, port int) PipeNodeConfigBuilder {
	b.outAddress = address
	b.outPort = port
	calculateOutModeType(b)
	return b
}

func (b *pipeNodeConfigBuilder) WithTLSCerts(certificate string, key string) PipeNodeConfigBuilder {
	cert, err := tls.LoadX509KeyPair(certificate, key)
	if err == nil {
		b.certificates = append(b.certificates, cert)
	}
	return b
}

func (b *pipeNodeConfigBuilder) MoreTLSCerts(certificate string, key string) PipeNodeConfigBuilder {
	cert, err := tls.LoadX509KeyPair(certificate, key)
	if err == nil {
		b.certificates = append(b.certificates, cert)
	}
	return b
}

func (b *pipeNodeConfigBuilder) WithRootCaCert(certificate string) PipeNodeConfigBuilder {
	caCert, err := ioutil.ReadFile(certificate)
	if err != nil {
		if b.rootCaPool == nil {
			b.rootCaPool = x509.NewCertPool()
		}
		_ = b.rootCaPool.AppendCertsFromPEM(caCert)
	}
	return b
}

func (b *pipeNodeConfigBuilder) WithClientCaCert(certificate string) PipeNodeConfigBuilder {
	caCert, err := ioutil.ReadFile(certificate)
	if err != nil {
		if b.caPool == nil {
			b.caPool = x509.NewCertPool()
		}
		_ = b.caPool.AppendCertsFromPEM(caCert)
	}
	return b
}

func (b *pipeNodeConfigBuilder) MoreClientCaCerts(certificate string) PipeNodeConfigBuilder {
	caCert, err := ioutil.ReadFile(certificate)
	if err != nil {
		if b.caPool == nil {
			b.caPool = x509.NewCertPool()
		}
		_ = b.caPool.AppendCertsFromPEM(caCert)
	}
	return b
}

func (b *pipeNodeConfigBuilder) MoreRootCaCerts(certificate string) PipeNodeConfigBuilder {
	caCert, err := ioutil.ReadFile(certificate)
	if err != nil {
		if b.rootCaPool == nil {
			b.rootCaPool = x509.NewCertPool()
		}
		_ = b.rootCaPool.AppendCertsFromPEM(caCert)
	}
	return b
}

func (b *pipeNodeConfigBuilder) WithCertificateManager(dir string) PipeNodeConfigBuilder {
	b.certManager = &autocert.Manager{
		Prompt: autocert.AcceptTOS,
		Cache: autocert.DirCache(dir),
	}
	return b
}

func (b *pipeNodeConfigBuilder) WithMinVersion(min uint16) PipeNodeConfigBuilder {
	b.minVersion = min
	return b
}

func (b *pipeNodeConfigBuilder) WithInsecureSkipVerify(insecure bool) PipeNodeConfigBuilder {
	b.insecure = insecure
	return b
}

func (b *pipeNodeConfigBuilder) WithRenegotiationSupport(renegotiation tls.RenegotiationSupport) PipeNodeConfigBuilder {
	b.renegotiation = renegotiation
	return b
}

func (b *pipeNodeConfigBuilder) WithClientSessionCache(cache tls.ClientSessionCache) PipeNodeConfigBuilder {
	b.cache = cache
	return b
}

func (b *pipeNodeConfigBuilder) MoreCipherSuites(cipherSuite uint16) PipeNodeConfigBuilder {
	b.cipherSuits = append(b.cipherSuits, cipherSuite)
	return b
}

func (b *pipeNodeConfigBuilder) MoreCurvePreferences(curve tls.CurveID) PipeNodeConfigBuilder {
	b.curvePref = append(b.curvePref, curve)
	return b
}

func (b *pipeNodeConfigBuilder) WithPreferServerCipherSuites(preferServerCipherSuites bool) PipeNodeConfigBuilder {
	b.preferServerCipherSuites = preferServerCipherSuites
	return b
}

func (b *pipeNodeConfigBuilder) Build() (model.PipeNodeConfig, error) {
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
	return model.PipeNodeConfig{
		Network: b.network,
		InHost: b.inAddress,
		InPort: b.inPort,
		OutHost: b.outAddress,
		OutPort: b.outPort,
		Type: b.pipeType,
		Config: tlsConfig,
	}, err
}

func NewPipeNodeConfigBuilder() PipeNodeConfigBuilder {
	return &pipeNodeConfigBuilder{
		network: "tcp",
		pipeType: model.NoTypeSelected,
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