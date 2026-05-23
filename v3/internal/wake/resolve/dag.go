package resolve

import (
	"fmt"
	"sort"

	"github.com/wailsapp/wails/v3/internal/wake/ast"
)

type DAG struct {
	Tasks    []*ast.Task
	Edges    map[string][]string
	InDegree map[string]int
	Order    []string
}

func BuildDAG(tf *ast.Taskfile, target string) (*DAG, error) {
	if _, ok := tf.Tasks[target]; !ok {
		return nil, fmt.Errorf("wake: task %q not found", target)
	}

	visited := make(map[string]bool)
	var collect func(name string) error
	collect = func(name string) error {
		if visited[name] {
			return nil
		}
		visited[name] = true
		task := tf.Tasks[name]
		if task == nil {
			return fmt.Errorf("wake: task %q not found", name)
		}
		for _, dep := range task.Deps {
			if err := collect(dep.Task); err != nil {
				return err
			}
		}
		return nil
	}

	if err := collect(target); err != nil {
		return nil, err
	}

	dag := &DAG{
		Edges:    make(map[string][]string),
		InDegree: make(map[string]int),
	}

	for name := range visited {
		task := tf.Tasks[name]
		dag.Tasks = append(dag.Tasks, task)
		dag.InDegree[name] = 0
	}

	for _, task := range dag.Tasks {
		for _, dep := range task.Deps {
			dag.Edges[dep.Task] = append(dag.Edges[dep.Task], task.Name)
			dag.InDegree[task.Name]++
		}
	}

	if err := detectCycle(dag); err != nil {
		return nil, err
	}

	dag.Order = topologicalSort(dag)
	return dag, nil
}

func detectCycle(dag *DAG) error {
	white, gray, black := 0, 1, 2
	color := make(map[string]int)
	for _, t := range dag.Tasks {
		color[t.Name] = white
	}

	var dfs func(string) error
	dfs = func(u string) error {
		color[u] = gray
		for _, v := range dag.Edges[u] {
			if color[v] == gray {
				return fmt.Errorf("wake: cycle detected involving task %q", v)
			}
			if color[v] == white {
				if err := dfs(v); err != nil {
					return err
				}
			}
		}
		color[u] = black
		return nil
	}

	for _, t := range dag.Tasks {
		if color[t.Name] == white {
			if err := dfs(t.Name); err != nil {
				return err
			}
		}
	}
	return nil
}

func topologicalSort(dag *DAG) []string {
	inDeg := make(map[string]int)
	for k, v := range dag.InDegree {
		inDeg[k] = v
	}

	var queue []string
	for name, deg := range inDeg {
		if deg == 0 {
			queue = append(queue, name)
		}
	}
	sort.Strings(queue)

	var order []string
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		order = append(order, u)

		for _, v := range dag.Edges[u] {
			inDeg[v]--
			if inDeg[v] == 0 {
				queue = append(queue, v)
			}
		}
	}

	return order
}
