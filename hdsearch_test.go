package main

import (
	"bufio"
	"encoding/hex"
	"log"
	"math/big"
	"math/bits"
	"os"

	"testing"
)

const (
	k         = 128     // bit長
	randsNum  = 1000000 // 乱数の数
	jaccardTh = 0.8     // jaccard係数の閾値
)

// ファイルから100万件の乱数を読み込む
func readRands(path string) [][]byte {
	rands := make([][]byte, 0, randsNum)

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		rand, err := hex.DecodeString(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		rands = append(rands, rand)
	}

	return rands
}

// []byte用の普通のpopcount
func popCount(bytes []byte) int {
	cnt := 0
	for _, b := range bytes {
		cnt += bits.OnesCount8(uint8(b))
	}

	return cnt
}

// 閾値より値が大きくなったら計算を打ち切るpopcount
func popCountKai(bytes []byte, threshold float64) (int, bool) {
	cnt := 0
	for _, b := range bytes {
		cnt += bits.OnesCount8(uint8(b))
		if float64(cnt) > threshold {
			return -1, false
		}
	}

	return cnt, true
}

// 普通のpopcountを用いた場合のベンチマーク
func BenchmarkHdAllSearch1(b *testing.B) {
	rands := readRands(`./rands10000.txt`)
	bb := big.NewInt(0)
	sampleRand := big.NewInt(0).SetBytes([]byte{8, 77, 202, 13, 237, 113, 66, 23, 18, 191, 17, 22, 153, 239, 171, 47})

	// ベンチマーク開始
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range rands {
			bb.SetBytes(v)
			popCnt := popCount(bb.Xor(sampleRand, bb).Bytes())
			jaccard := 2.0 * (float32(k-popCnt)/float32(k) - 0.5)
			if jaccard >= jaccardTh {
				_ = jaccard
			}
		}
	}
}

// 閾値ありpopcountを用いた場合のベンチマーク
func BenchmarkHdAllSearch2(b *testing.B) {
	rands := readRands(`./rands10000.txt`)
	bb := big.NewInt(0)
	sampleRand := big.NewInt(0).SetBytes([]byte{8, 77, 202, 13, 237, 113, 66, 23, 18, 191, 17, 22, 153, 239, 171, 47})

	hammingTh := k * (1 - jaccardTh) * 0.5

	// ベンチマーク開始
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range rands {
			bb.SetBytes(v)
			popCnt, ok := popCountKai(bb.Xor(sampleRand, bb).Bytes(), hammingTh)
			if ok {
				// jaccard
				_ = 2.0 * (float32(k-popCnt)/float32(k) - 0.5)
			}
		}
	}
}
