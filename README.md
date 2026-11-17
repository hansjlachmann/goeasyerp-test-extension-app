# Test Extension App

An extension app for [GoEasyERP](https://github.com/hansjlachmann/goeasyerp). It is meant as an
example other app authors can copy; the guide to extension apps is `EXTENSIONS.md` in the base
repository.

## What it adds

- A table **Customer Category** (50100) with Code and Description, and a list page
  **Customer Categories** (50100) to maintain them, reached from the main menu under
  Sales, after Customers.
- A field **Customer Category Code** on Customer (related to Customer Category), shown in the
  General section of the Customer Card and as a column in the Customer List.
- Texts in en-US, nb-NO and da-DK.

There is no business logic yet: `Register()` in `extensions/app.go` is empty.

## Layout

The app is the `extensions/` folder: `application.yaml` (identity, prefix `tea`, numbers
50100–50199, dependency on the base), `custom/` (own tables and pages), `tables/`, `pages/`
and `menus/` (additions to the base's), `translations/`, and the Go module (`go.mod`, `app.go`).

## Installing

In the installation's repository (see `deploy/installation-template/README.md` in the base):

1. In GoEasyERP, **Administration → Installed Apps → Install App**: this repository's URL, the
   release tag (e.g. `v1.0.0`) and the folder `extensions`.
2. **Build List**, commit the downloaded `apps.yaml` to the installation repository. Its
   workflow builds the image.
3. On the server, pull the image and restart.

After installing, create the categories on **Customer Categories** and set them on the customers.
The app can be turned off and on under **Installed Apps** (Disable App / Enable App); its data
stays.

## Developing

CI (`.github/workflows/ci.yml`) checks out the base release in `BASE_VERSION`, merges the app
into it, builds, runs the base's page, caption and installed-app checks, the app's tests, `go vet`
and `gofmt`. To do the same locally, follow `scripts/check-extension-fixture.sh` in the base
repository with this app's `extensions/` folder.

Releases are tags `vX.Y.Z` matching `version` in `extensions/application.yaml`.
Commits follow Conventional Commits.
