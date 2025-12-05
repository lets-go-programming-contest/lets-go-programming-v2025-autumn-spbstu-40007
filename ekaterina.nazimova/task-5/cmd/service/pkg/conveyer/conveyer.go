package conveyer

import (
    "context"
    "errors"
    "sync"
)

type conveyer interface {
    RegisterDecorator(fn func(ctx context.Context, input chan string, output chan string) error, input string, output string)
    RegisterMultiplexer(fn func(ctx context.Context, inputs []chan string, output chan string) error, inputs []string, output string)
    RegisterSeparator(fn func(ctx context.Context, input chan string, outputs []chan string) error, input string, outputs []string)
    Run(ctx context.Context) error
    Send(input string, data string) error
    Recv(output string) (string, error)
}

type ConveyerImpl struct {
    channelSize int
    channels    map[string]chan string
    runners     []func(ctx context.Context) error 
    mu          sync.RWMutex                     
}

func New(size int) *ConveyerImpl {
    return &ConveyerImpl{
        channelSize: size,
        channels:    make(map[string]chan string),
        runners:     make([]func(ctx context.Context) error, 0),
    }
}

func (c *ConveyerImpl) getOrCreateChannel(id string) chan string {
    c.mu.Lock()
    defer c.mu.Unlock()

    if ch, ok := c.channels[id]; ok {
        return ch
    }

    ch := make(chan string, c.channelSize)
    c.channels[id] = ch
    return ch
}

func (c *ConveyerImpl) RegisterDecorator(
    fn func(ctx context.Context, input chan string, output chan string) error,
    inputID string,
    outputID string,
) {
    inputCh := c.getOrCreateChannel(inputID)
    outputCh := c.getOrCreateChannel(outputID)

    runner := func(ctx context.Context) error {
        return fn(ctx, inputCh, outputCh)
    }

    c.mu.Lock()
    c.runners = append(c.runners, runner)
    c.mu.Unlock()
}

func (c *ConveyerImpl) RegisterMultiplexer(
    fn func(ctx context.Context, inputs []chan string, output chan string) error,
    inputIDs []string,
    outputID string,
) {
    inputChans := make([]chan string, len(inputIDs))
    for i, id := range inputIDs {
        inputChans[i] = c.getOrCreateChannel(id)
    }
    outputCh := c.getOrCreateChannel(outputID)

    runner := func(ctx context.Context) error {
        return fn(ctx, inputChans, outputCh)
    }

    c.mu.Lock()
    c.runners = append(c.runners, runner)
    c.mu.Unlock()
}

func (c *ConveyerImpl) RegisterSeparator(
    fn func(ctx context.Context, input chan string, outputs []chan string) error,
    inputID string,
    outputIDs []string,
) {
    inputCh := c.getOrCreateChannel(inputID)
    outputChans := make([]chan string, len(outputIDs))
    for i, id := range outputIDs {
        outputChans[i] = c.getOrCreateChannel(id)
    }

    runner := func(ctx context.Context) error {
        return fn(ctx, inputCh, outputChans)
    }

    c.mu.Lock()
    c.runners = append(c.runners, runner)
    c.mu.Unlock()
}

func (c *ConveyerImpl) Run(ctx context.Context) error {
    c.mu.RLock()
    numRunners := len(c.runners)
    c.mu.RUnlock()

    if numRunners == 0 {
        return nil
    }

    var wg sync.WaitGroup
    wg.Add(numRunners)

    errChan := make(chan error, numRunners)
    
    ctx, cancel := context.WithCancel(ctx)
    defer cancel() 

    c.mu.RLock()
    for _, runner := range c.runners {
        go func(r func(ctx context.Context) error) {
            defer wg.Done()
            if err := r(ctx); err != nil && !errors.Is(err, context.Canceled) {
                errChan <- err
                cancel()
            }
        }(runner)
    }
    c.mu.RUnlock()

    var runErr error
    select {
    case <-ctx.Done():
        runErr = ctx.Err()
    case err := <-errChan:
        runErr = err
    }
    
    wg.Wait()
    
    c.mu.Lock()
    defer c.mu.Unlock()
    for _, ch := range c.channels {
        close(ch)
    }

    if runErr == nil && errors.Is(ctx.Err(), context.Canceled) && len(errChan) > 0 {
         runErr = <-errChan
    }

    return runErr
}

func (c *ConveyerImpl) Send(inputID string, data string) error {
    c.mu.RLock()
    ch, ok := c.channels[inputID]
    c.mu.RUnlock()

    if !ok {
        return errors.New("chan not found")
    }

    ch <- data
    return nil
}

func (c *ConveyerImpl) Recv(outputID string) (string, error) {
    c.mu.RLock()
    ch, ok := c.channels[outputID]
    c.mu.RUnlock()

    if !ok {
        return "", errors.New("chan not found")
    }

    v, ok := <-ch
    if !ok {
        return "undefined", nil
    }
    return v, nil
}
