package form3

type Attributes struct {
	Country                 string   `json:"country,omitempty"`
	BaseCurrency            string   `json:"base_currency,omitempty"`
	BankID                  string   `json:"bank_id,omitempty"`
	BankIDCode              string   `json:"bank_id_code,omitempty"`
	AccountNumber           string   `json:"account_number,omitempty"`
	BIC                     string   `json:"bic,omitempty"`
	IBAN                    string   `json:"iban,omitempty"`
	CustomerID              string   `json:"customer_id,omitempty"`
	Name                    []string `json:"name,omitempty"`
	AlternativeNames        []string `json:"alternative_name,omitempty"`
	AccountClassification   string   `json:"account_classification,omitempty"`
	JoinAccount             bool     `json:"join_account,omitempty"`
	AccountMatchingOptOut   bool     `json:"account_matching_opt_out,omitempty"`
	SecondaryIdentification string   `json:"secondary_identification,omitempty"`
	Switched                bool     `json:"switched,omitempty"`
	Status                  string   `json:"status,omitempty"`
}

type Account struct {
	ID             string     `json:"id"`
	OrganisationID string     `json:"organisation_id"`
	Type           string     `json:"type"`
	Version        int        `json:"version"`
	CreatedOn      string     `json:"created_on"`
	ModifiedOn     string     `json:"modified_on"`
	Attributes     Attributes `json:"attributes"`
}
