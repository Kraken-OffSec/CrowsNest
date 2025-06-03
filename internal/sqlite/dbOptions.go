package sqlite

type IString interface {
	String() string
}

type DBOptions struct {
	Username              string
	Email                 string
	IPAddress             string
	Password              string
	HashedPassword        string
	Name                  string
	Vin                   string
	LicensePlate          string
	Address               string
	Phone                 string
	Social                string
	CryptoCurrencyAddress string
	Domain                string
	Limit                 int
	ExactMatch            bool
	NonEmptyFields        []string // Fields that should not be empty
	DisplayFields         []string // Fields to display in output
}

func (o *DBOptions) Empty() bool {
	return o.Username == "" && o.Email == "" && o.IPAddress == "" &&
		o.Password == "" && o.HashedPassword == "" && o.Name == "" &&
		o.Vin == "" && o.LicensePlate == "" && o.Address == "" &&
		o.Phone == "" && o.Social == "" && o.CryptoCurrencyAddress == "" && o.Domain == "" &&
		len(o.NonEmptyFields) == 0
}
