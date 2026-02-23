# Go Notes

## Context

In Go, `context.Context` is used to carry deadlines, cancellation signals, and request-scoped data across API boundaries. Think of it as metadata that travels with your operations.

At startup, you're not processing any user request, so there's no existing context. You create a "root" context with `context.Background()`.

## Data Types

### Struct

`struct` is a keyword that defines a composite data type (a collection of fields grouped together).

## Text Processing

- `unicode.IsLetter(ch) || unicode.IsDigit(ch)` - checks if the character is a letter or a digit
- `unicode.ToLower(ch)` - converts the character to lowercase
