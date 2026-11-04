@echo off
echo Testing release.go --create-release-notes functionality
echo ======================================================

REM Save current directory
set ORIGINAL_DIR=%CD%

REM Go to v3 root (where release_changelog.md should be)
cd ..\..

REM Backup existing release_changelog.md if it exists
if exist release_changelog.md (
    copy release_changelog.md release_changelog.md.backup > nul
)

REM Create a test release_changelog.md
echo # Unreleased Changes > release_changelog.md
echo. >> release_changelog.md
echo ^<^!-- >> release_changelog.md
echo This file is used to collect changelog entries for the next v3 alpha release. >> release_changelog.md
echo --^> >> release_changelog.md
echo. >> release_changelog.md
echo ## Added >> release_changelog.md
echo ^<^!-- New features, capabilities, or enhancements --^> >> release_changelog.md
echo - Add Windows dark theme support for menus and menubar >> release_changelog.md
echo - Add `--create-release-notes` flag to release script >> release_changelog.md
echo. >> release_changelog.md
echo ## Changed >> release_changelog.md
echo ^<^!-- Changes in existing functionality --^> >> release_changelog.md
echo - Update Go version to 1.23 in workflow >> release_changelog.md
echo - Improve error handling in release process >> release_changelog.md
echo. >> release_changelog.md
echo ## Fixed >> release_changelog.md
echo ^<^!-- Bug fixes --^> >> release_changelog.md
echo - Fix nightly release workflow changelog extraction >> release_changelog.md
echo - Fix Go cache configuration in GitHub Actions >> release_changelog.md
echo. >> release_changelog.md
echo ## Deprecated >> release_changelog.md
echo ^<^!-- Soon-to-be removed features --^> >> release_changelog.md
echo. >> release_changelog.md
echo ## Removed >> release_changelog.md
echo ^<^!-- Features removed in this release --^> >> release_changelog.md
echo. >> release_changelog.md
echo ## Security >> release_changelog.md
echo ^<^!-- Security-related changes --^> >> release_changelog.md
echo. >> release_changelog.md
echo --- >> release_changelog.md
echo. >> release_changelog.md
echo ### Example Entries: >> release_changelog.md

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

REM Restore original release_changelog.md if it exists
if exist release_changelog.md.backup (
    move /Y release_changelog.md.backup release_changelog.md > nul
)

cd %ORIGINAL_DIR%

echo.
echo ======================================================
echo Testing complete!