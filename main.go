package main

import (
	"fmt"
	"math"
	"sync"
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
	wait                                                       *sync.WaitGroup
	events                                                     []Event
	clientQueue, hamburguerQueue, sodaPriorityQueue, sodaQueue *SafeQueue
)

// Define experiment
const (
	QueueSize       = 15
	MaxQueueBurguer = QueueSize
	MaxQueueSoda    = QueueSize
	MaxQueueSodaP   = QueueSize
	MaxQueueRequest = QueueSize
	ArrivalTimeUnit = time.Millisecond
	ArrivalTime     = 110 //E[C]
	ArrivalQuantity = 1
)

func init() {
	var orderNumber uint32
	var duration time.Duration
	rand := randomNumberGenerator(1103515245, 12345, 1<<31, 17)
	simulationTime := time.Hour / 3600
	events = make([]Event, 0)
	for simulationTime > 0 {
		orderNumber = generateOrderNumber(rand())
		duration = generateRandomTime(rand(), ArrivalQuantity, ArrivalTime*ArrivalTimeUnit)
		events = append(events, Event{Order: generateOrder(orderNumber), Duration: duration})
		simulationTime = simulationTime - duration
	}
	fmt.Println(simulationTime)
}

func main() {
	wait = new(sync.WaitGroup)
	wait.Add(len(events))
	clientQueue = &SafeQueue{Queue: make([]Order, 0, MaxQueueRequest)}
	hamburguerQueue = &SafeQueue{Queue: make([]Order, 0, MaxQueueBurguer)}
	sodaPriorityQueue = &SafeQueue{Queue: make([]Order, 0, MaxQueueSodaP)}
	sodaQueue = &SafeQueue{Queue: make([]Order, 0, MaxQueueSoda)}

	//Consome pedidos de hamburguer
	go hamburguerHandler()
	//Consome os pedidos de refrigerante
	go sodaHandler()
	//Organiza pedidos da fila de clientes
	go handleClientRequests()
	//Gera pedidos
	generateClients()
	wait.Wait()
}

func generateClients() {
	for _, event := range events {
		//Boqueia fila de clientes para inserção
		clientQueue.mux.Lock()
		if len(clientQueue.Queue) < MaxQueueRequest {
			clientQueue.Queue = append(clientQueue.Queue, event.Order)
		} else {
			wait.Done()
		}
		//Desbloqueia fila de clientes
		clientQueue.mux.Unlock()
		time.Sleep(event.Duration)
		fmt.Printf("(%d, %d, %d, %d)\n", len(clientQueue.Queue), len(hamburguerQueue.Queue), len(sodaPriorityQueue.Queue), len(sodaQueue.Queue))
	}
	fmt.Println(wait)
}

func handleClientRequests() {
	rand := randomNumberGenerator(1103515245, 12345, 1<<31, 0)
	for {
		//Caso haja alguém na fila de clientes
		if len(clientQueue.Queue) > 0 && (clientQueue.Queue[0] != Order{}) {
			//Efetua o tempo de atendimento
			time.Sleep(generateRandomTime(rand(), 30, time.Hour))
			//Bloqueia a fila de cliente para leitura
			clientQueue.mux.RLock()
			processRequest(clientQueue.Queue[0])
			//Consome fila de clientes
			clientQueue.Queue = clientQueue.Queue[1:]
			//Desbloqueia a fila
			clientQueue.mux.RUnlock()
		}
	}
}

func hamburguerHandler() {
	rand := randomNumberGenerator(1103515245, 12345, 1<<31, 5)
	for {
		//Caso haja alguem na fila de hamburguer
		if len(hamburguerQueue.Queue) > 0 { // tem alguem na fila comum
			//Faz o hamburguer
			makeBurguer(hamburguerQueue.Queue[0], rand())
			//Bloqueia a fila de hamburguer para modificar a fila ou enviar para a fila de refrigerante com prioridade
			hamburguerQueue.mux.Lock()
			//Caso haja proximo pedido de refrigerante
			if hamburguerQueue.Queue[0].Next != nil {
				//Caso haja espaço na fila de prioridades de refrigerante
				if len(sodaPriorityQueue.Queue) < MaxQueueSodaP {
					//Bloqueia a fila para inserção
					sodaPriorityQueue.mux.Lock()
					sodaPriorityQueue.Queue = append(sodaPriorityQueue.Queue, *hamburguerQueue.Queue[0].Next)
					sodaPriorityQueue.mux.Unlock()
				} else {
					fmt.Println("Fila de prioridade de refrigerante cheia")
					//Encerra o pedido por descarte
					wait.Done()
				}
			} else { //Se não tiver proximo pedido
				//Encerra o pedido
				wait.Done()
			}
			//Anda com a fila de hamburgueres
			hamburguerQueue.Queue = hamburguerQueue.Queue[1:]
			hamburguerQueue.mux.Unlock()
		}
	}
}

func sodaHandler() {
	rand := randomNumberGenerator(1103515245, 12345, 1<<31, 5)
	for {
		//Caso haja alguem na fila de prioridade
		if len(sodaPriorityQueue.Queue) > 0 {
			//faz o refrigerante
			makeSoda(sodaPriorityQueue.Queue[0], rand())
			//bloqueia a fila para modificação
			sodaPriorityQueue.mux.Lock()
			sodaPriorityQueue.Queue = sodaPriorityQueue.Queue[1:]
			sodaPriorityQueue.mux.Unlock()
			//Encerra o pedido
			wait.Done()
		} else if len(sodaQueue.Queue) > 0 {
			//faz o refrigerante
			makeSoda(sodaQueue.Queue[0], rand())
			//bloqueia fila para modificação
			sodaQueue.mux.Lock()
			sodaQueue.Queue = sodaQueue.Queue[1:]
			sodaQueue.mux.Unlock()
			//Encerra pedido
			wait.Done()
		}
	}
}

func processRequest(order Order) {
	if order.Type == "hamburguer" {
		//Se houver lugar na fila de hamburguer
		if len(hamburguerQueue.Queue) < MaxQueueBurguer {
			//bloqueia a fila para inserção
			hamburguerQueue.mux.Lock()
			hamburguerQueue.Queue = append(hamburguerQueue.Queue, order)
			hamburguerQueue.mux.Unlock()
		} else {
			//Se não houver, encerra o pedido
			fmt.Println("Pedido de hamburguer descartado")
			wait.Done()
		}
	} else if len(sodaQueue.Queue) < MaxQueueSoda { // Se houver lugar na fila de refrigerante
		//Bloqueia a fila apra inserção
		sodaQueue.mux.Lock()
		sodaQueue.Queue = append(sodaQueue.Queue, order)
		sodaQueue.mux.Unlock()
	} else {
		//Se não houver, encerra o pedido
		fmt.Println("Pedido de Refrigerante descartado")
		wait.Done()
	}
}

func makeBurguer(order Order, rand float64) {
	time.Sleep(generateRandomTime(rand, 20, time.Hour))
	//fmt.Println(order.Type)
}

func makeSoda(order Order, rand float64) {
	time.Sleep(generateRandomTime(rand, 60, time.Hour))
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

func generateRandomTime(u, delta float64, duration time.Duration) time.Duration {
	return time.Duration(generateRandomExp(delta, u)*float64(duration)) / 3600
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
