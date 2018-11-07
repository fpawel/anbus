SET dir=%HOMEDRIVE%%HOMEPATH%\.anbus
buildmingw32 go build -o %dir%\anbus.exe -ldflags="-H windowsgui" github.com/fpawel/anbus/cmd
go build -o %dir%\runbus.exe -ldflags="-H windowsgui" github.com/fpawel/anbus/cmd/run
start %dir%
