# my-ls-1

**my-ls-1** is a custom implementation of the `ls` command in Go. It provides functionality to list files and folders in a directory, replicating the behavior of the original `ls` command with additional features.

## Features

- Display files and folders in a directory.
- Recursive listing of subdirectories.
- Include hidden files and folders.
- Sort by various criteria, such as modification time.
- Reverse the order of the listing.

## Installation

To install **my-ls-1**, you need to have Go installed on your system. Then, you can clone the repository and build the executable using the following steps:

1. Clone the repository:

   ```bash
   git clone https://learn.reboot01.com/git/emarei/my-ls-1.git
   ```

2. Change to the project directory:

   ```bash
   cd my-ls-1
   ```

3. Build the executable:

   ```bash
   go build -o my-ls-1
   ```

## Usage

To use **my-ls-1**, you can run the built executable with the desired flags and options. Here are some examples of how to use it:

- List files and folders in the current directory:

  ```bash
  ./my-ls-1
  ```

- List files and folders in a specific directory:

  ```bash
  ./my-ls-1 /path/to/directory
  ```

- List files and folders recursively:

  ```bash
  ./my-ls-1 -R
  ```

- Include hidden files and folders:

  ```bash
  ./my-ls-1 -a
  ```

- Sort by modification time, newest first:

  ```bash
  ./my-ls-1 -t
  ```

- Reverse the order of the listing:

  ```bash
  ./my-ls-1 -r
  ```

- Display long format listing:

  ```bash
  ./my-ls-1 -l
  ```

- Combine multiple flags:

  ```bash
  ./my-ls-1 -l -a -t
  ```

## Authurs
Emran marie (emarei)
Omar Abdulrahim (oabdulra)