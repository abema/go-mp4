package bufio

import (
	"errors"
	"io"
	"sort"
)

type cacheBlock struct {
	offset  int64
	data    []byte
	lastRef int
}

type ReadSeeker struct {
	reader      io.ReadSeeker
	blockSize   int
	historySize int
	blocks      []*cacheBlock
	curr        int64
	readCnt     int
}

func NewReadSeeker(r io.ReadSeeker, blockSize, historySize int) *ReadSeeker {
	blocks := make([]*cacheBlock, historySize)
	for i := range blocks {
		blocks[i] = &cacheBlock{
			offset: -1,
			data:   make([]byte, blockSize),
		}
	}
	return &ReadSeeker{
		reader:      r,
		blockSize:   blockSize,
		historySize: historySize,
		blocks:      blocks,
	}
}

func (r *ReadSeeker) Read(p []byte) (int, error) {
	r.readCnt++
	for _, block := range r.blocks {
		if block.offset < 0 {
			continue
		}
		if r.curr >= block.offset && r.curr < block.offset+int64(len(block.data)) {
			block.lastRef = r.readCnt
			r.sortBlocks()
			n, err := r.readFromBlock(p, block)
			return n, err
		}
	}
	block, err := r.cacheCurrentBlock()
	if err != nil {
		return 0, err
	}
	block.lastRef = r.readCnt
	r.sortBlocks()
	return r.readFromBlock(p, block)
}

func (r *ReadSeeker) readFromBlock(p []byte, block *cacheBlock) (int, error) {
	size := int(block.offset + int64(len(block.data)) - r.curr)
	if size == 0 {
		return 0, io.EOF
	}
	if size > len(p) {
		size = len(p)
	}
	begin := int(r.curr - block.offset)
	end := begin + size
	copy(p, block.data[begin:end])
	r.curr += int64(size)
	return size, nil
}

func (r *ReadSeeker) cacheCurrentBlock() (*cacheBlock, error) {
	block := r.blocks[len(r.blocks)-1]
	offset := r.curr - r.curr%int64(r.blockSize)
	data := block.data[:cap(block.data)]
	if _, err := r.reader.Seek(offset, io.SeekStart); err != nil {
		return nil, err
	}
	n, err := io.ReadFull(r.reader, data)
	if err == io.ErrUnexpectedEOF {
		data = data[:n]
	} else if err != nil {
		return nil, err
	}
	block.offset = r.curr - r.curr%int64(r.blockSize)
	block.data = data
	return block, nil
}

func (r *ReadSeeker) sortBlocks() {
	sort.Slice(r.blocks, func(i, j int) bool {
		return r.blocks[i].offset > r.blocks[j].offset
	})
}

func (r *ReadSeeker) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		r.curr = offset
	case io.SeekCurrent:
		r.curr += offset
	case io.SeekEnd:
		end, err := r.reader.Seek(0, io.SeekEnd)
		if err != nil {
			return 0, err
		}
		r.curr = end + offset
	default:
		return 0, errors.New("unknown whence")
	}
	return r.curr, nil
}
