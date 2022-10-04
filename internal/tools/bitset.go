package tools

import (
    "math/bits"
)

func NewBitset(size int) *Bitset {
    lastBucketSize, size := size&63, size>>6
    if lastBucketSize > 0 {
        size++
    }
    data := make([]uint64, size)
    return &Bitset{
        data:           data,
        lastBucketSize: lastBucketSize,
    }
}

type Bitset struct {
    data           []uint64
    lastBucketSize int
}

func (bs *Bitset) Size() int {
    if len(bs.data) == 0 {
        return 0
    }
    return (len(bs.data)-1)<<6 + bs.lastBucketSize
}

func (bs *Bitset) Xor(other *Bitset) bool {
    for i := range bs.data {
        if i >= len(other.data) {
            break
        }
        if bs.data[i]^other.data[i] != 0 {
            return true
        }
    }
    return false
}

func (bs *Bitset) Get(idx int) bool {
    bucket, ind := idx>>6, uint64(idx&63)
    mask := uint64(1) << ind
    return (bs.data[bucket] & mask) != 0
}

func (bs *Bitset) Fix(idx int) {
    bucket, ind := idx>>6, uint64(idx&63)
    mask := uint64(1) << ind
    bs.data[bucket] |= mask
}

func (bs *Bitset) Unfix(idx int) {
    bucket, ind := idx>>6, uint64(idx&63)
    mask := ^(uint64(1) << ind)
    bs.data[bucket] &= mask
}

func (bs *Bitset) Flip() {
    mask := ^uint64(0)
    for i := range bs.data {
        if i == len(bs.data)-1 && bs.lastBucketSize > 0 {
            mask = (uint64(1) << bs.lastBucketSize) - 1
        }
        bs.data[i] ^= mask
    }
}

func (bs *Bitset) All() bool {
    allOnes := ^uint64(0)
    for i, b := range bs.data {
        if b != allOnes {
            if bs.lastBucketSize == 0 || i != len(bs.data)-1 {
                return false
            }
            if b != (uint64(1)<<bs.lastBucketSize)-1 {
                return false
            }
        }
    }
    return true
}

func (bs *Bitset) Iterate() chan int {
    ch := make(chan int)
    go func() {
        for i := 0; i < (len(bs.data)-1)*64+bs.lastBucketSize; i++ {
            if bs.Get(i) {
                ch <- i
            }
        }
        close(ch)
    }()
    return ch
}

func (bs *Bitset) Count() int {
    cnt := 0
    for _, b := range bs.data {
        cnt += bits.OnesCount64(b)
    }
    return cnt
}
