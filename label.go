package historytree

import (
	"crypto"
	"crypto/sha512"
)

const (
	// LabelSize is the number of bytes in a label.
	LabelSize = sha512.Size256 // 32 bytes

	// LabelHash is the hash function used in this implementation.
	LabelHash = crypto.SHA512_256

	// confirm that MaxNodes labels would actually fit in a single file.
	_ = int64(LabelSize * MaxNodes) // 31 extra bytes available?
)

// Label represents the label of a node (a SHA-512/256 hash).
type Label [LabelSize]byte
