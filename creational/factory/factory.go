package factory

// Transport is the common interface for all types of transport
type Transport interface {
	Deliver() string
}

// Truck is a concrete implementation of Transport
type Truck struct{}

// Deliver implements the Transport interface for Truck
func (t *Truck) Deliver() string {
	return "Delivering by land in a truck"
}

// Ship is a concrete implementation of Transport
type Ship struct{}

// Deliver implements the Transport interface for Ship
func (s *Ship) Deliver() string {
	return "Delivering by sea in a ship"
}

// LogisticsService is the creator interface
type LogisticsService interface {
	CreateTransport() Transport
	PlanDelivery() string
}

// RoadLogistics is a concrete creator implementing LogisticsService
type RoadLogistics struct{}

// CreateTransport implements the factory method for RoadLogistics
func (r *RoadLogistics) CreateTransport() Transport {
	return &Truck{}
}

// PlanDelivery uses the factory method to create and use a transport
func (r *RoadLogistics) PlanDelivery() string {
	transport := r.CreateTransport()
	return "Road logistics: " + transport.Deliver()
}

// SeaLogistics is a concrete creator implementing LogisticsService
type SeaLogistics struct{}

// CreateTransport implements the factory method for SeaLogistics
func (s *SeaLogistics) CreateTransport() Transport {
	return &Ship{}
}

// PlanDelivery uses the factory method to create and use a transport
func (s *SeaLogistics) PlanDelivery() string {
	transport := s.CreateTransport()
	return "Sea logistics: " + transport.Deliver()
}

// CreateLogistics is a function to create different logistics services based on type
func CreateLogistics(logisticsType string) LogisticsService {
	switch logisticsType {
	case "road":
		return &RoadLogistics{}
	case "sea":
		return &SeaLogistics{}
	default:
		return &RoadLogistics{} // Default to road logistics
	}
}
