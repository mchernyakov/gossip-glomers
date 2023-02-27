package model

import "sync"

type NeibsState struct {
	neibs []string
	mu    sync.Mutex
}

func NewNeibsState() *NeibsState {
	return &NeibsState{
		neibs: []string{},
		mu:    sync.Mutex{},
	}
}

func (neibsState *NeibsState) AddNeibs(newNeibs []string) {
	neibsState.mu.Lock()
	defer neibsState.mu.Unlock()

	neibsState.neibs = append(neibsState.neibs, newNeibs...)
}

func (neibsState *NeibsState) GetNeibs() []string {
	neibsState.mu.Lock()
	defer neibsState.mu.Unlock()

	return append([]string(nil), neibsState.neibs...)
}
