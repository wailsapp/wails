# Spec review reproductions

Reviewed commit: 3452def9d.

Source: `spec-regression_test.go` is an isolated commands-package regression test. Copy it to `v3/internal/commands/zz_spec_ga_review_test.go`, run the following from `v3`, then remove the copied test:

```sh
go test ./internal/commands -run '^TestGAReviewSpec' -count=1 -v
```

The source is saved outside the production test tree so the review leaves tracked code unchanged. No Android SDK, Xcode, simulator or native application execution was used.

Expected contract assertions fail on this commit:

- Actual manifest loading, planning and asset generation produces Android Gradle `versionCode = 1`, `versionName = "1.0.0"`, `minSdk = 21`, `targetSdk = 36` despite HCL values 42, 2.4.1, 28, 35.
- Generated production Gradle still selects debug signing when the ambient legacy keystore environment is absent.
- Actual iOS generated Info.plist lacks the configured audio and remote-notification background modes.
- Actual development macOS app assembly selects `bin/app.app` and removes a sentinel file in the pre-existing production bundle.

`spec-all-repro.log` records the combined run; `spec-mobile-repro.log` and `spec-dev-repro.log` preserve the initial separate runs.

Platform evidence was checked against the September 4 macOS fix commit 8fe525f6d. That commit adds implementation fixes and automated tests for CGo, DMG packaging/signing, bundle executable metadata and dev readiness. It does not update platform-status.md or add native matrix evidence. The status document still identifies acdd11858 as its reviewed commit, so its uncompleted native/signing rows cannot be treated as accepted at HEAD.
