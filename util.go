package zkwasm

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"slices"
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
