@echo off
echo Testing edge cases for release.go
echo =================================

set ORIGINAL_DIR=%CD%
cd ..\..

REM Backup existing file
if exist release_changelog.md (
    copy release_changelog.md release_changelog.md.backup > nul
)

echo.
echo Test 1: Empty changelog (should fail)
echo -------------------------------------

REM Create empty changelog
echo # Unreleased Changes > release_changelog.md
echo. >> release_changelog.md
echo ## Added >> release_changelog.md
echo ^<^!-- New features --^> >> release_changelog.md
echo. >> release_changelog.md
echo ## Changed >> release_changelog.md
echo ^<^!-- Changes --^> >> release_changelog.md
echo. >> release_changelog.md
echo ## Fixed >> release_changelog.md
echo ^<^!-- Bug fixes --^> >> release_changelog.md
echo. >> release_changelog.md
echo --- >> release_changelog.md

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
echo # Unreleased Changes > release_changelog.md
echo. >> release_changelog.md
echo ## Added >> release_changelog.md
echo ^<^!-- This is just a comment --^> >> release_changelog.md
echo ^<^!-- Another comment --^> >> release_changelog.md
echo. >> release_changelog.md
echo --- >> release_changelog.md

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
echo # Unreleased Changes > release_changelog.md
echo. >> release_changelog.md
echo ## Added >> release_changelog.md
echo - Dash bullet point >> release_changelog.md
echo * Asterisk bullet point >> release_changelog.md
echo - Another dash >> release_changelog.md
echo. >> release_changelog.md
echo --- >> release_changelog.md

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
if exist release_changelog.md.backup (
    move /Y release_changelog.md.backup release_changelog.md > nul
)

cd %ORIGINAL_DIR%

echo.
echo =================================
echo Edge case testing complete!