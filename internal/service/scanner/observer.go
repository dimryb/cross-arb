package scanner

import (
	i "github.com/dimryb/cross-arb/internal/interface"
)

type OpportunityHandler func(Opportunity)

func LogOpportunities(logger i.Logger) OpportunityHandler {
	return func(opp Opportunity) {
		logger.Infof("Арбитраж %s: BUY %s @ %.4f → SELL %s @ %.4f (%.4f %%)",
			opp.Pair, opp.BuyOn, opp.BuyPrice, opp.SellOn, opp.SellPrice, opp.SpreadPct)
	}
}

func SubscribeAndHandle(
	scanner *Scanner,
	pairs []string,
	handler OpportunityHandler,
) {
	for _, pair := range pairs {
		ch, _ := scanner.Subscribe(pair, 10) // буфер 10 — можно тоже в конфиг
		go func(_ string, c <-chan Opportunity) {
			for opp := range c {
				handler(opp)
			}
		}(pair, ch)
	}
}
