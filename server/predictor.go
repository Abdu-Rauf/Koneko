package main

import (
	"time"
)

type Point struct {
	X, Y float64
	Ts   int64
}

type Predictor struct {
	History []Point
	Horizon float64 
}

func NewPredictor(horizon float64) *Predictor {
	return &Predictor{
		History: make([]Point, 0, 5),
		Horizon: horizon,
	}
}

func (p *Predictor) Predict(currX, currY float64) (float64, float64) {
	now := time.Now().UnixMilli()
	p.History = append(p.History, Point{X: currX, Y: currY, Ts: now})

	if len(p.History) < 2 {
		return currX, currY
	}

	// Keep only last 5 for stability
	if len(p.History) > 5 {
		p.History = p.History[1:]
	}

	oldest := p.History[0]
	newest := p.History[len(p.History)-1]
	
	dt := float64(newest.Ts - oldest.Ts)
	if dt <= 0 { return currX, currY }

	// Velocity = (NewPos - OldPos) / TimeElapsed
	vx := (newest.X - oldest.X) / dt
	vy := (newest.Y - oldest.Y) / dt

	// Future = Current + (Velocity * NetworkDelay)
	predX := currX + (vx * p.Horizon)
	predY := currY + (vy * p.Horizon)

	return predX, predY
}