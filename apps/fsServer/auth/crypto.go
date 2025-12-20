package auth

func hashAuthString(data string) string {
	return ""
}

type AuthCryptoModuleType struct {
	Hash    func(string) string
	Decrypt func(string) string
}

var AuthCryptoModule = AuthCryptoModuleType{
	Hash: func(s string) string {
		return ""
	},
	Decrypt: func(s string) string {
		return ""
	},
}
