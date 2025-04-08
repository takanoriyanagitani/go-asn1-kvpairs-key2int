package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"iter"
	"log"
	"os"
	"slices"
	"strconv"

	mi "github.com/takanoriyanagitani/go-asn1-kvpairs-key2int"
	. "github.com/takanoriyanagitani/go-asn1-kvpairs-key2int/util"
)

var envValByKey func(string) IO[string] = Lift(
	func(key string) (string, error) {
		val, found := os.LookupEnv(key)
		switch found {
		case true:
			return val, nil
		default:
			return "", fmt.Errorf("unknown env var: %s", key)
		}
	},
)

var str2int64 func(string) (int64, error) = mi.ComposeErr(
	strconv.Atoi,
	func(i int) (int64, error) { return int64(i), nil },
)

var originalMapSizeLimit IO[int64] = Bind(
	envValByKey("ENV_MAP_SIZE_LIMIT"),
	Lift(str2int64),
).Or(Of(int64(1048576)))

func ReaderToBytesLimited(limit int64) func(io.Reader) IO[[]byte] {
	return Lift(func(rdr io.Reader) ([]byte, error) {
		limited := &io.LimitedReader{
			R: rdr,
			N: limit,
		}
		var buf bytes.Buffer
		_, e := io.Copy(&buf, limited)
		return buf.Bytes(), e
	})
}

func Stdin2bytesLimited(limit int64) IO[[]byte] {
	return ReaderToBytesLimited(limit)(os.Stdin)
}

func FilenameToBytesLimited(limit int64) func(string) IO[[]byte] {
	return func(filename string) IO[[]byte] {
		return func(ctx context.Context) ([]byte, error) {
			f, e := os.Open(filename)
			if nil != e {
				return nil, e
			}
			defer f.Close()
			return ReaderToBytesLimited(limit)(f)(ctx)
		}
	}
}

var originalMapDerBytes IO[[]byte] = Bind(
	originalMapSizeLimit,
	Stdin2bytesLimited,
)

var originalPairs IO[[]mi.OriginalPair] = Bind(
	originalMapDerBytes,
	Lift(mi.BytesToOriginalPairs),
)

var originalPairIter IO[iter.Seq[mi.OriginalPair]] = Bind(
	originalPairs,
	Lift(func(pairs []mi.OriginalPair) (iter.Seq[mi.OriginalPair], error) {
		return slices.Values(pairs), nil
	}),
)

var str2intMapFilename IO[string] = envValByKey("ENV_STR2INT_MAP_DER_NAME")

var str2intMapDerBytes IO[[]byte] = Bind(
	originalMapSizeLimit,
	func(limit int64) IO[[]byte] {
		return Bind(
			str2intMapFilename,
			FilenameToBytesLimited(limit),
		)
	},
)

var str2intPairs IO[mi.StrToIntPairs] = Bind(
	str2intMapDerBytes,
	Lift(mi.StrToIntPairsFromDerBytes),
)

var str2intMap IO[mi.StrToInt] = Bind(
	str2intPairs,
	Lift(func(p mi.StrToIntPairs) (mi.StrToInt, error) {
		return p.ToMap(), nil
	}),
)

var normalized IO[iter.Seq2[mi.NormalizedPair, error]] = Bind(
	str2intMap,
	func(s2i mi.StrToInt) IO[iter.Seq2[mi.NormalizedPair, error]] {
		return Bind(
			originalPairIter,
			Lift(func(
				i iter.Seq[mi.OriginalPair],
			) (iter.Seq2[mi.NormalizedPair, error], error) {
				var op mi.OriginalPairs = mi.OriginalPairs(i)
				return op.Normalize(s2i), nil
			}),
		)
	},
)

var normalizedDerBytes IO[[]byte] = Bind(
	normalized,
	Lift(func(n iter.Seq2[mi.NormalizedPair, error]) ([]byte, error) {
		var s []mi.NormalizedPair
		for pair, e := range n {
			if nil != e {
				return nil, e
			}
			s = append(s, pair)
		}
		var np mi.NormalizedPairs = mi.NormalizedPairs(s)
		return np.ToDerBytes()
	}),
)

var bytes2stdout func([]byte) IO[Void] = func(der []byte) IO[Void] {
	return func(_ context.Context) (Void, error) {
		_, e := os.Stdout.Write(der)
		return Empty, e
	}
}

var original2normalized2der2stdout IO[Void] = Bind(
	normalizedDerBytes,
	bytes2stdout,
)

func main() {
	_, e := original2normalized2der2stdout(context.Background())
	if nil != e {
		log.Printf("%v\n", e)
	}
}
