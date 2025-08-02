@echo off
echo Testing release.go --create-release-notes functionality
echo ======================================================

REM Save current directory
set ORIGINAL_DIR=%CD%

REM Go to v3 root (where UNRELEASED_CHANGELOG.md should be)
cd ..\..

REM Backup existing UNRELEASED_CHANGELOG.md if it exists
if exist UNRELEASED_CHANGELOG.md (
    copy UNRELEASED_CHANGELOG.md UNRELEASED_CHANGELOG.md.backup > nul
)

REM Create a test UNRELEASED_CHANGELOG.md
echo # Unreleased Changes > UNRELEASED_CHANGELOG.md
echo. >> UNRELEASED_CHANGELOG.md
echo ^<^!-- >> UNRELEASED_CHANGELOG.md
echo This file is used to collect changelog entries for the next v3-alpha release. >> UNRELEASED_CHANGELOG.md
echo --^> >> UNRELEASED_CHANGELOG.md
echo. >> UNRELEASED_CHANGELOG.md
echo ## Added >> UNRELEASED_CHANGELOG.md
echo ^<^!-- New features, capabilities, or enhancements --^> >> UNRELEASED_CHANGELOG.md
echo - Add Windows dark theme support for menus and menubar >> UNRELEASED_CHANGELOG.md
echo - Add `--create-release-notes` flag to release script >> UNRELEASED_CHANGELOG.md
echo. >> UNRELEASED_CHANGELOG.md
echo ## Changed >> UNRELEASED_CHANGELOG.md
echo ^<^!-- Changes in existing functionality --^> >> UNRELEASED_CHANGELOG.md
echo - Update Go version to 1.23 in workflow >> UNRELEASED_CHANGELOG.md
echo - Improve error handling in release process >> UNRELEASED_CHANGELOG.md
echo. >> UNRELEASED_CHANGELOG.md
echo ## Fixed >> UNRELEASED_CHANGELOG.md
echo ^<^!-- Bug fixes --^> >> UNRELEASED_CHANGELOG.md
echo - Fix nightly release workflow changelog extraction >> UNRELEASED_CHANGELOG.md
echo - Fix Go cache configuration in GitHub Actions >> UNRELEASED_CHANGELOG.md
echo. >> UNRELEASED_CHANGELOG.md
echo ## Deprecated >> UNRELEASED_CHANGELOG.md
echo ^<^!-- Soon-to-be removed features --^> >> UNRELEASED_CHANGELOG.md
echo. >> UNRELEASED_CHANGELOG.md
echo ## Removed >> UNRELEASED_CHANGELOG.md
echo ^<^!-- Features removed in this release --^> >> UNRELEASED_CHANGELOG.md
echo. >> UNRELEASED_CHANGELOG.md
echo ## Security >> UNRELEASED_CHANGELOG.md
echo ^<^!-- Security-related changes --^> >> UNRELEASED_CHANGELOG.md
echo. >> UNRELEASED_CHANGELOG.md
echo --- >> UNRELEASED_CHANGELOG.md
echo. >> UNRELEASED_CHANGELOG.md
echo ### Example Entries: >> UNRELEASED_CHANGELOG.md

echo.
echo Test 1: Running with valid content
echo -----------------------------------

REM Run the release script
cd tasks\release
go run release.go --create-release-notes

if %ERRORLEVEL% EQU 0 (
    echo SUCCESS: Command succeeded
    
    REM Check if release_notes.md was created
    if exist "..\..\release_notes.md" (
        echo SUCCESS: release_notes.md was created
        echo.
        echo Content:
        echo --------
        type ..\..\release_notes.md
        echo.
        echo --------
    ) else (
        echo FAIL: release_notes.md was NOT created
    )
) else (
    echo FAIL: Command failed
)

echo.
echo Test 2: Check --check-only flag
echo --------------------------------

REM Test the check-only flag
go run release.go --check-only
if %ERRORLEVEL% EQU 0 (
    echo SUCCESS: --check-only detected content
) else (
    echo FAIL: --check-only did not detect content
)

echo.
echo Test 3: Check --extract-changelog flag
echo --------------------------------------

REM Test the extract-changelog flag
go run release.go --extract-changelog
if %ERRORLEVEL% EQU 0 (
    echo SUCCESS: --extract-changelog succeeded
) else (
    echo FAIL: --extract-changelog failed
)

REM Clean up
cd ..\..
if exist release_notes.md del release_notes.md

REM Restore original UNRELEASED_CHANGELOG.md if it exists
if exist UNRELEASED_CHANGELOG.md.backup (
    move /Y UNRELEASED_CHANGELOG.md.backup UNRELEASED_CHANGELOG.md > nul
)

cd %ORIGINAL_DIR%

echo.
echo ======================================================
echo Testing complete!