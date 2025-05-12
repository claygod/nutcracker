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

Получение изменений скорее всего через пул pull. Т.е. получение следующей порции данных по готовности
Это отвязывает от входного потока и решает проблему постоянной синхронизации
В то же время, получается, что модель мира работает с неким срезом "устаревших" данных
Важный момент - надо определиться, кто кого дергает и в каком режиме ????
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
	//
	// тут мы создаём объекты из входных групп и помещаем в репозиторий объектов
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

	// получаем объекты (текущий срез)
	// objs := m.objHub.GetActualObjects()

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

type ObjectsGroupWithPredictor struct {
	og                *ObjectsGroup
	mapObjTransducers *ObjectsGroupTransducers // слепок изменений
	// TODO: сюда надо и информацию по близости к проблемам
}

func NewObjectsGroupTransducers() *ObjectsGroupTransducers {
	return &ObjectsGroupTransducers{
		data: make(map[int][]int),
	}
}

type ObjectsGroupTransducers struct {
	data map[int][]int
}

func (o *ObjectsGroupTransducers) FingerPrint() string {
	return "" // TODO: implement me
}

func (o *ObjectsGroupTransducers) Set(objID int, transducersList []int) {
	o.data[objID] = transducersList
}

type BranchPredictor struct {
	trs map[int]*Transducer
}

func (b *BranchPredictor) Predict(og *ObjectsGroup) {
	/*
		Результат:
		набор цепочек с оценкой на каждом шагу
		(оценка появляется исходя из сравнения с проблемами)
		NOTE: есть ли смысл в том, что эти цепочки подрастают (просчитываются дальше)
	*/

	branchSets := make(map[int][][]int)

	for i, obj := range og.objs {
		listTr := make([]int, 0, len(b.trs)) // тут будет список изменений которые возможны

		// формируем список поддерживаемых изменений
		for trID, tr := range b.trs {
			if tr.IsObjSupported(obj) {
				listTr = append(listTr, trID)
			}
		}

		// генерируем набор вариантов изменений (возсожно стохастически)
		branchSets[i] = b.genChangeSets(listTr)
	}

	// это уже набор для генерации предполагаемых новых/дублированных и измененных объектов
	finBranchSets := b.genBranchSets(branchSets)

	ogwpList := make([]*ObjectsGroupWithPredictor, 0)

	for _, brSet := range finBranchSets {
		ogCopy := og.Copy()

		obt := NewObjectsGroupTransducers()

		for i, listTr := range brSet {
			for _, trID := range listTr {
				ogCopy.objs[i] = ogCopy.objs[i].Transformation(b.trs[trID])
			}

			obt.Set(ogCopy.objs[i].id, listTr)
		}

		ogwp := &ObjectsGroupWithPredictor{
			og:                ogCopy,
			mapObjTransducers: obt,
		}

		ogwpList = append(ogwpList, ogwp)
	}

	// по окончании в ogwpList набор вариантов для ОДНОГО шага
	// TODO: теперь надо сделать рекурсию и предусмотреть:
	// ------------ уровень вложенности
	// ----------- * указание конкретной
	// ----------- * можно менять гибко в процессе
	// ----------- * можно прерывать все например из-за найденной ситуации
	// ----------- * можно увеличивать глубину если хватает ресурсов
	//
	// Нужно теперь разобраться, как работать с коллизиями, взаимодействием с другими объектами
	// Это возможно на разных уровнях:
	// - Внутри объекта, тогда он должен знать про другие объекты,
	//   и может взаимодействовать только с теми что знает
	// - Снаружи, когда взаимодействие оторвано,
	//   т.е. оно может быть отдельным и унифицированным
	//
	// ПОИСК КОЛЛИЗИЙ!
	// - получение коллизий
	// - расчет результата коллизий

	// return nil // TODO: implement me
}

func (b *BranchPredictor) checkOGCollision(inList []*ObjectsGroupWithPredictor) [][]*ObjectsGroupWithPredictor {
	return nil // TODO: implement me
}

func (b *BranchPredictor) genChangeSets(inList []int) [][]int {
	outList := make([][]int, 0)

	// пока простая заглушка, допускающая только одно изменение
	// надо предусмотреть варианты невозможных одновременных изменений,
	// т.к. мы не можем одновременно повернуть вправо и влево
	for _, in := range inList {
		outList = append(outList, []int{in})
	}

	// при необходимости обрезаем список

	return outList // TODO: implement me
}

func (b *BranchPredictor) genBranchSets(inList map[int][][]int) []map[int][]int {
	outList := make([]map[int][]int, 0)

	// пока простая заглушка, допускающая только одно изменение
	for objNum, sets := range inList {
		item := make(map[int][]int)
		item[objNum] = sets[0]
		outList = append(outList, item)
	}

	// при необходимости обрезаем список

	return outList // TODO: implement me
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

func (o *ObjectsHub) GetActualObjects() []*Object {
	return nil // TODO: implement me
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

func (o *ObjectsGroup) Copy() *ObjectsGroup {
	return nil // TODO: implement me
}

type Object struct {
	id              int // пытаемся идентифицировать
	objType         ObjectType
	lifeBegin       int               // метка начала существования объекта
	lifeEnd         int               // метка конца существования объекта
	curPointsGroups *PointsGroupUnion // в объекте может быть несколько отдельных групп
	// prevPointsGroup *PointsGroup
	pointsGroupsChain []*PointsGroupUnion // TODO: возможно отсюда брать прогноз поведения (ЭТО ПРЕДЫДУЩАЯ ИСТОРИЯ)
}

func (o *Object) Merge(o2 *Object) *Object {
	// TODO: слияние двух объектов с одинаковым поведением или ещё что-то такое
	// TODO: в результате появится новый объект ИЛИ текущий обновится (пожирнеет)
	return nil
}

func (o *Object) Transformation(t *Transducer) *Object {
	// TODO: преобразование на основе знания своего типа и своих параметров
	//       возвращается копии самого себя

	return nil
}

func (o *Object) Copy() *Object {
	return nil // TODO: implement me
}

type Transducer struct {
	id       int           // надо как-то их различать, хотя возможно, идентификатор надо хранить "снаружи"
	objTypes []*ObjectType // преобразует только определенные типы

	// TODO: некий преобразователь, БАЗОВЫЙ!, то, что мы можем сделать (сдвиг влево-вправо)
	//       он естественно влияет только на некоторые или даже один объект (отбивалку),
	//       но привязку к нужному типу объектов надо найти самостоятельно
	// Оперирует наверно *PointsGroupUnion
	//
	// ОЧЕНЬ ВАЖНО: предполагаемое изменение делается с учетом предыдущей истории объекта!!
}

func (t *Transducer) Update(obj *Object) *Object {
	// например возвращаем объект только если он штменился
	return nil // TODO: implement me
}

func (t *Transducer) IsObjSupported(obj *Object) bool {
	// проверка на поддерживаемость для конкретного объекта (можем ли его изменить)
	return true // TODO: implement me
}

type ObjectType struct {
	// содержит варианты дельт или варианты действий (сдвиг, поворот и пр.)
	// TODO: ОДНО из двух, надо разобраться и определиться
	chs  []*Chainlet
	chcs []*ChainletContainer

	// achs []AtomicChanger // доступные изменятели

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
