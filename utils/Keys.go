package utils

type Keys struct {
	Name     string `json:"name"`
	KeyType  string `json:"type"`
	Address  string `json:"address"`
	PubKey   string `json:"pubkey"`
	Mnemonic string `json:"mnemonic, omitempty"`
}
