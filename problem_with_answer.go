package nutcracker

// Nutcracker
// Problem-based approach
// ProblemWithAnswer
// Copyright © 2022-2024 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"sort"
)

const (
	compareByStart compareBy = iota
	compareByEnd
	compareByDifference
)

type compareBy int

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
	deltaState  *State
	answers     []*ChainletContainer
}

func newProblemWithAnswer(curState *State, targetState *State, answers []*ChainletContainer) *ProblemWithAnswer {
	return &ProblemWithAnswer{
		curState:    curState,
		targetState: targetState,
		deltaState:  curState.Delta(targetState),
		answers:     answers,
	}
}

type ProblemWithAnswerContainer struct {
	repoID   int
	distance float64 // дистанция универсальна, по ситуации
	pwa      *ProblemWithAnswer
}

type repoIdent struct {
	repoID   int
	distance float64
}

type ProblemWithAnswerRepo struct {
	сomparer StateComparer
	list     []*ProblemWithAnswer

	// нужность этих полей под вопросом т.к. мы и так имеем возможность искать по сходству стейтов
	// они могут иметь смысл при большом (гигантском) количестве PWA в репозитории
	curStateSimilarity    map[int][]repoIdent
	targetStateSimilarity map[int][]repoIdent
	deltaStateSimilarity  map[int][]repoIdent
}

func (p *ProblemWithAnswerRepo) Add(pwa *ProblemWithAnswer) {
	p.list = append(p.list, pwa)
}

func (p *ProblemWithAnswerRepo) FindByState(state *State, comparison compareBy, rateSimilarity float64) []*ProblemWithAnswerContainer {
	out := make([]*ProblemWithAnswerContainer, 0)

	for repoID, pwa := range p.list {
		var distance float64

		switch comparison {
		case compareByStart:
			distance = p.сomparer.Comparison(pwa.curState, state)

		case compareByEnd:
			distance = p.сomparer.Comparison(pwa.targetState, state)

		default: // compareByDifference
			distance = p.сomparer.Comparison(pwa.deltaState, state)
			// дельту надо считать заранее distance = p.сomparer.Comparison(pwa.curState.Delta(pwa.targetState), state)
		}

		if distance < rateSimilarity {
			pwac := &ProblemWithAnswerContainer{repoID: repoID, distance: distance, pwa: pwa}
			out = append(out, pwac)
		}
	}

	// сортируем результат по дистанции
	sort.Slice(out, func(i, j int) bool {
		iv, jv := out[i], out[j]
		return iv.distance < jv.distance
	})

	return out
}
