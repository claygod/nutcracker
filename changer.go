package nutcracker

import (
	"math/rand"
	"sync"
)

// Nutcracker
// Problem-based approach
// Atomic changer (implementation)
// Copyright © 2022 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

/*
AtomicChangerRepository - имплементация интерфейса AtomicChangerRepo
*/
type AtomicChangerRepository struct {
	m       sync.Mutex
	counter uint64
	data    map[uint64]AtomicChanger
}

func NewAtomicChangerRepository() *AtomicChangerRepository {
	return &AtomicChangerRepository{
		data: make(map[uint64]AtomicChanger),
	}
}

/*
   GetRandom - берём случайную, это удобно для генерации случайного Chainlet-набора
*/
func (a *AtomicChangerRepository) GetRandom() (ID uint64, aChanger AtomicChanger) {
	a.m.Lock()
	defer a.m.Unlock()

	if a.counter == 0 {
		return 0, nil
	}

	id := uint64(rand.Int63n(int64(a.counter))-1) + 1

	return id, a.data[id]
}

/*
	Set - сначала добавляем действительно базовые возможности, а потом можно добавлять
	Chainlet-наборы, которые используются часто или которые короткие но эффективные
*/
func (a *AtomicChangerRepository) Set(aChanger AtomicChanger) (ID uint64) {
	a.m.Lock()
	defer a.m.Unlock()

	a.counter++
	a.data[a.counter] = aChanger

	return a.counter
}
