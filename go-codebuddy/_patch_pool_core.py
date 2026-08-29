from pathlib import Path

path = Path("internal/accounts/pool.go")
text = path.read_text(encoding="utf-8")

start = text.index("type Pool struct {")
end = text.index("func NormalizeStore(store Store) Store {")

new = r'''type Pool struct {
	path string

	mu       sync.RWMutex
	mem      Store
	loaded   bool
	dirty    bool
	dirReady bool

	persistMu sync.Mutex
	wakeCh    chan struct{}
	stopCh    chan struct{}
	doneCh    chan struct{}
	closeOnce sync.Once
}

func NewPool(path string) *Pool {
	p := &Pool{
		path:   path,
		mem:    EmptyStore(),
		wakeCh: make(chan struct{}, 1),
		stopCh: make(chan struct{}),
		doneCh: make(chan struct{}),
	}
	go p.flushLoop()
	return p
}

func (p *Pool) Path() string { return p.path }

func EmptyStore() Store {
	return Store{Version: 1, Provider: "codebuddy", NextIndex: 0, Accounts: []Account{}}
}

func cloneStore(store Store) Store {
	out := store
	if store.Accounts != nil {
		out.Accounts = append([]Account(nil), store.Accounts...)
	}
	return out
}

func (p *Pool) Read() (Store, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if err := p.ensureLoadedLocked(); err != nil {
		return Store{}, err
	}
	return cloneStore(p.mem), nil
}

func (p *Pool) Write(store Store) (Store, error) {
	normalized := NormalizeStore(store)
	p.mu.Lock()
	p.mem = cloneStore(normalized)
	p.loaded = true
	p.dirty = false
	snap := cloneStore(p.mem)
	p.mu.Unlock()
	if err := p.writeDisk(snap); err != nil {
		p.mu.Lock()
		p.dirty = true
		p.mu.Unlock()
		p.kickFlush()
		return Store{}, err
	}
	return cloneStore(normalized), nil
}

func (p *Pool) ensureLoadedLocked() error {
	if p.loaded {
		return nil
	}
	store, err := p.readDisk()
	if err != nil {
		return err
	}
	p.mem = store
	p.loaded = true
	p.dirty = false
	return nil
}

func (p *Pool) readDisk() (Store, error) {
	data, err := os.ReadFile(p.path)
	if errors.Is(err, os.ErrNotExist) {
		return EmptyStore(), nil
	}
	if err != nil {
		return Store{}, err
	}
	var raw Store
	if err := json.Unmarshal(data, &raw); err != nil {
		return Store{}, err
	}
	return NormalizeStore(raw), nil
}

func (p *Pool) writeDisk(store Store) error {
	p.persistMu.Lock()
	defer p.persistMu.Unlock()
	if !p.dirReady {
		if err := os.MkdirAll(filepath.Dir(p.path), 0o700); err != nil {
			return err
		}
		p.dirReady = true
	}
	payload, err := json.Marshal(store)
	if err != nil {
		return err
	}
	payload = append(payload, '\n')
	tmp := p.path + ".tmp"
	if err := os.WriteFile(tmp, payload, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, p.path)
}

func (p *Pool) markDirtyLocked() {
	p.dirty = true
	p.kickFlush()
}

func (p *Pool) kickFlush() {
	select {
	case p.wakeCh <- struct{}{}:
	default:
	}
}

func (p *Pool) Flush() error {
	p.mu.Lock()
	if err := p.ensureLoadedLocked(); err != nil {
		p.mu.Unlock()
		return err
	}
	if !p.dirty {
		p.mu.Unlock()
		return nil
	}
	snap := cloneStore(p.mem)
	p.dirty = false
	p.mu.Unlock()
	if err := p.writeDisk(snap); err != nil {
		p.mu.Lock()
		p.dirty = true
		p.mu.Unlock()
		return err
	}
	return nil
}

func (p *Pool) Close() error {
	var err error
	p.closeOnce.Do(func() {
		close(p.stopCh)
		<-p.doneCh
		err = p.Flush()
	})
	return err
}

func (p *Pool) flushLoop() {
	defer close(p.doneCh)
	for {
		select {
		case <-p.stopCh:
			_ = p.Flush()
			return
		case <-p.wakeCh:
			timer := time.NewTimer(250 * time.Millisecond)
			select {
			case <-timer.C:
			case <-p.stopCh:
				if !timer.Stop() {
					select {
					case <-timer.C:
					default:
					}
				}
				_ = p.Flush()
				return
			}
			for {
				select {
				case <-p.wakeCh:
				default:
					goto flushed
				}
			}
		flushed:
			_ = p.Flush()
		}
	}
}

'''

path.write_text(text[:start] + new + text[end:], encoding="utf-8")
print("core replaced")
