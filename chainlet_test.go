package nutcracker

// Nutcracker
// Problem-based approach
// Chainlet tests
// Copyright © 2022-2024 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"fmt"
	"testing"
)

func TestChainletGeneratorEasy(t *testing.T) {
	levelFunc := func(maxLevel int) int { // maxLevel для варианта с рандомным уровнем в перспективе
		return 0
	}
	achr := NewAtomicChangerRepository()
	achr.Set(newChangerIncrement(levelFunc, 0.000)) // бесполезный атомайзер
	achr.Set(newChangerIncrement(levelFunc, 0.0011))
	achr.Set(newChangerIncrement(levelFunc, 0.0021))
	achr.Set(newChangerIncrement(levelFunc, 0.0031))

	// --
	comparer := &EuclideanDistance{}

	if achr.counter != 4 {
		t.Errorf("Want 4, have %d", achr.counter)
	}

	chg := NewChainletGenerator(10, 300, achr, comparer)

	chg.Copy()
	state1 := &State{
		Data: []float64{0.000, 0.000},
	}

	state2 := &State{
		Data: []float64{0.007, 0.000},
	}

	fmt.Println(chg.GenChainlets(0.002, 0.001, state1, state2))
	containers := chg.GenChainlets(0.002, 0.001, state1, state2)

	for _, ctnr := range containers {
		fmt.Println("RATE: ", ctnr.Distance)
		fmt.Println("ШАГИ сколько: ", len(ctnr.Chainlet.Chain))
		fmt.Println("ШАГИ какие: ", ctnr.Chainlet.Chain)
	}
}
