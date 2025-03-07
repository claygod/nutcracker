package nutcracker

import (
	"fmt"
	"math"
)

// Nutcracker
// Problem-based approach
// Entities
// Copyright © 2022-2024 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

/*
Need - потребность
*/
// type Need struct {
// 	ID string
// }

// type TaskGenerator interface {
// 	GenTaskFromNeed(*Need)
// }

/*
Task - задача
Пока состояния (начальное, текущее и конечное) подразумеваются одного порядка,
т.е. длина массивов в них одинаковая, что упрощает их сравнение.
*/
type Task struct {
	ID            string
	maxSimilarity float64
	minSimilarity float64

	/* пояснение к нижележащим переменными:
	   мы генерируем findDirectsCount промежуточных целей, т.е. точек, церез которые попробуем найти путь к главной цели.
	   Напиример, рекурсия у нас 2 а каунт 5. Тогда мы генерируем 5 точек-целей, если из какой-то находим путь к главной цели,
	   то ура, всё ок, если нет, то из каждой из пяти чертим пять линий к новосгенерированнм пяти,
	   т.е. на втором уровне 25 промежуточных целей, из которых мы пытаемся дойти до главной цели
	*/
	recursionLevel int // количество рекурсий по поиску промежуточных шагов (поиск пути),
	// т.е. количество промежуточных шагов в неких направлениях
	findDirectsCount int // в одной рекурсии количество

	beginState  *State
	curState    *State
	targetState *State

	chlGen *ChainletGenerator
	// rStateGen IntermediateRandomStateGenerator
	sComparer StateComparer

	// Шаги, ведущие к цели
	// Steps []*ChainletContainer

	// ParentTasks []*Task
	// ChildTasks  []*Task
	// TODO: scope - контекст задачи. Возможно скоп, он подобен State (допустим это стартовый стейт при решении задачи).
	// Также возможно, что скоп снаружи, и возможный для использования генератор Chainlet-наборов уже относится к какому-то скопу).
}

func NewTask(
	id string,
	maxSimilarity float64,
	minSimilarity float64,
	recursionLevel int,
	findDirectsCount int,

	beginState *State,
	curState *State,
	targetState *State,

	chlGen *ChainletGenerator,
	// rStateGen IntermediateRandomStateGenerator,
	sComparer StateComparer,
) *Task {
	return &Task{
		ID:               id,
		maxSimilarity:    maxSimilarity,
		minSimilarity:    minSimilarity,
		recursionLevel:   recursionLevel,
		findDirectsCount: findDirectsCount,

		beginState:  beginState,
		curState:    curState,
		targetState: targetState,

		chlGen: chlGen,
		// rStateGen: rStateGen,
		sComparer: sComparer,

		// Steps:       make([]*ChainletContainer, 0),
		// ParentTasks: make([]*Task, 0),
		// ChildTasks:  make([]*Task, 0),
	}
}

func (t *Task) Copy() *Task {
	return &Task{
		ID:               t.ID,
		maxSimilarity:    t.maxSimilarity,
		minSimilarity:    t.minSimilarity,
		recursionLevel:   t.recursionLevel,
		findDirectsCount: t.findDirectsCount,

		beginState:  t.beginState,
		curState:    t.curState,
		targetState: t.targetState,

		chlGen: t.chlGen,
		// rStateGen: t.rStateGen,
		sComparer: t.sComparer,

		// Steps:       make([]*ChainletContainer, 0), // TODO: пока слайсы не копируем, не знаем, надо ли
		// ParentTasks: make([]*Task, 0),
		// ChildTasks:  make([]*Task, 0),
	}
}

func (t *Task) FindChainlets() []*ChainletContainer { // тут мы ищем оптимальный путь
	decisions := t.chlGen.GenChainlets(t.maxSimilarity, t.minSimilarity, t.curState, t.targetState)

	if len(decisions) == 0 && t.recursionLevel > 0 { // не найдено подходящих решений

		// === НИЖЕ РЕШЕНИЕ ПОКА ОТМЕНЕНО ===
		// и ещё можно создавать промежуточные шаги
		// генерируем новые (промежуточные) цели, которых можем добиться
		// и уже в каждой точке промежуточных целей пробуем заново добиться основной цели
		// (действуем рекурсивно)
		// === ВЫШЕ РЕШЕНИЕ ПОКА ОТМЕНЕНО ===

		// пробуем увеличить количество шагов для достижения цели
		chlGen := t.chlGen.Copy()

		for i := 0; i < t.findDirectsCount; i++ { // вместо i++ надо придумать более интересный механизм
			chlGen.MaxChainletLenght = t.chlGen.MaxVersionsCount + i

			if de2 := chlGen.GenChainlets(t.maxSimilarity, t.minSimilarity, t.curState, t.targetState); len(de2) > 0 {
				decisions = append(decisions, de2...)

				break
			}
		}

		// for i := 0; i < t.findDirectsCount; i++ {
		// 	newState := t.rStateGen.GenTask(t.curState, t.targetState, t.sComparer)

		// 	newTask := t.Copy()
		// 	newTask.recursionLevel = t.recursionLevel - 1
		// 	newTask.curState = newState

		// 	for _, dt := range newTask.FindChainlets() { // это получаем результаты к промежуточной цели
		// 		newTask2 := newTask.Copy()
		// 		newTask2.recursionLevel = t.recursionLevel - 2

		// 		for _, dt2 := range newTask2.FindChainlets() { // теперь из промежуточной точки пытамся добраться до основной цели
		// 			decisions = append(decisions, MergeChainletContainers(dt, dt2))
		// 		}
		// 	}
		// }
	}

	if len(decisions) > 0 {
		if achs, err := newAtomicChangerSyntheticFromChainlet(*decisions[0].Chainlet, t.chlGen.ChangersRepo); err == nil {
			fmt.Println("номер добавленного синтетического чейнжера: ", t.chlGen.ChangersRepo.Set(achs))
		}
	}

	return decisions
}

// type TaskCompletionCheck2 interface {
// 	CompletionCheck(*Task, *Task) float64 // оценка скорей всего от 0.0 до 1.0 CompletionCheck
// }

// type IntermediateRandomStateGenerator interface {
// 	GenTask(*State, *State, StateComparer) *State // генерирование некоторого состояния, находящегося где-то между начальной и конечной задачей
// }

/*
AtomicChanger - атомарный изменятель
Надо продумать возможность создания генерирующего атомарного изменятеля, т.е. когда он на лету генерирует вариант изменения
(случайный номер в массиве и к нему применяется или ещё что-то такое же) и возможность этот случайный вариант сохранять если решение
задачи успешное или наоборот выбрасывать, если неуспешное.
*/
type AtomicChanger interface { // минимальное атомарное изменение
	Change(*State) *State
	GetInnerSteps() int64 // количество внутренних встроенных AtomicChanger (для базовых AtomicChanger это всегда единица)
	GetName() string
}

type AtomicChangerRepo interface { // DONE: AtomicChangerRepository репо атомиков
	/*
	   GetByTarget - берём наиболее подходящее (близкое) по цели
	*/
	GetByTarget(targetState *State) (ID int64, aChanger AtomicChanger)

	/*
	   GetRandom - берём случайную, это удобно для генерации случайного Chainlet-набора
	*/
	GetRandom() (ID int64, aChanger AtomicChanger)
	GetByID(ID int64) (aChanger AtomicChanger, ok bool)

	/*
		SetRandom - сначала добавляем действительно базовые возможности, а потом можно добавлять
		Chainlet-наборы, которые используются часто или которые короткие но эффективные
	*/
	Set(aChanger AtomicChanger, opts ...bool) (ID int64)
	// NOTE: пока не требуется но возможно будет нужен Get(ID int64) (aChanger AtomicChanger)
}

type StateComparer interface { // DONE: EuclideanDistance сравниваем состояния (направление и координаты)
	Comparison(*State, *State) float64
}

/*
EuclideanDistance - вариант вычисления расстояния между векторами по методу Евклида
Понятное описание трёх популярных вариантов расчёта расстояния между векторами:
https://tproger.ru/translations/3-basic-distances-in-data-science/
*/
type EuclideanDistance struct {
}

func (e *EuclideanDistance) Comparison(st1, st2 *State) float64 { // отрицательное значение ответа как отрицательный результат
	// fmt.Println("STEP 500 -len- ", len(st1.Data), len(st2.Data))
	if len(st1.Data) < len(st2.Data) {
		return -9999999999999999.9 // TODO: проверить на проверку в использовании
	}

	var out float64

	for i := 0; i < len(st1.Data); i++ {
		// fmt.Println("STEP 501 -out- ", out)
		out += math.Pow((st1.Data[i] - st2.Data[i]), 2) // nolint: gomnd
		// fmt.Println("STEP 502 -out- ", out)
	}

	return math.Sqrt(out)
}

type State struct {
	// vector - coord. and direct(?)
	Data []float64
}

func NewState(data []float64) *State {
	return &State{
		Data: data,
	}
}

func (s *State) Add(in []float64) {
	s.Data = append(s.Data, in...)
}

func (s *State) Copy() *State {
	stateOut := &State{
		Data: make([]float64, len(s.Data)),
	}

	copy(stateOut.Data, s.Data)

	return stateOut
}

func (s *State) Delta(in *State) *State {
	stateOut := s.Copy()

	for i, _ := range stateOut.Data {
		stateOut.Data[i] = stateOut.Data[i] - in.Data[i]
	}

	return stateOut
}

// /*
// TaskResource -  учёт ресурсов, выделенных для решения задачи, обычно ресурсы только тратятся,
// но при каких-то определённых обстоятельствах ресурсы могут и повышаться
// (например найден Chainlet, достойный добавление в репо Changer-атомиков.
// */
// type TaskResource interface { // NOTE: возможно тут потребуются float
// 	Add(int64) int64
// 	Cut(int64) int64
// 	Total() int64
// 	ResetToZero() // напоминание о том, что у задачи может оказаться ситуация,
//                // когда точно надо остановить поиски путей (Chainlet) её выполнения
// }

type CalcChainletRater interface { // DONE: имплементированно в CalcChainletRate
	// пока вижу возможность считать исходя из длины цепочки (количества действий),
	// но если дать доступ к AtomicChangerRepo, а в нём внутри всем AtomicChanger назначить какие-то веса
	// (чтобы у первичных он был маленький, а для вторичных рос с количеством внутренних шагов, т.е. суммой внутренних операций)
	// TODO: в имплементации доступ к AtomicChangerRepo, где по идентификаторам в цепочке берём конкретный AtomicChanger.GetInnerSteps()
	CalcRate(*Chainlet) float64
}
