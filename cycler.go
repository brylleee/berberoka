package main

type Cycler struct {
	Charsets []string
	Current  string
	Steps    int

	LastFinish bool

	PreviousCycler *Cycler
}

func (cycler *Cycler) cycle() {
	if cycler.Steps+1 < len(cycler.Charsets) {
		cycler.Steps++
	} else {
		cycler.Steps = 0
		if cycler.PreviousCycler != nil {
			cycler.PreviousCycler.cycle()
		} else {
			cycler.LastFinish = true
		}
	}

	cycler.Current = cycler.Charsets[cycler.Steps]
}
