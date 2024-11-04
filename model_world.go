package nutcracker

// Nutcracker
// Problem-based approach
// Model of the world
// Copyright © 2024 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

/*
ModelWorld - модель мира
- Кластеризует объекты из входного массива данных
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

type Point [2]float64         // точка на экране - скорей всего будут использоваться исходные массивы, а это для образа
type PointsGroup [][2]float64 // группа точек, кластер
