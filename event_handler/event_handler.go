package event_handler

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	"../elevio"
	"../network/peers"
	"../order_distributor"
	"../sync/elevators"
	"../sync/nextfloor"
	"../sync/store"
)

var selfIP = peers.GetRelativeTo(peers.Self, 0)

// RunElevator Her skjer det
func RunElevator() {

	// First we start the server
	fmt.Println("Starting elevator server ...")
	err := (exec.Command("gnome-terminal", "-x", "/home/kristian/Dokumenter/Skole/sanntid/SimElevatorServer")).Run()
	if err != nil {
		fmt.Println("Something went wrong!")
		log.Fatal(err)
	}

	time.Sleep(time.Duration(1 * time.Second)) // To avoid crash due to not started sim
	elevio.Init("localhost:15657", numFloors)

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)
	nextFloor := make(chan int)

	go elevio.PollButtons(drv_buttons) // Etasje og hvilken type knapp som blir trykket
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)
	go nextfloor.SubscribeToDestinationUpdates(nextFloor)

	fmt.Println("Elevator server is running")

	// Initialize all elevators at the bottom when the program is first run.
	store.SetCurrentFloor(selfIP, store.NumFloors)

	goToFloor(0, drv_floors)

	for {
		select {
		case a := <-drv_buttons: // Just sets the button lamp, need to translate into calls
			fmt.Println("Noen trykket på en knapp, oh lø!")

			// Setter på lyset
			// light := store.DetermineLight(a.Floor, a.Button)
			// elevio.SetButtonLamp(a.Button, a.Floor, light)

			// Håndtere callen
			if a.Button == elevio.BT_Cab {
				store.AddCabCall(selfIP, a.Floor)
			} else {
				elevDir := btnDirToElevDir(a.Button)
				mostSuitedIP := store.MostSuitedElevator(a.Floor, elevDir)

				// Create and send HallCall
				hc := elevators.HallCall_s{Floor: a.Floor, Direction: elevDir}
				order_distributor.SendHallCall(mostSuitedIP, hc)
			}

		case floor := <-nextFloor:
			fmt.Print("dest:")
			fmt.Println(floor)
			/*go*/ goToFloor(floor, drv_floors)
		}
	}
}

func goToFloor(destinationFloor int, drv_floors <-chan int) { // Probably add a timeout'

	direction := elevators.DirectionIdle
	currentFloor, _ := store.GetCurrentFloor(selfIP)
	fmt.Print("curr: ")
	fmt.Println(currentFloor)
	if currentFloor < destinationFloor {
		direction = elevators.DirectionUp
	} else if currentFloor > destinationFloor {
		direction = elevators.DirectionDown
	}

	elevio.SetMotorDirection(direction)
	store.SetDirectionMoving(selfIP, direction)
	for {
		select {
		case floor := <-drv_floors: // Wait for elevator to reach floor
			fmt.Printf("Reaching floor: %d", floor)
			elevio.SetFloorIndicator(floor)
			if floor == destinationFloor {
				arrivedAtFloor(floor)
				return
			}
			break
		case <-time.After(5 * time.Second):
			fmt.Println("Didn't reach floor in time!")
			elevio.SetMotorDirection(e.DirectionIdle)
			//Do some shit
			return
			break
		}
	}
}

func arrivedAtFloor(floor int) {
	fmt.Printf("setting floor %d", floor)
	store.SetCurrentFloor(selfIP, floor)
	elevio.SetMotorDirection(elevators.DirectionIdle) // Stop elevator and set lamps and stuff
	store.SetDirectionMoving(selfIP, elevators.DirectionIdle)
	elevio.SetFloorIndicator(floor)
	elevio.SetButtonLamp(elevio.BT_HallUp, floor, false)
	elevio.SetButtonLamp(elevio.BT_HallDown, floor, false)
	elevio.SetButtonLamp(elevio.BT_Cab, floor, false)
	elevio.SetDoorOpenLamp(true)
	time.Sleep(3 * time.Second)
	elevio.SetDoorOpenLamp(false)
}

func btnDirToElevDir(btn elevio.ButtonType) elevators.Direction_e {
	switch btn {
	case elevio.BT_HallDown:
		return elevators.DirectionDown
	case elevio.BT_HallUp:
		return elevators.DirectionUp
	default:
		fmt.Println("Invalid use, must be either up or down")
		return elevators.DirectionIdle
	}
}
