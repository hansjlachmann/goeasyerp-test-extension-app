package testextensionapp

import (
	"database/sql"
	"sync"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	"github.com/hansjlachmann/goeasyerp/backend/business-logic/tables"
	"github.com/hansjlachmann/goeasyerp/backend/foundation/database"
	"github.com/hansjlachmann/goeasyerp/backend/foundation/extensions"
	ftables "github.com/hansjlachmann/goeasyerp/backend/foundation/tables"
	"github.com/hansjlachmann/goeasyerp/backend/foundation/types"
	gtables "github.com/hansjlachmann/goeasyerp/backend/generated/tables"
)

// appID is the app's id (application.yaml)
const appID = "c7c1e38d-51b4-453d-97ad-32fd54456d4b"

// installOnce: an app can be installed once per process
var installOnce sync.Once

// newSource creates a company's tables in an in-memory database and installs the app
func newSource(t *testing.T) ftables.Source {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { _ = db.Close() })
	err = tables.CreateCompanyTables(db, "test", database.DBTypeSQLite)
	if err != nil {
		t.Fatal(err)
	}
	src := ftables.DBSource(db, "test", database.DBTypeSQLite)
	installOnce.Do(func() { extensions.Install(appID, Register) })
	t.Cleanup(func() { extensions.SetDisabled() })
	return src
}

func insertCategory(t *testing.T, src ftables.Source, code string, description string) {
	t.Helper()
	cat := gtables.NewTEACustomerCategory(src)
	cat.Code = types.Code(code)
	cat.Description = types.Text(description)
	if !cat.Insert(true) {
		t.Fatalf("insert category %s: %v", code, cat.TriggerError())
	}
}

// A category is stored with its code and description, and its code is the primary key
func TestCategoryTable(t *testing.T) {
	src := newSource(t)
	insertCategory(t, src, "RETAIL", "Retail customers")

	cat := gtables.NewTEACustomerCategory(src)
	if !cat.Get("RETAIL") {
		t.Fatal("the category was not saved")
	}
	if string(cat.Description) != "Retail customers" {
		t.Errorf("description %q", cat.Description)
	}

	dup := gtables.NewTEACustomerCategory(src)
	dup.Code = "RETAIL"
	if dup.Insert(true) {
		t.Error("a second category with the same code was accepted")
	}
}

// The Customer's category code is saved with the customer, and must be a category that exists
func TestCustomerCategoryField(t *testing.T) {
	src := newSource(t)
	insertCategory(t, src, "RETAIL", "Retail customers")

	cust := gtables.NewCustomer(src)
	cust.No = "C1"
	if !cust.Insert(false) {
		t.Fatal("insert customer")
	}

	err := cust.Validate("tea_customer_category_code", "RETAIL")
	if err != nil {
		t.Fatalf("an existing category was refused: %v", err)
	}
	if !cust.Modify(true) {
		t.Fatalf("modify: %v", cust.TriggerError())
	}
	again := gtables.NewCustomer(src)
	if !again.Get("C1") {
		t.Fatal("get customer")
	}
	if string(again.Tea_customer_category_code) != "RETAIL" {
		t.Errorf("category code %q, want RETAIL", again.Tea_customer_category_code)
	}

	err = cust.Validate("tea_customer_category_code", "MISSING")
	if err == nil {
		t.Error("a category that does not exist was accepted")
	}
}

// Disabled (Installed Apps), the app's data stays and the customer can still be saved.
// (The category relation check still runs while the app is disabled: see NOTES.md.)
func TestDisabledAppKeepsData(t *testing.T) {
	src := newSource(t)
	insertCategory(t, src, "RETAIL", "Retail customers")

	cust := gtables.NewCustomer(src)
	cust.No = "C1"
	if !cust.Insert(false) {
		t.Fatal("insert customer")
	}
	err := cust.Validate("tea_customer_category_code", "RETAIL")
	if err != nil {
		t.Fatal(err)
	}
	if !cust.Modify(true) {
		t.Fatalf("modify: %v", cust.TriggerError())
	}

	extensions.SetDisabled(appID)
	cust.Name = "Changed while disabled"
	if !cust.Modify(true) {
		t.Fatalf("modify with the app disabled: %v", cust.TriggerError())
	}
	again := gtables.NewCustomer(src)
	if !again.Get("C1") {
		t.Fatal("get customer")
	}
	if string(again.Tea_customer_category_code) != "RETAIL" {
		t.Errorf("the app's data is gone: %q", again.Tea_customer_category_code)
	}
}
