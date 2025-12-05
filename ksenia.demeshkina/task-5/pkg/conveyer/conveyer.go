package conveyer

import (
    "context"
    "errors"
    "fmt"
    "sync"
)

type conveyer interface {
    RegisterDecorator(fn func(ctx context.Context, input chan string, output chan string) error, input string, output string) error
    RegisterMultiplexer(fn func(ctx context.Context, inputs []chan string, output chan string) error, inputs []string, output string) error
    RegisterSeparator(fn func(ctx context.Context, input chan string, outputs []chan string) error, input string, outputs []string) error
    Run(ctx context.Context) error
    Send(input string, data string) error
    Recv(output string) (string, error)
}

type HandlerType int

const (
    DecoratorType HandlerType = iota
    MultiplexerType
    SeparatorType
)

var (
    ErrChanNotFound = errors.New("chan not found")
)

type HandlerInfo struct {
    Type      HandlerType
    Func      interface{} 
    InputIDs  []string
    OutputIDs []string
}

type Conveyer struct {
    channelSize int
    channels    map[string]chan string
    handlers    []*HandlerInfo

    mu        sync.RWMutex
    wg        sync.WaitGroup
    runCtx    context.Context
    runCancel context.CancelFunc
}

func New(size int) *Conveyer {
    return &Conveyer{
        channelSize: size,
        channels:    make(map[string]chan string),
        handlers:    make([]*HandlerInfo, 0),
    }
}

func (c *Conveyer) getOrCreateChannel(id string) chan string {
    c.mu.RLock()
    ch, exists := c.channels[id]
    c.mu.RUnlock()

    if exists {
        return ch
    }

    c.mu.Lock()
    defer c.mu.Unlock()
    if ch, exists = c.channels[id]; exists {
        return ch
    }

    ch = make(chan string, c.channelSize)
    c.channels[id] = ch
    return ch
}

func (c *Conveyer) getChannel(id string) (chan string, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    ch, ok := c.channels[id]
    return ch, ok
}

func (c *Conveyer) RegisterDecorator(
    fn func(ctx context.Context, input chan string, output chan string) error,
    input string,
    output string,
) error {
    c.getOrCreateChannel(input)
    c.getOrCreateChannel(output)

    info := &HandlerInfo{
        Type: DecoratorType,
        Func: fn,
        InputIDs: []string{input},
        OutputIDs: []string{output},
    }
    c.handlers = append(c.handlers, info)
    return nil
}

func (c *Conveyer) RegisterMultiplexer(
    fn func(ctx context.Context, inputs []chan string, output chan string) error,
    inputs []string,
    output string,
) error {
    for _, id := range inputs {
        c.getOrCreateChannel(id)
    }
    c.getOrCreateChannel(output)

    info := &HandlerInfo{
        Type: MultiplexerType,
        Func: fn,
        InputIDs: inputs,
        OutputIDs: []string{output},
    }
    c.handlers = append(c.handlers, info)
    return nil
}

func (c *Conveyer) RegisterSeparator(
    fn func(ctx context.Context, input chan string, outputs []chan string) error,
    input string,
    outputs []string,
) error {
    c.getOrCreateChannel(input)
    for _, id := range outputs {
        c.getOrCreateChannel(id)
    }

    info := &HandlerInfo{
        Type: SeparatorType,
        Func: fn,
        InputIDs: []string{input},
        OutputIDs: outputs,
    }
    c.handlers = append(c.handlers, info)
    return nil
}

func (c *Conveyer) Run(ctx context.Context) error {
    c.runCtx, c.runCancel = context.WithCancel(ctx)
    
    for _, info := range c.handlers {
        c.wg.Add(1)
        go c.runHandler(info)
    }

    <-c.runCtx.Done()
    
    c.wg.Wait() 
    
    c.closeAllChannels() 

    return c.runCtx.Err()
}

func (c *Conveyer) runHandler(info *HandlerInfo) {
    defer c.wg.Done()

    var inputChs []chan string
    for _, id := range info.InputIDs {
        inputChs = append(inputChs, c.channels[id])
    }

    var outputChs []chan string
    for _, id := range info.OutputIDs {
        outputChs = append(outputChs, c.channels[id])
    }

    var err error
    switch info.Type {
    case DecoratorType:
        fn := info.Func.(func(ctx context.Context, input chan string, output chan string) error)
        err = fn(c.runCtx, inputChs[0], outputChs[0])

    case MultiplexerType:
        fn := info.Func.(func(ctx context.Context, inputs []chan string, output chan string) error)
        err = fn(c.runCtx, inputChs, outputChs[0])

    case SeparatorType:
        fn := info.Func.(func(ctx context.Context, input chan string, outputs []chan string) error)
        err = fn(c.runCtx, inputChs[0], outputChs)
    }

    if err != nil && err != context.Canceled {
        fmt.Printf("Handler error (%v): %v. Stopping conveyer.\n", info.Type, err)
        c.runCancel()
    }
}

func (c *Conveyer) closeAllChannels() {
    c.mu.Lock()
    defer c.mu.Unlock()
    for id, ch := range c.channels {
        select {
        case <-c.runCtx.Done(): 
            close(ch)
        default:
            close(ch) 
        }
        delete(c.channels, id)
    }
}

func (c *Conveyer) Send(inputID string, data string) error {
    ch, ok := c.getChannel(inputID)
    if !ok {
        return ErrChanNotFound
    }
    
    select {
    case ch <- data:
        return nil
    case <-c.runCtx.Done():
        return c.runCtx.Err()
    default:
        return errors.New("send failed: channel closed or full") 
    }
}

func (c *Conveyer) Recv(outputID string) (string, error) {
    ch, ok := c.getChannel(outputID)
    if !ok {
        return "", ErrChanNotFound
    }

    select {
    case data, ok := <-ch:
        if !ok {
            return "undefined", nil 
        }
        return data, nil
    case <-c.runCtx.Done():
        return "", c.runCtx.Err()
    }
}