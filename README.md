# beam-go

Beam-go interfaces with the `wallet-api` binary of the Beam ecosystem.

The goal is to provide robust request-response structures, using sensible defaults for optionals, and providing an architecture to resist future breaking changes.

## Usage

Non-idiomatic to Go, `beam-go` makes the use of pointers to indicate optionality on function inputs, and options bag inputs.

the `to` package, and all enum types feature `*ToPtr()` type functions for constant and literal values.

## Contribute

Please read the [Beam Contribution Guide](https://github.com/BeamMW/beam/wiki/Contribution-Guidelines).

### Refreshing test recordings

Beam-go records its tests using [`go-vcr`](https://github.com/dnaeon/go-vcr). 

Please, when submitting a PR that modifies the functionality of existing APIs, delete the relevant file in the `recordings` folder in the relevant package, and re-run tests yourself before submitting the PR.

Adding new tests does not require this. It is recommended to run `beam-node` with `--FakePoW 1 --storage testnode.db` and you _must_ have your `wallet-api` available on port 5000 to line up with the existing configuration.
