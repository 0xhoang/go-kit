package fsm

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

const (
	CreatingOrder     StateType = "CreatingOrder"
	OrderFailed       StateType = "OrderFailed"
	OrderPlaced       StateType = "OrderPlaced"
	ChargingCard      StateType = "ChargingCard"
	TransactionFailed StateType = "TransactionFailed"
	OrderShipped      StateType = "OrderShipped"

	CreateOrder     EventType = "CreateOrder"
	FailOrder       EventType = "FailOrder"
	PlaceOrder      EventType = "PlaceOrder"
	ChargeCard      EventType = "ChargeCard"
	FailTransaction EventType = "FailTransaction"
	ShipOrder       EventType = "ShipOrder"
)

type OrderCreationContext struct {
	items []string
	err   error
}

func (c *OrderCreationContext) String() string {
	return fmt.Sprintf("OrderCreationContext [ items: %s, err: %v ]", strings.Join(c.items, ","), c.err)
}

type CreatingOrderAction struct{}

func (a *CreatingOrderAction) Execute(eventCtx EventContext) (EventType, error) {
	order := eventCtx.(*OrderCreationContext)
	fmt.Println("Validating, order:", order)
	if len(order.items) == 0 {
		order.err = errors.New("Insufficient number of items in order")
		return FailOrder, nil
	}

	return PlaceOrder, nil
}

type OrderFailedAction struct{}

func (a *OrderFailedAction) Execute(eventCtx EventContext) (EventType, error) {
	order := eventCtx.(*OrderCreationContext)
	fmt.Println("Order failed, err:", order.err)
	return NoOp, nil
}

type OrderPlacedAction struct{}

func (a *OrderPlacedAction) Execute(eventCtx EventContext) (EventType, error) {
	order := eventCtx.(*OrderCreationContext)
	fmt.Println("Order placed, items:", order.items)
	return NoOp, nil
}

type OrderShipmentContext struct {
	cardNumber string
	address    string
	err        error
}

func (c *OrderShipmentContext) String() string {
	return fmt.Sprintf("OrderShipmentContext [ cardNumber: %s, address: %s, err: %v ]",
		c.cardNumber, c.address, c.err)
}

type ChargingCardAction struct{}

func (a *ChargingCardAction) Execute(eventCtx EventContext) (EventType, error) {
	shipment := eventCtx.(*OrderShipmentContext)
	fmt.Println("Validating card, shipment:", shipment)
	if shipment.cardNumber == "" {
		shipment.err = errors.New("Card number is invalid")
		return FailTransaction, nil
	}
	return ShipOrder, nil
}

type TransactionFailedAction struct{}

func (a *TransactionFailedAction) Execute(eventCtx EventContext) (EventType, error) {
	shipment := eventCtx.(*OrderShipmentContext)
	fmt.Println("Transaction failed, err:", shipment.err)
	return NoOp, nil
}

type OrderShippedAction struct{}

func (a *OrderShippedAction) Execute(eventCtx EventContext) (EventType, error) {
	shipment := eventCtx.(*OrderShipmentContext)
	fmt.Println("Order shipped, address:", shipment.address)
	return NoOp, nil
}

func NewOrderFSM() *StateMachine {
	return &StateMachine{
		States: States{
			Default: State{
				Events: Events{
					CreateOrder: CreatingOrder,
				},
			},
			CreatingOrder: State{
				Action: &CreatingOrderAction{},
				Events: Events{
					FailOrder:  OrderFailed,
					PlaceOrder: OrderPlaced,
				},
			},
			OrderFailed: State{
				Action: &OrderFailedAction{},
				Events: Events{
					CreateOrder: CreatingOrder,
				},
			},
			OrderPlaced: State{
				Action: &OrderPlacedAction{},
				Events: Events{
					ChargeCard: ChargingCard,
				},
			},
			ChargingCard: State{
				Action: &ChargingCardAction{},
				Events: Events{
					FailTransaction: TransactionFailed,
					ShipOrder:       OrderShipped,
				},
			},
			TransactionFailed: State{
				Action: &TransactionFailedAction{},
				Events: Events{
					ChargeCard: ChargingCard,
				},
			},
			OrderShipped: State{
				Action: &OrderShippedAction{},
			},
		},
	}
}
