package headers

import (
	"crypto/x509"
	"fmt"
	"encoding/base64"
	"strings"
	"net"
	utls "github.com/refraction-networking/utls"
	"golang.org/x/net/http2"
)

type ProxyAuth struct {
	Address  string
	Username string
	Password string
}

func BasicAuth(user, pass string) string {
	auth := fmt.Sprintf("%s:%s", user, pass)
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

func DialTLSOverProxy(proxy ProxyAuth, host string) (*http2.ClientConn, *utls.UConn, error) {
	// Connect to proxy
	conn, err := net.Dial("tcp", proxy.Address)
	if err != nil {
		return nil, nil, fmt.Errorf("proxy dial error: %w", err)
	}

	// Send CONNECT request
	connectReq := fmt.Sprintf("CONNECT %s:443 HTTP/1.1\r\nHost: %s:443\r\nProxy-Authorization: %s\r\n\r\n",
		host, host, BasicAuth(proxy.Username, proxy.Password))
	if _, err := conn.Write([]byte(connectReq)); err != nil {
		return nil, nil, fmt.Errorf("proxy CONNECT write failed: %w", err)
	}

	// Read proxy response
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, nil, fmt.Errorf("proxy CONNECT read failed: %w", err)
	}
	if !strings.Contains(string(buf[:n]), "200") {
		return nil, nil, fmt.Errorf("proxy CONNECT failed: %s", string(buf[:n]))
	}

	// Setup TLS config
	rootCAs, err := x509.SystemCertPool()
	if err != nil {
		return nil, nil, fmt.Errorf("cert pool error: %w", err)
	}
	tlsConf := &utls.Config{
		ServerName:         host,
		RootCAs:            rootCAs,
		NextProtos:         []string{"h2"},
		InsecureSkipVerify: false,
	}
	uconn := utls.UClient(conn, tlsConf, utls.HelloChrome_131)

	// Perform handshake
	if err := uconn.Handshake(); err != nil {
		return nil, nil, fmt.Errorf("uTLS handshake failed: %w", err)
	}
	if uconn.ConnectionState().NegotiatedProtocol != "h2" {
		return nil, nil, fmt.Errorf("ALPN negotiation failed: got %s", uconn.ConnectionState().NegotiatedProtocol)
	}

	// Create HTTP/2 connection
	transport := &http2.Transport{}
	clientConn, err := transport.NewClientConn(uconn)
	if err != nil {
		return nil, nil, fmt.Errorf("http2 conn error: %w", err)
	}

	return clientConn, uconn, nil
}

func DialTLSDirect(host string) (*http2.ClientConn, *utls.UConn, error) {
	conn, err := net.Dial("tcp", host+":443")
	if err != nil {
		return nil, nil, fmt.Errorf("dial error: %w", err)
	}

	rootCAs, err := x509.SystemCertPool()
	if err != nil {
		return nil, nil, fmt.Errorf("cert pool error: %w", err)
	}

	tlsConf := &utls.Config{
		ServerName:         host,
		RootCAs:            rootCAs,
		NextProtos:         []string{"h2"},
		InsecureSkipVerify: false,
	}
	uconn := utls.UClient(conn, tlsConf, utls.HelloChrome_131)

	if err := uconn.Handshake(); err != nil {
		return nil, nil, fmt.Errorf("handshake error: %w", err)
	}

	if uconn.ConnectionState().NegotiatedProtocol != "h2" {
		return nil, nil, fmt.Errorf("ALPN negotiation failed: got %s", uconn.ConnectionState().NegotiatedProtocol)
	}

	transport := &http2.Transport{}
	clientConn, err := transport.NewClientConn(uconn)
	if err != nil {
		return nil, nil, fmt.Errorf("http2 conn error: %w", err)
	}

	return clientConn, uconn, nil
}

func DialTLS(host string, proxy string) (*http2.ClientConn, *utls.UConn, error) {
	if proxy == "" {
		return DialTLSDirect(host)
	}

	parts := strings.Split(proxy, ":")
	if len(parts) != 4 {
		return nil, nil, fmt.Errorf("invalid proxy format (expected ip:port:user:pass)")
	}
	auth := ProxyAuth{
		Address:  parts[0] + ":" + parts[1],
		Username: parts[2],
		Password: parts[3],
	}
	return DialTLSOverProxy(auth, host)
}
