# Digital signature

Two library was used.

## unipdf (unidoc)

- It requires Unidoc's license, sign up on <https://cloud.unidoc.io>.
- The source code is encrypted, but can do a trick to bypass.
- How to:

```shell
export UNIDOC_LICENSE_API_KEY=<fake>
go run ./unidoc/pdf_sign_generate_key.go input.pdf ./unidoc/output.pdf
```

- The signature is inserted on the left.

## pdfsign

- Free to use.
- How to:

```shell
# sign
go run ./pdfsign/pdf_sign_generate_key.go sign input.pdf ./pdfsign/output.pdf
# verify
go run ./pdfsign/pdf_sign_generate_key.go verify ./pdfsign/output.pdf
```

- The signature isn't visible, but verify still works.
