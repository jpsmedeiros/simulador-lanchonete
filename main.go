package main

import (
	"fmt"
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
	BurguerTime  = 750
	IcecreamTime = 750
	SodaTime     = 750
)

func main() {
	wait = new(sync.WaitGroup)
	wait.Add(1)
	//cria um canal de pedidos
	orders = make(chan Order, 3)
	//Cria os canais de serviço
	hamburguers, icecreams, soda = make(chan Order, 5), make(chan Order, 5), make(chan Order, 5)

	//Gera pedidos
	//Precisa consumir de uma lista
	go func() {
		for {
			orders <- Order{Type: "hamburguer", Next: &Order{Type: "refrigerante", Next: &Order{Type: "sorvete"}}}
			time.Sleep(500 * time.Millisecond)
		}
	}()

	//Consome pedidos de hamburguer
	go hamburguerHandler()
	//Consome pedidos de sorvete
	go icecreamHandler()
	//Conso os pedidos de refrigerante
	go sodaHandler()

	//Fecha os canais de comunicação quando acabar
	go func() {
		//Lida com os pedidos para separa-los
		for item := range orders {
			wait.Add(1)
			switch item.Type {
			case "hamburguer":
				hamburguers <- item
			case "refrigerante":
				soda <- item
			case "sorvete":
				icecreams <- item
			default:
				wait.Done()
			}
		}
	}()

	wait.Wait()

}

func hamburguerHandler() {
	priority := make([]Order, 0, 5)
	normal := make([]Order, 0, 5)
	var h Order
	go func() {
		for {
			//Tem alguem na fila de prioridade
			if len(priority) > 0 {
				h = priority[0]
				priority = priority[1:]
			} else if len(normal) > 0 { // tem alguem na fila comum
				h = normal[0]
				normal = normal[1:]
			}
			if h.Type != "" {
				makeBurguer(h)
				if h.Next != nil {
					soda <- *h.Next
				} else {
					wait.Done()
				}
			}
		}
	}()

	for order := range hamburguers {
		if order.Priority {
			priority = append(priority, order)
		} else {
			normal = append(normal, order)
		}
	}
}

func icecreamHandler() {
	priority := make([]Order, 0, 5)
	normal := make([]Order, 0, 5)
	var i Order
	go func() {
		for {
			//Tem alguem na fila de prioridade
			if len(priority) > 0 {
				i = priority[0]
				priority = priority[1:]
			} else if len(normal) > 0 { // tem alguem na fila comum
				i = normal[0]
				normal = normal[1:]
			}
			if i.Type != "" {
				makeIcecream(i)
				wait.Done()
			}
		}
	}()

	for order := range icecreams {
		if order.Priority {
			priority = append(priority, order)
		} else {
			normal = append(normal, order)
		}
	}
}

func sodaHandler() {
	priority := make([]Order, 0, 5)
	normal := make([]Order, 0, 5)
	var s Order
	go func() {
		for {
			//Tem alguem na fila de prioridade
			if len(priority) > 0 {
				s = priority[0]
				priority = priority[1:]
			} else if len(normal) > 0 { // tem alguem na fila comum
				s = normal[0]
				normal = normal[1:]
			}
			if s.Type != "" {
				makeSoda(s)
				if s.Next != nil {
					icecreams <- *s.Next
				} else {
					wait.Done()
				}
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

func makeIcecream(order Order) {
	fmt.Println(order.Type)
	time.Sleep(IcecreamTime * time.Millisecond)
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
