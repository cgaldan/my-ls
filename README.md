# My ls

## Description
My ls project is a simple version of the unix `ls` command.
The `ls` command shows you the files and folders of the directory specified after the command.  By exclusion of this directory, it shows the files and folders of the present directory.

## Features
The behavior of this program is identical to the original `ls` command with the following flags :   
- `-l` : For long listing format with detailed info for each file
- `-R` : For Recursive directory listing
- `-a` : To include also the hidden files in the listing
- `-r` : To reverse the sort order
- `-t` : To sort by the modification time

## Usage
To run this program you need to install **golang**
1. Clone the repository :
    ```bash
    git clone https://platform.zone01.gr/git/cgkaldan/my-ls.git
2. Build the program :
    ```go
   make build
   ```
3. Run the program :
- **flag :** You can run the program with the flags provided. (optional)
- **filenames :** You can run the program with a name of a file that exists in the current directory. (optional)
    ```go
    ./myls [FLAGS]... [FILENAMES]...
    ```
# ToDo list
- Fixing dir headers format