#!/bin/bash

echo "======================================================="
echo "Run both my-ls and the system command ls with no arguments."
echo "======================================================="
echo
ls --color=always
echo
echo "======================================================="
echo "Run both my-ls and the system command ls with the arguments: \"<file name>\"."
echo "======================================================="
echo
ls --color=always main.go
echo
echo "======================================================="
echo "Run both my-ls and the system command ls with the arguments: \"<directory name>\"."
echo "======================================================="
echo
ls --color=always test
echo
echo "======================================================="
echo "Run both my-ls and the system command ls with the flag: \"-l\"."
echo "======================================================="
echo
ls --color=always -l
echo
echo "======================================================="
echo "Run both my-ls and the system command ls with the arguments: \"-l <file name>\"."
echo "======================================================="
echo
ls --color=always -l main.go
echo
echo "======================================================="
echo "Run both my-ls and the system command ls with the arguments: \"-l <directory name>\"."
echo "======================================================="
echo
ls --color=always -l test
echo
echo "======================================================="
echo "Run both my-ls and the system command ls with the flag: \"-R\", in a directory with folders in it."
echo "======================================================="
echo
ls --color=always -R test
echo
echo "======================================================="
echo "Run both my-ls and the system command ls with the flag: \"-a\"."
echo "======================================================="
echo
ls --color=always -a
echo
echo "======================================================="
echo "Run both my-ls and the system command ls with the flag: \"-r\"."
echo "======================================================="
echo
ls --color=always -r
echo
echo "======================================================="
echo "Run both my-ls and the system command ls with the flag: \"-t\"."
echo "======================================================="
echo
ls --color=always -t
echo
echo "======================================================="
echo "Run both my-ls and the system command ls with the flag: \"-la\"."
echo "======================================================="
echo
ls --color=always -la
echo
echo "======================================================="
echo "Run both my-ls and the system command ls with the arguments: \"-l -t <directory name>\"."
echo "======================================================="
echo
ls --color=always -l -t test
echo
echo "======================================================="
echo "Run both my-ls and the system command ls with the arguments: \"-lRr <directory name>\", in which the directory chosen contains folders."
echo "======================================================="
echo
ls --color=always -lRr test
echo
echo "======================================================="
echo "Run both my-ls and the system command ls with the arguments: \"-l <directory name> -a <file name>\"."
echo "======================================================="
echo
ls --color=always -l test -a main.go
echo
echo "======================================================="
echo "Run both my-ls and the system command ls with the arguments: \"-lR <directory name>///<sub directory name>/// <directory name>/<sub directory name>/\""
echo "======================================================="
echo
ls --color=always -lR test///test_dir_00/// -- "-/test_folder/"
echo
echo "======================================================="
echo "Run both my-ls and the system command ls with the arguments: \"-alRrt <directory name>\", in which the directory chosen contains folders and files within folders."
echo "======================================================="
echo
ls --color=always -alRrt test
echo
echo "======================================================="
echo "Create directory with - name and run both my-ls and the system command ls with the arguments: \"-\""
echo "======================================================="
echo
ls --color=always "-"
echo
echo "======================================================="