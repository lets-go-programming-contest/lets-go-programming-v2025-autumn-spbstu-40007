package handlers

func PrefixDecoratorFunc(
	ctx context.Context,
	input chan string,
	output chan string,
) error