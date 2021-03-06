package core

import "github.com/enjoypi/gostatechart"

type ActorType = int64
type Message = gostatechart.Event
type MessageQueue chan Message

const (
	EvStopped = iota + 1
	EvPanic
)
