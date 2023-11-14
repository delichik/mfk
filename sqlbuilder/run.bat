@echo off


echo ">>> prepare go.mod"
rem go mod verify
rem go get github.com/delichik/mfk
go mod tidy

echo ">>> compiling collector"
go build -o .\.mfk\__mfk_sqlbuilder_collector.exe github.com/delichik/mfk/sqlbuilder/collector

echo ">>> collector running"
.\.mfk\__mfk_sqlbuilder_collector.exe -i %1 -o ./.mfk/gen/main.go

echo ">>> compiling generator"
go build -o ./.mfk/__mfk_sqlbuilder_generator.exe ./.mfk/gen/main.go

echo ">>> generator running"
.\.mfk\__mfk_sqlbuilder_generator.exe

echo ">>> clean up"
del .\.mfk\__mfk_sqlbuilder_collector.exe
del .\.mfk\__mfk_sqlbuilder_generator.exe
del .\.mfk\gen\main.go
