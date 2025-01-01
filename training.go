package nutcracker

// Nutcracker
// Problem-based approach
// Learning by example
// Copyright © 2022-2024 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

type School struct { // Alphabet
}

func (s *School) CreateABCAtomicChanger(name string, curState, targetState *State) AtomicChanger {
	return &ABCChanger{
		name:  name,
		delta: targetState.Delta(curState),
	}
}

func (s *School) CreateDriftAtomicChanger(name string, delta, sample *State, cmpr StateComparer, rate float64, max int) AtomicChanger {
	return &DriftChanger{
		name:            name,
		delta:           delta,
		sample:          sample,
		comparer:        cmpr,
		rate:            rate,
		maxCountChanges: max,
	}
}

type ABCChanger struct {
	name  string
	delta *State
}

func (a *ABCChanger) Change(in *State) *State {
	out := in.Copy()

	for i := range out.Data {
		if i >= len(a.delta.Data) {
			break
		}

		out.Data[i] = out.Data[i] + a.delta.Data[i]
	}

	return out
}

func (a *ABCChanger) GetInnerSteps() int64 {
	return 1
}

func (a *ABCChanger) GetName() string {
	return a.name
}

type DriftChanger struct {
	name            string
	delta           *State
	sample          *State
	comparer        StateComparer
	rate            float64
	maxCountChanges int
}

func (d *DriftChanger) Change(in *State) *State {
	out := in.Copy()

	shift := len(out.Data) - len(d.delta.Data)

	if shift <= 0 {
		return out
	}

	for u, c := 0, 0; u < shift && c < d.maxCountChanges; u, c = u+1, c+1 {
		item := NewState(out.Data[u : u+len(d.delta.Data)])

		distance := d.comparer.Comparison(item, d.sample)
		if distance < d.rate {
			for i := range item.Data {
				item.Data[i] = item.Data[i] + d.delta.Data[i]
			}

			u += len(d.delta.Data) - 1
		}
	}

	return out
}

func (d *DriftChanger) GetInnerSteps() int64 {
	return int64((len(d.delta.Data) / 2) + 1)
}

func (d *DriftChanger) GetName() string {
	return d.name
}
