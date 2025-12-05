package handlers

import (
    "context"
    "errors"
    "strings"
    "sync"
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case data, ok := <-input:
            if !ok {
                return nil
            }
            if strings.Contains(data, "no decorator") {
                return errors.New("can't be decorated")
            }
            newData := data
            if !strings.HasPrefix(data, "decorated: ") {
                newData = "decorated: " + data
            }
            select {
            case <-ctx.Done():
                return ctx.Err()
            case output <- newData:
            }
        }
    }
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
    var wg sync.WaitGroup
    transfer := make(chan string)

    for _, in := range inputs {
        wg.Add(1)
        go func(c chan string) {
            defer wg.Done()
            for {
                select {
                case <-ctx.Done():
                    return
                case val, ok := <-c:
                    if !ok {
                        return
                    }
                    select {
                    case <-ctx.Done():
                        return
                    case transfer <- val:
                    }
                }
            }
        }(in)
    }

    go func() {
        wg.Wait()
        close(transfer)
    }()

    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case data, ok := <-transfer:
            if !ok {
                return nil
            }
            if strings.Contains(data, "no multiplexer") {
                continue
            }
            select {
            case <-ctx.Done():
                return ctx.Err()
            case output <- data:
            }
        }
    }
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
    var i int
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case data, ok := <-input:
            if !ok {
                return nil
            }
            targetIdx := i % len(outputs)
            i++
            select {
            case <-ctx.Done():
                return ctx.Err()
            case outputs[targetIdx] <- data:
            }
        }
    }
}
