package conveyer

import (
    "context"
    "errors"
    "fmt"
    "sync"

    "golang.org/x/sync/errgroup"
)

var (
    ErrChanNotFound = errors.New("канал не найден")
)

const undefined = "undefined"

type Conveyer struct {
    size int

    mu       sync.RWMutex
    channels map[string]chan string
    handlers []func(ctx context.Context) error
}

func New(size int) *Conveyer {
    return &Conveyer{
        size:     size,
        channels: make(map[string]chan string),
        handlers: make([]func(context.Context) error, 0),
    }
}

func (c *Conveyer) register(name string) chan string {
    if ch, ok := c.channels[name]; ok {
        return ch
    }

    ch := make(chan string, c.size)
    c.channels[name] = ch

    return ch
}

func (c *Conveyer) RegisterDecorator(
    fn func(ctx context.Context, input chan string, output chan string) error,
    input string,
    output string,
) {
    c.mu.Lock()
    defer c.mu.Unlock()

    in := c.register(input)
    out := c.register(output)

    c.handlers = append(c.handlers, func(ctx context.Context) error {
        return fn(ctx, in, out)
    })
}

func (c *Conveyer) RegisterMultiplexer(
    fn func(ctx context.Context, inputs []chan string, output chan string) error,
    inputs []string,
    output string,
) {
    c.mu.Lock()
    defer c.mu.Unlock()

    inChans := make([]chan string, 0, len(inputs))
    for _, name := range inputs {
        inChans = append(inChans, c.register(name))
    }

    out := c.register(output)

    c.handlers = append(c.handlers, func(ctx context.Context) error {
        return fn(ctx, inChans, out)
    })
}

func (c *Conveyer) RegisterSeparator(
    fn func(ctx context.Context, input chan string, outputs []chan string) error,
    input string,
    outputs []string,
) {
    c.mu.Lock()
    defer c.mu.Unlock()

    in := c.register(input)

    outChans := make([]chan string, 0, len(outputs))
    for _, name := range outputs {
        outChans = append(outChans, c.register(name))
    }

    c.handlers = append(c.handlers, func(ctx context.Context) error {
        return fn(ctx, in, outChans)
    })
}

func closeChannelSafe(ch chan string) {
    defer func() {
        _ = recover()
    }()

    close(ch)
}

func (c *Conveyer) closeAllChannels() {
    c.mu.Lock()
    defer c.mu.Unlock()

    for _, ch := range c.channels {
        closeChannelSafe(ch)
    }
}

func (c *Conveyer) Run(ctx context.Context) error {
    c.mu.RLock()
    handlers := make([]func(context.Context) error, len(c.handlers))
    copy(handlers, c.handlers)
    c.mu.RUnlock()

    group, ctxWithCancel := errgroup.WithContext(ctx)

    for _, h := range handlers {
        fn := h

        group.Go(func() error {
            return fn(ctxWithCancel)
        })
    }

    err := group.Wait()

    c.closeAllChannels()

    if err != nil && !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
        return fmt.Errorf("конвейер приостановлен: %w", err)
    }

    return nil
}

func (c *Conveyer) Send(input string, data string) error {
    c.mu.RLock()
    ch, ok := c.channels[input]
    c.mu.RUnlock()

    if !ok {
        return ErrChanNotFound
    }

    ch <- data
    return nil
}

func (c *Conveyer) Recv(output string) (string, error) {
    c.mu.RLock()
    ch, ok := c.channels[output]
    c.mu.RUnlock()

    if !ok {
        return "", ErrChanNotFound
    }

    v, ok := <-ch
    if !ok {
        return undefined, nil
    }

    return v, nil
}
