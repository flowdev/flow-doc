package functyps

import (
	"fmt"
)

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

//flowdev:flow
func simpleFunc() {
	msg := "vim-go"
	fmt.Println(msg)
}

/*
//flowdev:flow
func (ba *BankAccount) doAccountingMagic(newHolder string, newType string) (iban, bic string) {
	ba.HolderName = newHolder
	ba.AccountType = newType

	return ba.IBAN, ba.BIC
}

//flowdev:flow
func (sba *SpecialBankAccount) doSpecialAccountingMagic(newHolder string, newTitle string,
) (portOut1 string, portSpecialOut *SpecialBankAccount, err error) {

	sba.doAccountingMagic(newHolder, "special")
	sba.HolderTitle = newTitle
	return sba.IBAN, sba, nil
}

//flowdev:flow
func funcWithEllipsis_in2(i, j int, names ...string) ([]string, error) {
	var addNames []string
	for k := i; k < j; k++ {
		addNames = append(addNames, names[k])
	}

	return append(names, addNames...), nil
}

//flowdev:flow
func funcWithErrorOnly() error {
	return nil
}

//flowdev:flow
func funcWithoutResults(
	d *tool.Data,
	m map[string]*SpecialBankAccount,
	s string,
) {

	print(s)
}
*/
