@echo off
rem UWP applications disallow loopbacks. For development purpose it is necessary to run the following command in an administrative command prompt
CheckNetIsolation.exe LoopbackExempt -a -n="Microsoft.Win32WebViewHost_cw5n1h2txyewy"
