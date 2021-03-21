# beam-go

Beam-go interfaces with the `wallet-api` and `explorer-node` binaries of the Beam ecosystem.

The goal is to provide robust request-response structures, using sensible defaults for optionals, and providing an architecture to resist future breaking changes.

## Usage

Non-idiomatic to Go, `beam-go` makes the use of pointers to indicate optionality on function inputs, and options bag inputs.

the `to` package, and all enum types feature `*ToPtr()` type functions for constant and literal values.

## Contribute

Please read the [Beam Contribution Guide](https://github.com/BeamMW/beam/wiki/Contribution-Guidelines).