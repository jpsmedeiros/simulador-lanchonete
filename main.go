package main

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

type Order struct {
	Type              string
	Next              *Order
	sodaPriorityQueue bool
}

type Event struct {
	Order    Order
	Duration time.Duration
}

type SafeQueue struct {
	Queue []Order
	mux   sync.RWMutex
}

var (
	events                                                     []Event
	clientQueue, hamburguerQueue, sodaPriorityQueue, sodaQueue *SafeQueue
	eventsCount, descarteCliente, totalEvents, totalFila       int64
	utilization                                                time.Duration
	mux                                                        *sync.Mutex
)

// Define experiment
const (
	QueueSize         = 15
	MaxQueueBurguer   = QueueSize
	MaxQueueSoda      = QueueSize
	MaxQueueSodaP     = QueueSize
	MaxQueueRequest   = QueueSize
	ArrivalTimeUnit   = time.Millisecond
	ArrivalTime       = 110 //E[C]
	ArrivalQuantity   = 1
	ServiceTimeUnit   = time.Millisecond
	ServiceTime       = 90 //E[C]
	ServiceQuantity   = 1
	TimeFasterMult    = 60
	MaxSimulationTime = time.Hour / TimeFasterMult
)

func init() {
	mux = &sync.Mutex{}
	//Inicializa eventos
	randDuration := randomNumberGenerator(1103515245, 12345, 1<<31, 3)
	randOrder := randomNumberGenerator(1103515245, 12345, 1<<31, 4)
	simulationTime := MaxSimulationTime
	events = make([]Event, 0)
	var duration time.Duration
	var orderNumber uint32
	for simulationTime > 0 {
		orderNumber = generateOrderNumber(randOrder())
		duration = generateRandomTime(randDuration(), ArrivalQuantity, ArrivalTime*ArrivalTimeUnit)
		if duration < simulationTime {
			events = append(events, Event{Order: generateOrder(orderNumber), Duration: duration})
			simulationTime = simulationTime - duration
		} else {
			events = append(events, Event{Order: generateOrder(orderNumber), Duration: simulationTime})
			simulationTime = 0
		}
	}
	//Inicializa variavel de controle
	eventsCount = int64(len(events))
	//Inicializa variaveis de analise
	totalEvents = eventsCount
	//Inicializa fila
	clientQueue = &SafeQueue{Queue: make([]Order, 0, MaxQueueRequest)}
}

func main() {
	//Gera pedidos
	go generateEvents()
	handleEvents()
	for eventsCount > 0 {
	}
	fmt.Printf("(Descarte: %.2f%%, Utilização: %.0f%%, Media fila: %.0f)\n", float64(descarteCliente)/float64(totalEvents), (float64(utilization)/float64(MaxSimulationTime))*100, float64(totalFila)/float64(totalEvents))
}

func generateEvents() {
	for _, event := range events {
		//Boqueia fila de clientes para inserção
		if len(clientQueue.Queue) < MaxQueueRequest {
			mux.Lock()
			clientQueue.Queue = append(clientQueue.Queue, event.Order)
			mux.Unlock()
		} else {
			mux.Lock()
			atomic.AddInt64(&eventsCount, -1)
			atomic.AddInt64(&descarteCliente, 1)
			mux.Unlock()
		}
		totalFila += int64(len(clientQueue.Queue))
		//Desbloqueia fila de cliente
		time.Sleep(event.Duration)
	}
	eventsCount = atomic.LoadInt64(&eventsCount)
	descarteCliente = atomic.LoadInt64(&descarteCliente)
}

func handleEvents() {
	go func() {
		randTime := randomNumberGenerator(1103515245, 12345, 1<<31, 5)
		for eventsCount > 0 {
			if len(clientQueue.Queue) > 0 {
				t := generateRandomTime(randTime(), 1, ServiceTime*ServiceTimeUnit)
				time.Sleep(t)
				clientQueue.mux.Lock()
				utilization = utilization + t
				clientQueue.Queue = clientQueue.Queue[1:]
				atomic.AddInt64(&eventsCount, -1)
				clientQueue.mux.Unlock()
				fmt.Println("BURGAO: ", eventsCount)
			}
		}
	}()
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

func generateRandomTime(u, lambda float64, duration time.Duration) time.Duration {
	return time.Duration(generateRandomExp(lambda, u)*float64(duration)) / TimeFasterMult
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
