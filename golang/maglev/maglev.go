package main

import (
	"errors"
	"math/big"
	"sort"
	"sync"

	"github.com/dchest/siphash"
)

func getHashKey(obj string) uint64 {
	return siphash.Hash(0xdeadbabe, 0, []byte(obj))
}

// Maglev
type Maglev struct {
	n           uint64 // Number of backends
	m           uint64 // Size of the lookup table
	permutation [][]uint64
	lookup      []int64
	nodeList    []string
	lock        *sync.RWMutex
}

// NewMaglev
func NewMaglev(backends []string, m uint64) (*Maglev, error) {
	if !big.NewInt(0).SetUint64(m).ProbablyPrime(1) {
		return nil, errors.New("lookup table size is not a prime number")
	}

	maglev := &Maglev{m: m, lock: &sync.RWMutex{}}
	if err := maglev.Set(backends); err != nil {
		return nil, err
	}

	return maglev, nil
}

func (m *Maglev) Set(backends []string) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	n := uint64(len(backends))
	if m.m < n {
		return errors.New("number of backends is greater than lookup table")
	}

	m.nodeList = make([]string, n)
	// Copy to avoid modifying original input afterwards
	copy(m.nodeList, backends)
	m.n = n
	m.generatePopulation()
	m.populate()
	return nil
}

// Remove removes a given backend from the existing list
func (m *Maglev) Remove(backend string) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	index := sort.SearchStrings(m.nodeList, backend)
	if index == len(m.nodeList) {
		return errors.New("backend is not present")
	}

	m.nodeList = append(m.nodeList[:index], m.nodeList[index+1:]...)
	m.n = uint64(len(m.nodeList))
	m.generatePopulation()
	m.populate()
	return nil
}

// Add adds new backend to Maglev
func (m *Maglev) Add(backend string) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	for _, v := range m.nodeList {
		if v == backend {
			return errors.New("backend already exists")
		}
	}

	if m.m == m.n {
		return errors.New("number of backend would be greater than lookup table")
	}

	m.nodeList = append(m.nodeList, backend)
	m.n = uint64(len(m.nodeList))
	m.generatePopulation()
	m.populate()
	return nil
}

// Clear resets Maglev to initial state
func (m *Maglev) Clear() {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.nodeList = nil
	m.permutation = nil
	m.lookup = nil
}

// Get gets node name using object string
func (m *Maglev) Get(obj string) (string, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if len(m.nodeList) == 0 {
		return "", errors.New("backend list is empty")
	}

	key := getHashKey(obj)
	return m.nodeList[m.lookup[key%m.m]], nil
}

func (m *Maglev) generatePopulation() {
	m.permutation = nil
	if len(m.nodeList) == 0 {
		return
	}

	sort.Strings(m.nodeList)

	for i := 0; i < len(m.nodeList); i++ {
		bData := []byte(m.nodeList[i])
		// Hexspeak: https://en.wikipedia.org/wiki/Hexspeak
		offset := siphash.Hash(0xdeadbabe, 0, bData) % m.m
		skip := (siphash.Hash(0xdeadbeef, 0, bData) % (m.m - 1)) + 1

		iRow := make([]uint64, m.m)
		var j uint64
		for j = 0; j < m.m; j++ {
			iRow[j] = (offset + uint64(j)*skip) % m.m
		}

		m.permutation = append(m.permutation, iRow)
	}
}

func (m *Maglev) populate() {
	if len(m.nodeList) == 0 {
		return
	}

	var i, j uint64
	next := make([]uint64, m.n)
	entry := make([]int64, m.m)
	for j = 0; j < m.m; j++ {
		entry[j] = -1
	}

	var n uint64

	for {
		for i = 0; i < m.n; i++ {
			c := m.permutation[i][next[i]]
			for entry[c] >= 0 {
				next[i] = next[i] + 1
				c = m.permutation[i][next[i]]
			}

			entry[c] = int64(i)
			next[i] = next[i] + 1
			n++

			if n == m.m {
				m.lookup = entry
				return
			}
		}
	}
}
