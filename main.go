package main

import (
	"fmt"
	"sync"
	"time"
)

type Order struct {
	Type string
	Next *Order
}

var (
	orders, hamburguers, icecreams, soda chan Order
	wait                                 *sync.WaitGroup
)

func main() {
	wait = new(sync.WaitGroup)
	wait.Add(1)
	//cria um canal de pedidos
	orders = make(chan Order)
	//Cria os canais de serviço
	hamburguers, icecreams, soda = make(chan Order, 5), make(chan Order, 5), make(chan Order, 5)

	//Gera pedidos
	//Precisa consumir de uma lista
	go func() {
		for {
			orders <- Order{Type: "hamburguer", Next: &Order{Type: "sorvete", Next: &Order{Type: "refrigerante"}}}
			time.Sleep(1000 * time.Millisecond)
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
	for order := range hamburguers {
		fmt.Println(order)
		time.Sleep(2000 * time.Millisecond)
		if order.Next != nil {
			orders <- *order.Next
		}
		wait.Done()
	}
}

func icecreamHandler() {
	for order := range icecreams {
		fmt.Println(order)
		time.Sleep(2000 * time.Millisecond)
		if order.Next != nil {
			orders <- *order.Next
		}
		wait.Done()
	}
}

func sodaHandler() {
	for order := range soda {
		fmt.Println(order)
		time.Sleep(2000 * time.Millisecond)
		if order.Next != nil {
			orders <- *order.Next
		}
		wait.Done()
	}
}
