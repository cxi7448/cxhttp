package clIP2Region

import (
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var shiftIndex = []int{24, 16, 8, 0}

const (
	HeaderInfoLength      = 256
	VectorIndexRows       = 256
	VectorIndexCols       = 256
	VectorIndexSize       = 8
	SegmentIndexBlockSize = 14
)

func CheckIP(ip string) (uint32, error) {
	var ps = strings.Split(strings.TrimSpace(ip), ".")
	if len(ps) != 4 {
		return 0, fmt.Errorf("invalid ip address `%s`", ip)
	}

	var val = uint32(0)
	for i, s := range ps {
		d, err := strconv.Atoi(s)
		if err != nil {
			return 0, fmt.Errorf("the %dth part `%s` is not an integer", i, s)
		}

		if d < 0 || d > 255 {
			return 0, fmt.Errorf("the %dth part `%s` should be an integer bettween 0 and 255", i, s)
		}

		val |= uint32(d) << uint32(shiftIndex[i])
	}

	// convert the ip to integer
	return val, nil
}

// --- Index policy define

type IndexPolicy int

const (
	VectorIndexPolicy IndexPolicy = 1
	BTreeIndexPolicy  IndexPolicy = 2
)

// LoadHeader load the header info from the specified handle
func LoadHeader(handle *os.File) (*Header, error) {
	_, err := handle.Seek(0, 0)
	if err != nil {
		return nil, fmt.Errorf("seek to the header: %w", err)
	}

	var buff = make([]byte, HeaderInfoLength)
	rLen, err := handle.Read(buff)
	if err != nil {
		return nil, err
	}

	if rLen != len(buff) {
		return nil, fmt.Errorf("incomplete read: readed bytes should be %d", len(buff))
	}

	return NewHeader(buff)
}

// LoadHeaderFromFile load header info from the specified db file path
func LoadHeaderFromFile(dbFile string) (*Header, error) {
	handle, err := os.OpenFile(dbFile, os.O_RDONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf("open xdb file `%s`: %w", dbFile, err)
	}

	defer func(handle *os.File) {
		_ = handle.Close()
	}(handle)

	header, err := LoadHeader(handle)
	if err != nil {
		return nil, err
	}

	return header, nil
}

// LoadHeaderFromBuff wrap the header info from the content buffer
func LoadHeaderFromBuff(cBuff []byte) (*Header, error) {
	return NewHeader(cBuff[0:256])
}

// LoadVectorIndex util function to load the vector index from the specified file handle
func LoadVectorIndex(handle *os.File) ([]byte, error) {
	// load all the vector index block
	_, err := handle.Seek(HeaderInfoLength, 0)
	if err != nil {
		return nil, fmt.Errorf("seek to vector index: %w", err)
	}

	var buff = make([]byte, VectorIndexRows*VectorIndexCols*VectorIndexSize)
	rLen, err := handle.Read(buff)
	if err != nil {
		return nil, err
	}

	if rLen != len(buff) {
		return nil, fmt.Errorf("incomplete read: readed bytes should be %d", len(buff))
	}

	return buff, nil
}

// LoadVectorIndexFromFile load vector index from a specified file path
func LoadVectorIndexFromFile(dbFile string) ([]byte, error) {
	handle, err := os.OpenFile(dbFile, os.O_RDONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf("open xdb file `%s`: %w", dbFile, err)
	}

	defer func() {
		_ = handle.Close()
	}()

	vIndex, err := LoadVectorIndex(handle)
	if err != nil {
		return nil, err
	}

	return vIndex, nil
}

// LoadContent load the whole xdb content from the specified file handle
func LoadContent(handle *os.File) ([]byte, error) {
	// get file size
	fi, err := handle.Stat()
	if err != nil {
		return nil, fmt.Errorf("stat: %w", err)
	}

	size := fi.Size()

	// seek to the head of the file
	_, err = handle.Seek(0, 0)
	if err != nil {
		return nil, fmt.Errorf("seek to get xdb file length: %w", err)
	}

	var buff = make([]byte, size)
	rLen, err := handle.Read(buff)
	if err != nil {
		return nil, err
	}

	if rLen != len(buff) {
		return nil, fmt.Errorf("incomplete read: readed bytes should be %d", len(buff))
	}

	return buff, nil
}

// LoadContentFromFile load the whole xdb content from the specified db file path
func LoadContentFromFile(dbFile string) ([]byte, error) {
	handle, err := os.OpenFile(dbFile, os.O_RDONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf("open xdb file `%s`: %w", dbFile, err)
	}

	defer func() {
		_ = handle.Close()
	}()

	cBuff, err := LoadContent(handle)
	if err != nil {
		return nil, err
	}

	return cBuff, nil
}

func (i IndexPolicy) String() string {
	switch i {
	case VectorIndexPolicy:
		return "VectorIndex"
	case BTreeIndexPolicy:
		return "BtreeIndex"
	default:
		return "unknown"
	}
}

// --- Header define

type Header struct {
	// data []byte
	Version       uint16
	IndexPolicy   IndexPolicy
	CreatedAt     uint32
	StartIndexPtr uint32
	EndIndexPtr   uint32
}

func NewHeader(input []byte) (*Header, error) {
	if len(input) < 16 {
		return nil, fmt.Errorf("invalid input buffer")
	}

	return &Header{
		Version:       binary.LittleEndian.Uint16(input),
		IndexPolicy:   IndexPolicy(binary.LittleEndian.Uint16(input[2:])),
		CreatedAt:     binary.LittleEndian.Uint32(input[4:]),
		StartIndexPtr: binary.LittleEndian.Uint32(input[8:]),
		EndIndexPtr:   binary.LittleEndian.Uint32(input[12:]),
	}, nil
}

// --- searcher implementation

type Searcher struct {
	handle *os.File

	// header info
	header  *Header
	ioCount int

	// use it only when this feature enabled.
	// Preload the vector index will reduce the number of IO operations
	// thus speedup the search process
	vectorIndex []byte

	// content buffer.
	// running with the whole xdb file cached
	contentBuff []byte
}

func baseNew(dbFile string, vIndex []byte, cBuff []byte) (*Searcher, error) {
	var err error

	// content buff first
	if cBuff != nil {
		return &Searcher{
			vectorIndex: nil,
			contentBuff: cBuff,
		}, nil
	}

	// open the xdb binary file
	handle, err := os.OpenFile(dbFile, os.O_RDONLY, 0600)
	if err != nil {
		return nil, err
	}

	return &Searcher{
		handle:      handle,
		vectorIndex: vIndex,
	}, nil
}

func NewWithFileOnly(dbFile string) (*Searcher, error) {
	return baseNew(dbFile, nil, nil)
}

func NewWithVectorIndex(dbFile string, vIndex []byte) (*Searcher, error) {
	return baseNew(dbFile, vIndex, nil)
}

func NewWithBuffer(cBuff []byte) (*Searcher, error) {
	return baseNew("", nil, cBuff)
}

func (s *Searcher) Close() {
	if s.handle != nil {
		err := s.handle.Close()
		if err != nil {
			return
		}
	}
}

// GetIOCount return the global io count for the last search
func (s *Searcher) GetIOCount() int {
	return s.ioCount
}

// SearchByStr find the region for the specified ip string
func (s *Searcher) SearchByStr(str string) (string, error) {
	ip, err := CheckIP(str)
	if err != nil {
		return "", err
	}

	return s.Search(ip)
}

// Search find the region for the specified long ip
func (s *Searcher) Search(ip uint32) (string, error) {
	// reset the global ioCount
	s.ioCount = 0

	// locate the segment index block based on the vector index
	var il0 = (ip >> 24) & 0xFF
	var il1 = (ip >> 16) & 0xFF
	var idx = il0*VectorIndexCols*VectorIndexSize + il1*VectorIndexSize
	var sPtr, ePtr = uint32(0), uint32(0)
	if s.vectorIndex != nil {
		sPtr = binary.LittleEndian.Uint32(s.vectorIndex[idx:])
		ePtr = binary.LittleEndian.Uint32(s.vectorIndex[idx+4:])
	} else if s.contentBuff != nil {
		sPtr = binary.LittleEndian.Uint32(s.contentBuff[HeaderInfoLength+idx:])
		ePtr = binary.LittleEndian.Uint32(s.contentBuff[HeaderInfoLength+idx+4:])
	} else {
		// read the vector index block
		var buff = make([]byte, VectorIndexSize)
		err := s.read(int64(HeaderInfoLength+idx), buff)
		if err != nil {
			return "", fmt.Errorf("read vector index block at %d: %w", HeaderInfoLength+idx, err)
		}

		sPtr = binary.LittleEndian.Uint32(buff)
		ePtr = binary.LittleEndian.Uint32(buff[4:])
	}

	// fmt.Printf("sPtr=%d, ePtr=%d", sPtr, ePtr)

	// binary search the segment index to get the region
	var dataLen, dataPtr = 0, uint32(0)
	var buff = make([]byte, SegmentIndexBlockSize)
	var l, h = 0, int((ePtr - sPtr) / SegmentIndexBlockSize)
	for l <= h {
		m := (l + h) >> 1
		p := sPtr + uint32(m*SegmentIndexBlockSize)
		err := s.read(int64(p), buff)
		if err != nil {
			return "", fmt.Errorf("read segment index at %d: %w", p, err)
		}

		// decode the data step by step to reduce the unnecessary operations
		sip := binary.LittleEndian.Uint32(buff)
		if ip < sip {
			h = m - 1
		} else {
			eip := binary.LittleEndian.Uint32(buff[4:])
			if ip > eip {
				l = m + 1
			} else {
				dataLen = int(binary.LittleEndian.Uint16(buff[8:]))
				dataPtr = binary.LittleEndian.Uint32(buff[10:])
				break
			}
		}
	}

	//fmt.Printf("dataLen: %d, dataPtr: %d", dataLen, dataPtr)
	if dataLen == 0 {
		return "", nil
	}

	// load and return the region data
	var regionBuff = make([]byte, dataLen)
	err := s.read(int64(dataPtr), regionBuff)
	if err != nil {
		return "", fmt.Errorf("read region at %d: %w", dataPtr, err)
	}

	return string(regionBuff), nil
}

// do the data read operation based on the setting.
// content buffer first or will read from the file.
// this operation will invoke the Seek for file based read.
func (s *Searcher) read(offset int64, buff []byte) error {
	if s.contentBuff != nil {
		cLen := copy(buff, s.contentBuff[offset:])
		if cLen != len(buff) {
			return fmt.Errorf("incomplete read: readed bytes should be %d", len(buff))
		}
	} else {
		_, err := s.handle.Seek(offset, 0)
		if err != nil {
			return fmt.Errorf("seek to %d: %w", offset, err)
		}

		s.ioCount++
		rLen, err := s.handle.Read(buff)
		if err != nil {
			return fmt.Errorf("handle read: %w", err)
		}

		if rLen != len(buff) {
			return fmt.Errorf("incomplete read: readed bytes should be %d", len(buff))
		}
	}

	return nil
}
