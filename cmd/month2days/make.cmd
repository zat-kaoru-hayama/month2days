setlocal
set GOARCH=386
call :"%1"
endlocal
exit /b

:""
    go build -ldflags "-s -w"
    upx *.exe
    exit /b

:"zip"
    for %%I in (%CD%) do set NAME=%%~nI
    set /P "VERSION=Please input version:"
    zip %NAME%-%VERSION%-windows-%GOARCH%.zip *.exe
    exit /b
