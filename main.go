package main

import (
	"fmt"
	"log"
)

type Step interface {
	Do() error
	Compensate()
}

type Saga struct {
	Steps []Step
}

type PaymentService struct {
	Balance float64
	Cost    float64
}

type InventoryService struct {
	Products   int
	OrderCount int
}

type ShippingService struct {
	AvailableMethods []string
	Method           string
}

func (s *PaymentService) Do() error {
	if s.Cost > s.Balance {
		return fmt.Errorf("not enough money")
	}

	s.Balance -= s.Cost
	return nil
}

func (s *PaymentService) Compensate() {
	s.Balance += s.Cost
}

func (s *InventoryService) Do() error {
	if s.OrderCount > s.Products {
		return fmt.Errorf("not enough products")
	}

	s.Products -= s.OrderCount
	return nil
}

func (s *InventoryService) Compensate() {
	s.Products += s.OrderCount
}

func (s *ShippingService) Do() error {
	for i, v := range s.AvailableMethods {
		if v == s.Method {
			s.AvailableMethods = append(s.AvailableMethods[:i], s.AvailableMethods[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("method is not supported")
}

func (s *ShippingService) Compensate() {
	s.AvailableMethods = append(s.AvailableMethods, s.Method)
}

func (s *Saga) Rollback(errIndex int) {
	for i := errIndex - 1; i >= 0; i-- {
		log.Printf("rollback of completed step %d\n", i+1)
		s.Steps[i].Compensate()
	}
}

func (s *Saga) Execute() error {
	for i, step := range s.Steps {
		if err := step.Do(); err != nil {
			log.Printf("error occurred during execution of step %d: %v\n", i+1, err)
			s.Rollback(i)
			return err
		}
		log.Printf("step %d completed successfully", i+1)
	}
	log.Println("all steps completed successfully")
	return nil
}

func main() {
	paymentService := &PaymentService{Balance: 1100, Cost: 150}
	inventoryService := &InventoryService{Products: 15, OrderCount: 10}
	shippingService := &ShippingService{AvailableMethods: []string{"air", "plane", "truck"}, Method: "truck"}

	saga := Saga{[]Step{paymentService, inventoryService, shippingService}}

	log.Printf("Before:\n%v\n%v\n%v", paymentService, inventoryService, shippingService)
	if err := saga.Execute(); err != nil {
		log.Println("workflow ended unsuccessfully")
	} else {
		log.Println("workflow successfully completed")
	}
	log.Printf("After:\n%v\n%v\n%v", paymentService, inventoryService, shippingService)
}
