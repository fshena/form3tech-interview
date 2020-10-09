# Form3 Take Home Exercise

## Candidate
Florian Shena

## Requirements

In order to run the project the following requirements must be met:

- [Docker Engine](https://docs.docker.com/engine/installation/)
- [Docker Compose](https://docs.docker.com/compose/)
- [make tool](http://www.gnu.org/software/make/) (optional)

If _**make**_ tool is not available, the commands inside the "Makefile" can be copied and run directly in the terminal.

## Run
In order to run the application type the command 
```bash
make serve
```

Once the Postgresql and Account API containers have been created, the tests for the Client will run on
a separate containers and, that container, will exit once the tests have finished.

## Usage 

Create a client, then call commands on it.

```go
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
```
The example above can be run locally by serving the application with `make serve` and `make example`, 
which will create a binary in the root of the project called `account-api-client`. Afterwards run the aforementioned 
binary, and the newly created Account will be logged.

## Commands
A Makefile is included containing some basic commands for convenience:
- `make serve` runs docker compose "up" command in order to build the images and create the containers.
- `make stop` stops all running containers.
- `make test` runs all phpunit test.
- `make tests-all` runs all phpunit test and integration tests.
- `make ci` runs tests in a separate container.
- `make example` runs the clients __Create__ functionality as an example

## Technicals
Moved the docker-compose file on a separate folder in order to abide to go's standard folder structure. Extracted
the http client as a dependency in order to make it possible to mock it for the unit tests. The integration tests
where placed in a separate file in order to be able to exclude them from the build. Unit tests where written using
Table-Driven approach for better test organisation.
