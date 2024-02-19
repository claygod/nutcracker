package nutcracker

// Nutcracker
// Problem-based approach
// Chainlet
// Copyright © 2022-2024 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"log"
	"sort"
	"strings"
	"sync"
)

const (
	separator = "#"
)

type Chainlet struct { // цепочка действий имеющая удовленворяющий результат (смысл)
	// ID uint64 // возможно снаружи
	// Rate float64
	countSteps int64
	ChainIDs   []int64 // храним идентификаторы а не ссылки чтобы сравнивать цепочки на похожесть
	ChainNames []string
}

func NewChainlet() *Chainlet {
	return &Chainlet{
		ChainIDs:   make([]int64, 0),
		ChainNames: make([]string, 0),
	}
}

func (c *Chainlet) Add(chID int64, chName string, steps int64) {
	c.countSteps += steps
	c.ChainIDs = append(c.ChainIDs, chID)
	c.ChainNames = append(c.ChainNames, chName)
}

func (c *Chainlet) MergeChainlet(ch *Chainlet) {
	c.countSteps += ch.GetCountSteps()
	c.ChainIDs = append(c.ChainIDs, ch.ChainIDs...)
	c.ChainNames = append(c.ChainNames, ch.ChainNames...)
}

func (c *Chainlet) GetCountSteps() int64 {
	return c.countSteps
}

func (c *Chainlet) GetMultiName() string {
	return strings.Join(c.ChainNames, separator)
}

/*
ChainletContainer - контейнер нужен для того, чтобы иметь возможность сравнить
*/
type ChainletContainer struct {
	// ID uint64 // возможно снаружи
	Distance float64 // исчисляется исходя только из коэффициэнта сравнения state! но не длины цепочки (количества действий)
	Chainlet *Chainlet
}

/*
GetChainletStepsCount - для дополнительной сортировки по количеству шагов при одинаковой Distance
*/
func (c *ChainletContainer) GetChainletStepsCount() int64 {
	return c.Chainlet.GetCountSteps()
}

/*
ProblemWithAnswer - структура которую можно сформировать по результату поиска решения задачи
С помощью репозитория таких структур можно группировать их по:
- входному состоянию
- выходному состоянию
- дельте между входным и выходным состояниям
Т.е. можно создать этакую карту Кохонена, по которой можно искать решение задачи/проблемы помимо простого перебора атомарных чейнжеров.
А каждая секция карты становится похожей на колонку нейронов в мозговой нейросети.
*/
type ProblemWithAnswer struct {
	curState    *State
	targetState *State
	answers     []*ChainletContainer
}

// func MergeChainletContainers222(c1, c2 *ChainletContainer) *ChainletContainer { // возвращаем НОВЫЙ экземпляр!
// 	chLetOut := &Chainlet{
// 		Chain: append(c1.Chainlet.Chain, c2.Chainlet.Chain...),
// 	}

// 	out := &ChainletContainer{
// 		Distance: rateCalc.CalcRate(chLetOut),
// 		Chainlet: chLetOut,
// 	}

// 	return out // TODO: реализация возможно упрощена, можно будет доработать
// }

// type ChainletRepo interface { // репо цепочек
// 	SetNewChainlet(*Chainlet) (ID uint64)
// }

// var rateCalc CalcChainletRater // TODO: пока проще сделать автономной сущностью, для которой потом найду место
// var rateCalc = &CalcChainletRate{}

// type CalcChainletRate struct { // имплементация пока не юзанного CalcChainletRater
// 	// в перспективе пригодится AtomicChangerRepo (см. описание интерфейса CalcChainletRater)
// }

/*
CalcRate - по сути анализируем, насколько быстрая (эффективная) эта цепочка
Чем длиней, тем хуже, т.к. потребуется больше шагов для решения задачи
*/
// func (c *CalcChainletRate) CalcRate(chl *Chainlet) float64 {
// 	// var sum int64 = 1

// 	// for i, k := range chl.Chain {
// 	// 	sum += int64(i) * k
// 	// }

// 	// return 1.0 / float64(sum)
// 	return 1.0 / float64(len(chl.Chain))
// }

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

func (c *ChainletGenerator) GenChainlets(rateSimilarity, minSimilarity float64, curState, targetState *State) []*ChainletContainer {
	wg := sync.WaitGroup{}
	wg.Add(c.MaxVersionsCount)

	out := make([]*ChainletContainer, c.MaxVersionsCount)

	//var emptyChainlets int64

	for i := 0; i < c.MaxVersionsCount; i++ {
		num := i

		go func() {
			resp := c.GenChainlet(rateSimilarity, curState, targetState)
			out[num] = resp

			// if resp == nil {
			// 	atomic.AddInt64(&emptyChainlets, 1)
			// }

			wg.Done()
		}()
	}

	wg.Wait()

	// проверяем на содержимое nil в списке
	outWithoutNil := make([]*ChainletContainer, 0, len(out))

	for _, item := range out {
		if item != nil {
			outWithoutNil = append(outWithoutNil, item)
		}
	}
	//fmt.Println(out)
	out = outWithoutNil
	//fmt.Println(out)
	// if emptyChainlets > 0 {
	// 	return make([]*ChainletContainer, 0)
	// }

	// сортируем и обрезаем по minSimilarity
	// sort.Slice(out, func(i, j int) bool {
	// 	return out[i].Distance > out[j].Distance
	// })

	sort.Slice(out, func(i, j int) bool {
		iv, jv := out[i], out[j]
		switch {
		case iv.Distance != jv.Distance:
			return iv.Distance < jv.Distance
		default:
			return iv.GetChainletStepsCount() < jv.GetChainletStepsCount()
		}
	})
	//fmt.Println(out)

	// обрезаем по minSimilarity
	outMinSimilarity := make([]*ChainletContainer, 0, len(out))

	for _, item := range out {
		if item.Distance < minSimilarity {
			outMinSimilarity = append(outMinSimilarity, item)
		}
	}

	// по условиям можем первое (лучшее) решение добавлять в репо атомиков

	// for i, chCon := range out {
	// 	if chCon.Distance < minSimilarity {
	// 		out = out[:i]

	// 		break
	// 	}
	// }

	return out
}

func (c *ChainletGenerator) SetChainletAsAtomicChanger(ch *ChainletContainer) {
	// TODO: определиться, подходит ли такое решение (насколько оно идеально)
	// пока добавляем всегда
	achs, err := newAtomicChangerSyntheticFromChainlet(*ch.Chainlet, c.ChangersRepo)
	if err != nil {
		log.Println(err)

		return
	}

	c.ChangersRepo.Set(achs)
}

/*
GenChainlet - генерируем цепочку (один из вариантов набора последовательности действий)
*/
func (c *ChainletGenerator) GenChainlet(rateSimilarity float64, curState, targetState *State) *ChainletContainer {
	// fmt.Println("STEP 301 ", rateSimilarity)
	out := &ChainletContainer{
		Distance: 0.0,
		Chainlet: NewChainlet(),
	}

	for i := 0; i < c.MaxChainletLenght; i++ {
		chID, chGer := c.ChangersRepo.GetRandom() // каждый раз берём случайное действие
		if chID == -1 {                           // ноль означает полное отсутствие цепочек в репе, не из чего выбирать
			return nil
		}
		// fmt.Println("STEP 303 ", chID, chGer)
		out.Chainlet.Add(chID, chGer.GetName(), chGer.GetInnerSteps())
		curState = chGer.Change(curState)
		//fmt.Println("STEP 304 -измененный текущий статус- ", curState)
		out.Distance = c.Comparer.Comparison(curState, targetState)
		//fmt.Println("STEP 305 -похожесть- ", out.Rate)
		if out.Distance < rateSimilarity {
			//fmt.Println("STEP 306 -ПОДОШЛО!!- ", out.Rate)
			break
		}
	}

	return out
}
