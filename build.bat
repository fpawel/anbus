SET dir=%HOMEDRIVE%%HOMEPATH%\.anbus
buildmingw32 go build -o %dir%\anbus.exe -ldflags="-H windowsgui" github.com/fpawel/anbus/cmd
go build -o %dir%\runankat.exe -ldflags="-H windowsgui" github.com/fpawel/ankat/run
start %dir%
