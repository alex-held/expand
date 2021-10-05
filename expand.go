package expand

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	mapset "github.com/deckarep/golang-set"

	"github.com/alex-held/expand/internal/graph"
	"github.com/alex-held/expand/internal/parse"
)

// Expander expands registered variables recursively
type Expander interface {
	Expand() (resolved Expansions, err error)
}

// Expansions contains the expanded values accessible by its id
type Expansions map[string]string

// Get returns the value for an id and whether a value was found
func (r Expansions) Get(id string) (res string, ok bool) { res, ok = r[id]; return res, ok }

// MustGet returns the value for an id
func (r Expansions) MustGet(id string) string { return r[id] }

// Contains returns whether Expansions contain a value for an id
func (r Expansions) Contains(id string) bool { _, ok := r.Get(id); return ok }

type resolver struct {
	env        map[string]string
	unexpanded map[string]*parse.Unexpanded
}

// Expand returns expanded Expansions and maybe an error
func (r *resolver) Expand() (resolved Expansions, err error) {
	var nodes []*graph.Node
	for key, _ := range r.env {
		nodes = append(nodes, graph.NewNode(key))
	}
	for key, u := range r.unexpanded {
		nodes = append(nodes, graph.NewNode(key, u.Depends...))
	}

	g := graph.New(nodes...)
	resolvedGraph, err := g.Resolve()
	if err != nil {
		return r.env, err
	}

	resolvedNodes := resolvedGraph.Nodes()

	log.Println("resolve order:")
	for i, node := range resolvedNodes {
		log.Printf("[%d] %s\n", i, node.Name)
	}
	envSet := mapset.NewSet()
	for key := range r.env {
		envSet.Add(key)
	}

	for i, node := range resolvedNodes {

		if val, ok := r.env[node.Name]; ok {
			log.Printf("[%d] resolved '%s' with value '%s'\n", i, node.Name, val)
			continue
		}

		if val, ok := r.unexpanded[node.Name]; ok {

			depsSet := mapset.NewSet()
			for _, dep := range val.Depends {
				depsSet.Add(dep)
			}

			if !depsSet.IsSubset(envSet) {
				log.Printf("[ERROR ] unable to resolve '%s'", node.Name)
				return r.env, errors.New(fmt.Sprintf("unable to resolve '%s'", node.Name))
			}

			log.Printf("resolving '%s' with raw value '%s' and deps %v\n", val.ID, val.RawValue, val.Depends)
			resolved := val.Expand(r.env)
			log.Printf("[%d] resolved '%s' with value '%s'\n", i, node.Name, resolved)

			envSet.Add(val.ID)
			r.env[val.ID] = resolved

			continue
		}
		log.Printf("[ERROR] unable to find var with key '%s'; available=%v", node.Name, r.env)
		return r.env, errors.New(fmt.Sprintf("[ERROR] unable to find var with key '%s'; available=%v", node.Name, r.env))
	}

	return r.env, nil
}

// NewExpanderWithEnvironment create a new Expander with variables from os.Environ and the provided variables
//
// The provided variables overwrite the os.Environ variables
func NewExpanderWithEnvironment(m map[string]string) Expander {
	vars := map[string]string{}
	for _, env := range os.Environ() {
		i := strings.Index(env, "=")
		vars[env[:i]] = env[i+1:]
	}
	for key, val := range m {
		vars[key] = val
	}

	return NewExpander(vars)
}

// NewExpander create a new Expander with the provided variables
func NewExpander(vars map[string]string) Expander {
	r := &resolver{
		env:        map[string]string{},
		unexpanded: map[string]*parse.Unexpanded{},
	}

	for key, val := range vars {
		if parse.IsExpanded(val) {
			log.Printf("adding initial expanded: key=%s, val=%s\n", key, val)
			r.env[key] = val
		}

		unresolved := parse.ParseUnexpanded(val)
		u := parse.Unexpanded{
			ID:       key,
			RawValue: val,
			Depends:  unresolved,
		}
		log.Printf("adding unres key=%s; val=%s; deps:%v\n", u.ID, u.RawValue, u.Depends)
		r.unexpanded[key] = &u
	}

	return r
}

// ExpandWithEnvironment expands vars using os.Environ and returns Expansions and an error
func ExpandWithEnvironment(vars map[string]string) (Expansions, error) {
	r := NewExpanderWithEnvironment(vars)
	return r.Expand()
}

// Expand expands vars as Expansions and an error
func Expand(vars map[string]string) (Expansions, error) {
	r := NewExpander(vars)
	return r.Expand()
}
