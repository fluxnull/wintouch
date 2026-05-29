@echo off
setlocal EnableExtensions EnableDelayedExpansion

REM ============================================================================
REM WinTouch CMD-only behavior runner v3
REM No PowerShell.
REM Constructs individual command lines, runs each through a temporary .cmd file, and saves all output.
REM ============================================================================

chcp 65001 >nul

set "SCRIPT_DIR=%~dp0"
set "TOUCH_EXE=%~1"

if "%TOUCH_EXE%"=="" (
    if exist "%SCRIPT_DIR%touch.exe" set "TOUCH_EXE=%SCRIPT_DIR%touch.exe"
)
if "%TOUCH_EXE%"=="" (
    if exist "%SCRIPT_DIR%bin\touch.exe" set "TOUCH_EXE=%SCRIPT_DIR%bin\touch.exe"
)
if "%TOUCH_EXE%"=="" (
    for %%A in (touch.exe) do (
        if not "%%~$PATH:A"=="" set "TOUCH_EXE=%%~$PATH:A"
    )
)

if not exist "%TOUCH_EXE%" (
    echo ERROR: touch.exe not found.
    echo Usage:
    echo   test-wintouch-cmd.cmd F:\path\to\touch.exe
    exit /b 1
)

set "STAMP=%DATE%_%TIME%"
set "STAMP=%STAMP:/=-%"
set "STAMP=%STAMP::=-%"
set "STAMP=%STAMP:.=-%"
set "STAMP=%STAMP: =0%"
set "STAMP=%STAMP:,=-%"

set "OUT_ROOT=%SCRIPT_DIR%test-output"
set "CASE_ROOT=%OUT_ROOT%\wintouch_cmd_cases_%STAMP%"
set "REPORT=%OUT_ROOT%\wintouch_cmd_report_%STAMP%.txt"

mkdir "%OUT_ROOT%" >nul 2>nul
mkdir "%CASE_ROOT%" >nul 2>nul

cd /d "%CASE_ROOT%" || exit /b 1

call :Header

REM ============================================================================
REM HELP AND VERSION
REM ============================================================================

call :Section "HELP AND VERSION"
set "CMDLINE="%TOUCH_EXE%" /?"
call :Run "help_slash_question"
set "CMDLINE="%TOUCH_EXE%" /H"
call :Run "help_slash_h"
set "CMDLINE="%TOUCH_EXE%" /HELP"
call :Run "help_slash_help"
set "CMDLINE="%TOUCH_EXE%" -h"
call :Run "help_dash_h"
set "CMDLINE="%TOUCH_EXE%" -help"
call :Run "help_dash_help"
set "CMDLINE="%TOUCH_EXE%" --help"
call :Run "help_long_help"
set "CMDLINE="%TOUCH_EXE%" -v"
call :Run "version_dash_v"
set "CMDLINE="%TOUCH_EXE%" --version"
call :Run "version_long"

REM ============================================================================
REM DEFAULT CREATE / NO-CREATE / PARENTS
REM ============================================================================

call :Section "DEFAULT CREATE / NO-CREATE / PARENTS"

set "CMDLINE="%TOUCH_EXE%" default_create.txt"
call :Run "default_create"
call :State "default_create.txt"

set "CMDLINE="%TOUCH_EXE%" default_create.txt"
call :Run "default_update_existing"
call :State "default_create.txt"

set "CMDLINE="%TOUCH_EXE%" -n missing_no_create.txt"
call :Run "no_create_missing"
call :State "missing_no_create.txt"

set "CMDLINE="%TOUCH_EXE%" -p logs\build\output.txt"
call :Run "parents_create"
call :State "logs\build\output.txt"

set "CMDLINE="%TOUCH_EXE%" -np logs_np\build\output.txt"
call :Run "parents_no_create_cluster_np"
call :State "logs_np\build\output.txt"

REM ============================================================================
REM TIMESTAMP SELECTION USING -d
REM ============================================================================

call :Section "TIMESTAMP SELECTION USING -d"

call :MakeFile selector_a.txt
set "CMDLINE="%TOUCH_EXE%" -a -d "2026-05-28 21:30:00" selector_a.txt"
call :Run "selector_a_accessed_only"
call :State "selector_a.txt"

call :MakeFile selector_m.txt
set "CMDLINE="%TOUCH_EXE%" -m -d "2026-05-28 21:30:00" selector_m.txt"
call :Run "selector_m_modified_only"
call :State "selector_m.txt"

call :MakeFile selector_c.txt
set "CMDLINE="%TOUCH_EXE%" -c -d "2026-05-28 21:30:00" selector_c.txt"
call :Run "selector_c_created_only"
call :State "selector_c.txt"

call :MakeFile selector_am.txt
set "CMDLINE="%TOUCH_EXE%" -am -d "2026-05-28 21:30:00" selector_am.txt"
call :Run "selector_am_accessed_modified"
call :State "selector_am.txt"

call :MakeFile selector_ac.txt
set "CMDLINE="%TOUCH_EXE%" -ac -d "2026-05-28 21:30:00" selector_ac.txt"
call :Run "selector_ac_accessed_created"
call :State "selector_ac.txt"

call :MakeFile selector_mc.txt
set "CMDLINE="%TOUCH_EXE%" -mc -d "2026-05-28 21:30:00" selector_mc.txt"
call :Run "selector_mc_modified_created"
call :State "selector_mc.txt"

call :MakeFile selector_amc.txt
set "CMDLINE="%TOUCH_EXE%" -amc -d "2026-05-28 21:30:00" selector_amc.txt"
call :Run "selector_amc_all"
call :State "selector_amc.txt"

REM ============================================================================
REM BOOLEAN CLUSTERS
REM ============================================================================

call :Section "BOOLEAN CLUSTERS"

call :MakeFile cluster_np.txt
set "CMDLINE="%TOUCH_EXE%" -np cluster_np.txt"
call :Run "cluster_np"
call :State "cluster_np.txt"

call :MakeFile cluster_rn.txt
set "CMDLINE="%TOUCH_EXE%" -rn cluster_rn.txt"
call :Run "cluster_rn_no_create_runas_if_needed"
call :State "cluster_rn.txt"

call :MakeFile cluster_amcr.txt
set "CMDLINE="%TOUCH_EXE%" -amcr -d "2026-05-28 21:30:00" cluster_amcr.txt"
call :Run "cluster_amcr"
call :State "cluster_amcr.txt"

call :MakeFile cluster_amcnp.txt
set "CMDLINE="%TOUCH_EXE%" -amcnp -d "2026-05-28 21:30:00" cluster_amcnp.txt"
call :Run "cluster_amcnp"
call :State "cluster_amcnp.txt"

call :MakeFile cluster_amcnpf.txt
set "CMDLINE="%TOUCH_EXE%" -amcnpf -d "2026-05-28 21:30:00" cluster_amcnpf.txt"
call :Run "cluster_amcnpf"
call :State "cluster_amcnpf.txt"

call :MakeFile cluster_amcnpi.txt
set "CMDLINE="%TOUCH_EXE%" -amcnpi -d "2026-05-28 21:30:00" cluster_amcnpi.txt"
call :Run "cluster_amcnpi"
call :State "cluster_amcnpi.txt"

call :MakeFile cluster_amcnpr.txt
set "CMDLINE="%TOUCH_EXE%" -amcnpr -d "2026-05-28 21:30:00" cluster_amcnpr.txt"
call :Run "cluster_amcnpr"
call :State "cluster_amcnpr.txt"

REM ============================================================================
REM VALUE OPTION ORDER
REM ============================================================================

call :Section "VALUE OPTION ORDER"

call :MakeFile order_bool_then_value.txt
set "CMDLINE="%TOUCH_EXE%" -amc -d "2026-05-28 21:30:00" order_bool_then_value.txt"
call :Run "order_boolean_before_value"
call :State "order_bool_then_value.txt"

call :MakeFile order_value_then_bool.txt
set "CMDLINE="%TOUCH_EXE%" -d "2026-05-28 21:30:00" -amc order_value_then_bool.txt"
call :Run "order_value_before_boolean"
call :State "order_value_then_bool.txt"

call :MakeFile order_value_between_bool.txt
set "CMDLINE="%TOUCH_EXE%" -m -d "2026-05-28 21:30:00" -c order_value_between_bool.txt"
call :Run "order_value_between_boolean"
call :State "order_value_between_bool.txt"

REM ============================================================================
REM REJECT ATTACHED VALUE FORMS
REM ============================================================================

call :Section "REJECT ATTACHED VALUE FORMS"

set "CMDLINE="%TOUCH_EXE%" -sreference.txt bad_s.txt"
call :Run "reject_attached_s"
set "CMDLINE="%TOUCH_EXE%" -d2026-05-28T21:30:00 bad_d.txt"
call :Run "reject_attached_d"
set "CMDLINE="%TOUCH_EXE%" -t202605282130.00 bad_t.txt"
call :Run "reject_attached_t"
set "CMDLINE="%TOUCH_EXE%" -e0130 bad_e.txt"
call :Run "reject_attached_e"
set "CMDLINE="%TOUCH_EXE%" -amct bad_cluster_t.txt"
call :Run "reject_value_inside_cluster_t"
set "CMDLINE="%TOUCH_EXE%" -amce bad_cluster_e.txt"
call :Run "reject_value_inside_cluster_e"

REM ============================================================================
REM DATE FORMATS
REM ============================================================================

call :Section "DATE FORMATS"

call :MakeFile date_space.txt
set "CMDLINE="%TOUCH_EXE%" -d "2026-05-28 21:30:00" date_space.txt"
call :Run "date_space"
call :State "date_space.txt"

call :MakeFile date_T.txt
set "CMDLINE="%TOUCH_EXE%" -d "2026-05-28T21:30:00" date_T.txt"
call :Run "date_T"
call :State "date_T.txt"

call :MakeFile date_fraction_space.txt
set "CMDLINE="%TOUCH_EXE%" -d "2026-05-28 21:30:00.123" date_fraction_space.txt"
call :Run "date_fraction_space"
call :State "date_fraction_space.txt"

call :MakeFile date_fraction_T.txt
set "CMDLINE="%TOUCH_EXE%" -d "2026-05-28T21:30:00.123" date_fraction_T.txt"
call :Run "date_fraction_T"
call :State "date_fraction_T.txt"

call :MakeFile date_zulu.txt
set "CMDLINE="%TOUCH_EXE%" -d "2026-05-28T21:30:00Z" date_zulu.txt"
call :Run "date_zulu"
call :State "date_zulu.txt"

REM ============================================================================
REM BSD COMPACT TIMESTAMP -t
REM ============================================================================

call :Section "BSD COMPACT TIMESTAMP -t"

call :MakeFile timestamp_12.txt
set "CMDLINE="%TOUCH_EXE%" -t 202605282130.00 timestamp_12.txt"
call :Run "timestamp_12_digit"
call :State "timestamp_12.txt"

call :MakeFile timestamp_10.txt
set "CMDLINE="%TOUCH_EXE%" -t 2605282130.00 timestamp_10.txt"
call :Run "timestamp_10_digit"
call :State "timestamp_10.txt"

call :MakeFile timestamp_8.txt
set "CMDLINE="%TOUCH_EXE%" -t 05282130.00 timestamp_8.txt"
call :Run "timestamp_8_digit_current_year"
call :State "timestamp_8.txt"

call :MakeFile timestamp_am.txt
set "CMDLINE="%TOUCH_EXE%" -am -t 202605282130.00 timestamp_am.txt"
call :Run "timestamp_am_selector"
call :State "timestamp_am.txt"

REM ============================================================================
REM SYNC -s
REM ============================================================================

call :Section "SYNC -s"

call :MakeFile reference.txt
set "CMDLINE="%TOUCH_EXE%" -amc -d "2027-06-29 22:31:01" reference.txt"
call :Run "set_reference_times"
call :State "reference.txt"

call :MakeFile sync_all.txt
set "CMDLINE="%TOUCH_EXE%" -s reference.txt sync_all.txt"
call :Run "sync_all"
call :State "sync_all.txt"

call :MakeFile sync_am.txt
set "CMDLINE="%TOUCH_EXE%" -am -s reference.txt sync_am.txt"
call :Run "sync_am"
call :State "sync_am.txt"

call :MakeFile sync_with_edit.txt
set "CMDLINE="%TOUCH_EXE%" -s reference.txt -e +0130 sync_with_edit.txt"
call :Run "sync_with_edit_plus"
call :State "sync_with_edit.txt"

REM ============================================================================
REM EDIT OFFSET -e
REM ============================================================================

call :Section "EDIT OFFSET -e"

call :MakeFile edit_plus.txt
set "CMDLINE="%TOUCH_EXE%" -amc -d "2026-05-28 21:30:00" edit_plus.txt"
call :Run "seed_edit_plus"
set "CMDLINE="%TOUCH_EXE%" -e +0130 edit_plus.txt"
call :Run "edit_plus"
call :State "edit_plus.txt"

call :MakeFile edit_unsigned.txt
set "CMDLINE="%TOUCH_EXE%" -amc -d "2026-05-28 21:30:00" edit_unsigned.txt"
call :Run "seed_edit_unsigned"
set "CMDLINE="%TOUCH_EXE%" -e 0130 edit_unsigned.txt"
call :Run "edit_unsigned"
call :State "edit_unsigned.txt"

call :MakeFile edit_minus_m.txt
set "CMDLINE="%TOUCH_EXE%" -m -d "2026-05-28 21:30:00" edit_minus_m.txt"
call :Run "seed_edit_minus_m"
set "CMDLINE="%TOUCH_EXE%" -m -e -0100 edit_minus_m.txt"
call :Run "edit_minus_m"
call :State "edit_minus_m.txt"

call :MakeFile edit_date_plus.txt
set "CMDLINE="%TOUCH_EXE%" -amc -d "2026-05-28 21:30:00" -e +0130 edit_date_plus.txt"
call :Run "edit_date_plus"
call :State "edit_date_plus.txt"

set "CMDLINE="%TOUCH_EXE%" -e +0130 edit_missing_create.txt"
call :Run "edit_missing_create"
call :State "edit_missing_create.txt"

REM ============================================================================
REM CONFLICTS
REM ============================================================================

call :Section "CONFLICTS"

set "CMDLINE="%TOUCH_EXE%" -i -f conflict_if.txt"
call :Run "conflict_ignore_follow_reparse"
set "CMDLINE="%TOUCH_EXE%" -s reference.txt -d "2026-05-28 21:30:00" conflict_sd.txt"
call :Run "conflict_sync_date"
set "CMDLINE="%TOUCH_EXE%" -d "2026-05-28 21:30:00" -t 202605282130.00 conflict_dt.txt"
call :Run "conflict_date_timestamp"
set "CMDLINE="%TOUCH_EXE%" -s reference.txt -t 202605282130.00 conflict_st.txt"
call :Run "conflict_sync_timestamp"

REM ============================================================================
REM WILDCARDS
REM ============================================================================

call :Section "WILDCARDS"

call :MakeFile wild_one.go
call :MakeFile wild_two.go
set "CMDLINE="%TOUCH_EXE%" -m -d "2026-05-28 21:30:00" *.go"
call :Run "wildcard_go_modified"
call :State "wild_one.go"
call :State "wild_two.go"
set "CMDLINE="%TOUCH_EXE%" -n *.nomatch"
call :Run "wildcard_no_match_no_create"

REM ============================================================================
REM DIRECTORIES
REM ============================================================================

call :Section "DIRECTORIES"

mkdir dir_target >nul 2>nul
set "CMDLINE="%TOUCH_EXE%" -m -d "2026-05-28 21:30:00" dir_target"
call :Run "directory_modified"
call :State "dir_target"

REM ============================================================================
REM UNICODE PATH
REM ============================================================================

call :Section "UNICODE PATH"

call :MakeFile "unicode_test_é_测试.txt"
set "CMDLINE="%TOUCH_EXE%" -m -d "2026-05-28 21:30:00" "unicode_test_é_测试.txt""
call :Run "unicode_path"
call :State "unicode_test_é_测试.txt"

REM ============================================================================
REM LONG PATH
REM ============================================================================

call :Section "LONG PATH"

set "LONGDIR=longpath"
for /L %%N in (1,1,7) do set "LONGDIR=!LONGDIR!\aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
set "LONGFILE=!LONGDIR!\long_target.txt"
set "CMDLINE="%TOUCH_EXE%" -p "!LONGFILE!""
call :Run "long_path_create"
call :State "!LONGFILE!"

REM ============================================================================
REM REPARSE / SYMLINK
REM ============================================================================

call :Section "REPARSE / SYMLINK"

call :MakeFile reparse_target.txt
cmd /c mklink reparse_link.txt reparse_target.txt >> "%REPORT%" 2>&1
echo MKLINK_ERRORLEVEL=%ERRORLEVEL%>> "%REPORT%"
if exist reparse_link.txt (
    set "CMDLINE="%TOUCH_EXE%" -m -d "2026-05-28 21:30:00" reparse_link.txt"
    call :Run "reparse_default_follow"
    call :State "reparse_link.txt"
    call :State "reparse_target.txt"

    set "CMDLINE="%TOUCH_EXE%" -f -m -d "2027-06-29 22:31:01" reparse_link.txt"
    call :Run "reparse_explicit_follow"
    call :State "reparse_link.txt"
    call :State "reparse_target.txt"

    set "CMDLINE="%TOUCH_EXE%" -i -m -d "2028-07-30 23:32:02" reparse_link.txt"
    call :Run "reparse_ignore"
    call :State "reparse_link.txt"
    call :State "reparse_target.txt"
) else (
    echo SKIP reparse tests: mklink failed. Enable Developer Mode or run elevated.>> "%REPORT%"
)

REM ============================================================================
REM RUNAS
REM ============================================================================

call :Section "RUNAS"

echo NOTE: -r / --runas can trigger UAC. This script does not force an interactive UAC prompt.>> "%REPORT%"
echo Manual runas command examples:>> "%REPORT%"
echo   "%TOUCH_EXE%" -r "C:\Windows\System32\drivers\etc\hosts">> "%REPORT%"
echo   "%TOUCH_EXE%" -rm -d "2026-05-28 21:30:00" protected.txt>> "%REPORT%"

REM ============================================================================
REM FINISH
REM ============================================================================

call :Section "SUMMARY"
echo Report:   "%REPORT%">> "%REPORT%"
echo CaseRoot: "%CASE_ROOT%">> "%REPORT%"
echo.>> "%REPORT%"
echo Finished: %DATE% %TIME%>> "%REPORT%"

echo.
echo DONE
echo Report:   "%REPORT%"
echo CaseRoot: "%CASE_ROOT%"
exit /b 0

REM ============================================================================
REM SUBROUTINES
REM ============================================================================

:Header
(
echo WINTOUCH CMD-ONLY TEST REPORT
echo Generated:    %DATE% %TIME%
echo TouchExe:     "%TOUCH_EXE%"
echo CaseRoot:     "%CASE_ROOT%"
echo Report:       "%REPORT%"
echo ComputerName: %COMPUTERNAME%
echo UserName:     %USERNAME%
echo.
)> "%REPORT%"
exit /b 0

:Section
echo.>> "%REPORT%"
echo =============================================================================>> "%REPORT%"
echo %~1>> "%REPORT%"
echo =============================================================================>> "%REPORT%"
echo %~1
exit /b 0

:Run
set "TEST_NAME=%~1"
echo.>> "%REPORT%"
echo RUN: %TEST_NAME%>> "%REPORT%"
echo CMD: !CMDLINE!>> "%REPORT%"

set "RUNFILE=%CASE_ROOT%\__wintouch_run_%RANDOM%_%RANDOM%.cmd"
(
    echo @echo off
    echo cd /d "%CASE_ROOT%"
    echo !CMDLINE! ^>^> "%REPORT%" 2^>^&1
    echo exit /b %%ERRORLEVEL%%
)> "!RUNFILE!"

cmd /d /c ""!RUNFILE!""
set "RC=!ERRORLEVEL!"

del /q "!RUNFILE!" >nul 2>nul

echo ERRORLEVEL: !RC!>> "%REPORT%"
exit /b 0

:State
set "TARGET=%~1"
echo.>> "%REPORT%"
echo STATE: "%TARGET%">> "%REPORT%"
if exist "%TARGET%" (
    echo DIR /T:C>> "%REPORT%"
    dir /a /t:c "%TARGET%" >> "%REPORT%" 2>&1
    echo DIR /T:W>> "%REPORT%"
    dir /a /t:w "%TARGET%" >> "%REPORT%" 2>&1
    echo DIR /T:A>> "%REPORT%"
    dir /a /t:a "%TARGET%" >> "%REPORT%" 2>&1
    echo ATTRIB>> "%REPORT%"
    attrib "%TARGET%" >> "%REPORT%" 2>&1
    echo FSUTIL REPARSEPOINT QUERY>> "%REPORT%"
    fsutil reparsepoint query "%TARGET%" >> "%REPORT%" 2>&1
) else (
    echo MISSING>> "%REPORT%"
)
exit /b 0

:MakeFile
set "MF=%~1"
for %%D in ("%MF%") do (
    if not "%%~dpD"=="" mkdir "%%~dpD" >nul 2>nul
)
break > "%MF%"
echo MAKEFILE: "%MF%">> "%REPORT%"
exit /b 0
