package nutcracker

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

/*
PointsGroup - результат кластеризации (один из кластеров)
*/
type PointsGroup struct {
	points [][2]float64
}

func (p *PointsGroup) GetFingerPrint() *FingerPrint {
	// TODO: implement me
	// приводим к начальной системе координат и генерируем отпечаток группы
	return nil
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

func (p *PointsGroupRepo) Grep(pg *PointsGroup) []int {
	out := make([]int, 0, len(p.data))
	// TODO: implement me
	// тут проходим for-ом
	return out
}

/*
FingerPrint - универсально описание группы точек по которому можно искать похожие
*/
type FingerPrint struct {
	// TODO: implement me
	left    float64
	right   float64
	top     float64
	bottom  float64
	center  float64
	density float64 // плотность
}

func (f *FingerPrint) Finger() []float64 {
	// TODO: implement me
	return nil
}
