package nutcracker

// Nutcracker
// Problem-based approach
// Learning by example (tests)
// Copyright © 2022-2024 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"testing"
)

func TestTrainingABC(t *testing.T) {
	sch := &School{}

	state1 := &State{
		Data: []float64{0.000, 0.000},
	}

	state2 := &State{
		Data: []float64{0.000, 0.001},
	}

	ach := sch.CreateABCAtomicChanger("increment-0001", state1, state2)

	if steps := ach.GetInnerSteps(); steps != 1 {
		t.Errorf("Want 1, have %d", steps)
	}

	stateNov := ach.Change(state1)

	if len(stateNov.Data) != len(state2.Data) {
		t.Errorf("Want %d, have %d", len(state2.Data), len(stateNov.Data))

		return
	}

	for i := range stateNov.Data {
		if stateNov.Data[i] != state2.Data[i] {
			t.Errorf("Want %f, have %f", state2.Data[i], stateNov.Data[i])
		}
	}
}

func TestTrainingDrift(t *testing.T) {
	sch := &School{}

	delta := &State{
		Data: []float64{0.001},
	}

	sample := &State{
		Data: []float64{0.002},
	}

	state1 := &State{
		Data: []float64{0.000, 0.003, 0.001},
	}

	state2 := &State{
		Data: []float64{0.000, 0.004, 0.001},
	}

	comparer := &EuclideanDistance{}

	ach := sch.CreateDriftAtomicChanger("incr-drift-0001", delta, sample, comparer, 0.0011, 5)

	if steps := ach.GetInnerSteps(); steps != 1 {
		t.Errorf("Want 1, have %d", steps)
	}

	stateNov := ach.Change(state1)

	if len(stateNov.Data) != len(state2.Data) {
		t.Errorf("Want %d, have %d", len(state2.Data), len(stateNov.Data))

		return
	}

	for i := range stateNov.Data {
		if stateNov.Data[i] != state2.Data[i] {
			t.Errorf("Want %f, have %f", state2.Data[i], stateNov.Data[i])
		}
	}
}
