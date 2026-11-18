# Changelog

## 1.0.0

First release. Needs GoEasyERP base application 0.1.161 or later.

### Features

- Table **Customer Category** (50100) with Code and Description, and the list page
  **Customer Categories**, in the main menu under Sales after Customers.
- Field **Customer Category Code** on Customer, related to Customer Category, on the Customer Card
  (General section) and as a column in the Customer List.
- Texts in en-US, nb-NO and da-DK.

### Tests and CI

- Go tests for the category table, the Customer field and the disabled app.
- CI merges the app into the base release, builds it and runs the base's checks, the app's tests,
  `go vet` and `gofmt`.
