package router

import "github.com/go-webauthn/webauthn/protocol"

type passkeyRegisterCreationResponse struct {
	SessionID string                       `json:"session_id"`
	Data      *protocol.CredentialCreation `json:"data"`
}

type passkeyLoginAssertionResponse struct {
	SessionID string                        `json:"session_id"`
	Data      *protocol.CredentialAssertion `json:"data"`
}
