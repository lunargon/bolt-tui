package app

type state int

const (
	stateBuckets state = iota
	stateCreateBucket
	stateCreateKey
	stateEditValue
	stateEditBucket
	stateEditKey
	stateConfirmDelete
	stateConfirmDeleteBucket
	stateSettings
)

// DisplayMode represents different ways to display data
type DisplayMode int

const (
	DisplayString DisplayMode = iota
	DisplayBase64
	DisplayBase58
	DisplayHex
)

func (d DisplayMode) String() string {
	switch d {
	case DisplayString:
		return "String"
	case DisplayBase64:
		return "Base64"
	case DisplayBase58:
		return "Base58"
	case DisplayHex:
		return "Hex"
	default:
		return "String"
	}
}
