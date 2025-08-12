package types

type User struct {
	AccountName string
	PersonaName string
	Steam64     string // Steam64 ID
	AutoLogin   bool   // Whether user is likely to have a valid Connect Cache JWT saved
}
