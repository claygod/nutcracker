package nutcracker

// Nutcracker
// Problem-based approach
// Chainlet
// Copyright © 2022 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"sort"
	"sync"
	"sync/atomic"
)

type Chainlet struct { // цепочка действий имеющая удовленворяющий результат (смысл)
	// ID uint64 // возможно снаружи
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
ChainletContainer - контейнер нужен для того, чтобы иметь возможность сравнить
*/
type ChainletContainer struct {
	// ID uint64 // возможно снаружи
	Rate     float64 // исчисляется исходя не только из коэффициэнта сравнения state, но и длины цепочки (количества действий)
	Chainlet *Chainlet
}

func MergeChainletContainers(c1, c2 *ChainletContainer) *ChainletContainer { // возвращаем НОВЫЙ экземпляр!
	chLetOut := &Chainlet{
		Chain: append(c1.Chainlet.Chain, c2.Chainlet.Chain...),
	}

	out := &ChainletContainer{
		Rate:     rateCalc.CalcRate(chLetOut),
		Chainlet: chLetOut,
	}

	return out // TODO: реализация возможно упрощена, можно будет доработать
}

// type ChainletRepo interface { // репо цепочек
// 	SetNewChainlet(*Chainlet) (ID uint64)
// }

var rateCalc CalcChainletRate // TODO: пока проще сделать автономной сущностью, для которой потом найду место

type CalcChainletRate interface {
	// пока вижу возможность считать исходя из длины цепочки (количества действий),
	// но если дать доступ к AtomicChangerRepo, а в нём внутри всем AtomicChanger назначить какие-то веса
	// (чтобы у первичных он был маленький, а для вторичных рос с количеством внутренних шагов, т.е. суммой внутренних операций)
	// TODO: в имплементации доступ к AtomicChangerRepo, где по идентификаторам в цепочке берём конкретный AtomicChanger.GetInnerSteps()
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

func NewChainletGenerator(maxChainletLenght, maxVersionsCount int, changersRepo AtomicChangerRepo, comparer StateComparer) *ChainletGenerator {
	return &ChainletGenerator{
		MaxChainletLenght: maxChainletLenght,
		MaxVersionsCount:  maxVersionsCount,
		// TODO: Parallelism
		ChangersRepo: changersRepo,
		Comparer:     comparer,
	}
}

func (c *ChainletGenerator) Copy() *ChainletGenerator {
	return &ChainletGenerator{
		MaxChainletLenght: c.MaxChainletLenght,
		MaxVersionsCount:  c.MaxVersionsCount,
		ChangersRepo:      c.ChangersRepo,
		Comparer:          c.Comparer,
	}
}

func (c *ChainletGenerator) GenChainlets(maxSimilarity, minSimilarity float64, curState, targetState *State) []*ChainletContainer {
	wg := sync.WaitGroup{}
	wg.Add(c.MaxVersionsCount)

	out := make([]*ChainletContainer, c.MaxVersionsCount)

	var emptyChainlets int64

	for i := 0; i < c.MaxVersionsCount; i++ {
		num := i

		go func() {
			resp := c.GenChainlet(maxSimilarity, curState, targetState)
			out[num] = resp

			if resp == nil {
				atomic.AddInt64(&emptyChainlets, 1)
			}

			wg.Done()
		}()
	}

	wg.Wait()

	// проверяем на содержимое nil в списке
	if emptyChainlets > 0 {
		return make([]*ChainletContainer, 0)
	}

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
func (c *ChainletGenerator) GenChainlet(maxSimilarity float64, curState, targetState *State) *ChainletContainer {
	out := &ChainletContainer{
		Rate:     0.0,
		Chainlet: NewChainlet(),
	}

	for i := 0; i < c.MaxChainletLenght; i++ {
		chID, chGer := c.ChangersRepo.GetRandom() // каждый раз берём случайное действие
		if chID == 0 {                            // ноль означает полное отсутствие цепочек в репе, не из чего выбирать
			return nil
		}

		out.Chainlet.Add(chID)
		curState = chGer.Change(curState)

		if out.Rate = c.Comparer.Comparison(curState, targetState); out.Rate >= maxSimilarity {
			break
		}
	}

	return out
}
