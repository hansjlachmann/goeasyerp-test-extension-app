# Notes for the base application's maintainer

Findings made while building this app that belong in the base repository (not worked around here).

## A disabled app's table relation check still runs

Reported to the base. Not a blocker for 1.0.0: the app works as described while it is disabled,
the only difference is that the relation check of the Customer field is not skipped.

`EXTENSIONS.md` §11 says that while an app is disabled "its business logic, event subscribers and
field checks do not run". Field `tea_customer_category_code` (Customer) has a `table_relation`; the
generated `Customer.Validate` calls `tables.CheckRelation` without consulting the disabled-app
switch, so with the app disabled a value that is not a Customer Category is still refused
(`ERR_RELATED_NOT_FOUND`). The `validation: {min: …}` check of the test app's fields is skipped when
disabled, the relation check is not.

- **Where:** the `// TableRelation:` block of every generated `Validate` for an extension field
  (`backend/generated/tables/*_gen.go`, from the table generator), and `tables.CheckRelation`
  (`backend/foundation/tables/relations.go`).
- **Needed:** skip the relation check for a field of a disabled app, like the min/max check.
- **Test:** install an app that adds a related field, `extensions.SetDisabled(appID)`, `Validate`
  the field with a value that does not exist in the related table: it must be accepted.
