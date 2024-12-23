package nutcracker

import (
	"sort"
)

// Nutcracker
// Problem-based approach
// Clusterise
// Copyright © 2022-2024 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

/*
Кластеризация происходит примерно так:
- в тике генерируются группы точек (кластеры)
					по сути это обьекты
- проводится анализ по похожести с предыдущим тиком
		похожие увязываются в цепочку
- формируются новые объекты или к имеющимся привязываются свежие данные
	(мигающие объекты - искать по предыдущим и более старым на историческую глубину(должна быть/формироваться)
- объекты могут формировать дельты, которые уже могут сохраняться и применяться для прогнозирования
*/

const (
	coordX = iota
	coordY
)

func NewPointsGroup() *PointsGroup {
	pg := &PointsGroup{
		points: make([][2]float64, 0),
	}

	// pg.fingerPrint = pg.genFingerPrint()

	return pg
}

/*
PointsGroup - результат кластеризации (один из кластеров) группа точек, кластер
*/
type PointsGroup struct {
	points [][2]float64 // при создании проверять что точки есть (не пустой слайс)
	// fingerPrint *FingerPrint
}

func (p *PointsGroup) Compare(in *PointsGroup) float64 {
	// TODO: implement me
	// поучаем некий finger print
	// перобразуем его в State и сравниваем
	return 0.0
}

func (p *PointsGroup) Add(point [2]float64) {
	p.points = append(p.points, point)
}

func (p *PointsGroup) GetFingerPrint() *FingerPrint {
	return p.genFingerPrint() // p.fingerPrint
}

func (p *PointsGroup) genFingerPrint() *FingerPrint {
	fdv := &FingerData{}

	xList := make([]float64, len(p.points))
	yList := make([]float64, len(p.points))

	for i, point := range p.points {
		xList[i] = point[coordX]
		yList[i] = point[coordY]

	}

	sort.Float64s(xList)
	sort.Float64s(yList)

	// вычисляем центр
	fdv.medianaX = medianForSorted(xList)
	fdv.medianaY = medianForSorted(yList)
	// fdv.centerX = fdv.medianaX
	// fdv.centerY = fdv.medianaY

	// находим габариты
	fdv.minX = xList[0]
	fdv.maxX = xList[len(xList)-1]

	fdv.minY = yList[0]
	fdv.maxY = yList[len(yList)-1]

	// приводим к началу координат
	fdt := &FingerData{
		minX:     0,
		maxX:     fdv.maxX - fdv.minX,
		minY:     0,
		maxY:     fdv.maxY - fdv.minY,
		medianaX: fdv.medianaX - fdv.minX,
		medianaY: fdv.medianaY - fdv.minY,
	}

	return &FingerPrint{
		typeData:  fdt,
		valueData: fdv,
	}
}

// type PointsGroupRepoRepo struct {
// 	data []*PointsGroupRepo
// }

func NewTaktWrap(prevTakt *TaktWrap, repo *PointsGroupRepo) *TaktWrap {
	/*
		TODO: тут ГЕНЕРАЦИЯ объектов !!!!!!!!!!!!!!!!!

			! каждая pointsGroup это объект

		Но желательно связать каждый объект с объектами предыдущего такта, а при необходимости и пред-предыдущего (м.б. нескольких)

		Надо определить тип	и искать по типу в предыдущем такте
		А из выбранных по типу выбирать по значению

		Это вообще длинный путь, и вроде ка можно сразу по значению искать
			Плюсы: отфильтровываем другие типы, всё точнее
			Минусы: дольше и фильтр по типу может сработать хуже прямого сравнения

		Можно делать прямой поиск, т.е. объект, и искать ему предков, и обратный, предкам искать текущуюее состояние
	*/

	newTakt := &TaktWrap{
		repo: repo,
		// objs         []*Object
		previousTakt: prevTakt,
	}

	return newTakt
}

type TaktWrap struct {
	repo         *PointsGroupRepo
	objs         []*Object
	previousTakt *TaktWrap
}

// type Object struct {
// 	chain []*ObjChainLink
// }

/*
PointsGroupRepo - набор групп точек одного тика
*/
type PointsGroupRepo struct {
	data map[int]*PointsGroup
}

type ObjChainLink struct {
	taktNum int
	pg      *PointsGroup
}

func (p *PointsGroupRepo) Add(id int, pGroup *PointsGroup) {
	p.data[id] = pGroup
}

func (p *PointsGroupRepo) Get(id int) *PointsGroup {
	return p.data[id]
}

func (p *PointsGroupRepo) Grep(pgIn *PointsGroup) []CompareState {
	out := make([]CompareState, len(p.data))

	for i, pg := range p.data {
		out[i] = CompareState{
			pgID:       i,
			similarity: pg.Compare(pgIn),
		}
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].similarity < out[j].similarity
	})

	return out
}

type CompareState struct {
	pgID       int
	similarity float64
}

/*
FingerPrint - универсально описание группы точек по которому можно искать похожие
*/
type FingerPrint struct {
	typeData  *FingerData
	valueData *FingerData
}

type FingerData struct {
	minX, maxX         float64
	minY, maxY         float64
	medianaX, medianaY float64
	// centerX, centerY   float64 // координаты медианы
	// variance float64 // дисперсия
	// density float64 // плотность
}

func (f *FingerPrint) Value() []float64 {
	// TODO: implement me
	return []float64{
		f.valueData.medianaX, f.valueData.medianaY,
		f.valueData.maxX, f.valueData.minX,
		f.valueData.maxY, f.valueData.minY,
	}
}

func (f *FingerPrint) Type() []float64 {
	// TODO: implement me
	return []float64{
		f.typeData.medianaX, f.typeData.medianaY,
		f.typeData.maxX, f.typeData.minX,
		f.typeData.maxY, f.typeData.minY,
	}
}

func medianForSorted(in []float64) float64 {
	ln := len(in)

	if ln%2 == 1 {
		return in[(ln-1)/2]
	} else {
		return (in[ln/2] + in[(ln/2)-1]) / 2
	}

	return 0.0
}
