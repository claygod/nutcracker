# nutcracker

[![API documentation](https://godoc.org/github.com/claygod/nutcracker?status.svg)](https://godoc.org/github.com/claygod/nutcracker)

Problem-based approach

## Задачный подход

Решатель, который может искать цепочку возможных шагов для того чтобы приблизиться из стартового состояния к целевому.
Проблема постановки задачи остаётся "снаружи".

## Процесс

Создаём решатель, учим каким-то базовым, хардкорным умениям и затем ставим зандачи и смотрим что получится. 
Если в процессе решения задачи будет найден какой-то весьма эффективный ход (цепочка ходов), то этот ход может быть добавлен к базовым.

Решатель ищет множество решений, и из этого пула решений какое-то будет выбрано исходя из неких критериев.
Если решения не найдены, то предусмотрен вариант с поиском случайного хода, после которого можно будет найти решение, т.е. поиск кружного пути.

## ToDo

- [ ] имплементация CalcChainletRate (нужна хоть какая-нибудь, базовая)
- [ ] имплементация AtomicChanger
- [ ] имплементация AtomicChangerRepo
- [ ] имплементация StateComparer (возможно снаружи передается при начальном создании)


## Визуализация пакета

- устанавливаем https://github.com/davidschlachter/embedded-struct-visualizer/tree/main
- в корне нужного пакета запускаем embedded-struct-visualizer -out sh.txt ./
- на https://dreampuf.github.io/GraphvizOnline/ визуализируем полученный текстовый файл

# Copyright

Copyright © 2022-2023 Eduard Sesigin. All rights reserved. Contacts: claygod@yandex.ru
