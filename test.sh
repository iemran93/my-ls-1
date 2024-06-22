#!/bin/bash

# Function to compare the outputs of my-ls and ls
compare_outputs() {
    local description="$1"
    local cmd_args="$2"

    ./my-ls $cmd_args > my_ls_output.txt
    ls $cmd_args > ls_output.txt

    if diff -q my_ls_output.txt ls_output.txt > /dev/null; then
        echo "[PASS] $description"
    else
        echo "[FAIL] $description"
        echo "my-ls output:"
        cat my_ls_output.txt
        echo "ls output:"
        cat ls_output.txt
    fi
}

# Test cases
compare_outputs "No arguments" ""
compare_outputs "Specific file" "<file_name_here>"
compare_outputs "Specific directory" "<directory_name_here>"
compare_outputs "Long listing format" "-l"
compare_outputs "Long format for specific file" "-l <file_name_here>"
compare_outputs "Long format for specific directory" "-l <directory_name_here>"
compare_outputs "Long format for /usr/bin" "-l /usr/bin"
compare_outputs "Recursive listing" "-R <directory_with_folders_here>"
compare_outputs "All files including hidden" "-a"
compare_outputs "Reverse order" "-r"
compare_outputs "Sort by modification time" "-t"
compare_outputs "Long and all files" "-la"
compare_outputs "Long and sorted by time in specific directory" "-l -t <directory_name_here>"
compare_outputs "Recursive, reverse order in directory with folders" "-lRr <directory_name_here>"
compare_outputs "Combined long and all files with specific file" "-l <directory_name_here> -a <file_name_here>"
compare_outputs "Recursive listing with multiple slashes" "-lR <directory_name_here>///<sub_directory_name_here>/// <directory_name_here>/<sub_directory_name_here>/"
compare_outputs "Long and all files in /dev" "-la /dev"
# Test with time constraints and large directories
compare_outputs "All flags on large directory" "-alRrt <directory_name_here>"
# Test with special file and directory names
mkdir -- "-name_test_dir"
compare_outputs "Directory with '-' as name" "-name_test_dir"
rm -r -- "-name_test_dir"
# Test with symbolic links
touch testfile
ln -s testfile symlinkfile
compare_outputs "Symbolic link file with /" "-l symlinkfile/"
compare_outputs "Symbolic link file" "-l symlinkfile"
mkdir testdir
touch testdir/testfile
ln -s testdir symlinkdir
compare_outputs "Symbolic link directory with /" "-l symlinkdir/"
compare_outputs "Symbolic link directory" "-l symlinkdir"
# Clean up
rm -f testfile symlinkfile
rm -rf testdir symlinkdir

# Colors and performance tests
compare_outputs "Color output" "--color=auto"
time ./my-ls -R ~ > my_ls_output.txt
time ls -R ~ > ls_output.txt

echo "Tests completed."