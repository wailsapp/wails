@echo off
echo Testing edge cases for release.go
echo =================================

set ORIGINAL_DIR=%CD%
cd ..\..

REM Backup existing file
if exist UNRELEASED_CHANGELOG.md (
    copy UNRELEASED_CHANGELOG.md UNRELEASED_CHANGELOG.md.backup > nul
)

echo.
echo Test 1: Empty changelog (should fail)
echo -------------------------------------

REM Create empty changelog
echo # Unreleased Changes > UNRELEASED_CHANGELOG.md
echo. >> UNRELEASED_CHANGELOG.md
echo ## Added >> UNRELEASED_CHANGELOG.md
echo ^<^!-- New features --^> >> UNRELEASED_CHANGELOG.md
echo. >> UNRELEASED_CHANGELOG.md
echo ## Changed >> UNRELEASED_CHANGELOG.md
echo ^<^!-- Changes --^> >> UNRELEASED_CHANGELOG.md
echo. >> UNRELEASED_CHANGELOG.md
echo ## Fixed >> UNRELEASED_CHANGELOG.md
echo ^<^!-- Bug fixes --^> >> UNRELEASED_CHANGELOG.md
echo. >> UNRELEASED_CHANGELOG.md
echo --- >> UNRELEASED_CHANGELOG.md

cd tasks\release
go run release.go --create-release-notes 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo SUCCESS: Command failed as expected for empty changelog
) else (
    echo FAIL: Command should have failed for empty changelog
)

echo.
echo Test 2: Only comments (should fail)
echo -----------------------------------

cd ..\..
echo # Unreleased Changes > UNRELEASED_CHANGELOG.md
echo. >> UNRELEASED_CHANGELOG.md
echo ## Added >> UNRELEASED_CHANGELOG.md
echo ^<^!-- This is just a comment --^> >> UNRELEASED_CHANGELOG.md
echo ^<^!-- Another comment --^> >> UNRELEASED_CHANGELOG.md
echo. >> UNRELEASED_CHANGELOG.md
echo --- >> UNRELEASED_CHANGELOG.md

cd tasks\release
go run release.go --create-release-notes 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo SUCCESS: Command failed as expected for comment-only changelog
) else (
    echo FAIL: Command should have failed for comment-only changelog
)

echo.
echo Test 3: Mixed bullet styles
echo ---------------------------

cd ..\..
echo # Unreleased Changes > UNRELEASED_CHANGELOG.md
echo. >> UNRELEASED_CHANGELOG.md
echo ## Added >> UNRELEASED_CHANGELOG.md
echo - Dash bullet point >> UNRELEASED_CHANGELOG.md
echo * Asterisk bullet point >> UNRELEASED_CHANGELOG.md
echo - Another dash >> UNRELEASED_CHANGELOG.md
echo. >> UNRELEASED_CHANGELOG.md
echo --- >> UNRELEASED_CHANGELOG.md

cd tasks\release
go run release.go --create-release-notes
if %ERRORLEVEL% EQU 0 (
    echo SUCCESS: Mixed bullet styles handled
    echo Content:
    type ..\..\release_notes.md
) else (
    echo FAIL: Mixed bullet styles should work
)

echo.
echo Test 4: Custom output path
echo --------------------------

go run release.go --create-release-notes ..\..\custom_notes.md
if %ERRORLEVEL% EQU 0 (
    if exist "..\..\custom_notes.md" (
        echo SUCCESS: Custom path works
        del ..\..\custom_notes.md
    ) else (
        echo FAIL: Custom path file not created
    )
) else (
    echo FAIL: Custom path should work
)

REM Clean up
cd ..\..
if exist release_notes.md del release_notes.md
if exist UNRELEASED_CHANGELOG.md.backup (
    move /Y UNRELEASED_CHANGELOG.md.backup UNRELEASED_CHANGELOG.md > nul
)

cd %ORIGINAL_DIR%

echo.
echo =================================
echo Edge case testing complete!