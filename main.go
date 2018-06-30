package main

import (
	"fmt"
	"math"
	"sync"
	"time"
)

type Order struct {
	Type     string
	Next     *Order
	Priority bool
}

var (
	orders, hamburguers, icecreams, soda chan Order
	wait                                 *sync.WaitGroup
	totalTime                            uint
)

const (
	BurguerTime = 750
	SodaTime    = 750
)

func main() {
	wait = new(sync.WaitGroup)
	wait.Add(1)
	//cria um canal de pedidos
	//Cria os canais de serviço
	orders, hamburguers, soda = make(chan Order, 3), make(chan Order, 5), make(chan Order, 5)

	//Gera pedidos
	//Precisa consumir de uma lista
	go func() {
		var duration float64
		rand := generateRandomNumber(1103515245, 12345, 1<<31, 3)
		for {
			duration = generateRandomExp(5, rand()) * float64(time.Minute)
			orders <- Order{Type: "hamburguer", Next: &Order{Type: "refrigerante"}}
			time.Sleep(time.Duration(duration))
		}
	}()

	//Consome pedidos de hamburguer
	go hamburguerHandler()
	//Conso os pedidos de refrigerante
	go sodaHandler()

	//Fecha os canais de comunicação quando acabar
	go func() {
		//Lida com os pedidos para separa-los
		for item := range orders {
			switch item.Type {
			case "hamburguer":
				hamburguers <- item
			case "refrigerante":
				soda <- item
			default:
				fmt.Println("Aguardando pedidos")
			}
		}
	}()

	wait.Wait()

}

func hamburguerHandler() {
	normal := make([]Order, 0, 5)
	go func() {
		for {
			//Tem alguem na fila de prioridade
			if len(normal) > 0 { // tem alguem na fila comum
				if normal[0].Type != "" {
					makeBurguer(normal[0])
					if normal[0].Next != nil {
						normal[0].Priority = true
						soda <- *normal[0].Next
					}
					normal = normal[1:]
				}
			}
		}
	}()

	for order := range hamburguers {
		normal = append(normal, order)
	}
}

func sodaHandler() {
	priority := make([]Order, 0, 5)
	normal := make([]Order, 0, 5)
	go func() {
		for {
			//Tem alguem na fila de prioridade
			if len(priority) > 0 {
				makeSoda(priority[0])
				priority = priority[1:]
			} else if len(normal) > 0 { // tem alguem na fila comum
				makeSoda(normal[0])
				normal = normal[1:]
			}
		}
	}()
	for order := range soda {
		if order.Priority {
			priority = append(priority, order)
		} else {
			normal = append(normal, order)
		}
	}
}

func makeBurguer(order Order) {
	fmt.Println(order.Type)
	time.Sleep(BurguerTime * time.Millisecond)
}

func makeSoda(order Order) {
	fmt.Println(order.Type)
	time.Sleep(SodaTime * time.Millisecond)
}

func generateRandomNumber(a, c, m, seed uint32) func() float64 {
	// usage: rand := generateRandomNumber(1103515245, 12345, 1<<31, 0)
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
