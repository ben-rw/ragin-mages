package wizarena

func (w *WizArena) CleanUp() {
	w.audioPlayer.Close()
}
