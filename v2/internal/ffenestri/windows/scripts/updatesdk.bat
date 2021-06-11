@echo off
IF %1.==. GOTO NoVersion
nuget install microsoft.web.webview2 -Version %1 -OutputDirectory . >NUL || goto :eof
echo Downloaded microsoft.web.webview2.%1

set sdk_version=%1
set native_dir="%~dp0\microsoft.web.webview2.%sdk_version%\build\native"
copy "%native_dir%\include\*.h" .. >NUL
copy "%native_dir%\x64\WebView2Loader.dll" "..\x64" >NUL
@rd /S /Q "microsoft.web.webview2.%sdk_version%"
del /s version.txt  >nul 2>&1
echo The version of WebView2 used: %sdk_version% > version.txt
echo SDK updated to %sdk_version%
goto :eof

:NoVersion
  echo Please provide a version number, EG: 1.0.664.37
  goto :eof
