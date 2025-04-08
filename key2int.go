package key2int

import (
	"encoding/asn1"
	"errors"
	"fmt"
	"iter"
	"maps"
)

var ErrNotFound error = errors.New("not found")

type OriginalPair struct {
	Key string `asn1:"utf8"`
	Val string `asn1:"utf8"`
}

type NormalizedPair struct {
	Key int64
	Val string `asn1:"utf8"`
}

type StrToInt map[string]int64

type StrToIntPairs iter.Seq2[string, int64]

type StrToIntPair struct {
	Key        string `asn1:"utf8"`
	MapdSerial int64
}

func (p StrToIntPairs) ToMap() StrToInt {
	return maps.Collect(iter.Seq2[string, int64](p))
}

func StrToIntPairsFromDerBytes(der []byte) (StrToIntPairs, error) {
	var pairs []StrToIntPair
	_, e := asn1.Unmarshal(der, &pairs)
	var i iter.Seq2[string, int64] = func(
		yield func(string, int64) bool,
	) {
		for _, pair := range pairs {
			var originalKey string = pair.Key
			var mapdKey int64 = pair.MapdSerial
			if !yield(originalKey, mapdKey) {
				return
			}
		}
	}
	return StrToIntPairs(i), e
}

func (c StrToInt) Normalize(o OriginalPair) (NormalizedPair, error) {
	var rawKey string = o.Key
	key, found := c[rawKey]

	ret := NormalizedPair{}
	ret.Key = key
	ret.Val = o.Val

	switch found {
	case true:
		return ret, nil
	default:
		return ret, fmt.Errorf("%w: %s", ErrNotFound, rawKey)
	}
}

func BytesToOriginalPairs(der []byte) ([]OriginalPair, error) {
	var ret []OriginalPair
	_, e := asn1.Unmarshal(der, &ret)
	return ret, e
}

type OriginalPairs iter.Seq[OriginalPair]

func (o OriginalPairs) Normalize(
	s2i StrToInt,
) iter.Seq2[NormalizedPair, error] {
	return func(yield func(NormalizedPair, error) bool) {
		for original := range o {
			normalized, e := s2i.Normalize(original)
			if !yield(normalized, e) {
				return
			}
		}
	}
}

type NormalizedPairs []NormalizedPair

func (n NormalizedPairs) ToDerBytes() ([]byte, error) {
	var raw []NormalizedPair = n
	return asn1.Marshal(raw)
}
