package typ

import (
	"fmt"

	"github.com/flowdev/ea-flow-doc/data/testdata/typ/tool"
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
	fmt.Println("vim-go")
}

//flowdev:flow
func (ba *BankAccount) doAccountingMagic(newHolder string, newType string) (iban, bic string) {
	ba.HolderName = newHolder
	ba.AccountType = newType

	return ba.IBAN, ba.BIC
}

//flowdev:flow
func (sba *SpecialBankAccount) doSpecialAccountingMagic(newHolder string, newTitle string, ba *BankAccount,
) (portOut1 string, portSpecialOut *SpecialBankAccount, err error) {

	sba.doAccountingMagic(newHolder, "special")
	sba.HolderTitle = newTitle
	return sba.IBAN, sba, nil
}

//flowdev:flow
func funcWithEllipsis_in2(i, j int, names [16]string, moreNames []string, mostNames ...bool,
) (newNames []string, err error) {
	var addNames []string
	for k := i; k < j; k++ {
		addNames = append(addNames, names[k])
	}

	return append(addNames, moreNames...), nil
}

//flowdev:flow
func funcWithErrorOnly() error {
	return nil
}

//flowdev:flow
func funcWithoutResults(
	d *tool.Data,
	m map[string]*SpecialBankAccount,
	m2 []map[string]map[*int][]*tool.Data,
	s string,
) {

	print(s)
}

//flowdev:flow
func funcWithTooComplexData(s string, f func(int) int, c chan<- int, i int) {
	fmt.Println(s, i)
}
