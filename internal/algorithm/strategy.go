package algorithm

import "github.com/snakoner/go-crypto-robot/internal/models"

type Strategy struct {
	Node []*StrategyElement
}

type StrategyElement struct {
	Name string
	Func func([]*models.MarketPoint) (bool, bool)
}

func NewStrategy(algos []string) *Strategy {
	strategy := &Strategy{}
	for _, name := range algos {
		element := &StrategyElement{
			Name: name,
		}

		switch name {
		case "divergence":
			element.Func = nil
			break
		case "rsi":
			element.Func = Rsi
			break
		default:
			return nil
		}

		strategy.Node = append(strategy.Node, element)
	}

	return strategy
}

func (s *Strategy) String() string {
	first := true
	ret := "["

	for _, node := range s.Node {
		if first {
			first = false
			ret += node.Name
			continue
		}
		ret += ", " + node.Name
	}

	ret += "]"

	return ret
}

func (s *Strategy) Calculate(tracker *models.TokenTracker) bool {
	for _, node := range s.Node {
		res := false
		if tracker.Long {
			res, _ = node.Func(tracker.MarketPoints)
		} else {
			_, res = node.Func(tracker.MarketPoints)
		}
		if !res {
			return false
		}
	}

	return true
}
