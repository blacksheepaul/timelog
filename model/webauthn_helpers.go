package model

import (
	"strings"

	"github.com/go-webauthn/webauthn/protocol"
)

func serializeCredentialTransport(transports []protocol.AuthenticatorTransport) string {
	if len(transports) == 0 {
		return ""
	}

	values := make([]string, len(transports))
	for i, transport := range transports {
		values[i] = string(transport)
	}

	return strings.Join(values, ",")
}

func parseCredentialTransport(raw string) []protocol.AuthenticatorTransport {
	if raw == "" {
		return nil
	}

	parts := strings.Split(raw, ",")
	transports := make([]protocol.AuthenticatorTransport, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		transports = append(transports, protocol.AuthenticatorTransport(trimmed))
	}

	return transports
}

func parseAuthenticatorAttachment(raw string) protocol.AuthenticatorAttachment {
	switch raw {
	case string(protocol.Platform):
		return protocol.Platform
	case string(protocol.CrossPlatform):
		return protocol.CrossPlatform
	default:
		return protocol.AuthenticatorAttachment(raw)
	}
}
