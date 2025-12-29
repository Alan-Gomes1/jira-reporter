package model

// User representa os dados do usuário/empresa para o relatório.
type User struct {
	CompanyName string
	CNPJ        string
	Username    string
}

// NewUser cria uma nova instância de User.
func NewUser(companyName, cnpj, username string) *User {
	return &User{
		CompanyName: companyName,
		CNPJ:        cnpj,
		Username:    username,
	}
}
