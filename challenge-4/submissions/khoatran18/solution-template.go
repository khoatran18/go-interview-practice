package main

import (
	"fmt"
	"sync"
)

func ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int {
	if numWorkers <= 0 {
		return map[int][]int{}
	}

	results := make(map[int][]int)
	jobs := make(chan int)
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	worker := func() {
		defer wg.Done()
		for start := range jobs {
			result := bfs(graph, start)
			mu.Lock()
			results[start] = result
			mu.Unlock()
		}
	}

	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go worker()
	}

	go func() {
		for _, query := range queries {
			jobs <- query
		}
		close(jobs)
	}()

	wg.Wait()
	return results

}

func bfs(graph map[int][]int, start int) []int {

	if len(graph[start]) == 0 {
		return []int{start}
	}

	visited := make(map[int]bool)
	queue := []int{start}
	order := []int{}

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		order = append(order, node)
		visited[node] = true

		neighbors := graph[node]
		for _, neighbor := range neighbors {
			if !visited[neighbor] {
				queue = append(queue, neighbor)
				visited[neighbor] = true
			}
		}
	}
	return order
}

func main() {
	graph := map[int][]int{
		0: {1, 2},
		1: {2, 3},
		2: {3},
		3: {4},
		4: {},
	}
	queries := []int{0, 1, 2}
	numWorkers := 2

	results := ConcurrentBFSQueries(graph, queries, numWorkers)
	for k, v := range results {
		fmt.Println(k, ":", v)
	}
}
