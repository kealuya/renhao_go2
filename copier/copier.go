package copier

import (
	"fmt"
	"github.com/jinzhu/copier"
)

type Address struct {
	City    string
	Country string
}

type Contact struct {
	Email  string
	Phones []string
}

type Employee struct {
	Name      string
	Age       int32
	Addresses []Address
	Contact   *Contact
}

type Manager struct {
	Name            string `copier:"must"`
	Age             int32  `copier:"must,nopanic"`
	ManagedCities   []string
	Contact         *Contact `copier:"override"`
	SecondaryEmails []string
}

func Copier() {
	employee := Employee{
		Name: "John Doe",
		Age:  30,
		Addresses: []Address{
			{City: "New York", Country: "USA"},
			{City: "San Francisco", Country: "USA"},
		},
		Contact: nil,
	}

	manager := Manager{
		ManagedCities: []string{"Los Angeles", "Boston"},
		Contact: &Contact{
			Email:  "john.doe@example.com",
			Phones: []string{"123-456-7890", "098-765-4321"},
		}, // since override is set this should be overridden with nil
		SecondaryEmails: []string{"secondary@example.com"},
	}

	err := copier.CopyWithOption(&manager, &employee, copier.Option{IgnoreEmpty: true, DeepCopy: true})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Manager: %#v\n", manager)
	// Output: Manager struct showcasing copied fields from Employee,
	// including overridden and deeply copied nested slices.
}
