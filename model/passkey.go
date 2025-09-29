package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"one-api/common"
	"strings"
	"time"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"gorm.io/gorm"
)

var (
	ErrPasskeyNotFound         = errors.New("passkey credential not found")
	ErrFriendlyPasskeyNotFound = errors.New("Passkey 验证失败，请重试或联系管理员")
)

type PasskeyCredential struct {
	ID              int            `json:"id" gorm:"primaryKey"`
	UserID          int            `json:"user_id" gorm:"uniqueIndex;not null"`
	CredentialID    []byte         `json:"credential_id" gorm:"type:blob;uniqueIndex;not null"`
	PublicKey       []byte         `json:"public_key" gorm:"type:blob;not null"`
	AttestationType string         `json:"attestation_type" gorm:"type:varchar(255)"`
	AAGUID          []byte         `json:"aaguid" gorm:"type:blob"`
	SignCount       uint32         `json:"sign_count" gorm:"default:0"`
	CloneWarning    bool           `json:"clone_warning"`
	UserPresent     bool           `json:"user_present"`
	UserVerified    bool           `json:"user_verified"`
	BackupEligible  bool           `json:"backup_eligible"`
	BackupState     bool           `json:"backup_state"`
	Transports      string         `json:"transports" gorm:"type:text"`
	Attachment      string         `json:"attachment" gorm:"type:varchar(32)"`
	LastUsedAt      *time.Time     `json:"last_used_at"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

func (p *PasskeyCredential) TransportList() []protocol.AuthenticatorTransport {
	if p == nil || strings.TrimSpace(p.Transports) == "" {
		return nil
	}
	var transports []string
	if err := json.Unmarshal([]byte(p.Transports), &transports); err != nil {
		return nil
	}
	result := make([]protocol.AuthenticatorTransport, 0, len(transports))
	for _, transport := range transports {
		result = append(result, protocol.AuthenticatorTransport(transport))
	}
	return result
}

func (p *PasskeyCredential) SetTransports(list []protocol.AuthenticatorTransport) {
	if len(list) == 0 {
		p.Transports = ""
		return
	}
	stringList := make([]string, len(list))
	for i, transport := range list {
		stringList[i] = string(transport)
	}
	encoded, err := json.Marshal(stringList)
	if err != nil {
		return
	}
	p.Transports = string(encoded)
}

func (p *PasskeyCredential) ToWebAuthnCredential() webauthn.Credential {
	flags := webauthn.CredentialFlags{
		UserPresent:    p.UserPresent,
		UserVerified:   p.UserVerified,
		BackupEligible: p.BackupEligible,
		BackupState:    p.BackupState,
	}

	return webauthn.Credential{
		ID:              p.CredentialID,
		PublicKey:       p.PublicKey,
		AttestationType: p.AttestationType,
		Transport:       p.TransportList(),
		Flags:           flags,
		Authenticator: webauthn.Authenticator{
			AAGUID:       p.AAGUID,
			SignCount:    p.SignCount,
			CloneWarning: p.CloneWarning,
			Attachment:   protocol.AuthenticatorAttachment(p.Attachment),
		},
	}
}

func NewPasskeyCredentialFromWebAuthn(userID int, credential *webauthn.Credential) *PasskeyCredential {
	if credential == nil {
		return nil
	}
	passkey := &PasskeyCredential{
		UserID:          userID,
		CredentialID:    credential.ID,
		PublicKey:       credential.PublicKey,
		AttestationType: credential.AttestationType,
		AAGUID:          credential.Authenticator.AAGUID,
		SignCount:       credential.Authenticator.SignCount,
		CloneWarning:    credential.Authenticator.CloneWarning,
		UserPresent:     credential.Flags.UserPresent,
		UserVerified:    credential.Flags.UserVerified,
		BackupEligible:  credential.Flags.BackupEligible,
		BackupState:     credential.Flags.BackupState,
		Attachment:      string(credential.Authenticator.Attachment),
	}
	passkey.SetTransports(credential.Transport)
	return passkey
}

func (p *PasskeyCredential) ApplyValidatedCredential(credential *webauthn.Credential) {
	if credential == nil || p == nil {
		return
	}
	p.CredentialID = credential.ID
	p.PublicKey = credential.PublicKey
	p.AttestationType = credential.AttestationType
	p.AAGUID = credential.Authenticator.AAGUID
	p.SignCount = credential.Authenticator.SignCount
	p.CloneWarning = credential.Authenticator.CloneWarning
	p.UserPresent = credential.Flags.UserPresent
	p.UserVerified = credential.Flags.UserVerified
	p.BackupEligible = credential.Flags.BackupEligible
	p.BackupState = credential.Flags.BackupState
	p.Attachment = string(credential.Authenticator.Attachment)
	p.SetTransports(credential.Transport)
}

func GetPasskeyByUserID(userID int) (*PasskeyCredential, error) {
	if userID == 0 {
		common.SysLog("GetPasskeyByUserID: empty user ID")
		return nil, ErrFriendlyPasskeyNotFound
	}
	var credential PasskeyCredential
	if err := DB.Where("user_id = ?", userID).First(&credential).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			common.SysLog(fmt.Sprintf("GetPasskeyByUserID: passkey not found for user %d", userID))
			return nil, ErrFriendlyPasskeyNotFound
		}
		common.SysLog(fmt.Sprintf("GetPasskeyByUserID: database error for user %d: %v", userID, err))
		return nil, ErrFriendlyPasskeyNotFound
	}
	return &credential, nil
}

func GetPasskeyByCredentialID(credentialID []byte) (*PasskeyCredential, error) {
	if len(credentialID) == 0 {
		common.SysLog("GetPasskeyByCredentialID: empty credential ID")
		return nil, ErrFriendlyPasskeyNotFound
	}

	var credential PasskeyCredential
	if err := DB.Where("credential_id = ?", credentialID).First(&credential).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			common.SysLog(fmt.Sprintf("GetPasskeyByCredentialID: passkey not found for credential ID length %d", len(credentialID)))
			return nil, ErrFriendlyPasskeyNotFound
		}
		common.SysLog(fmt.Sprintf("GetPasskeyByCredentialID: database error for credential ID: %v", err))
		return nil, ErrFriendlyPasskeyNotFound
	}

	return &credential, nil
}

func UpsertPasskeyCredential(credential *PasskeyCredential) error {
	if credential == nil {
		common.SysLog("UpsertPasskeyCredential: nil credential provided")
		return fmt.Errorf("Passkey 保存失败，请重试")
	}
	return DB.Transaction(func(tx *gorm.DB) error {
		// 使用Unscoped()进行硬删除，避免唯一索引冲突
		if err := tx.Unscoped().Where("user_id = ?", credential.UserID).Delete(&PasskeyCredential{}).Error; err != nil {
			common.SysLog(fmt.Sprintf("UpsertPasskeyCredential: failed to delete existing credential for user %d: %v", credential.UserID, err))
			return fmt.Errorf("Passkey 保存失败，请重试")
		}
		if err := tx.Create(credential).Error; err != nil {
			common.SysLog(fmt.Sprintf("UpsertPasskeyCredential: failed to create credential for user %d: %v", credential.UserID, err))
			return fmt.Errorf("Passkey 保存失败，请重试")
		}
		return nil
	})
}

func DeletePasskeyByUserID(userID int) error {
	if userID == 0 {
		common.SysLog("DeletePasskeyByUserID: empty user ID")
		return fmt.Errorf("删除失败，请重试")
	}
	// 使用Unscoped()进行硬删除，避免唯一索引冲突
	if err := DB.Unscoped().Where("user_id = ?", userID).Delete(&PasskeyCredential{}).Error; err != nil {
		common.SysLog(fmt.Sprintf("DeletePasskeyByUserID: failed to delete passkey for user %d: %v", userID, err))
		return fmt.Errorf("删除失败，请重试")
	}
	return nil
}
