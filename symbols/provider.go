package symbols

import (
	"fmt"
	"math/rand"
	"sync"
)

var (
	ErrSymbolNotFound = fmt.Errorf("symbol not found")
)

type Provider struct {
	symbols   []string
	tickSizes map[string]float64

	sync.RWMutex
}

func NewProvider() *Provider {
	return &Provider{
		symbols:   make([]string, 0),
		tickSizes: make(map[string]float64),
	}
}

func (p *Provider) Symbols() ([]string, error) {
	if err := p.fetchIfNeeded(); err != nil {
		return nil, err
	}
	return p.symbols, nil
}

func (p *Provider) TickSize(symbol string) (float64, error) {
	if err := p.fetchIfNeeded(); err != nil {
		return 0, err
	}

	p.RLock()
	defer p.RUnlock()
	if tickSize, ok := p.tickSizes[symbol]; !ok {
		return 0, ErrSymbolNotFound
	} else {
		return tickSize, nil
	}
}

func (p *Provider) RandomSymbol() (string, error) {
	if err := p.fetchIfNeeded(); err != nil {
		return "", err
	}

	p.RLock()
	defer p.RUnlock()
	return p.symbols[rand.Intn(len(p.symbols))], nil
}

func (p *Provider) RandomSymbolNot(symbol string) (string, error) {
	if err := p.fetchIfNeeded(); err != nil {
		return "", err
	}

	p.RLock()
	defer p.RUnlock()
	for {
		s := p.symbols[rand.Intn(len(p.symbols))]
		if s != symbol {
			return s, nil
		}
	}
}

func (p *Provider) fetchIfNeeded() error {
	p.Lock()
	defer p.Unlock()

	if len(p.tickSizes) > 0 {
		return nil
	}

	syms, err := fetch()
	if err != nil {
		return err
	}

	p.symbols = make([]string, 0, len(syms))
	p.tickSizes = make(map[string]float64, len(syms))
	for _, sym := range syms {
		p.symbols = append(p.symbols, sym.Symbol)
		p.tickSizes[sym.Symbol] = sym.TickSize
	}
	return nil
}
