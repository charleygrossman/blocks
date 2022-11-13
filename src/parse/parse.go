package parse

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

type BlockFileParser struct {
	file          *os.File
	endian        binary.ByteOrder
	magicBytes    []byte
	headerOffsets []int64
}

func NewBlockFileParser(filepath string, endian binary.ByteOrder) (*BlockFileParser, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file at %s: %v", filepath, err)
	}
	p := &BlockFileParser{
		file:   f,
		endian: endian,
	}
	if err := p.setMagicBytes(); err != nil {
		return nil, err
	}
	if err := p.setHeaderOffsets(); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *BlockFileParser) BlockHeader(height uint64) ([]byte, error) {
	if height > uint64(len(p.headerOffsets)-1) {
		return nil, fmt.Errorf("height higher than block count: max=%d got=%d", len(p.headerOffsets)-1, height)
	}
	header := make([]byte, 80)
	n, err := p.file.ReadAt(header, p.headerOffsets[height])
	if err != nil {
		return nil, fmt.Errorf("failed to read block header at height %d: %v", height, err)
	}
	if n != 80 {
		return nil, fmt.Errorf("failed to read block header at height %d: expect 80 bytes, found %d", height, n)
	}
	return header, nil
}

func (p *BlockFileParser) BlockCount() uint64 {
	return uint64(len(p.headerOffsets))
}

func (p *BlockFileParser) MagicBytes() []byte {
	return p.magicBytes
}

func (p *BlockFileParser) setMagicBytes() error {
	bts := make([]byte, 4)
	n, err := p.file.ReadAt(bts, 0)
	if err != nil {
		return fmt.Errorf("failed to set magic bytes: %v", err)
	}
	if n != 4 {
		return fmt.Errorf("failed to set magic bytes: expect 4 bytes at beginning of file, found %d", n)
	}
	p.magicBytes = bts
	return nil
}

// TODO: Error handling for reads.
func (p *BlockFileParser) setHeaderOffsets() error {
	var (
		wantMbts = fmt.Sprintf("%x", p.magicBytes)
		offset   int64
		offsets  = make([]int64, 0, 100)
	)
	for {
		mbts := make([]byte, 4)
		p.file.ReadAt(mbts, offset)
		if fmt.Sprintf("%x", mbts) != wantMbts {
			break
		}
		offset += 4

		sbts := make([]byte, 4)
		p.file.ReadAt(sbts, offset)
		size, err := p.parseInt32(sbts)
		if err != nil {
			return fmt.Errorf("failed to set header offsets: %v", err)
		}
		offset += 4
		offsets = append(offsets, offset)
		offset += int64(size)
	}
	p.headerOffsets = offsets
	return nil
}

func (p *BlockFileParser) parseInt32(b []byte) (int32, error) {
	if len(b) != 4 {
		return 0, fmt.Errorf("unexpected byte count for int32: want=4 got=%d", len(b))
	}
	var result int32
	if err := binary.Read(bytes.NewReader(b), p.endian, &result); err != nil {
		return 0, fmt.Errorf("failed to parse bytes as int32: %v", err)
	}
	return result, nil
}
