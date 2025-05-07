package core

type Graph map[string][]string

func NewGraph() Graph {
	return Graph{}
}

func (this Graph) AddNode(node string) {
	if this.NodeExist(node) {
		return
	}
	this[node] = make([]string, 0)
}

func (this Graph) NodeExist(node string) bool {
	_, ok := this[node]
	return ok
}

func (this Graph) AddEdge(from, to string) {
	if this.NodeExist(from) && this.NodeExist(to) {
		this[from] = append(this[from], to)
		this[to] = append(this[to], from)
	}
}

func (this Graph) GetNodes() []string {
	nodes := make([]string, 0)
	for k := range this {
		nodes = append(nodes, k)
	}
	return nodes
}

func (this Graph) ConnectedComponents() [][]string {
	components := make([][]string, 0)
	visited := make(map[string]bool)
	for node := range this {
		if !visited[node] {
			nodes := make([]string, 0)
			this.dfs(node, visited, &nodes)
			components = append(components, nodes)
		}
	}
	return components
}

func (g Graph) SubGraph(node string) Graph {
	subNodes := make([]string, 0)
	visited := make(map[string]bool)

	g.dfs(node, visited, &subNodes)

	subGraph := NewGraph()
	for _, n := range subNodes {
		subGraph.AddNode(n)
		subGraph[n] = g[n]
	}
	return subGraph
}

func (g Graph) dfs(node string, visited map[string]bool, nodes *[]string) {
	if visited[node] {
		return
	}

	visited[node] = true
	*nodes = append(*nodes, node)

	for _, neighbor := range g[node] {
		g.dfs(neighbor, visited, nodes)
	}
}
