package cmd

import (
	"blocks/src/parse"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	Cmd = cobra.Command{
		Use:   "blocks",
		Short: "print the header of a bitcoin block data file at the given height",
		RunE:  runCmd,
	}
	endian string
	height uint64
)

func init() {
	Cmd.PersistentFlags().StringVarP(&endian, "endian", "e", "l", "endian byte order (choices: b (big), l (little))")
	Cmd.PersistentFlags().Uint64Var(&height, "height", 0, "height of block to print header")
}

func runCmd(c *cobra.Command, args []string) error {
	fpath, ok := os.LookupEnv("FILEPATH")
	if !ok {
		return errors.New("must set FILEPATH")
	}
	var e binary.ByteOrder
	switch strings.ToLower(strings.TrimSpace(endian)) {
	case "b":
		e = binary.BigEndian
	case "l":
		e = binary.LittleEndian
	default:
		return fmt.Errorf("invalid endian: want=[b,l] got=%s", endian)
	}
	parser, err := parse.NewBlockFileParser(fpath, e)
	if err != nil {
		return err
	}
	h, err := parser.BlockHeader(height)
	if err != nil {
		return err
	}
	fmt.Printf("%x\n", h)
	return nil
}
