# Smoke tests for the Performance Platform

## Running the tests

Use the provided script:

```bash
./run_tests.sh
```

### In development:

1. Deploy the branch you want to test to preview.
2. Used the provided script:

```bash
SIGNON_USERNAME=See manual for signon credentials SIGNON_PASSWORD=See manual for signon credentials ./run_tests.sh
```


## Environment

By default these tests run against the Performance Platform's preview environment.

You can override this by specifying the following environment variables:

- `PP_APP_DOMAIN`, the public hostname
- `PP_FULL_APP_DOMAIN`, the internal hostname
- `GOVUK_APP_DOMAIN`, the GOV.UK hostname

## Integration with GOV.UK Signon

These tests make use of [GOV.UK Signon][signon] to work through the OAuth flow
to log in.

Environment variables (`SIGNON_USERNAME` and `SIGNON_PASSWORD`) are passed in when
these tests are run to allow authentication using the Signon interface.

[signon]: https://github.com/alphagov/signonotron2/
