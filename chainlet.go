package nutcracker

// Nutcracker
// Problem-based approach
// Chainlet
// Copyright © 2022 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"sort"
	"sync"
)

type Chainlet struct { // цепочка действий имеющая удовленворяющий результат (смысл)
	//ID uint64 // возможно снаружи
	// Rate float64
	Chain []uint64 // храним идентификаторы а не ссылки чтобы сравнивать цепочки на похожесть
}

func NewChainlet() *Chainlet {
	return &Chainlet{
		Chain: make([]uint64, 0),
	}
}

func (c *Chainlet) Add(chID uint64) {
	c.Chain = append(c.Chain, chID)
}

func (c *Chainlet) MergeChainlet(ch *Chainlet) {
	c.Chain = append(c.Chain, ch.Chain...)
}

/*
hainletContainer - контейнер нужен для того, чтобы иметь возможность сравнить
*/
type ChainletContainer struct {
	// ID uint64 // возможно снаружи
	Rate     float64 // исчисляется исходя не только из коэффициэнта сравнения state, но и длины цепочки (количества действий)
	Chainlet *Chainlet
}

func MergeChainletContainers(c1 *ChainletContainer, c2 *ChainletContainer) *ChainletContainer { // возвращаем НОВЫЙ экземпляр!
	chOut := append(c1.Chainlet.Chain, c2.Chainlet.Chain...)

	chLetOut := &Chainlet{
		Chain: chOut,
	}

	out := &ChainletContainer{
		Rate:     rateCalc.CalcRate(chLetOut),
		Chainlet: chLetOut,
	}

	return out // TODO: реализация возможно упрощена, можно будет доработать
}

type ChainletRepo interface { // репо цепочек
	SetNewChainlet(*Chainlet) (ID uint64)
}

var rateCalc CalcChainletRate // TODO: пока проще сделать автономной сущностью, для которой потом найду место

type CalcChainletRate interface {
	CalcRate(*Chainlet) float64
}

/*
ChainletGenerator - генерирует набор цепочек действий, которые можно провести с текущим состоянием
*/
type ChainletGenerator struct {
	MaxChainletLenght int
	MaxVersionsCount  int
	// TODO: Parallelism
	ChangersRepo AtomicChangerRepo
	Comparer     StateComparer
}

func (c *ChainletGenerator) GenChainlets(maxSimilarity float64, minSimilarity float64, curState *State, targetState *State) []*ChainletContainer {
	wg := sync.WaitGroup{}
	wg.Add(c.MaxVersionsCount)

	out := make([]*ChainletContainer, c.MaxVersionsCount)

	for i := 0; i < c.MaxVersionsCount; i++ {
		num := i

		go func() {
			out[num] = c.GenChainlet(maxSimilarity, curState, targetState)

			wg.Done()
		}()
	}

	wg.Wait()

	// сортируем и обрезаем по minSimilarity
	sort.Slice(out, func(i, j int) bool {
		return out[i].Rate < out[j].Rate
	})

	// обрезаем по minSimilarity
	for i, chCon := range out {
		if chCon.Rate < minSimilarity {
			out = out[:i]

			break
		}
	}

	return out
}

/*
GenChainlet - генерируем цепочку (один из вариантов набора последовательности действий)
*/
func (c *ChainletGenerator) GenChainlet(maxSimilarity float64, curState *State, targetState *State) *ChainletContainer {
	out := &ChainletContainer{
		Rate:     0.0,
		Chainlet: NewChainlet(),
	}

	for i := 0; i < c.MaxChainletLenght; i++ {
		chID, chGer := c.ChangersRepo.GetRandom() // каждый раз берём случайное действие
		out.Chainlet.Add(chID)
		curState = chGer.Change(curState)

		if out.Rate = c.Comparer.Comparison(curState, targetState); out.Rate >= maxSimilarity {
			break
		}
	}

	return out
}
