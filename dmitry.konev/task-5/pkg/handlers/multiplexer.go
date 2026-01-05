package handlers

func MultiplexerFunc(
	ctx context.Context,
	inputs []chan string,
	output chan string,
) error
