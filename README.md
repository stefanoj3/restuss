## Restuss
This is a Go package that support part of the Nessus rest(kind of) api.

```
go get github.com/stefanoj3/restuss
```

The amount of method supported is extremely small, I actually needed it for a personal project and decided to make a library out of it, so others can reuse it.

Example usage of the library:
```go
auth := restuss.NewKeyAuthProvider(
    "fa74fcdd10db53bf54cd1467c11547efd70af0b526eb0d2b347b1050e1cab639",
    "71ebdf108b4d2fa9fef8c895096ff1a2a7732c215c64edc8d6495053488004d6",
)

c, err := restuss.NewClient(auth, "https://127.0.0.1:8834", true)

if err != nil {
    log.Fatal(err.Error())
}

var lastModificationDate int64 = 0

res, err := c.GetScans(lastModificationDate)

if err != nil {
    log.Fatal(err.Error())
}
```

Support for basic auth is also planned but not a priority.

For now the available calls are: create scan, launch scan, stop scan, list scans, list scan's templates.

This package has zero dependencies and I plan to keep it like that.

Tested with Nessus 6.10.1

## Important note
Since it's in early development the next versions could contain breaking changes.