package elevators

import "fmt"

type HallCall_s struct {
	Floor     int
	Direction Direction_e
}

type Direction_e int

const (
	DirectionUp   Direction_e = 1
	DirectionDown Direction_e = -1
	DirectionBoth Direction_e = 2 // Only used for HallCalls
	DirectionIdle Direction_e = 0
)

type Elevator_s struct {
	ip              string
	currentFloor    int
	NumFloors       int
	prevFloor       int
	directionMoving Direction_e
	hallCalls       []HallCall_s
	cabCalls        []bool
}

func New(peerIP string, numFloors int, currentFloor int) Elevator_s {
	elevator := Elevator_s{
		ip:              peerIP,
		currentFloor:    currentFloor,
		prevFloor:       currentFloor, // initialize with same floor
		NumFloors:       numFloors,
		directionMoving: DirectionIdle,
		hallCalls:       make([]HallCall_s, numFloors),
		cabCalls:        make([]bool, numFloors),
	}

	return elevator
}

func (e Elevator_s) GetIP() string {
	return e.ip
}

func (e Elevator_s) GetCurrentFloor() int {
	return e.currentFloor
}

func (e *Elevator_s) SetCurrentFloor(newCurrentFloor int) {
	fmt.Print("Elevator set: ")
	fmt.Println(newCurrentFloor)
	e.currentFloor = newCurrentFloor
}

func (e Elevator_s) GetDirectionMoving() Direction_e {
	return e.directionMoving
}
func (e Elevator_s) GetPreviousFloor() int {
	return e.prevFloor
}

func (e *Elevator_s) SetDirectionMoving(newDirection Direction_e) {
	e.directionMoving = newDirection
}

func (e Elevator_s) GetAllHallCalls() []HallCall_s {
	return e.hallCalls
}

func (e *Elevator_s) AddHallCall(hallCall HallCall_s) {
	if hallCall.Direction == DirectionUp && e.hallCalls[hallCall.Floor].Direction == DirectionDown {
		e.hallCalls[hallCall.Floor].Direction = DirectionBoth
	} else if hallCall.Direction == DirectionDown && e.hallCalls[hallCall.Floor].Direction == DirectionUp {
		e.hallCalls[hallCall.Floor].Direction = DirectionBoth
	} else {
		e.hallCalls[hallCall.Floor].Direction = hallCall.Direction
	}
}

func (e *Elevator_s) RemoveHallCalls(floor int) {
	e.hallCalls[floor].Direction = DirectionIdle
}

func (e *Elevator_s) AddCabCall(floor int) {
	e.cabCalls[floor] = true
}

func (e *Elevator_s) RemoveCabCall(floor int) {
	e.cabCalls[floor] = false
}

func (e Elevator_s) GetAllCabCalls() []bool {
	return e.cabCalls
}
