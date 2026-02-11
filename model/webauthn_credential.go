package model

import (
	"time"

	"github.com/go-webauthn/webauthn/webauthn"
	"gorm.io/gorm"
)

type WebAuthnCredential struct {
	ID                            uint           `gorm:"primaryKey" json:"id"`
	CredentialID                  []byte         `gorm:"column:credential_id;not null;unique" json:"credential_id"`
	PublicKey                     []byte         `gorm:"column:public_key;not null" json:"public_key"`
	AttestationType               string         `gorm:"column:attestation_type;not null" json:"attestation_type"`
	Transport                     string         `gorm:"column:transport;not null" json:"transport"`
	DeviceName                    string         `gorm:"column:device_name" json:"device_name"`
	UserPresent                   bool           `gorm:"column:user_present;not null" json:"user_present"`
	UserVerified                  bool           `gorm:"column:user_verified;not null" json:"user_verified"`
	BackupEligible                bool           `gorm:"column:backup_eligible;not null" json:"backup_eligible"`
	BackupState                   bool           `gorm:"column:backup_state;not null" json:"backup_state"`
	AuthenticatorAAGUID           []byte         `gorm:"column:authenticator_aaguid;not null" json:"authenticator_aaguid"`
	AuthenticatorSignCount        uint32         `gorm:"column:authenticator_sign_count;not null" json:"authenticator_sign_count"`
	AuthenticatorCloneWarning     bool           `gorm:"column:authenticator_clone_warning;not null" json:"authenticator_clone_warning"`
	AuthenticatorAttachment       string         `gorm:"column:authenticator_attachment;not null" json:"authenticator_attachment"`
	AttestationClientDataJSON     []byte         `gorm:"column:attestation_client_data_json;not null" json:"attestation_client_data_json"`
	AttestationClientDataHash     []byte         `gorm:"column:attestation_client_data_hash;not null" json:"attestation_client_data_hash"`
	AttestationAuthenticatorData  []byte         `gorm:"column:attestation_authenticator_data;not null" json:"attestation_authenticator_data"`
	AttestationPublicKeyAlgorithm int64          `gorm:"column:attestation_public_key_algorithm;not null" json:"attestation_public_key_algorithm"`
	AttestationObject             []byte         `gorm:"column:attestation_object;not null" json:"attestation_object"`
	CreatedAt                     time.Time      `json:"created_at"`
	UpdatedAt                     time.Time      `json:"updated_at"`
	DeletedAt                     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (WebAuthnCredential) TableName() string {
	return "webauthn_credentials"
}

func CreateWebAuthnCredential(db *gorm.DB, credential *WebAuthnCredential) error {
	return db.Create(credential).Error
}

func GetWebAuthnCredentialByCredentialID(db *gorm.DB, credentialID []byte) (*WebAuthnCredential, error) {
	var credential WebAuthnCredential
	err := db.Where("credential_id = ?", credentialID).First(&credential).Error
	return &credential, err
}

func ListWebAuthnCredentials(db *gorm.DB) ([]WebAuthnCredential, error) {
	var credentials []WebAuthnCredential
	err := db.Order("created_at DESC").Find(&credentials).Error
	return credentials, err
}

func DeleteWebAuthnCredential(db *gorm.DB, id uint) error {
	return db.Delete(&WebAuthnCredential{}, id).Error
}

func UpdateWebAuthnCredentialAuth(db *gorm.DB, credentialID []byte, credential *webauthn.Credential) error {
	if credential == nil {
		return nil
	}

	return db.Model(&WebAuthnCredential{}).
		Where("credential_id = ?", credentialID).
		Updates(map[string]interface{}{
			"authenticator_sign_count":    credential.Authenticator.SignCount,
			"authenticator_clone_warning": credential.Authenticator.CloneWarning,
			"user_present":                credential.Flags.UserPresent,
			"user_verified":               credential.Flags.UserVerified,
			"backup_eligible":             credential.Flags.BackupEligible,
			"backup_state":                credential.Flags.BackupState,
		}).Error
}

func (w WebAuthnCredential) ToCredential() webauthn.Credential {
	return webauthn.Credential{
		ID:              w.CredentialID,
		PublicKey:       w.PublicKey,
		AttestationType: w.AttestationType,
		Transport:       parseCredentialTransport(w.Transport),
		Flags: webauthn.CredentialFlags{
			UserPresent:    w.UserPresent,
			UserVerified:   w.UserVerified,
			BackupEligible: w.BackupEligible,
			BackupState:    w.BackupState,
		},
		Authenticator: webauthn.Authenticator{
			AAGUID:       w.AuthenticatorAAGUID,
			SignCount:    w.AuthenticatorSignCount,
			CloneWarning: w.AuthenticatorCloneWarning,
			Attachment:   parseAuthenticatorAttachment(w.AuthenticatorAttachment),
		},
		Attestation: webauthn.CredentialAttestation{
			ClientDataJSON:     w.AttestationClientDataJSON,
			ClientDataHash:     w.AttestationClientDataHash,
			AuthenticatorData:  w.AttestationAuthenticatorData,
			PublicKeyAlgorithm: w.AttestationPublicKeyAlgorithm,
			Object:             w.AttestationObject,
		},
	}
}

func WebAuthnCredentialFromCredential(credential *webauthn.Credential) *WebAuthnCredential {
	if credential == nil {
		return nil
	}

	return &WebAuthnCredential{
		CredentialID:                  credential.ID,
		PublicKey:                     credential.PublicKey,
		AttestationType:               credential.AttestationType,
		Transport:                     serializeCredentialTransport(credential.Transport),
		UserPresent:                   credential.Flags.UserPresent,
		UserVerified:                  credential.Flags.UserVerified,
		BackupEligible:                credential.Flags.BackupEligible,
		BackupState:                   credential.Flags.BackupState,
		AuthenticatorAAGUID:           credential.Authenticator.AAGUID,
		AuthenticatorSignCount:        credential.Authenticator.SignCount,
		AuthenticatorCloneWarning:     credential.Authenticator.CloneWarning,
		AuthenticatorAttachment:       string(credential.Authenticator.Attachment),
		AttestationClientDataJSON:     credential.Attestation.ClientDataJSON,
		AttestationClientDataHash:     credential.Attestation.ClientDataHash,
		AttestationAuthenticatorData:  credential.Attestation.AuthenticatorData,
		AttestationPublicKeyAlgorithm: credential.Attestation.PublicKeyAlgorithm,
		AttestationObject:             credential.Attestation.Object,
	}
}
