#!/bin/bash

echo "======================================================="
echo "Run both my-ls and the system command ls with no arguments."
echo "======================================================="
echo
go run .
echo
echo "======================================================="
echo "Run both my-ls and the system command ls with the arguments: \"<file name>\"."
echo "======================================================="
echo
go run . main.go
echo
echo "======================================================="
echo "Run both my-ls and the system command ls with the arguments: \"<directory name>\"."
echo "======================================================="
echo
go run . test
echo
echo "======================================================="
echo "Run both my-ls and the system command ls with the flag: \"-l\"."
echo "======================================================="
echo
go run . -l
echo
echo "======================================================="
echo "Run both my-ls and the system command ls with the arguments: \"-l <file name>\"."
echo "======================================================="
echo
go run . -l main.go
echo
echo "======================================================="
echo "Run both my-ls and the system command ls with the arguments: \"-l <directory name>\"."
echo "======================================================="
echo
go run . -l test
echo
echo "======================================================="
echo "Run both my-ls and the system command ls with the flag: \"-R\", in a directory with folders in it."
echo "======================================================="
echo
go run . -R test
echo
echo "======================================================="
echo "Run both my-ls and the system command ls with the flag: \"-a\"."
echo "======================================================="
echo
go run . -a
echo
echo "======================================================="
echo "Run both my-ls and the system command ls with the flag: \"-r\"."
echo "======================================================="
echo
go run . -r
echo
echo "======================================================="
echo "Run both my-ls and the system command ls with the flag: \"-t\"."
echo "======================================================="
echo
go run . -t
echo
echo "======================================================="
echo "Run both my-ls and the system command ls with the flag: \"-la\"."
echo "======================================================="
echo
go run . -la
echo
echo "======================================================="
echo "Run both my-ls and the system command ls with the arguments: \"-l -t <directory name>\"."
echo "======================================================="
echo
go run . -l -t test
echo
echo "======================================================="
echo "Run both my-ls and the system command ls with the arguments: \"-lRr <directory name>\", in which the directory chosen contains folders."
echo "======================================================="
echo
go run . -lRr test
echo
echo "======================================================="
echo "Run both my-ls and the system command ls with the arguments: \"-l <directory name> -a <file name>\"."
echo "======================================================="
echo
go run . -l test -a main.go
echo
echo "======================================================="
echo "Run both my-ls and the system command ls with the arguments: \"-lR <directory name>///<sub directory name>/// <directory name>/<sub directory name>/\""
echo "======================================================="
echo
go run . -lR test///test_dir_00/// -- "-/test_folder/"
echo
echo "======================================================="
echo "Run both my-ls and the system command ls with the arguments: \"-alRrt <directory name>\", in which the directory chosen contains folders and files within folders."
echo "======================================================="
echo
go run . -alRrt test
echo
echo "======================================================="
echo "Create directory with - name and run both my-ls and the system command ls with the arguments: \"-\""
echo "======================================================="
echo
go run . "-"
echo
echo "======================================================="