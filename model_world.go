package nutcracker

import (
	"github.com/lfritz/clustering/dbscan"
	"github.com/lfritz/clustering/index"
)

// Nutcracker
// Problem-based approach
// Model of the world
// Copyright © 2024 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

/*
ModelWorld - модель мира
- Кластеризует объекты из входного массива данных
- Входных потоков несколько и их надо синхронизировать (или это поручить роутеру)
- Полученные объекты пытается классифицировать (определить тип)
- Пытается предугадать будущее
- Из предположений о будущем могут формироваться Проблемы
- Предлагает следующий шаг (с характеристикой "уверенность")

Получение изменений скорее всего через пуш. Т.е. получение следующей порции данных по готовности
Это отвязывает от входного потока и решает проблему постоянной синхронизации
*/
type ModelWorld struct {
	// NOTE: ProblemWithAnswer скорей всего внутри, используется для/как ?? ОПРЕДЕЛИТЬСЯ
	curObjects []*Object
	objHub     *ObjectsHub
	// TODO:		главное, это как и куда результаты работы отправляются.
	// 			Т.е. в некий модуль оценки предложений от потенциально нескольких/многих MW
}

/*
PushCurPointsGroups - основной рабочий метод
Инфраструктурная часть проводит анализ текущего среза с целью сформировать кластеры (группы точек)
После формирования кластеров они попадают на вход модели мира (моделей может быть несколько)
Модель должна полученные кластеры преобразовать в объекты, провести анализ и сделать предсказание(я)
Даже если предсказаний несколько, всё равно, выбирается одно, на основании которого формируется реакция (действие)
Действие скорей всего отправляется вовне, а не реализуется самомстоятельно (предположительно их надо разнести по разным структурам).

Аргумент *PointsGroupUnion тут скорее технический, по факту это []*PointsGroup т.е. группы, которые потом расползуться по объектам
*/
func (m *ModelWorld) PushCurPointsGroups(pgu *PointsGroupUnion) {
	// NOTE: получение кластеров, работа с объектами, прогнозирование, принятие решений - всё это может быть асинхронным!

	// полученые кластеры надо проанализировать, отдать объектам, если они к ним относятся, или создать новые объекты
	// некоторые объекты могут исчезнуть. Скорее всего мы некоторое время их "помним" и учитываем (???)
	// После проведения работы с объектами мы имеем актуализованное состояние (мы в моменте)

	// Дальнейшая работа с актуализованным срезом/состоянием:
	// TODO:		анализ ситуации и прогнозирование, т.е. из объектов берем предположение о будущем состоянии и анализируем его
	// 			прогноз - это возмодно целая цепочка, и анализ по шагам вперед
	// 			при этом для каждого шага просчитываем варианты действия (дерево вариантов) и из этого дерева выбираем вариант
}

/*
WorkerObjGen - воркер генерации обьектов
*/
func (m *ModelWorld) WorkerObjGen() {
	// TODO: implement me
}

/*
WorkerSituationAnalysis - воркер анализа ситуации и прогнозирование
*/
func (m *ModelWorld) WorkerSituationAnalysis() {
	// TODO: implement me
}

/*
WorkerSelectAction - воркер выбора решения
*/
func (m *ModelWorld) WorkerSelectAction() {
	// TODO: implement me
}

/*
WorkerSend - воркер отправки решения
*/
func (m *ModelWorld) WorkerSend() {
	// TODO: implement me
}

/*
ObjectsHub - верхнеуровневая структура, выполняющая работу по поддержанию их в актуальном состоянии
*/
type ObjectsHub struct {
	// TODO: objRepo
	// TODO: objTypeRepo
}

func (o *ObjectsHub) PushCurPointsGroups(pgu *PointsGroupUnion) {
}

func (o *ObjectsHub) GetActualObjects() {
}

/*
Problem - предположительно проблема, это:
- потенциальная цепочка действий (и/или её результат)
- конечное состояние нескольких объектов (это признак проблемы)
- окрас (отрицательный-положительный) и величина

Проблема может иметь и отрицательный и положительный окрас
*/
type Problem struct {
	// TODO:
	og *ObjectsGroup
}

/*
ObjectsGroup - набор объектов, который может служить как некий признак
Пример: мяч и ворота с сильно близкими координатами
*/
type ObjectsGroup struct {
	objs []*Object
}

func (o *ObjectsGroup) Compare(o2 *ObjectsGroup) float64 {
	// Сравнение может быть как в абсолютных, так и в относительных координатах
	return 0.0 // TODO: признак близости, далеко - это ноль
}

type Object struct { // Alphabet
	objType         ObjectType
	lifeBegin       int               // метка начала существования объекта
	lifeEnd         int               // метка конца существования объекта
	curPointsGroups *PointsGroupUnion // в объекте может быть несколько отдельных групп
	// prevPointsGroup *PointsGroup
	pointsGroupsChain []*PointsGroupUnion // TODO: возможно отсюда брать прогноз поведения
}

func (o *Object) Merge(o2 *Object) *Object {
	// TODO: слияние двух объектов с одинаковым поведением или ещё что-то такое
	// TODO: в результате появится новый объект ИЛИ текущий обновится (пожирнеет)
	return nil
}

// func (o *Object) CheckSimilarity(o2 *Object) float64 {
// 	// TODO: пока просто сравниваем текущие группы, в перспективе бы сравнивать историю
// 	return o.curPointsGroups.Compare(o2.curPointsGroups)
// }

type ObjectType struct { // Alphabet
	// содержит варианты дельт или варианты действий (сдвиг, поворот и пр.)
	// TODO: ОДНО из двух, надо разобраться и определиться
	chs  []*Chainlet
	chcs []*ChainletContainer

	// содержит признаки похожести, по которым можно сделать вывод о похожести объектов
	// ИЛИ же по списку объектов можно пробежаться и сделать вывод о похожести
	pgus []*PointsGroupUnion

	// количество всех объектов такого типа (м.б. ссылки)
	// количество живых объектов такого типа (м.б. ссылки)
	// TODO: под вопросом objs []*Object
}

func (o *ObjectType) CheckSimilarity(oIn *Object) float64 {
	// TODO: пока просто перебираем и отдаём самое похожее, но в перспективе более старые должны иметь меньшую значимость
	resp := 0.0

	// NOTE: вместо сравнения объектов лучше сравнить группы - так можно и историю зацепить
	// for _, obj := range o.objs {
	// 	if sim := oIn.CheckSimilarity(obj); sim > resp {
	// 		resp = sim
	// 	}
	// }

	for _, pgu := range o.pgus {
		if sim := oIn.curPointsGroups.Compare(pgu); sim > resp {
			resp = sim
		}
	}

	return resp
}

// type Point [2]float64         // точка на экране - скорей всего будут использоваться исходные массивы, а это для образа

/*
PointsGroupFormer - формирование групп точек (кластеров)
используется алгоритм DBSCAN - https://en.m.wikipedia.org/wiki/DBSCAN
*/
type PointsGroupFormer struct {
	eps    float64 //  distance to the nearest neighbor
	minPts int     // min points count
	groups []*PointsGroup
	// rawData
}

func NewPointsGroupFormer(eps float64, minPts int) *PointsGroupFormer {
	return &PointsGroupFormer{
		eps:    eps,
		minPts: minPts,
		groups: make([]*PointsGroup, 0),
	}
}

func (p *PointsGroupFormer) AddRawCluster(in [][2]float64) {
	indexPoints := index.NewTrivialIndex(in)
	cl := dbscan.Dbscan(indexPoints, p.eps, p.minPts)

	if len(cl) > 0 {
		curObjNum := cl[0]
		curCluster := NewPointsGroup()

		for i, objNum := range cl {
			if curObjNum != objNum {
				curObjNum = objNum
				p.groups = append(p.groups, curCluster)
				// nov obj
				curCluster = NewPointsGroup()
			}

			curCluster.Add(in[i])
		}

		p.groups = append(p.groups, curCluster)
	}
}

func (p *PointsGroupFormer) LoadGroups() []*PointsGroup {
	return p.groups
}
