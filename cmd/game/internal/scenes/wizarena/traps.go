package wizarena

func (w *WizArena) UpdateTraps() {
	w.trapTimer -= 1
	if w.trapTimer <= 0 {
		w.trapsUp = !w.trapsUp
		if w.trapsUp {
			w.trapTimer = 50 // 50 ticks up
		} else {
			w.trapTimer = 100 //100 ticks down
		}
	}
}
