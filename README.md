# Test Extension App

An extension app for [GoEasyERP](https://github.com/hansjlachmann/goeasyerp). It is meant as an
example other app authors can copy; the guide to extension apps is `EXTENSIONS.md` in the base
repository.

**Needs:** GoEasyERP base application **0.1.161** or later (the release with Install App,
Disable App and Enable App).

## What it adds

- A table **Customer Category** (50100) with Code and Description, and a list page
  **Customer Categories** (50100) to maintain them, reached from the main menu under
  Sales, after Customers.
- A field **Customer Category Code** on Customer (related to Customer Category), shown in the
  General section of the Customer Card and as a column in the Customer List.
- Texts in en-US, nb-NO and da-DK.

There is no business logic: `Register()` in `extensions/app.go` is empty.

## Layout

The app is the `extensions/` folder: `application.yaml` (identity, prefix `tea`, numbers
50100–50199, dependency on the base), `custom/` (own tables and pages), `tables/`, `pages/`
and `menus/` (additions to the base's), `translations/`, and the Go module (`go.mod`, `app.go`,
`app_test.go`).

## Installing

1. In GoEasyERP, **Administration → Installed Apps → Install App** (SUPER only), with
   - repository: `https://github.com/hansjlachmann/goeasyerp-test-extension-app`
   - version: `v1.0.0`
   - folder: `extensions`

   The app shows as **Pending Install**.
2. **Build List** downloads `apps.yaml`. Commit it to the installation repository (see
   `deploy/installation-template/README.md` in the base); its workflow builds the backend image
   with the app compiled in.
3. Deploy the image: on the server `scripts/prod.sh pull && scripts/prod.sh up -d`. On startup the
   app shows as **Installed**.

After installing, create the categories on **Customer Categories** and set them on the customers.

## Disabling, enabling and uninstalling

On **Installed Apps** (SUPER only):

- **Disable App** switches the app off without removing it: the page parts and menu entry it added
  are left out, and its data stays. **Enable App** brings it back (servers pick the switch up
  within 30 seconds).
- **Uninstall App** marks the app **Pending Removal**. Then **Build List**, commit the new
  `apps.yaml` and deploy the image as when installing. The app is gone from the program; its data
  stays in the database. **Cancel Pending** takes back a pending install, update or removal.

## Developing

CI (`.github/workflows/ci.yml`) checks out the base release in `BASE_VERSION`, merges the app
into it, builds, runs the base's page, caption and installed-app checks, the app's tests, `go vet`
and `gofmt`. To do the same locally, follow `scripts/check-extension-fixture.sh` in the base
repository with this app's `extensions/` folder.

Releases are tags `vX.Y.Z` matching `version` in `extensions/application.yaml`; see
`CHANGELOG.md`. Commits follow Conventional Commits.
