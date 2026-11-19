# F3: Products need a status

**Area:** SQL. **Type:** build. **Points:** 10. Core.

## The request

From the product team:

> Some products are no longer sold. We want to mark a product as retired instead of deleting it.

## What to build

Add a `status` to products, by a database migration:

- A product's status is `active` or `retired`. Nothing else can be stored, and it is never empty.
- Every product that already exists becomes `active`.

## What makes it interesting

We release without downtime. While your migration is applied, and for some time after, **the previous release
of the API is still running**. It still adds products, and it knows nothing about a status. It must keep
working, and the products it adds are `active`.

## Done means

- `go run ./cmd/devshelf migrate` applies your change to a database that already holds data.
- Running it a second time changes nothing and reports no error.

The API does not need to return or accept the status yet. That is a later ticket.
