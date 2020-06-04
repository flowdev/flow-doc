package find_flows

import "fmt"

// BankAccount represents a bank account.
type BankAccount struct {
	HolderName  string
	IBAN        string
	BIC         string
	AccountType string
}

// SpecialBankAccount represents a special bank account.
type SpecialBankAccount struct {
	BankAccount
	HolderTitle string
}

// SimpleFunc does bla.
//flowdev:flow
func SimpleFunc() {
	fmt.Println("vim-go")
}

// DoAccountingMagic func does bla.
//
//flowdev:flow
func (ba *BankAccount) DoAccountingMagic(newHolder string, newType string) (iban, bic string) {
	ba.HolderName = newHolder
	ba.AccountType = newType

	return ba.IBAN, ba.BIC
}

//flowdev:flow
// DoSpecialAccountingMagic func does bla.
func (sba *SpecialBankAccount) DoSpecialAccountingMagic(newHolder string, newTitle string) (iban string) {
	sba.DoAccountingMagic(newHolder, "special")
	sba.HolderTitle = newTitle
	return sba.IBAN
}

//flowdev:flow
func funcWithEllipsis(i, j int, names ...string) []string {
	var addNames []string
	for k := i; k < j; k++ {
		addNames = append(addNames, names[k])
	}

	return append(names, addNames...)
}

// funcWithoutAFlow returns foo.
func funcWithoutAFlow() string {
	return "foo"
}

func funcWithoutComment() string {
	return "bar"
}
