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
- Предлагает следующий шаг (с характеристикой "уверенность")

Получение изменений скорее всего через пуш. Т.е. получение следующей порции данных по готовности
Это отвязывает от входного потока и решает проблему постоянной синхронизации
*/
type ModelWorld struct {
	curObjects []*Object
}

type Object struct { // Alphabet
	objType         ObjectType
	lifeBegin       int // метка начала существования объекта
	lifeEnd         int // метка конца существования объекта
	curPointsGroup  *PointsGroup
	prevPointsGroup *PointsGroup
}

type ObjectType struct { // Alphabet
	// содержит варианты дельт или варианты действий (сдвиг, поворот и пр.)

	// сожержит признаки похожести, по которым можно сделать вывод о похожести объектов

	// количество всех объектов такого типа (м.б. ссылки)
	// количество живых объектов такого типа (м.б. ссылки)
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
	indexPoints := index.NewTrivialIndex(testPongPoints)
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
