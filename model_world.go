package nutcracker

import (
	"sort"

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

	// TODO: надо научиться генерировать Transducer'ы
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
	// - тут должен быть доступ к состояниям (положительным, отрицательным)
	//   чтобы было с чем сравнивать, т.е. делать выводы и оценивать

	// цикл (с рекурсией):
	// - исходное состояние - набор объектов
	// - возможные изменения (набор вариантов)
	// -- для каждого варианта пробуем преобразовать (спрогнозировать) набор объектов (дубликатов)
	// --- новые наборы проверяем на коллизии, если надо, какой-то объект отктываем (или игнорируем всю вероятность)
	// --- новые наборы сравниваем с набором проблем и при некотором близком значении
	//     в коллектор результатов пишем возможную возможность (с указанием/расчётом с учетом уровня рекурсии и м.б.цвета)

	// результат - набор проблем: ключ - действие, значение - список проблем с вероятностями

}

/*
WorkerSelectAction - воркер выбора решения (м.б. это часть WorkerSituationAnalysis ??)
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
	og     *ObjectsGroup
	rating float64 // если объектов много, то могло бы быть так,
	// что проблема отрацательна для чего-то одного и положительна для другого
	// в шахматах мы жертвуем пешку чтобы получить коня, т.е. есть и плюсы и минусы
	// но возможно, две рядом существующие проблемы - это выход, одна положительная, другая отрицательная
}

func (p *Problem) Copy() *Problem {
	return nil // TODO: implement me
}

/*
ProblemHub - верхнеуровневая структура, выполняющая работу по поддержанию Problem в актуальном состоянии
*/
type ProblemHub struct {
	pList      []*Problem
	lowerLimit float64
}

func (p *ProblemHub) PushProblem(prb *Problem) {
}

func (p *ProblemHub) GetActualProblem(og *ObjectsGroup) []*Problem {
	out := make([]*Problem, 0)

	for _, prb := range p.pList {
		rate := prb.og.Compare(og)
		if rate > p.lowerLimit {
			pCopy := prb.Copy()
			pCopy.rating *= rate
			out = append(out, prb)
		}
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].rating < out[j].rating
	})

	return out
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

type Object struct {
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

func (o *Object) Transformation(t *Transducer) *Object {
	// TODO: преобразование на основе знания своего типа и своих параметров
	//       возвращается копия самого себя

	return nil
}

type Transducer struct {
	objTypes []*ObjectType // преобразует только определенные типы

	// TODO: некий преобразователь, БАЗОВЫЙ!, то, что мы можем сделать (сдвиг влево-вправо)
	//       он естественно влияет только на некоторые или даже один объект (отбивалку),
	//       но привязку к нужному типу объектов надо найти самостоятельно
	// Оперирует наверно *PointsGroupUnion
}

func (t *Transducer) Update(pgu *PointsGroupUnion) *PointsGroupUnion {
	return nil // TODO: implement me
}

type ObjectType struct {
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
