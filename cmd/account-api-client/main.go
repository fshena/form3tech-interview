package main

import (
	"fmt"
	"form3/internal/app/form3"
	"log"
	"net/http"
)

func main() {
	c := form3.NewClient(&http.Client{}, "http://localhost:8080")

	a := form3.Account{
		ID:             "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
		OrganisationID: "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
		Type:           "accounts",
		Attributes: form3.Attributes{
			Country:      "GB",
			BaseCurrency: "GBP",
			BankID:       "400300",
			BankIDCode:   "GBDSC",
			BIC:          "NWBKGB22",
		},
	}

	newAccount, err := c.Create(a)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v", newAccount)
}
