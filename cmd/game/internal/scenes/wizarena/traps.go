package wizarena

var trapTimer = 100

func (w *WizArena) UpdateTraps() {
	trapTimer -= 1
	if trapTimer <= 0 {
		w.trapsUp = !w.trapsUp
		if w.trapsUp {
			trapTimer = 50 // 60 ticks up
		} else {
			trapTimer = 100 //120 ticks down
		}
	}
}
