# blocks
Read Bitcoin block data (`.dat`) files.

## Usage
Build with `make build`.

`blocks` takes a single `height` argument, which represents the height of the block in the `.dat` file pointed to by the `FILEPATH` environment variable for which to print its 80-byte block header.