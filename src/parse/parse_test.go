package parse

import (
	"encoding/binary"
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestParser(t *testing.T) {
	suite.Run(t, new(ParserTest))
}

type ParserTest struct {
	suite.Suite
	assert *assert.Assertions
	parser *BlockFileParser
}

func (t *ParserTest) SetupSuite() {
	t.assert = assert.New(t.T())

	cwd, err := os.Getwd()
	t.Suite.NoError(err)
	parser, err := NewBlockFileParser(
		path.Join(cwd, "testdata/blk00000.dat"),
		binary.LittleEndian,
	)
	t.Suite.NoError(err)
	t.parser = parser
}

func (t *ParserTest) Test_MagicBytes() {
	t.assert.Equal("f9beb4d9", fmt.Sprintf("%x", t.parser.MagicBytes()))
}

func (t *ParserTest) Test_headerOffsets() {
	t.assert.Equal(uint64(120018), t.parser.BlockCount())
}

func (t *ParserTest) Test_BlockHeader() {
	want1 := "01000000000000000000000000" +
		"00000000000000000000000000000000000" +
		"000000000003ba3edfd7a7b12b27ac72c3e" +
		"67768f617fc81bc3888a51323a9fb8aa4b1" +
		"e5e4a29ab5f49ffff001d1dac2b7c"
	got1, err := t.parser.BlockHeader(0)
	t.assert.NoError(err)
	t.assert.Equal(want1, fmt.Sprintf("%x", got1))

	want2 := "010000006fe28c0ab6f1b372c1" +
		"a6a246ae63f74f931e8365e15a089c68d6" +
		"190000000000982051fd1e4ba744bbbe68" +
		"0e1fee14677ba1a3c3540bf7b1cdb606e8" +
		"57233e0e61bc6649ffff001d01e36299"
	got2, err := t.parser.BlockHeader(1)
	t.assert.NoError(err)
	t.assert.Equal(want2, fmt.Sprintf("%x", got2))
}
