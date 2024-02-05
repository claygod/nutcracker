package nutcracker

// Nutcracker
// Problem-based approach
// Changer tests
// Copyright © 2024 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"testing"
)

func TestChangerEasy(t *testing.T) {
	achr := NewAtomicChangerRepository()
	achr.Set(newChangerIncrement(nil, 0.001))

	if achr.counter != 1 {
		t.Errorf("Want 1, have %d", achr.counter)
	}
}

func TestChangerEasy2(t *testing.T) {
	state1 := &State{
		Data: []float64{0.1, 0.1},
	}

	levelFunc := func(maxLevel int) int { // maxLevel для варианта с рандомным уровнем в перспективе
		return 0
	}
	achr := NewAtomicChangerRepository()
	achr.Set(newChangerIncrement(levelFunc, 0.001))

	chID, chCur := achr.GetRandom()
	if chID != 0 {
		t.Errorf("Want 0, have %d", chID)
	}

	state2 := chCur.Change(state1)
	if state2.Data[0] != 0.101 {
		t.Errorf("Want 0.101, have %v", state2.Data[0])
	}
}

// ================== Changer example ==================

func newChangerIncrement(levelFunc func(int) int, delta float64) *ChangerIncrement {
	return &ChangerIncrement{
		levelFunc: levelFunc,
		delta:     delta,
	}
}

type ChangerIncrement struct {
	levelFunc func(int) int
	// level int
	delta float64
}

func (c *ChangerIncrement) Change(stateIn *State) *State {
	level := c.levelFunc(len(stateIn.Data) - 1)

	if len(stateIn.Data) < level { // возвращаем оригинал
		return stateIn
	}

	stateOut := stateIn.Copy()

	stateOut.Data[level] = stateIn.Data[level] + c.delta

	return stateOut
}

func (c *ChangerIncrement) GetInnerSteps() int64 {
	return 1 // NOTE: для базовых чейнжеров как правило это 1, а остальные должны уметь вычислять число шагов
}
