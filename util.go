package zkwasm

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"slices"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

var (
	ErrUnsupportedInputType = errors.New("UnsupportedInputType")
)

func BuildInputsString(input []any) ([]string, error) {
	if len(input) == 0 {
		return nil, nil
	}

	arr := make([]string, 0, len(input))

	for _, i := range input {
		switch ti := i.(type) {
		case json.Number:
			iti, err := ti.Int64()
			if err == nil {
				arr = append(arr, fmt.Sprintf("%d:i64", iti))
				continue
			}

			return nil, ErrUnsupportedInputType
		case int64:
			arr = append(arr, fmt.Sprintf("%d:i64", ti))
		case uint64:
			arr = append(arr, fmt.Sprintf("%d:i64", ti))
		case int:
			arr = append(arr, fmt.Sprintf("%d:i64", ti))
		case uint:
			arr = append(arr, fmt.Sprintf("%d:i64", ti))
		case int32:
			arr = append(arr, fmt.Sprintf("%d:i64", ti))
		case uint32:
			arr = append(arr, fmt.Sprintf("%d:i64", ti))
		case []byte:
			hexStr := hex.EncodeToString(ti)
			arr = append(arr, fmt.Sprintf("0x%s:bytes", hexStr))
		case []uint64:
			buf := make([]byte, 0, len(ti)*8)
			for _, u := range ti {
				buf = binary.LittleEndian.AppendUint64(buf, u)
			}
			hexStr := hex.EncodeToString(buf)
			arr = append(arr, fmt.Sprintf("0x%s:bytes-packed", hexStr))
		default:
			return nil, ErrUnsupportedInputType
		}
	}

	return arr, nil
}

func ChunkSlice[T any](slice []T, chunkSize int) [][]T {
	var chunks = [][]T{}
	for {
		if len(slice) == 0 {
			break
		}

		// necessary check to avoid slicing beyond
		// slice capacity
		if len(slice) < chunkSize {
			chunkSize = len(slice)
		}

		chunks = append(chunks, slice[0:chunkSize])
		slice = slice[chunkSize:]
	}

	return chunks
}

func ByteSliceToBigIntSlice(input []byte, le bool) []*big.Int {
	chunkSize := 256 / 8

	chunks := ChunkSlice(input, chunkSize)

	output := make([]*big.Int, len(chunks))

	if len(chunks) == 0 {
		return output
	}

	for i := range chunks {
		bi := new(big.Int)
		if !le {
			bi.SetBytes(chunks[i])
		} else {
			tmp := make([]byte, chunkSize)
			copy(tmp, chunks[i])
			slices.Reverse(tmp)
			bi.SetBytes(tmp)
		}
		output[i] = bi
	}

	return output
}

func ParseInputsString(inputs []string) ([]uint64, error) {
	re := make([]uint64, 0, len(inputs))

	for _, i := range inputs {
		fmtErr := fmt.Errorf("illegal input string: %s", i)

		arr := strings.SplitN(i, ":", 2)
		if len(arr) != 2 {
			return nil, fmtErr
		}

		v := arr[0]
		t := arr[1]

		switch t {
		case "i64":
			if strings.HasPrefix(v, "0x") {
				d, err := hexutil.DecodeUint64(v)
				if err != nil {
					return nil, fmtErr
				}
				re = append(re, d)
			} else {
				d, err := strconv.ParseUint(v, 10, 64)
				if err != nil {
					return nil, fmtErr
				}
				re = append(re, d)
			}
		case "bytes":
			if !strings.HasPrefix(v, "0x") {
				return nil, fmtErr
			}
			bytes, err := hexutil.Decode(v)
			if err != nil {
				return nil, fmtErr
			}
			for _, b := range bytes {
				re = append(re, uint64(b))
			}
		case "bytes-packed":
			if !strings.HasPrefix(v, "0x") {
				return nil, fmtErr
			}
			bytes, err := hexutil.Decode(v)
			if err != nil {
				return nil, fmtErr
			}

			bytesArr := ChunkSlice(bytes, 8)
			for _, n := range bytesArr {
				re = append(re, binary.LittleEndian.Uint64(n))
			}
		// TODO: support "file"
		default:
			return nil, fmtErr
		}
	}

	return re, nil
}
