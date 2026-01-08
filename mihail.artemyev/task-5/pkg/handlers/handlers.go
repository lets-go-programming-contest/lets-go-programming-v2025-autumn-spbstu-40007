package handlers

import (
    "context"
    "errors"
    "fmt"
    "strings"
    "sync"
)

var ErrPrefixDecoratorCantBeDecorated = errors.New("can't be decorated")

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
    defer close(output)

    const prefix = "decorated: "

    for {
        select {
        case <-ctx.Done():
            return nil

        case data, ok := <-input:
            if !ok {
                return nil
            }

            if strings.Contains(data, "no decorator") {
                return fmt.Errorf("%w", ErrPrefixDecoratorCantBeDecorated)
            }

            if !strings.HasPrefix(data, prefix) {
                data = prefix + data
            }

            select {
            case <-ctx.Done():
                return nil
            case output <- data:
            }
        }
    }
}

func SeparatorFunc(ctx context.Context, input chan string, outputs []chan string) error {
    if len(outputs) == 0 {
        return nil
    }


    idx := 0

    for {
        select {
        case <-ctx.Done():
            return nil

        case data, ok := <-input:
            if !ok {
                return nil
            }

            target := outputs[idx%len(outputs)]
            idx++

            select {
            case <-ctx.Done():
                return nil
            case target <- data:
            }
        }
    }
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
    if len(inputs) == 0 {
        return nil
    }

    type message struct {
        value string
        ok    bool
    }

    fanIn := make(chan message)
    var wg sync.WaitGroup

    wg.Add(len(inputs))
    for _, in := range inputs {
        ch := in

        go func(c chan string) {
            defer wg.Done()

            for {
                select {
                case <-ctx.Done():
                    return
                case v, ok := <-c:
                    if !ok {
                        return
                    }

                    fanIn <- message{value: v, ok: true}
                }
            }
        }(ch)
    }

    go func() {
        wg.Wait()
        close(fanIn)
    }()

    for {
        select {
        case <-ctx.Done():
            return nil

        case msg, ok := <-fanIn:
            if !ok {
                return nil
            }

            if !msg.ok {
                continue
            }

            if strings.Contains(msg.value, "no multiplexer") {
                continue
            }

            select {
            case <-ctx.Done():
                return nil
            case output <- msg.value:
            }
        }
    }
}
