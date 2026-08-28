# JetBrains Modern Go Guidelines (Go <= 1.26)

Source: official JetBrains `go-modern-guidelines` / `use-modern-go` skill.

This project targets **Go 1.26+**. Prefer these idioms in all new/edited code.

## `new_expression` (since Go 1.26)

Use new(value) for pointer fields or arguments instead of generic/type-specific pointer helper functions or temporary variables used only for &value.

new(value) creates a *T from a value expression. In struct literals, prefer Field: new(value) for pointer fields over helper calls whose only purpose is returning &value; keep helpers only when they add behavior.

**Before:**
```go
func Pointer[T any](value T) *T {
  return &value
}

cfg := Config{
  Timeout: Pointer(30),
  Debug:   Pointer(true),
}
```

**After:**
```go
cfg := Config{
  Timeout: new(30),
  Debug:   new(true),
}
```

## `errors_as_type` (since Go 1.26)

Use errors.AsType[T](err) when checking whether an error matches a specific type.

errors.AsType returns the matched error value and a boolean directly. It avoids a separate temporary variable and the pointer-to-target pattern required by errors.As.

**Before:**
```go
var pathErr *os.PathError
if errors.As(err, &pathErr) {
  handle(pathErr)
}
```

**After:**
```go
if pathErr, ok := errors.AsType[*os.PathError](err); ok {
  handle(pathErr)
}
```

## `sync_waitgroup_go` (since Go 1.25)

Use wg.Go when spawning goroutines tracked by a sync.WaitGroup.

WaitGroup.Go starts a goroutine and handles the matching Add and Done calls. Use it when the goroutine's lifetime is exactly what the WaitGroup should track.

**Before:**
```go
var wg sync.WaitGroup
for _, item := range items {
  wg.Add(1)
  go func() {
    defer wg.Done()
    process(item)
  }()
}
wg.Wait()
```

**After:**
```go
var wg sync.WaitGroup
for _, item := range items {
  wg.Go(func() {
    process(item)
  })
}
wg.Wait()
```

## `testing_t_context` (since Go 1.24)

Use t.Context() when a test function needs a context tied to the test lifetime.

t.Context returns a context tied to the test lifetime. It removes manual background context setup when helper work should stop as the test is ending.

**Before:**
```go
func TestFoo(t *testing.T) {
  ctx, cancel := context.WithCancel(context.Background())
  defer cancel()
  result := doSomething(ctx)
}
```

**After:**
```go
func TestFoo(t *testing.T) {
  ctx := t.Context()
  result := doSomething(ctx)
}
```

## `json_omitzero` (since Go 1.24)

Use omitzero on JSON-tagged bool, numeric, struct, and time fields whose zero value should be omitted; keep omitempty for empty strings, slices, and maps.

When adding or editing JSON struct tags, use omitzero when the Go zero value means the field is absent, such as false, 0, or a zero time. Use omitempty when JSON empty-value semantics are intended, especially empty strings, slices, and maps.

**Before:**
```go
type CacheEntry struct {
  Name      string    `json:"name,omitempty"`
  Warm      bool      `json:"warm,omitempty"`
  Hits      int64     `json:"hits,omitempty"`
  ExpiresAt time.Time `json:"expiresAt,omitempty"`
}
```

**After:**
```go
type CacheEntry struct {
  Name      string    `json:"name,omitempty"`
  Warm      bool      `json:"warm,omitzero"`
  Hits      int64     `json:"hits,omitzero"`
  ExpiresAt time.Time `json:"expiresAt,omitzero"`
}
```

## `testing_b_loop` (since Go 1.24)

Use b.Loop() for the main loop in benchmark functions.

b.Loop is the modern benchmark loop. It manages benchmark iteration mechanics for you and can remove the need for manual timer control in the benchmark body.

**Before:**
```go
func BenchmarkFoo(b *testing.B) {
  for i := 0; i < b.N; i++ {
    doWork()
  }
}
```

**After:**
```go
func BenchmarkFoo(b *testing.B) {
  for b.Loop() {
    doWork()
  }
}
```

## `strings_split_seq` (since Go 1.24)

Use strings or bytes SplitSeq and FieldsSeq helpers when iterating over split results.

SplitSeq and FieldsSeq stream substrings instead of allocating a slice of all parts. The same idea applies to the matching bytes package helpers.

**Before:**
```go
for _, part := range strings.Split(s, ",") {
  process(part)
}
```

**After:**
```go
for part := range strings.SplitSeq(s, ",") {
  process(part)
}
```

## `maps_keys_values_iter` (since Go 1.23)

Use maps.Keys or maps.Values directly as iterators instead of manually looping over a map.

maps.Keys and maps.Values can be ranged over directly when a slice is not needed. This avoids allocating a temporary collection just to iterate.

**Before:**
```go
for k := range m {
  process(k)
}
```

**After:**
```go
for k := range maps.Keys(m) {
  process(k)
}
```

## `slices_collect` (since Go 1.23)

Use slices.Collect to build a slice from an iterator.

slices.Collect materializes an iterator into a slice when a slice is actually required. Prefer ranging over the iterator directly when streaming is enough.

**Before:**
```go
keys := make([]string, 0, len(m))
for k := range m {
  keys = append(keys, k)
}
```

**After:**
```go
keys := slices.Collect(maps.Keys(m))
```

## `slices_sorted` (since Go 1.23)

Use slices.Sorted to collect and sort iterator values in one step.

slices.Sorted collects iterator values and sorts them in one call. It is useful for deterministic output from unordered sources such as map keys.

**Before:**
```go
keys := make([]string, 0, len(m))
for k := range m {
  keys = append(keys, k)
}
slices.Sort(keys)
```

**After:**
```go
keys := slices.Sorted(maps.Keys(m))
```

## `time_tick_gc` (since Go 1.23)

Use time.Tick when it fits; Go 1.23 can recover unreferenced tickers without requiring Stop for GC.

Go 1.23 made unreferenced tickers recoverable by the garbage collector. Use time.Tick for simple forever loops, and use time.NewTicker when you need Stop or Reset.

**Before:**
```go
ticker := time.NewTicker(time.Second)
defer ticker.Stop()
for range ticker.C {
  poll()
}
```

**After:**
```go
for range time.Tick(time.Second) {
  poll()
}
```

## `range_over_int` (since Go 1.22)

Use for i := range n when iterating from 0 to n-1.

Ranging over an integer is the concise Go 1.22 form for zero-based count loops. Keep a traditional for loop when you need a non-zero start, custom step, or changing bound.

**Before:**
```go
for i := 0; i < len(items); i++ {
  process(items[i])
}
```

**After:**
```go
for i := range len(items) {
  process(items[i])
}
```

## `loopvar_capture` (since Go 1.22)

Do not add redundant loop-variable copies before closures or taking addresses; Go 1.22 gives each iteration its own variables.

Go 1.22 gives each loop iteration its own variables, so defensive copies like v := v are usually redundant before closures, goroutines, deferred functions, or appending &v. Keep an explicit copy only when it serves another purpose. If you need a pointer to the original slice element rather than the per-iteration copy, use &slice[i].

**Before:**
```go
for _, item := range items {
  item := item
  go func() {
    process(item)
  }()
}
```

**After:**
```go
for _, item := range items {
  go func() {
    process(item)
  }()
}
```

## `cmp_or` (since Go 1.22)

Use cmp.Or to pick the first non-zero value from a fallback chain.

cmp.Or returns the first non-zero value from its arguments. It is concise for simple fallback chains, but remember that all arguments are evaluated before the call.

**Before:**
```go
name := os.Getenv("NAME")
if name == "" {
  name = "default"
}
```

**After:**
```go
name := cmp.Or(os.Getenv("NAME"), "default")
```

## `reflect_type_for` (since Go 1.22)

Use reflect.TypeFor[T]() instead of reflect.TypeOf((*T)(nil)).Elem().

reflect.TypeFor returns the reflect.Type for a compile-time type parameter or type argument. It avoids nil-pointer tricks and is clearer for interface types.

**Before:**
```go
typ := reflect.TypeOf((*T)(nil)).Elem()
```

**After:**
```go
typ := reflect.TypeFor[T]()
```

## `http_servemux_patterns` (since Go 1.22)

Use method-aware ServeMux patterns and r.PathValue for path parameters.

The modern ServeMux pattern syntax can include an HTTP method and named path wildcards. r.PathValue retrieves wildcard values without manual path trimming.

**Before:**
```go
mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodGet {
    http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
    return
  }
  id := strings.TrimPrefix(r.URL.Path, "/api/")
  handleID(w, r, id)
})
```

**After:**
```go
mux.HandleFunc("GET /api/{id}", func(w http.ResponseWriter, r *http.Request) {
  handleID(w, r, r.PathValue("id"))
})
```

## `min_max` (since Go 1.21)

Use built-in min and max instead of handwritten comparisons.

The min and max built-ins express ordered comparisons directly. They remove branch boilerplate when the logic is simply choosing the smaller or larger value.

**Before:**
```go
if b > a {
  a = b
}
```

**After:**
```go
a = max(a, b)
```

## `clear` (since Go 1.21)

Use clear(m) to delete all map entries or clear(s) to zero slice elements.

clear is the idiomatic built-in for removing all map entries or zeroing the elements of a slice. For slices it preserves the slice length and capacity.

**Before:**
```go
for k := range m {
  delete(m, k)
}
```

**After:**
```go
clear(m)
```

## `slices_contains` (since Go 1.21)

Use slices.Contains instead of a manual search loop.

slices.Contains is the direct membership test for slices of comparable values. It replaces a manual search loop whose only result is whether the element is present.

**Before:**
```go
found := false
for _, item := range items {
  if item == x {
    found = true
    break
  }
}
```

**After:**
```go
found := slices.Contains(items, x)
```

## `slices_index` (since Go 1.21)

Use slices.Index to find the index of an element, returning -1 when absent.

slices.Index is the direct index lookup for slices of comparable values. It returns -1 when the value is absent, matching the common handwritten convention.

**Before:**
```go
index := -1
for i, item := range items {
  if item == x {
    index = i
    break
  }
}
```

**After:**
```go
index := slices.Index(items, x)
```

## `slices_index_func` (since Go 1.21)

Use slices.IndexFunc to find an element by predicate.

slices.IndexFunc is the standard predicate-based lookup. Use it when the match depends on fields, derived values, or any condition other than direct equality.

**Before:**
```go
index := -1
for i, item := range items {
  if item.ID == id {
    index = i
    break
  }
}
```

**After:**
```go
index := slices.IndexFunc(items, func(item Item) bool {
  return item.ID == id
})
```

## `slices_sort_func` (since Go 1.21)

Use slices.SortFunc with cmp.Compare instead of sort.Slice for typed comparisons.

slices.SortFunc compares typed elements directly, avoiding sort.Slice closures that index back into the slice. cmp.Compare is the usual comparator for ordered fields.

**Before:**
```go
sort.Slice(items, func(i, j int) bool {
  return items[i].X < items[j].X
})
```

**After:**
```go
slices.SortFunc(items, func(a, b Item) int {
  return cmp.Compare(a.X, b.X)
})
```

## `slices_sort` (since Go 1.21)

Use slices.Sort for slices of ordered values.

slices.Sort is the generic sort for slices of ordered values. It replaces older type-specific helpers and simple sort.Slice calls.

**Before:**
```go
sort.Ints(values)
```

**After:**
```go
slices.Sort(values)
```

## `slices_max_min` (since Go 1.21)

Use slices.Max and slices.Min instead of manual loops over ordered values.

slices.Max and slices.Min capture the common non-empty-slice scan for an ordered maximum or minimum. Keep explicit logic when empty slices or custom comparison rules need special handling.

**Before:**
```go
maxValue := values[0]
for _, value := range values[1:] {
  if value > maxValue {
    maxValue = value
  }
}
```

**After:**
```go
maxValue := slices.Max(values)
```

## `slices_reverse` (since Go 1.21)

Use slices.Reverse instead of a manual swap loop.

slices.Reverse performs an in-place reversal without hand-written index math. It is clearer and less error-prone than open-coded swap loops.

**Before:**
```go
for i, j := 0, len(items)-1; i < j; i, j = i+1, j-1 {
  items[i], items[j] = items[j], items[i]
}
```

**After:**
```go
slices.Reverse(items)
```

## `slices_compact` (since Go 1.21)

Use slices.Compact to remove consecutive duplicates in place.

slices.Compact removes consecutive duplicate values in place. If the goal is to remove all duplicates, sort or otherwise group equal values first.

**Before:**
```go
out := values[:0]
for i, value := range values {
  if i == 0 || value != values[i-1] {
    out = append(out, value)
  }
}
values = out
```

**After:**
```go
values = slices.Compact(values)
```

## `slices_clip` (since Go 1.21)

Use slices.Clip to remove unused capacity.

slices.Clip limits a slice's capacity to its length. Use it to stop retaining excess backing-array capacity or to prevent later appends from reusing hidden capacity.

**Before:**
```go
s = s[:len(s):len(s)]
```

**After:**
```go
s = slices.Clip(s)
```

## `slices_clone` (since Go 1.21)

Use slices.Clone to copy a slice.

slices.Clone returns a shallow copy of a slice using the standard helper. It is clearer than append-copy boilerplate and preserves nil slices.

**Before:**
```go
copied := append([]T(nil), values...)
```

**After:**
```go
copied := slices.Clone(values)
```

## `maps_clone` (since Go 1.21)

Use maps.Clone instead of manual map iteration.

maps.Clone returns a shallow copy of a map using the standard helper. It preserves nil maps and makes the copy operation explicit.

**Before:**
```go
copied := make(map[string]int, len(src))
for k, v := range src {
  copied[k] = v
}
```

**After:**
```go
copied := maps.Clone(src)
```

## `maps_copy` (since Go 1.21)

Use maps.Copy to copy entries from one map into another.

maps.Copy copies all entries from one map into another, overwriting existing keys in the destination. It replaces loops whose body only assigns dst[k] = v.

**Before:**
```go
for k, v := range src {
  dst[k] = v
}
```

**After:**
```go
maps.Copy(dst, src)
```

## `maps_delete_func` (since Go 1.21)

Use maps.DeleteFunc to delete map entries that match a predicate.

maps.DeleteFunc deletes entries selected by a predicate. It keeps the filtering condition in one callback and avoids hand-written delete loops.

**Before:**
```go
for k, v := range m {
  if shouldDelete(k, v) {
    delete(m, k)
  }
}
```

**After:**
```go
maps.DeleteFunc(m, func(k string, v int) bool {
  return shouldDelete(k, v)
})
```

## `sync_once_func` (since Go 1.21)

Use sync.OnceFunc instead of sync.Once plus a wrapper closure.

sync.OnceFunc wraps a function so that the returned function runs it at most once. It is a compact fit for idempotent cleanup or one-time initialization hooks.

**Before:**
```go
var once sync.Once
cleanup := func() {
  once.Do(func() {
    close(ch)
  })
}
```

**After:**
```go
cleanup := sync.OnceFunc(func() {
  close(ch)
})
```

## `sync_once_value` (since Go 1.21)

Use sync.OnceValue to memoize a computed value.

sync.OnceValue memoizes a computed value safely across goroutines. It replaces a sync.Once plus a separate result variable and getter closure.

**Before:**
```go
var once sync.Once
var value T
getter := func() T {
  once.Do(func() {
    value = computeValue()
  })
  return value
}
```

**After:**
```go
getter := sync.OnceValue(func() T {
  return computeValue()
})
```

## `context_after_func` (since Go 1.21)

Use context.AfterFunc to run cleanup when a context is canceled.

context.AfterFunc registers work to run after cancellation and returns a stop function. It avoids starting a goroutine whose only job is to wait on ctx.Done.

**Before:**
```go
go func() {
  <-ctx.Done()
  cleanup()
}()
```

**After:**
```go
stop := context.AfterFunc(ctx, cleanup)
defer stop()
```

## `context_timeout_deadline_cause` (since Go 1.21)

Use timeout and deadline contexts with causes when callers need to inspect the cancellation reason.

Timeout and deadline causes let callers distinguish the reason for cancellation with context.Cause. Use them when ctx.Err is too coarse for diagnostics or control flow.

**Before:**
```go
ctx, cancel := context.WithTimeout(parent, d)
defer cancel()
```

**After:**
```go
ctx, cancel := context.WithTimeoutCause(parent, d, errTimeout)
defer cancel()
```

## `strings_clone` (since Go 1.20)

Use strings.Clone to copy a string without retaining shared backing memory.

strings.Clone forces a string copy. Use it when retaining a small string derived from a much larger string would otherwise keep the larger backing memory alive.

**Before:**
```go
copied := string([]byte(s))
```

**After:**
```go
copied := strings.Clone(s)
```

## `bytes_clone` (since Go 1.20)

Use bytes.Clone to copy a byte slice.

bytes.Clone communicates that a new byte slice is required. It avoids append-copy boilerplate and preserves the standard library's nil-slice behavior.

**Before:**
```go
copied := append([]byte(nil), b...)
```

**After:**
```go
copied := bytes.Clone(b)
```

## `strings_cut_prefix_suffix` (since Go 1.20)

Use strings.CutPrefix or strings.CutSuffix when you need both the trimmed result and whether it matched.

CutPrefix and CutSuffix combine the match check and trimming operation. They avoid writing HasPrefix or HasSuffix followed by a second operation that repeats the same condition.

**Before:**
```go
if strings.HasPrefix(s, "pre:") {
  rest := strings.TrimPrefix(s, "pre:")
  use(rest)
}
```

**After:**
```go
if rest, ok := strings.CutPrefix(s, "pre:"); ok {
  use(rest)
}
```

## `errors_join` (since Go 1.20)

Use errors.Join to combine multiple errors while preserving error matching.

errors.Join combines multiple non-nil errors while preserving error matching through errors.Is and errors.As. It returns nil when there are no non-nil errors.

**Before:**
```go
if err1 != nil && err2 != nil {
  return fmt.Errorf("%v; %w", err1, err2)
}
```

**After:**
```go
return errors.Join(err1, err2)
```

## `context_cancel_cause` (since Go 1.20)

Use context.WithCancelCause and context.Cause when cancellation needs to carry an error cause.

Cancellation causes let code record why a context was canceled. Callers can inspect context.Cause(ctx) instead of only seeing the broad ctx.Err result.

**Before:**
```go
ctx, cancel := context.WithCancel(parent)
cancel()
```

**After:**
```go
ctx, cancel := context.WithCancelCause(parent)
cancel(err)
cause := context.Cause(ctx)
```

## `fmt_appendf` (since Go 1.19)

Use fmt.Appendf when appending formatted text to a byte slice and an intermediate fmt.Sprintf string is unnecessary.

fmt.Appendf writes formatted output directly into a byte slice. It is useful when code is already accumulating []byte and the extra string allocation from fmt.Sprintf is not needed.

**Before:**
```go
buf = append(buf, []byte(fmt.Sprintf("x=%d", x))...)
```

**After:**
```go
buf = fmt.Appendf(buf, "x=%d", x)
```

## `atomic_types` (since Go 1.19)

Use typed atomics such as atomic.Bool, atomic.Int64, and atomic.Pointer[T] instead of untyped atomic functions.

Typed atomic wrapper values keep the storage and the atomic operations together. They make the value type visible, reduce accidental non-atomic access, and avoid old pointer-alignment pitfalls.

**Before:**
```go
var enabled int32
atomic.StoreInt32(&enabled, 1)
if atomic.LoadInt32(&enabled) != 0 {
  run()
}
```

**After:**
```go
var enabled atomic.Bool
enabled.Store(true)
if enabled.Load() {
  run()
}
```

## `any` (since Go 1.18)

Use any instead of interface{}.

any is the built-in alias for interface{} introduced with generics. Use it for unconstrained values and type parameters so the code reads in current Go terminology.

**Before:**
```go
func Decode(v interface{}) error {
  return nil
}
```

**After:**
```go
func Decode(v any) error {
  return nil
}
```

## `bytes_cut` (since Go 1.18)

Use bytes.Cut instead of bytes.Index plus manual slicing.

bytes.Cut returns before, after, and found in one call. It avoids duplicated Index checks and separator-length arithmetic while returning slices of the original byte slice.

**Before:**
```go
i := bytes.Index(b, sep)
if i < 0 {
  return nil, nil, false
}
before, after := b[:i], b[i+len(sep):]
```

**After:**
```go
before, after, found := bytes.Cut(b, sep)
```

## `strings_cut` (since Go 1.18)

Use strings.Cut instead of strings.Index plus manual slicing.

strings.Cut returns before, after, and found in one call. It makes split-at-first-separator logic explicit and removes manual Index checks and slicing.

**Before:**
```go
i := strings.Index(s, ":")
if i < 0 {
  return "", "", false
}
key, value := s[:i], s[i+1:]
```

**After:**
```go
key, value, found := strings.Cut(s, ":")
```

## `errors_is` (since Go 1.13)

Use errors.Is(err, target) instead of err == target so wrapped errors are handled correctly.

errors.Is walks wrapped errors and honors custom Is methods. Direct equality only matches the exact sentinel value and misses errors wrapped with fmt.Errorf or custom wrappers.

**Before:**
```go
if err == os.ErrNotExist {
  return nil
}
```

**After:**
```go
if errors.Is(err, os.ErrNotExist) {
  return nil
}
```

## `time_until` (since Go 1.8)

Use time.Until(deadline) instead of deadline.Sub(time.Now()).

time.Until is the standard spelling for remaining time before a deadline. It is equivalent to deadline.Sub(time.Now()) but reads in the direction the caller cares about.

**Before:**
```go
remaining := deadline.Sub(time.Now())
```

**After:**
```go
remaining := time.Until(deadline)
```

## `time_since` (since Go 1.0)

Use time.Since(start) instead of time.Now().Sub(start).

time.Since is the standard spelling for elapsed time from a recorded start. It keeps the time.Now call inside the time package helper and makes the intent easier to scan.

**Before:**
```go
elapsed := time.Now().Sub(start)
```

**After:**
```go
elapsed := time.Since(start)
```
