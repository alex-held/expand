package expand

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	mapset "github.com/deckarep/golang-set"

	"github.com/alex-held/expand/graph"
)

const (
	dOLLAR        = "$"
	dOUBLE_DOLLAR = "$$"
)

const (
	dOLLAR_TOKEN_REPLACEMENT        = "##<&DOLLAR_TOKEN&>##"
	dOUBLE_DOLLAR_TOKEN_REPLACEMENT = "##<&DOUBLE_DOLLAR_TOKEN&>##"
)

type Resolver interface {
	Resolve() (resolved map[string]string, err error)
}

func isExpanded(val string) bool {
	val = strings.ReplaceAll(val, dOUBLE_DOLLAR, dOUBLE_DOLLAR_TOKEN_REPLACEMENT)
	if strings.Contains(val, dOLLAR) {
		return false
	}
	return true
}

func parseUnexpanded(raw string) (unexpanded []string) {
	unexpanded = []string{}
	val := strings.ReplaceAll(raw, dOUBLE_DOLLAR, dOUBLE_DOLLAR_TOKEN_REPLACEMENT)
	req := make(map[string]struct{})

	if strings.Contains(val, dOLLAR) {
		r := regexp.MustCompile(`\$([\w_]*)`)
		unmatched := r.FindAllString(val, -1)
		for _, u := range unmatched {
			if _, ok := req[u]; ok {
				continue
			}
			unexpanded = append(unexpanded, strings.TrimPrefix(u, "$"))
			req[u] = struct{}{}
		}
	}

	return unexpanded
}

type resolver struct {
	env        map[string]string
	unresolved map[string]*unres
}

func (r *resolver) Resolve() (resolved map[string]string, err error) {
	var nodes []*graph.Node
	for key, _ := range r.env {
		nodes = append(nodes, graph.NewNode(key))
	}
	for key, u := range r.unresolved {
		nodes = append(nodes, graph.NewNode(key, u.depends...))
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

		if val, ok := r.unresolved[node.Name]; ok {

			depsSet := mapset.NewSet()
			for _, dep := range val.depends {
				depsSet.Add(dep)
			}

			if !depsSet.IsSubset(envSet) {
				log.Printf("[ERROR ] unable to resolve '%s'", node.Name)
				return r.env, errors.New(fmt.Sprintf("unable to resolve '%s'", node.Name))
			}

			log.Printf("resolving '%s' with raw value '%s' and deps %v\n", val.id, val.val, val.depends)
			resolved := val.Resolve(r.env)
			log.Printf("[%d] resolved '%s' with value '%s'\n", i, node.Name, resolved)

			envSet.Add(val.id)
			r.env[val.id] = resolved

			continue
		}
		log.Printf("[ERROR] unable to find var with key '%s'; available=%v", node.Name, r.env)
		return r.env, errors.New(fmt.Sprintf("[ERROR] unable to find var with key '%s'; available=%v", node.Name, r.env))
	}

	return r.env, nil
}

func (u *unres) Resolve(vals map[string]string) string {
	if isExpanded(u.val) {
		return u.val
	}

	val := strings.ReplaceAll(u.val, dOUBLE_DOLLAR, dOUBLE_DOLLAR_TOKEN_REPLACEMENT)

	for _, dep := range u.depends {
		if subst, ok := vals[dep]; ok {
			newVal := strings.ReplaceAll(val, fmt.Sprintf("$%s", dep), subst)
			val = newVal
		}
		newVal := strings.ReplaceAll(val, fmt.Sprintf("$%s", dep), "")
		val = newVal
	}

	val = strings.ReplaceAll(val, dOUBLE_DOLLAR_TOKEN_REPLACEMENT, dOUBLE_DOLLAR)
	return val
}

type unres struct {
	id      string
	val     string
	depends []string
}

// NewResolverWithEnvironment create a new Resolver with variables from os.Environ and the provided variables
//
// The provided variables overwrite the os.Environ variables
func NewResolverWithEnvironment(m map[string]string) Resolver {
	vars := map[string]string{}
	for _, env := range os.Environ() {
		i := strings.Index(env, "=")
		vars[env[:i]] = env[i+1:]
	}
	for key, val := range m {
		vars[key] = val
	}

	return NewResolver(vars)
}

// NewResolver create a new Resolver with the provided variables
func NewResolver(vars map[string]string) Resolver {
	r := &resolver{
		env:        map[string]string{},
		unresolved: map[string]*unres{},
	}

	for key, val := range vars {
		if isExpanded(val) {
			log.Printf("adding initial expanded: key=%s, val=%s\n", key, val)
			r.env[key] = val
		}

		unresolved := parseUnexpanded(val)
		u := unres{
			id:      key,
			val:     val,
			depends: unresolved,
		}
		log.Printf("adding unres key=%s; val=%s; deps:%v\n", u.id, u.val, u.depends)
		r.unresolved[key] = &u
	}

	return r
}
