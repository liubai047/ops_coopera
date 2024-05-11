// Package bloomFilter 布隆过滤器
package bloomFilter

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/binary"
	"hash"
	"math"
	"math/rand"
	"strconv"
	"sync"
)

type bloom struct {
	bitNum   uint64
	hashFunc []func(val []byte) uint64 // 哈希函数数组
	bit      []uint64                  // bloom位数组
	lock     sync.RWMutex              // 读写锁
	num      int                       // 布隆中元素个数
}

var hashBaseFunc = []hash.Hash{md5.New(), sha1.New(), sha256.New()}
var minHashNum = uint64(len(hashBaseFunc))
var maxHashNum = uint64(15)

/*
NewBloom
param maxN 预计元素数量
param p 期望误差概率
*/
func NewBloom(maxN uint64, p float64) *bloom {
	var (
		bitNum  = optimalM(maxN, p)
		hashNum = optimalK(bitNum, maxN)
		bloom   = &bloom{
			bitNum:   bitNum,
			hashFunc: make([]func(val []byte) uint64, 0),
			bit:      make([]uint64, (bitNum+63)/64),
			lock:     sync.RWMutex{},
			num:      0,
		}
	)
	if hashNum < minHashNum {
		hashNum = minHashNum
	}
	if hashNum > maxHashNum {
		hashNum = maxHashNum
	}
	bloom.baseHash()
	for _, salt := range createSalt(hashNum - minHashNum) {
		bloom.hashFunc = append(bloom.hashFunc, bloom.createHash(salt))
	}
	return bloom
}

// 生成hash函数createHash
func (b *bloom) createHash(salt []byte) func(val []byte) uint64 {
	return func(val []byte) uint64 {
		hs := sha512.New()
		res := hs.Sum(append(val, salt...))
		return binary.BigEndian.Uint64(res[:8]) % b.bitNum
	}
}

// Check 判断数据是否在布隆过滤器中(存在true 不存在false)
func (b *bloom) Check(val []byte) bool {
	for _, hs := range b.hashFunc {
		hsVal := hs(val)
		k := hsVal / 64  // 对应数组下标在哪
		mk := hsVal % 64 // 在对应数组下标的第几个偏移量中
		b.lock.RLock()
		var t = b.bit[k] & (1 << mk)
		b.lock.RUnlock()
		if t == 0 { // 找到bitmap中对应下标，对应偏移量的位，判断是否0
			return false
		}
	}
	return true
}

// Add 数据插入布隆过滤器(不可删除) error is always return nil
func (b *bloom) Add(val []byte) error {
	for _, hs := range b.hashFunc {
		hsVal := hs(val)
		k := hsVal / 64  // 对应数组下标在哪
		mk := hsVal % 64 // 在对应数组下标的第几个偏移量中
		b.lock.Lock()
		b.bit[k] = b.bit[k] | (1 << mk)
		b.lock.Unlock()
	}
	return nil
}

// 基础的hash函数(不加盐)
func (b *bloom) baseHash() {
	for _, hs := range hashBaseFunc {
		var f = func(val []byte) uint64 {
			res := hs.Sum(val)
			return binary.BigEndian.Uint64(res[:8]) % b.bitNum
		}
		b.hashFunc = append(b.hashFunc, f)
	}
}

// 计算过滤器-所需的位数量
func optimalM(maxN uint64, p float64) uint64 {
	return uint64(math.Ceil(-float64(maxN) * math.Log(p) / (math.Ln2 * math.Ln2)))
}

// 计算过滤器-所需的最大哈希函数
func optimalK(m, maxN uint64) uint64 {
	return uint64(math.Ceil(float64(m) * math.Ln2 / float64(maxN)))
}

// 生成指定位数的salt
func createSalt(num uint64) [][]byte {
	var res = make([][]byte, 0)
	for i := uint64(0); i < num; i++ {
		res = append(res, []byte(string(rand.Int31())+strconv.FormatUint(i, 10)))
	}
	return res
}
