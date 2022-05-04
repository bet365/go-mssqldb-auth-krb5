// Package krb5 implements the mssql.Auth interface in order to provide kerberos/active directory (Windows) based authentication.
package krb5

import (
	"fmt"
	"net"
	"strings"

	"github.com/microsoft/go-mssqldb/integratedauth"

	"github.com/jcmturner/gokrb5/v8/client"
	"github.com/jcmturner/gokrb5/v8/gssapi"
	"github.com/jcmturner/gokrb5/v8/spnego"
)

var (
	_ integratedauth.Provider                = &AuthProvider{}
	_ integratedauth.IntegratedAuthenticator = &krbAuth{}
)

// NewAuthProvider creates an instance of the Auth interface for kerberos authentication
func NewAuthProvider(client *client.Client) *AuthProvider {
	a := &AuthProvider{
		client: client,
	}

	return a
}

// AuthProvider wraps gokrb5/v8/client and can in turn allow instances of mssql.Auth to be created which then handle kerberos/active directory logins.
type AuthProvider struct {
	client *client.Client
}

// GetAuth returns an instance of the mssql.Auth interface. That is then responsible for kerberos Service Provider Negotiation.
func (a AuthProvider) GetIntegratedAuthenticator(user, _, service, _ string) (integratedauth.IntegratedAuthenticator, bool) {
	// If we've got a user, assume SQL Authentication
	if user != "" {
		return nil, false
	}

	spnegoClient := spnego.SPNEGOClient(a.client, canonicalize(service))
	return &krbAuth{client: spnegoClient}, true
}

// responsible for transforming network CNames into their actual Hostname.
// For cases where service tickets can only be bound to hostnames, not cnames.
func canonicalize(service string) string {
	parts := strings.SplitAfterN(service, "/", 2)
	if len(parts) != 2 {
		return service
	}
	host, port, err := net.SplitHostPort(parts[1])
	if err != nil {
		return service
	}
	cname, err := net.LookupCNAME(strings.ToLower(host))
	if err != nil {
		return service
	}
	// Put service back together with cname (stripped of trailing .) and port
	return parts[0] + net.JoinHostPort(cname[:len(cname)-1], port)
}

// krbAuth implements the mssql.Auth interface. It is responsible for kerberos Service Provider Negotiation.
type krbAuth struct {
	client *spnego.SPNEGO
}

func (k *krbAuth) InitialBytes() ([]byte, error) {
	tkn, err := k.client.InitSecContext()
	if err != nil {
		return nil, err
	}
	return tkn.Marshal()
}

func (k *krbAuth) NextBytes(bytes []byte) ([]byte, error) {
	var resp spnego.SPNEGOToken
	if err := resp.Unmarshal(bytes); err != nil {
		return nil, err
	}

	ok, status := resp.Verify()
	if ok { // we're ok, done
		return nil, nil
	}

	switch status.Code {
	case gssapi.StatusContinueNeeded:
		return nil, nil
	// case gssapi.StatusComplete: // could also be ok
	default:
		return nil, fmt.Errorf("bad status: %+v", status)
	}
}

func (k *krbAuth) Free() {}
