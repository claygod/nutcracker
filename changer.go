package nutcracker

import (
	"fmt"
	"math/rand"
	"sync"
)

// Nutcracker
// Problem-based approach
// Atomic changer (implementation)
// Copyright © 2022-2024 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

const errID int64 = 9e+18 // 9000000000000000000 можно 9223372036854775807

type AtomicChangerOverRepository struct {
	baseRepo  AtomicChangerRepo // основное ререпо, доступное всем
	innerRepo AtomicChangerRepo // кастомное репо
}

func NewAtomicChangerOverRepository(baseRepo AtomicChangerRepo) *AtomicChangerOverRepository {
	return &AtomicChangerOverRepository{
		baseRepo:  baseRepo,
		innerRepo: NewAtomicChangerRepository(),
	}
}

func (a *AtomicChangerOverRepository) GetByTarget(targetState *State) (ID int64, aChanger AtomicChanger) {
	panic("implement me")

	return 0, nil
}

func (a *AtomicChangerOverRepository) GetRandom() (int64, AtomicChanger) {
	if rand.Int63n(1) == 0 {
		return a.baseRepo.GetRandom()
	}

	iID, iAch := a.innerRepo.GetRandom()

	return -iID, iAch
}

func (a *AtomicChangerOverRepository) GetByID(id int64) (AtomicChanger, bool) {
	if id >= 0 {
		return a.baseRepo.GetByID(id)
	}

	return a.innerRepo.GetByID(-id)
}

func (a *AtomicChangerOverRepository) Set(aChanger AtomicChanger, opts ...bool) int64 {
	if len(opts) == 0 || opts[0] == false {
		return a.baseRepo.Set(aChanger)
	}

	id := a.innerRepo.Set(aChanger)

	if id != errID {
		id *= -1
	}

	return id
}

/*
AtomicChangerRepository - имплементация интерфейса AtomicChangerRepo
*/
type AtomicChangerRepository struct {
	m       sync.Mutex
	counter int64
	names   map[string]struct{}
	data    map[int64]AtomicChanger
}

func NewAtomicChangerRepository() *AtomicChangerRepository {
	return &AtomicChangerRepository{
		m:     sync.Mutex{},
		names: make(map[string]struct{}),
		data:  make(map[int64]AtomicChanger),
	}
}

/*
GetByTarget - берём наиболее подходящее (близкое) по цели
*/
func (a *AtomicChangerRepository) GetByTarget(targetState *State) (ID int64, aChanger AtomicChanger) {
	a.m.Lock()
	defer a.m.Unlock()

	panic("implement me")

	return 0, nil // TODO: надо имплементировать, если потребуется
}

/*
GetRandom - берём случайную, это удобно для генерации случайного Chainlet-набора
*/
func (a *AtomicChangerRepository) GetRandom() (int64, AtomicChanger) {
	a.m.Lock()
	defer a.m.Unlock()

	if a.counter == 0 {
		return 0, nil
	}

	id := rand.Int63n(int64(a.counter))

	return id, a.data[id]
}

func (a *AtomicChangerRepository) GetByID(id int64) (AtomicChanger, bool) {
	a.m.Lock()
	defer a.m.Unlock()

	ch, ok := a.data[id]

	return ch, ok
}

/*
Set - сначала добавляем действительно базовые возможности, а потом можно добавлять
Chainlet-наборы, которые используются часто или которые короткие но эффективные
*/
func (a *AtomicChangerRepository) Set(aChanger AtomicChanger, _ ...bool) (ID int64) {
	a.m.Lock()
	defer a.m.Unlock()

	if _, ok := a.names[aChanger.GetName()]; ok {
		return errID // такой чейнжер есть (возможный вариант для синтетических ченжеров)
	}

	a.counter++
	a.names[aChanger.GetName()] = struct{}{}
	a.data[a.counter-1] = aChanger

	return a.counter
}

func newAtomicChangerSyntheticFromChainlet(ch Chainlet, chRepo AtomicChangerRepo) (*AtomicChangerSynthetic, error) {
	chList := make([]AtomicChanger, 0, ch.countSteps)

	// перебираем чтобы получить вместо иденитификаторов сами чейнжеры
	// кроме того, это защищает от вероятности, что в цепочке будет что-то неидентифицируемое
	for _, id := range ch.ChainIDs {
		if ch, ok := chRepo.GetByID(id); ok {
			chList = append(chList, ch)
		} else {
			return nil, fmt.Errorf("changer %d not found in repo", id)
		}
	}

	return &AtomicChangerSynthetic{
		name:   ch.GetMultiName(),
		steps:  ch.GetCountSteps(),
		chList: chList,
	}, nil
}

type AtomicChangerSynthetic struct {
	name   string
	steps  int64
	chList []AtomicChanger
}

func (a *AtomicChangerSynthetic) Change(stIn *State) *State {
	stOut := stIn.Copy()

	for _, ch := range a.chList {
		stOut = ch.Change(stOut)
	}

	return stOut
}

func (a *AtomicChangerSynthetic) GetInnerSteps() int64 {
	return a.steps
}

func (a *AtomicChangerSynthetic) GetName() string {
	return a.name
}
