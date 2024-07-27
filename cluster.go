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

/*
PointsGroup - результат кластеризации (один из кластеров)
*/
type PointsGroup struct {
	points      [][2]float64 // при создании проверять что точки есть (не пустой слайс)
	fingerPrint *FingerPrint
}

func (p *PointsGroup) Compare(in *PointsGroup) float64 {
	// TODO: implement me
	// поучаем некий finger print
	// перобразуем его в State и сравниваем
	return 0.0
}

func (p *PointsGroup) GetFingerPrint() *FingerPrint {
	return p.fingerPrint
}

func (p *PointsGroup) genFingerPrint() *FingerPrint {
	// TODO: implement me
	// приводим к начальной системе координат и генерируем отпечаток группы
	fp := &FingerPrint{}

	xList := make([]float64, len(p.points))
	yList := make([]float64, len(p.points))

	for i, point := range p.points {
		xList[i] = point[coordX]
		yList[i] = point[coordY]

	}

	sort.Float64s(xList)
	sort.Float64s(yList)

	// вычисляем центр
	fp.medianaX = medianForSorted(xList)
	fp.medianaY = medianForSorted(yList)
	fp.centerX = fp.medianaX
	fp.centerY = fp.medianaY

	// находим габариты
	fp.minX = xList[0]
	fp.maxX = xList[len(xList)-1]

	fp.minY = yList[0]
	fp.maxY = yList[len(yList)-1]

	// приволим к началу координат (не всё)
	fp.medianaX -= fp.minX
	fp.maxX -= fp.minX
	fp.minX = 0

	fp.medianaY -= fp.minY
	fp.maxY -= fp.minY
	fp.minY = 0

	return fp
}

/*
PointsGroupRepo - набор групп точек одного тика
*/
type PointsGroupRepo struct {
	data map[int]*PointsGroup
}

func (p *PointsGroupRepo) Add(id int, pGroup *PointsGroup) {
	// TODO: implement me
}

func (p *PointsGroupRepo) Get(id int) *PointsGroup {
	return p.data[id]
}

func (p *PointsGroupRepo) Grep(pgIn *PointsGroup) []CompareState {
	// TODO: implement me
	// тут проходим for-ом

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
	minX, maxX         float64
	minY, maxY         float64
	medianaX, medianaY float64
	centerX, centerY   float64 // координаты медианы
	// variance float64 // дисперсия
	// density float64 // плотность
}

func (f *FingerPrint) Finger() []float64 {
	// TODO: implement me
	return []float64{
		f.centerX, f.centerY,
		f.medianaX, f.medianaY,
		f.maxX, f.minX,
		f.maxY, f.minY,
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
