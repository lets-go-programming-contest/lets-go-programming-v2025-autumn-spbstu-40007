package handlers

import (
    "context"
    "fmt"
    "strings"
    "sync"
)

func PrefixDecoratorFunc(ctx context.Context, input chan string, output chan string) error {
    prefix := "decorated: "
    errorSubstring := "can't be decorated"

    for {
        select {
        case data, ok := <-input:
            if !ok {
                return nil
            }
            
            if strings.Contains(data, "no decorator") {
                return fmt.Errorf("data contains 'no decorator': %s", errorSubstring)
            }

            if strings.HasPrefix(data, prefix) {
                continue 
            }
            
            newData := prefix + data
            
            select {
            case output <- newData:
            case <-ctx.Done():
                return ctx.Err()
            }

        case <-ctx.Done():
            return ctx.Err()
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
        case data, ok := <-input:
            if !ok {
                return nil
            }
            
            targetCh := outputs[idx % len(outputs)]
            
            select {
            case targetCh <- data:
                idx++ 
            case <-ctx.Done():
                return ctx.Err()
            }
            
        case <-ctx.Done():
            return ctx.Err()
        }
    }
}

func MultiplexerFunc(ctx context.Context, inputs []chan string, output chan string) error {
    if len(inputs) == 0 {
        return nil
    }
    
    mergedInput := make(chan string) 
    var wg sync.WaitGroup
    
    for _, ch := range inputs {
        wg.Add(1)
        go func(inCh chan string) {
            defer wg.Done()
            for data := range inCh {
                select {
                case mergedInput <- data:
                case <-ctx.Done():
                    return
                }
            }
        }(ch)
    }
    
    go func() {
        wg.Wait()
        close(mergedInput) 
    }()

    for {
        select {
        case data, ok := <-mergedInput:
            if !ok {
                return nil
            }
            
            if strings.Contains(data, "no multiplexer") {
                continue 
            }
            
            select {
            case output <- data:
            case <-ctx.Done():
                return ctx.Err()
            }
            
        case <-ctx.Done():
            return ctx.Err()
        }
    }
}