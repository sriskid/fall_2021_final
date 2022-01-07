package game

type priorityPos struct {
	Pos
	priority int
}

type pQueue []priorityPos

func (pq pQueue) Push(pos Pos, priority int) pQueue {
	newNode := priorityPos{pos, priority}
	pq = append(pq, newNode)
	newNodeindex := len(pq) - 1
	parentIndex, parent := pq.Parent(newNodeindex)
	for newNode.priority < parent.priority && newNodeindex != 0 {
		pq.Swap(newNodeindex, parentIndex)
		newNodeindex = parentIndex
		parentIndex, parent = pq.Parent(newNodeindex)
	}
	return pq
}

func (pq pQueue) Pop() (pQueue, Pos) {
	result := pq[0].Pos
	pq[0] = pq[len(pq)-1]
	pq = pq[:len(pq)-1]

	if len(pq) == 0 {
		return pq, result
	}

	index := 0
	node := pq[index]

	leftExists, leftIndex, left := pq.Left(index)
	rightExists, rightIndex, right := pq.Right(index)

	for (leftExists && node.priority > left.priority) || (rightExists && node.priority > right.priority) {
		if !rightExists || left.priority <= right.priority {
			pq.Swap(index, leftIndex)
			index = leftIndex
		} else {
			pq.Swap(index, rightIndex)
			index = rightIndex
		}
		leftExists, leftIndex, left = pq.Left(index)
		rightExists, rightIndex, right = pq.Right(index)

	}

	return pq, result
}

func (pq pQueue) Swap(i, j int) {
	temp := pq[i]
	pq[i] = pq[j]
	pq[j] = temp
}

func (pq pQueue) Parent(i int) (int, priorityPos) {
	index := (i - 1) / 2
	return index, pq[index]
}

func (pq pQueue) Left(i int) (bool, int, priorityPos) {
	index := i*2 + 1
	if index < len(pq) {
		return true, index, pq[index]
	}
	return false, 0, priorityPos{}
}

func (pq pQueue) Right(i int) (bool, int, priorityPos) {
	index := i*2 + 2
	if index < len(pq) {
		return true, index, pq[index]
	}
	return false, 0, priorityPos{}
}
