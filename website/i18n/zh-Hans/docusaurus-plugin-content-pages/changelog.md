# 更新日志

此项目的所有重要变更都将记录在此文件中。

格式基于 [维护更新日志](https://keepachangelog.com/en/1.0.0/)， 并且该项目遵循 [语义化版本](https://semver.org/spec/v2.0.0.html)。

## [Unreleased][]

## [v2.0.0-beta.43] - 2022-08-08

### Added

* NSIS Installer support for creating installers for Windows applications - Thanks [@stffabi](https://github.com/stffabi) 🎉
* New frontend:dev:watcher command to spin out 3rd party watchers when using wails dev - Thanks [@stffabi](https://github.com/stffabi)🎉

### Fixed

* Fix formatting of some error messages by @stffabi in https://github.com/wailsapp/wails/pull/1665
* Windows dialogs now work when window is not visible yet by @leaanthony in https://github.com/wailsapp/wails/pull/1662
* Multiple fixes for MacOS asset requests by @stffabi in https://github.com/wailsapp/wails/pull/1668 and https://github.com/wailsapp/wails/pull/1681
* Fix for Go 1.19 by @stffabi in https://github.com/wailsapp/wails/pull/1689
* Removed Linux warnings by @leaanthony in https://github.com/wailsapp/wails/pull/1656
* Better support for doubleclick events in drag regions by @leaanthony in https://github.com/wailsapp/wails/pull/1704
* Allow MacOS frameless window to be miniturisable by @leaanthony in https://github.com/wailsapp/wails/pull/1705

### What's Changed

* add wails-sveltekit-template by @h8gi in https://github.com/wailsapp/wails/pull/1671
* wails doctor now reports correct MacOS os id by @stffabi in https://github.com/wailsapp/wails/pull/1673
* Update application-development.mdx by @SamHennessy in https://github.com/wailsapp/wails/pull/1682
* Move SetMin/Max calls to main thread by @leaanthony in https://github.com/wailsapp/wails/pull/1684
* Change `frontend:dev` to `frontend:dev:build` by @LGiki in https://github.com/wailsapp/wails/pull/1691
* Build frontend only before starting the dev watcher command by @stffabi in https://github.com/wailsapp/wails/pull/1694
* Improve error message for auto dev server discovery without a dev watcher by @stffabi in https://github.com/wailsapp/wails/pull/1711

### New Contributors

* @h8gi made their first contribution in https://github.com/wailsapp/wails/pull/1671
* @SamHennessy made their first contribution in https://github.com/wailsapp/wails/pull/1682

## [v2.0.0-beta.42] - 2022-07-25

### Added

* Added `options.NewRGBA` and `options.NewRGB` functions to create `*options.RGBA` by @leaanthony

### Fixed

* Fixed initial build of frontend when using `wails dev` on new projects by @leaanthony in https://github.com/wailsapp/wails/pull/1650
* Ignore empty install command when running `wails dev` by @stffabi in https://github.com/wailsapp/wails/pull/1651
* Fixed error reporting in templates
* BackgroundColour documentation fix
* Generalize manual compile steps [Documentation] by @acheong08 in https://github.com/wailsapp/wails/pull/1644

## [v2.0.0-beta.40] - 2022-07-24

### Added

* Add Show() and Hide() to runtime to show/hide application by @leaanthony in https://github.com/wailsapp/wails/pull/1599
* Override target platform/arch using GOOS and GOARCH environment variables by @leaanthony in https://github.com/wailsapp/wails/pull/1618
* Add option to skip frontend rebuild in dev mode by @leaanthony in https://github.com/wailsapp/wails/pull/1632

### Fixed

* Update svelte templates to use vite 3 by @leaanthony in https://github.com/wailsapp/wails/pull/1643
* Fix plain template by @stffabi in https://github.com/wailsapp/wails/pull/1609
* Fix Website layout by @leaanthony in https://github.com/wailsapp/wails/pull/1616
* Fixed typo in documentation page docs/howdoesitwork by @MyNameIsAres in https://github.com/wailsapp/wails/pull/1636
* Use scaling when setting min/max window by @leaanthony in https://github.com/wailsapp/wails/pull/1557

### Changed

* Install dev dependencies before starting dev mode by @leaanthony in https://github.com/wailsapp/wails/pull/1615
* Translate and fix website text by @misitebao in https://github.com/wailsapp/wails/pull/1525
* docs: add MyNameIsAres as a contributor for doc by @allcontributors in https://github.com/wailsapp/wails/pull/1638
* Deprecate Fullscreen appoption by @acheong08 in https://github.com/wailsapp/wails/pull/1640

### Deprecated

* The `Fullscreen` application option is deprecated. The `Fullscreen` application option is deprecated. Please use [`WindowStartState`](https://wails.io/docs/reference/options#windowstartstate) instead.

### New Contributors

* @MyNameIsAres made their first contribution in https://github.com/wailsapp/wails/pull/1636

## [v2.0.0-beta.39.2] - 2022-07-20

### Added

* Update example for macOS Menu by @acheong08 in https://github.com/wailsapp/wails/pull/1600

### Fixed

* Reinstate Go 1.17 compatibility by @leaanthony in https://github.com/wailsapp/wails/pull/1605

## [v2.0.0-beta.39] - 2022-07-19

### Added

* New screen dimensions runtime API by @skamensky in https://github.com/wailsapp/wails/pull/1519
* Auto discover vite devserver port by @leaanthony in https://github.com/wailsapp/wails/pull/1547
* Add nixpkgs support to doctor command. by @ianmjones in https://github.com/wailsapp/wails/pull/1551 by @ianmjones in https://github.com/wailsapp/wails/pull/1551
* New pre-build hooks feature by @leaanthony in https://github.com/wailsapp/wails/pull/1578
* New production log level option by @leaanthony in https://github.com/wailsapp/wails/pull/1555

### Fixed

* Fix stack corruption in Windows when using ICoreWebView2HttpHeadersCollectionIterator by @stffabi in https://github.com/wailsapp/wails/pull/1589
* Move WindowGet* to main thread by @leaanthony in https://github.com/wailsapp/wails/pull/1464
* Allow -appargs flag to pass flags to binary. by @ianmjones in https://github.com/wailsapp/wails/pull/1534 by @ianmjones in https://github.com/wailsapp/wails/pull/1534
* Fix checking for installed apt package in none English session. by @ianmjones in https://github.com/wailsapp/wails/pull/1548 by @ianmjones in https://github.com/wailsapp/wails/pull/1548
* Fix OnBeforeClose code for Mac by @leaanthony in https://github.com/wailsapp/wails/pull/1558
* Support Maps in TS conversion by @leaanthony in https://github.com/wailsapp/wails/pull/1435
* Check for line length when scanning for local devserver url by @leaanthony in https://github.com/wailsapp/wails/pull/1566
* Remove usage of unsafe.Pointer in winc by @stffabi and @leaanthony in https://github.com/wailsapp/wails/pull/1556

### Changed

* Rename WindowSetRGBA -> WindowSetBackgroundColour by @leaanthony in https://github.com/wailsapp/wails/pull/1506
* Improvements to the dev command by @stffabi in https://github.com/wailsapp/wails/pull/1510
* Update vscode template by @leaanthony in https://github.com/wailsapp/wails/pull/1398
* Bump svelte from 3.42.2 to 3.49.0 in /v2/internal/frontend/runtime/dev by @dependabot in https://github.com/wailsapp/wails/pull/1572
* Bump svelte from 3.42.5 to 3.49.0 in /v2/internal/frontend/runtime by @dependabot in https://github.com/wailsapp/wails/pull/1573
* Add troubleshooting for `Not Found` error by @acheong08 in https://github.com/wailsapp/wails/pull/1586
* Docs/better homepage by @leaanthony in https://github.com/wailsapp/wails/pull/1591

### New Contributors

* @skamensky made their first contribution in https://github.com/wailsapp/wails/pull/1519
* @acheong08 made their first contribution in https://github.com/wailsapp/wails/pull/1586

**Full Changelog**: https://github.com/wailsapp/wails/compare/v2.0.0-beta.38...v2.0.0-beta.39

## [v2.0.0-beta.38] - 2022-06-27

### Added

* Add race detector to build & dev by @Lyimmi in https://github.com/wailsapp/wails/pull/1426
* [linux] Support `linux/arm` architecture by @Lyimmi in https://github.com/wailsapp/wails/pull/1427
* Create gitignore when using `-g` option by @jaesung9507 in https://github.com/wailsapp/wails/pull/1430
* [windows] Add Suspend/Resume callback support by @leaanthony in https://github.com/wailsapp/wails/pull/1474
* Add runtime function `WindowSetAlwaysOnTop` by @chenxiao1990 in https://github.com/wailsapp/wails/pull/1442
* [windows] Allow setting browser path by @NanoNik in https://github.com/wailsapp/wails/pull/1448

### Fixed

* [linux] Improve switching to main thread for callbacks by @stffabi in https://github.com/wailsapp/wails/pull/1392
* [windows] Fix WebView2 minimum runtime version check by @stffabi in https://github.com/wailsapp/wails/pull/1456
* [linux] Fix apt command syntax (#1458) by @abtin in https://github.com/wailsapp/wails/pull/1461
* [windows] Set Window Background colour if provided + debounce redraw option by @leaanthony in https://github.com/wailsapp/wails/pull/1466
* Fix small typo in docs by @LukenSkyne in https://github.com/wailsapp/wails/pull/1449
* Fix the url to surge by @andywenk in https://github.com/wailsapp/wails/pull/1460
* Fixed theme change at runtime by @leaanthony in https://github.com/wailsapp/wails/pull/1473
* Fix: Don't stop if unable to remove temporary bindings build by @leaanthony in https://github.com/wailsapp/wails/pull/1465
* [windows] Pass the correct installationStatus to the webview installation strategy by @stffabi in https://github.com/wailsapp/wails/pull/1483
* [windows] Make `SetBackgroundColour` compatible for `windows/386` by @stffabi in https://github.com/wailsapp/wails/pull/1493
* Fix lit-ts template by @Orijhins in https://github.com/wailsapp/wails/pull/1494

### Changed

* [windows] Load WebView2 loader from embedded only by @stffabi in https://github.com/wailsapp/wails/pull/1432
* Add showcase entry for October + update homepage carousel entry for October by @marcus-crane in https://github.com/wailsapp/wails/pull/1436
* Always use return in wrapped method by @leaanthony in https://github.com/wailsapp/wails/pull/1410
* [windows] Unlock OSThread after native calls have been finished by @stffabi in https://github.com/wailsapp/wails/pull/1441
* Add `BackgroundColour` and deprecate `RGBA` by @leaanthony in https://github.com/wailsapp/wails/pull/1475
* AssetsHandler remove retry logic in dev mode by @stffabi in https://github.com/wailsapp/wails/pull/1479
* Add Solid JS template to docs by @sidwebworks in https://github.com/wailsapp/wails/pull/1492
* Better signal handling by @leaanthony in https://github.com/wailsapp/wails/pull/1488
* Chore/react 18 create root by @tomanagle in https://github.com/wailsapp/wails/pull/1489

### New Contributors

* @jaesung9507 made their first contribution in https://github.com/wailsapp/wails/pull/1430
* @LukenSkyne made their first contribution in https://github.com/wailsapp/wails/pull/1449
* @andywenk made their first contribution in https://github.com/wailsapp/wails/pull/1460
* @abtin made their first contribution in https://github.com/wailsapp/wails/pull/1461
* @chenxiao1990 made their first contribution in https://github.com/wailsapp/wails/pull/1442
* @NanoNik made their first contribution in https://github.com/wailsapp/wails/pull/1448
* @sidwebworks made their first contribution in https://github.com/wailsapp/wails/pull/1492
* @tomanagle made their first contribution in https://github.com/wailsapp/wails/pull/1489

## [v2.0.0-beta.37] - 2022-05-26

### Added

* Add `nogen` flag in wails dev command by @mondy in https://github.com/wailsapp/wails/pull/1413
* Initial support for new native translucency in Windows Preview by @leaanthony in https://github.com/wailsapp/wails/pull/1400

### Fixed

* Bugfix/incorrect bindings by @leaanthony in https://github.com/wailsapp/wails/pull/1383
* Fix runtime.js events by @polikow in https://github.com/wailsapp/wails/pull/1369
* Fix docs formatting by @antimatter96 in https://github.com/wailsapp/wails/pull/1372
* Events | fixes #1388 by @lambdajack in https://github.com/wailsapp/wails/pull/1390
* bugfix: correct typo by @tmclane in https://github.com/wailsapp/wails/pull/1391
* Fix typo in docs by @LGiki in https://github.com/wailsapp/wails/pull/1393
* Fix typo bindings.js to ipc.js by @rayshoo in https://github.com/wailsapp/wails/pull/1406
* Make sure to execute the menu callbacks on a new goroutine by @stffabi in https://github.com/wailsapp/wails/pull/1403
* Update runtime.d.ts & templates by @Yz4230 in https://github.com/wailsapp/wails/pull/1421
* Add missing className to input in React and Preact templates by @edwardbrowncross in https://github.com/wailsapp/wails/pull/1419

### Changed
* Improve multi-platform builds by @stffabi in https://github.com/wailsapp/wails/pull/1373
* During wails dev only use reload logic if no AssetsHandler are in use by @stffabi in https://github.com/wailsapp/wails/pull/1385
* Update events.mdx by @Junkher in https://github.com/wailsapp/wails/pull/1387
* Add Next.js template by @LGiki in https://github.com/wailsapp/wails/pull/1394
* Add docs on wails generate module by @TechplexEngineer in https://github.com/wailsapp/wails/pull/1414
* Add macos custom menu EditMenu tips by @daodao97 in https://github.com/wailsapp/wails/pull/1423

### New Contributors
* @polikow made their first contribution in https://github.com/wailsapp/wails/pull/1369
* @antimatter96 made their first contribution in https://github.com/wailsapp/wails/pull/1372
* @Junkher made their first contribution in https://github.com/wailsapp/wails/pull/1387
* @lambdajack made their first contribution in https://github.com/wailsapp/wails/pull/1390
* @LGiki made their first contribution in https://github.com/wailsapp/wails/pull/1393
* @rayshoo made their first contribution in https://github.com/wailsapp/wails/pull/1406
* @TechplexEngineer made their first contribution in https://github.com/wailsapp/wails/pull/1414
* @mondy made their first contribution in https://github.com/wailsapp/wails/pull/1413
* @Yz4230 made their first contribution in https://github.com/wailsapp/wails/pull/1421
* @daodao97 made their first contribution in https://github.com/wailsapp/wails/pull/1423
* @edwardbrowncross made their first contribution in https://github.com/wailsapp/wails/pull/1419


## [v2.0.0-beta.36] - 2022-04-27

### Fixed
- [v2] Validate devServer property to be of the correct form by [@stffabi](https://github.com/stffabi) in https://github.com/wailsapp/wails/pull/1359
- [v2, darwin] Initialize native variables on stack to prevent segfault by [@stffabi](https://github.com/stffabi) in https://github.com/wailsapp/wails/pull/1362
- Vue-TS template fix

### Changed
- Added `OnStartup` method back to default templates

## [v2.0.0-beta.35] - 2022-04-27

### Breaking Changes

- When data was sent to the `EventsOn` callback, it was being sent as a slice of values, instead of optional parameters to the method. `EventsOn` now works as expected, but you will need to update your code if you currently use this. [More information](https://github.com/wailsapp/wails/issues/1324) `EventsOn` now works as expected, but you will need to update your code if you currently use this. [More information](https://github.com/wailsapp/wails/issues/1324)
- The broken `bindings.js` and `bindings.d.ts` files have been replaced by a new JS/TS code generation system. More details [here](https://wails.io/docs/howdoesitwork#calling-bound-go-methods) More details [here](https://wails.io/docs/howdoesitwork#calling-bound-go-methods)

### Added

- **New Templates**: Svelte, React, Vue, Preact, Lit and Vanilla templates, both JS and TS versions. `wails init -l` for more info. `wails init -l` for more info.
- Default templates now powered by [Vite](https://vitejs.dev). This enables lightning fast reloads when you use `wails dev`! Default templates now powered by [Vite](https://vitejs.dev). This enables lightning fast reloads when you use `wails dev`! This enables lightning fast reloads when you use `wails dev`!
- Add support for external frontend development servers. Add support for external frontend development servers. See `frontend:dev:serverUrl` in the [project config](https://wails.io/docs/reference/project-config) - [@stffabi](https://github.com/stffabi)
- [Fully configurable dark mode](https://wails.io/docs/reference/options#theme) for Windows.
- Hugely improved [WailsJS generation](https://wails.io/docs/howdoesitwork#calling-bound-go-methods) (both Javascript and Typescript)
- Wails doctor now reports information about the wails installation - [@stffabi](https://github.com/stffabi)
- Added docs for [code-signing](https://wails.io/docs/guides/signing) and [NSIS installer](https://wails.io/docs/guides/windows-installer) - [@gardc](https://github.com/gardc)
- Add support for `-trimpath` [build flag](https://wails.io/docs/reference/cli#build)
- Add support for a default AssetsHandler - [@stffabi](https://github.com/stffabi)

### Fixed

- Improved mimetype detection for BOM marker and comments - [@napalu](https://github.com/napalu)
- Remove duplicate mimetype entries - [@napalu](https://github.com/napalu)
- Remove duplicate Typescript imports in generated definition files - [@adalessa](https://github.com/adalessa)
- Add missing method declaration - [@adalessa](https://github.com/adalessa)
- Fix Linux sigabrt on start - [@napalu](https://github.com/napalu)
- Double Click event now works on elements with `data-wails-drag` attribute - [@jicg](https://github.com/jicg)
- Suppress resizing during minimize of a frameless window - [@stffabi](https://github.com/stffabi)
- Fixed TS/JS generation for Go methods with no returns
- Fixed WailsJS being generated in project directory

### Changed

- Website docs are now versioned
- Improved `runtime.Environment` call
- Improve the close action for Mac
- A bunch of dependabot security updates
- Improved website content - [@misitebao](https://github.com/misitebao)
- Upgrade issue template - [@misitebao](https://github.com/misitebao)
- Convert documents that don't require version management to individual pages
  - [@misitebao](https://github.com/misitebao)
- Website now using Algolia search

## [v2.0.0-beta.34] - 2022-03-26

### Added

- Add support for 'DomReady' callback on linux by [@napalu](https://github.com/napalu) in #1249
- MacOS - Show extension by default by [@leaanthony](https://github.com/leaanthony)  in #1228

### Fixed

- [v2, nsis] Seems like / as path separator works only for some directives in a cross platform way by [@stffabi](https://github.com/stffabi) in #1227
- import models on binding definition by [@adalessa](https://github.com/adalessa) in #1231
- Use local search on website by [@leaanthony](https://github.com/leaanthony)  in #1234
- Ensure binary resources can be served by [@napalu](https://github.com/napalu) in #1240
- Only retry loading assets when loading from disk by [@leaanthony](https://github.com/leaanthony)  in #1241
- [v2, windows] Fix maximised start state by [@stffabi](https://github.com/stffabi) in #1243
- Ensure Linux IsFullScreen uses GDK_WINDOW_STATE_FULLSCREEN bitmask appropriately. by [@ianmjones](https://github.com/ianmjones) in #1245 by [@ianmjones](https://github.com/ianmjones) in #1245
- Fix memory leak in ExecJS for Mac by [@leaanthony](https://github.com/leaanthony)  in #1230
- Fix, or at least a workaround, for (#1232) by [@BillBuilt](https://github.com/BillBuilt) in #1247
- [v2] Use os.Args[0] for self starting wails by [@stffabi](https://github.com/stffabi) in #1258
- [v2, windows] Windows switch scheme: https -> http by @stefpap in #1255
- Ensure Focus is regained by Webview2 when tabbing by [@leaanthony](https://github.com/leaanthony)  in #1257
- Try to focus window when Show() is called. by [@leaanthony](https://github.com/leaanthony)  in #1212 by [@leaanthony](https://github.com/leaanthony)  in #1212
- Check system for user installed Linux dependencies by [@leaanthony](https://github.com/leaanthony)  in #1180

### Changed

- feat(website): sync documents and add content by [@misitebao](https://github.com/misitebao) in #1215
- refactor(cli): optimize default templates by [@misitebao](https://github.com/misitebao) in #1214
- Run watcher after initial build by [@leaanthony](https://github.com/leaanthony)  in #1216
- Feature/docs update by [@leaanthony](https://github.com/leaanthony)  in #1218
- feat(website): optimize website and sync documents by [@misitebao](https://github.com/misitebao) in #1219
- docs: sync documents by [@misitebao](https://github.com/misitebao) in #1224
- Default index page by [@leaanthony](https://github.com/leaanthony)  in #1229
- Build added win32 compatibility by [@fengweiqiang](https://github.com/fengweiqiang) in #1238
- docs: sync documents by [@misitebao](https://github.com/misitebao) in #1260

## [v2.0.0-beta.33][] - 2022-03-05

### Added

- NSIS Installer support for creating installers for Windows applications - Thanks [@stffabi](https://github.com/stffabi) 🎉
- New frontend:dev:watcher command to spin out 3rd party watchers when using wails dev - Thanks [@stffabi](https://github.com/stffabi)🎉
- Remote templates now support version tags - Thanks [@misitebao](https://github.com/misitebao) 🎉

### Fixed

- A number of fixes for ARM Linux providing a huge improvement - Thanks [@ianmjones](https://github.com/ianmjones) 🎉
- Fixed potential Nil reference when discovering the path to `index.html`
- Fixed crash when using `runtime.Log` methods in a production build
- Improvements to internal file handling meaning webworkers will now work on Windows - Thanks [@stffabi](https://github.com/stffabi)🎉

### Changed

- The Webview2 bootstrapper is now run as a normal user and doesn't require admin rights
- The docs have been improved and updated
- Added troubleshooting guide

[Unreleased]: https://github.com/wailsapp/wails/compare/v2.0.0-beta.33...HEAD
[v2.0.0-beta.33]: https://github.com/wailsapp/wails/compare/v2.0.0-beta.32...v2.0.0-beta.33
