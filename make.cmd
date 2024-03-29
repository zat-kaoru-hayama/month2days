@echo off
setlocal
set "DIRS=cmd\month2days cmd\month2daysweb cmd\month2daysgui"
call :"%1"
endlocal
exit /b

:"all"
    set GOARCH=386
    for %%I in ( %DIRS% ) do (
        pushd %%I
        go build -ldflags "-s -w"
        popd %%I
    )
    exit /b

:"package"
    set /P "VERSION=Version ?"
    for %%I in ( %DIRS% ) do zip %%~nI-%VERSION%.zip %%I\*.exe
    exit /b
