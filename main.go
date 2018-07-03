package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"sync"
	"time"
)

type Order struct {
	Duration time.Time
}

type Event struct {
	Order    Order
	Duration time.Duration
}

type SafeQueue struct {
	Queue []Order
	mux   *sync.RWMutex
}

var (
	events                                                     []Event
	clientQueue, hamburguerQueue, sodaPriorityQueue, sodaQueue *SafeQueue
	eventsCount, descarteCliente, totalEvents, totalFila       int64
	finished                                                   []time.Duration

	maxQueueRequest, queueSize uint

	wg *sync.WaitGroup

	serviceTime, arrivalTime uint
	utilization, meanTime    time.Duration
)

// Define experiment
const (
	ArrivalTimeUnit = 1000000
	ArrivalQuantity = 1

	ServiceTimeUnit = 1000000
	ServiceQuantity = 1

	TimeFasterMult    = 60
	MaxSimulationTime = time.Hour / TimeFasterMult

	Seed = 0
)

func init() {
	var err error
	temp, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("erro ao executar")
		return
	}
	maxQueueRequest = uint(temp)

	temp, err = strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("erro ao executar")
		return
	}
	serviceTime = uint(temp)

	temp, err = strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Println("erro ao executar")
		return
	}
	arrivalTime = uint(temp)

	//Inicializa eventos
	randDuration := randomNumberGenerator(1103515245, 12345, 1<<31, Seed)
	simulationTime := MaxSimulationTime
	events = make([]Event, 0)
	var duration time.Duration
	for simulationTime > 0 {
		duration = generateRandomTime(randDuration(), ArrivalQuantity, time.Duration(arrivalTime)*ArrivalTimeUnit)
		if duration < simulationTime {
			events = append(events, Event{Duration: duration})
			simulationTime = simulationTime - duration
		} else {
			events = append(events, Event{Duration: simulationTime})
			simulationTime = 0
		}
		fmt.Println(simulationTime)
	}
	//Inicializa variavel de controle
	eventsCount = int64(len(events))
	//Inicializa variaveis de analise
	totalEvents = eventsCount
	//Inicializa fila
	clientQueue = &SafeQueue{Queue: make([]Order, 0), mux: &sync.RWMutex{}}

	finished = make([]time.Duration, 0, totalEvents)
	wg = &sync.WaitGroup{}
}

func main() {
	// InitExp()
	wg.Add(1)
	//Gera pedidosi
	go generateEvents()
	handleEvents()
	wg.Wait()
	for _, time := range finished {
		meanTime += time
	}
	fmt.Printf("(Descarte: %.5f%%, Utilização: %.5f%%, Media fila: %.5f, Media no sistema: %.2f)\n", float64(descarteCliente)/float64(totalEvents)*100, (float64(utilization)/float64(MaxSimulationTime))*100, float64(totalFila)/float64(totalEvents), (meanTime.Seconds()/1000)/float64(len(finished)))
}

func generateEvents() {
	for _, event := range events {
		//Boqueia fila de clientes para inserção
		if uint(len(clientQueue.Queue)) < maxQueueRequest {
			clientQueue.mux.Lock()
			clientQueue.Queue = append(clientQueue.Queue, event.Order)
			clientQueue.mux.Unlock()
		} else {
			clientQueue.mux.Lock()
			eventsCount--
			descarteCliente++
			clientQueue.mux.Unlock()
		}
		totalFila += int64(len(clientQueue.Queue))
		//Desbloqueia fila de cliente
		time.Sleep(event.Duration)
	}
}

func handleEvents() {
	go func() {
		randTime := randomNumberGenerator(1103515245, 12345, 1<<31, Seed+1)
		for eventsCount > 0 {
			if len(clientQueue.Queue) > 0 {
				t := generateRandomTime(randTime(), ServiceQuantity, time.Duration(serviceTime)*ServiceTimeUnit)
				utilization = utilization + t

				clientQueue.mux.RLock()
				finished = append(finished, time.Now().Sub(clientQueue.Queue[0].Duration))
				clientQueue.mux.RUnlock()

				time.Sleep(t)

				clientQueue.mux.Lock()
				clientQueue.Queue = clientQueue.Queue[1:]
				clientQueue.mux.Unlock()

				eventsCount--
				fmt.Println("BURGAO: ", eventsCount)
			}
		}
		wg.Done()
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
	return Order{Duration: time.Now()}
}
