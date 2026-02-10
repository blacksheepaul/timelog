package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"

	"github.com/blacksheepaul/timelog/model"
)

var webAuthnInstance *webauthn.WebAuthn

func InitWebAuthn() error {
	if cfg == nil {
		return errors.New("config not initialized")
	}

	if cfg.Passkey.RPID == "" || cfg.Passkey.RPName == "" || len(cfg.Passkey.RPOrigins) == 0 {
		return errors.New("passkey config missing rp_id/rp_name/rp_origins")
	}

	instance, err := webauthn.New(&webauthn.Config{
		RPID:          cfg.Passkey.RPID,
		RPDisplayName: cfg.Passkey.RPName,
		RPOrigins:     cfg.Passkey.RPOrigins,
		AuthenticatorSelection: protocol.AuthenticatorSelection{
			UserVerification: protocol.VerificationPreferred,
		},
	})
	if err != nil {
		return err
	}

	webAuthnInstance = instance
	return nil
}

func GetWebAuthn() *webauthn.WebAuthn {
	return webAuthnInstance
}

func StorePasskeySession(sessionID string, session *webauthn.SessionData, ttlSeconds int64) error {
	if session == nil {
		return errors.New("session is nil")
	}
	dao := model.GetDao()
	// Namespace the key to prevent confusion with auth tokens
	dao.WriteCache("passkey_session:"+sessionID, session, ttlSeconds)
	return nil
}

func LoadPasskeySession(sessionID string) (*webauthn.SessionData, error) {
	dao := model.GetDao()
	// Use namespaced key to retrieve passkey session
	raw, ok := dao.GetCache("passkey_session:" + sessionID)
	if !ok {
		return nil, errors.New("session not found")
	}

	session, ok := raw.(*webauthn.SessionData)
	if !ok || session == nil {
		return nil, errors.New("invalid session data")
	}

	return session, nil
}

func CreatePasskeyCredential(credential *webauthn.Credential, deviceName string) (*model.WebAuthnCredential, error) {
	dao := model.GetDao()
	record := model.WebAuthnCredentialFromCredential(credential)
	if record == nil {
		return nil, errors.New("credential is nil")
	}
	record.DeviceName = deviceName
	if err := model.CreateWebAuthnCredential(dao.Db(), record); err != nil {
		return nil, err
	}
	return record, nil
}

func ListPasskeyCredentials() ([]model.WebAuthnCredential, error) {
	dao := model.GetDao()
	return model.ListWebAuthnCredentials(dao.Db())
}

func DeletePasskeyCredential(id uint) error {
	dao := model.GetDao()
	return model.DeleteWebAuthnCredential(dao.Db(), id)
}

func LoadPasskeyCredentialByID(rawID []byte) (*model.WebAuthnCredential, error) {
	dao := model.GetDao()
	return model.GetWebAuthnCredentialByCredentialID(dao.Db(), rawID)
}

func GenerateSessionToken() (string, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(tokenBytes), nil
}

func StoreSessionToken(token string, ttlSeconds int64) error {
	dao := model.GetDao()
	// Namespace the key to distinguish from passkey sessions
	dao.WriteCache("auth_token:"+token, true, ttlSeconds)
	return nil
}

func GenerateTempPassword() (string, string, error) {
	raw := make([]byte, 16)
	if _, err := rand.Read(raw); err != nil {
		return "", "", err
	}

	password := hex.EncodeToString(raw)
	hashBytes := sha256.Sum256([]byte(password))
	return password, hex.EncodeToString(hashBytes[:]), nil
}

func CreateTempPassword(ttlSeconds int) (*model.TempPassword, string, error) {
	password, hash, err := GenerateTempPassword()
	if err != nil {
		return nil, "", err
	}

	dao := model.GetDao()
	record := &model.TempPassword{
		PasswordHash: hash,
		ExpiresAt:    time.Now().Add(time.Duration(ttlSeconds) * time.Second),
	}
	if err := model.CreateTempPassword(dao.Db(), record); err != nil {
		return nil, "", err
	}

	return record, password, nil
}

func ListTempPasswords() ([]model.TempPassword, error) {
	dao := model.GetDao()
	return model.ListTempPasswords(dao.Db())
}

func DeleteTempPassword(id uint) error {
	dao := model.GetDao()
	return model.DeleteTempPassword(dao.Db(), id)
}

func CleanupExpiredTempPasswords() error {
	dao := model.GetDao()
	return model.DeleteExpiredTempPasswords(dao.Db(), time.Now())
}

func ValidateTempPassword(password string) (*model.TempPassword, error) {
	hashBytes := sha256.Sum256([]byte(password))
	hash := hex.EncodeToString(hashBytes[:])

	dao := model.GetDao()
	return model.GetTempPasswordByHash(dao.Db(), hash, time.Now())
}
