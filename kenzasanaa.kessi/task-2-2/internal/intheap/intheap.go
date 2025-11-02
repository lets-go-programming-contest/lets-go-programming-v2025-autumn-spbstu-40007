package intheap

type CustomHeap []int

func (h *CustomHeap) Size() int {
	return len(*h)
}

func (h *CustomHeap) Compare(i, j int) bool {
	if i < 0 || j < 0 || i >= h.Size() || j >= h.Size() {
		panic("index beyond heap boundaries")
	}

	return (*h)[i] > (*h)[j]
}

func (h *CustomHeap) Exchange(i, j int) {
	if i < 0 || j < 0 || i >= h.Size() || j >= h.Size() {
		panic("index beyond heap boundaries")
	}

	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *CustomHeap) Add(value interface{}) {
	intValue, valid := value.(int)
	if !valid {
		panic("non-integer value provided")
	}

	*h = append(*h, intValue)
}

func (h *CustomHeap) Remove() interface{} {
	if h.Size() == 0 {
		return nil
	}

	current := *h
	lastIndex := len(current) - 1
	element := current[lastIndex]
	*h = current[:lastIndex]

	return element
}

func (h *CustomHeap) Len() int {
	return len(*h)
}

func (h *CustomHeap) Less(i, j int) bool {
	return (*h)[i] > (*h)[j]
}

func (h *CustomHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *CustomHeap) Push(x interface{}) {
	h.Add(x)
}

func (h *CustomHeap) Pop() interface{} {
	return h.Remove()
}
