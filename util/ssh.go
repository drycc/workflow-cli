package util

import (
	"fmt"
	"regexp"
)

var (
	sshPubKeyRegex = regexp.MustCompile("^(ssh-...|ecdsa-[^ ]+) ([^ ]+) ?(.*)")
)

// SSHPubKeyInfo contains the information on an SSH public key
type SSHPubKeyInfo struct {
	ID     string
	Public string
}

// ErrInvalidSSHPubKey is the error returned when an SSH public key is unrecognizable
type ErrInvalidSSHPubKey struct {
	pubKey []byte
}

// Error is the error interface implementation
func (e ErrInvalidSSHPubKey) Error() string {
	return fmt.Sprintf("invalid SSH public key %s", string(e.pubKey))
}

// ErrUnknownSSHPubKeyID is the error returned when an SSH public key is not the expected format
type ErrUnknownSSHPubKeyID struct {
	pubKey []byte
}

func (e ErrUnknownSSHPubKeyID) Error() string {
	return fmt.Sprintf("unknown SSH public key ID for %s", string(e.pubKey))
}

// ParseSSHPubKey parses a byte slice representation of an SSH Public Key into an
// SSHPubKeyInfo struct. If it cannot find the key ID from the pubKey byte slice itself,
// it uses backupKeyID instead. Returns an appropriate error if parsing failed.
func ParseSSHPubKey(backupKeyID string, pubKey []byte) (*SSHPubKeyInfo, error) {
	if !sshPubKeyRegex.Match(pubKey) {
		return nil, ErrInvalidSSHPubKey{pubKey: pubKey}
	}
	capture := sshPubKeyRegex.FindStringSubmatch(string(pubKey))
	if len(capture) < 4 || capture[3] == "" {
		return &SSHPubKeyInfo{ID: backupKeyID, Public: string(pubKey)}, nil
	}
	return &SSHPubKeyInfo{ID: capture[3], Public: string(pubKey)}, nil
}
