package ssh

import (
	"fmt"
	"regexp"
)

var (
	pubKeyRegex = regexp.MustCompile("^(ssh-...|ecdsa-[^ ]+) ([^ ]+) ?(.*)")
)

// PubKeyInfo contains the information on an SSH public key
type PubKeyInfo struct {
	ID     string
	Public string
}

// ErrInvalidPubKey is the error returned when an SSH public key is unrecognizable
type ErrInvalidPubKey struct {
	pubKey []byte
}

// Error is the error interface implementation
func (e ErrInvalidPubKey) Error() string {
	return fmt.Sprintf("invalid SSH public key %s", string(e.pubKey))
}

// ErrUnknownPubKeyID is the error returned when an SSH public key is not the expected format
type ErrUnknownPubKeyID struct {
	pubKey []byte
}

func (e ErrUnknownPubKeyID) Error() string {
	return fmt.Sprintf("unknown SSH public key ID for %s", string(e.pubKey))
}

// ParsePubKey parses a byte slice representation of an SSH Public Key into an
// SSHPubKeyInfo struct. If it cannot find the key ID from the pubKey byte slice itself,
// it uses backupKeyID instead. Returns an appropriate error if parsing failed.
func ParsePubKey(backupKeyID string, pubKey []byte) (*PubKeyInfo, error) {
	if !pubKeyRegex.Match(pubKey) {
		return nil, ErrInvalidPubKey{pubKey: pubKey}
	}
	capture := pubKeyRegex.FindStringSubmatch(string(pubKey))
	if len(capture) < 4 || capture[3] == "" {
		return &PubKeyInfo{ID: backupKeyID, Public: string(pubKey)}, nil
	}
	return &PubKeyInfo{ID: capture[3], Public: string(pubKey)}, nil
}
