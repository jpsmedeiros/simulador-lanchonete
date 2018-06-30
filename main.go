package main

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/jpsmedeiros/simulador-lanchonete/simulation"
)

type Order struct {
	Type     string
	Next     *Order
	Priority bool
}

type Event struct {
	Order    Order
	Duration time.Duration
}

var (
	orders, hamburguers, icecreams, soda chan Order
	wait                                 *sync.WaitGroup
	totalTime                            uint
	simSystem                            simulation.System
	events                               []Event
)

const (
	MaxQueueBurguer = 5
	MaxQueueSoda    = 5
	MaxQueueSodaP   = 5
	MaxQueueRequest = 5
)

func init() {
	var orderNumber uint32
	seed := uint32(17)
	rand := randomNumberGenerator(1103515245, 12345, 1<<31, seed)
	simulationTime := time.Hour / 3600
	events = make([]Event, 0)
	for simulationTime > 0 {
		orderNumber = generateOrderNumber(rand())
		duration := generateRandomTime(3, 17, time.Hour)
		events = append(events, Event{Order: generateOrder(orderNumber), Duration: duration})
		simulationTime = simulationTime - duration
	}
	fmt.Println(simulationTime)
}

func main() {
	wait = new(sync.WaitGroup)
	wait.Add(1)
	//cria um canal de pedidos
	//Cria os canais de serviço
	orders, hamburguers, soda = make(chan Order, 3), make(chan Order, 5), make(chan Order, 5)

	//Gera pedidos
	//TODO: Precisa consumir de uma lista
	go generateClients()
	//Consome pedidos de hamburguer
	go hamburguerHandler()
	//Consome os pedidos de refrigerante
	go sodaHandler()
	//Organiza pedidos da fila de clientes
	go handleClientRequests()
	wait.Wait()

}

func generateClients() {
	for _, event := range events {
		newEvent := event
		orders <- newEvent.Order
		time.Sleep(event.Duration)
		fmt.Println(simSystem)
	}
	for simSystem != (simulation.System{}) {
	}
	wait.Done()
}

func handleClientRequests() {
	clientQueue := make([]Order, 0, MaxQueueRequest)
	go func() {
		for {
			if len(clientQueue) > 0 && (clientQueue[0] != Order{}) {
				processRequest(clientQueue[0])
				clientQueue = clientQueue[1:]
				simSystem.RemoveClient()
			}
		}
	}()
	//Lida com os pedidos para separa-los
	for item := range orders {
		if len(clientQueue) < MaxQueueRequest {
			clientQueue = append(clientQueue, item)
			simSystem.AddClient()
		} else {
			fmt.Println("Fila de pedidos cheia. Pedido descartado")
		}
	}
}

func hamburguerHandler() {
	normal := make([]Order, 0, MaxQueueBurguer)
	go func() {
		for {
			//Tem alguem na fila de prioridade
			if len(normal) > 0 && (normal[0] != Order{}) { // tem alguem na fila comum
				makeBurguer(normal[0])
				simSystem.RemoveHamburguer()
				if normal[0].Next != nil {
					normal[0].Priority = true
					soda <- *normal[0].Next
				}
				normal = normal[1:]
			}
		}
	}()

	for order := range hamburguers {
		if len(normal) <= MaxQueueBurguer {
			normal = append(normal, order)
			simSystem.AddHamburguer()
		} else {
			fmt.Println("Pedido de hamburguer descartado")
		}
	}
}

func sodaHandler() {
	priority := make([]Order, 0, MaxQueueSodaP)
	normal := make([]Order, 0, MaxQueueSoda)
	go func() {
		for {
			//Tem alguem na fila de prioridade
			if len(priority) > 0 && (priority[0] != Order{}) {
				makeSoda(priority[0])
				simSystem.RemoveSodaP()
				priority = priority[1:]
			} else if len(normal) > 0 && (normal[0] != Order{}) { // tem alguem na fila comum
				makeSoda(normal[0])
				simSystem.RemoveSoda()
				normal = normal[1:]
			}
		}
	}()
	for order := range soda {
		if order.Priority {
			if len(priority) <= MaxQueueSodaP {
				priority = append(priority, order)
				simSystem.AddSodaP()
			} else {
				fmt.Println("Pedido de refrigerante com prioridade descartado")
			}
		} else {
			if len(normal) <= MaxQueueSoda {
				normal = append(normal, order)
				simSystem.AddSoda()
			} else {
				fmt.Println("Pedido de refrigerante descartado")
			}
		}
	}
}

func processRequest(order Order) {
	time.Sleep(generateRandomTime(0, 30, time.Hour))
	switch order.Type {
	case "hamburguer":
		hamburguers <- order
	case "soda":
		soda <- order
	default:
		fmt.Println("lixo")
	}
}

func makeBurguer(order Order) {
	time.Sleep(generateRandomTime(5, 20, time.Hour))
	//fmt.Println(order.Type)
}

func makeSoda(order Order) {
	time.Sleep(generateRandomTime(3, 60, time.Hour))
	//fmt.Println(order.Type)
}

func randomNumberGenerator(a, c, m, seed uint32) func() float64 {
	// usage: rand := randomNumberGenerator(1103515245, 12345, 1<<31, 0)
	// then call rand()
	r := seed
	return func() float64 {
		r = (a*r + c) % m
		return float64(r) / float64(m)
	}
}

func generateRandomExp(lambda float64, u float64) float64 {
	return (-1 / lambda) * math.Log(1-u)
}

func generateRandomTime(seed uint32, delta float64, duration time.Duration) time.Duration {
	rand := randomNumberGenerator(1103515245, 12345, 1<<31, seed)
	return time.Duration(generateRandomExp(delta, rand())*float64(time.Minute)) / 3600
}

func generateOrderNumber(number float64) uint32 {
	orderNumber := uint32(number * 10)
	if isInRange(0, 3, orderNumber) {
		return 1 // hamburguer
	} else if isInRange(4, 5, orderNumber) {
		return 2 // soda
	} else {
		return 3 // hamburguer + soda
	}
}

func isInRange(x, y, value uint32) bool {
	return x <= value && value <= y
}

func generateOrder(orderNumber uint32) Order {
	if orderNumber == 1 {
		return Order{Type: "hamburguer"}
	} else if orderNumber == 2 {
		return Order{Type: "soda"}
	} else {
		return Order{Type: "hamburguer", Next: &Order{Type: "soda"}}
	}
}
