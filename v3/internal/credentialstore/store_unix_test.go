//go:build linux || freebsd || openbsd || netbsd || dragonfly

package credentialstore

import (
	"bufio"
	"errors"
	"os/exec"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/godbus/dbus/v5"
)

const testCollection dbus.ObjectPath = "/org/freedesktop/secrets/collection/login"
const testSession dbus.ObjectPath = "/org/freedesktop/secrets/session/test"
const testPrompt dbus.ObjectPath = "/org/freedesktop/secrets/prompt/test"

type testVault struct {
	itemPrompt                                                  bool
	lookups                                                     []string
	mu                                                          sync.Mutex
	conn                                                        *dbus.Conn
	items                                                       map[dbus.ObjectPath]*testItem
	next                                                        int
	alias, locked, dismiss, denyUnlock, badSecret, createPrompt bool
	opened, closed                                              int
	promptResult                                                dbus.Variant
	promptOperation                                             func()
}

type testItem struct {
	vault      *testVault
	path       dbus.ObjectPath
	attributes map[string]string
	value      []byte
	locked     bool
}

type testCollectionObject struct{ vault *testVault }
type testSessionObject struct{ vault *testVault }
type testPromptObject struct{ vault *testVault }
type testProperties struct {
	vault *testVault
	path  dbus.ObjectPath
}

func newTestVault(t testing.TB) *testVault {
	t.Helper()
	daemon, err := exec.LookPath("dbus-daemon")
	if err != nil {
		t.Skip("dbus-daemon is required for isolated Secret Service tests")
	}
	cmd := exec.Command(daemon, "--session", "--nofork", "--nopidfile", "--print-address=1")
	output, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err = cmd.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { cmd.Process.Kill(); cmd.Wait() })
	address, err := bufio.NewReader(output).ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("DBUS_SESSION_BUS_ADDRESS", strings.TrimSpace(address))
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { conn.Close() })
	return exportTestVault(t, conn)
}

func exportTestVault(t testing.TB, conn *dbus.Conn) *testVault {
	t.Helper()
	reply, err := conn.RequestName(secretService, dbus.NameFlagDoNotQueue)
	if err != nil || reply != dbus.RequestNameReplyPrimaryOwner {
		t.Fatalf("own isolated service: %v, %v", reply, err)
	}
	v := &testVault{conn: conn, items: make(map[dbus.ObjectPath]*testItem), alias: true}
	v.mu.Lock()
	defer v.mu.Unlock()
	exports := []struct {
		value any
		path  dbus.ObjectPath
		iface string
	}{
		{v, servicePath, serviceInterface},
		{&testCollectionObject{v}, testCollection, collectionInterface},
		{&testCollectionObject{v}, "/org/freedesktop/secrets/aliases/default", collectionInterface},
		{&testSessionObject{v}, testSession, "org.freedesktop.Secret.Session"},
		{&testPromptObject{v}, testPrompt, promptInterface},
	}
	for _, e := range exports {
		if err = conn.Export(e.value, e.path, e.iface); err != nil {
			t.Fatal(err)
		}
	}
	for _, path := range []dbus.ObjectPath{servicePath, testCollection, "/org/freedesktop/secrets/aliases/default"} {
		if err = conn.Export(&testProperties{v, path}, path, "org.freedesktop.DBus.Properties"); err != nil {
			t.Fatal(err)
		}
	}
	return v
}

func (v *testVault) OpenSession(algorithm string, input dbus.Variant) (dbus.Variant, dbus.ObjectPath, *dbus.Error) {
	v.mu.Lock()
	defer v.mu.Unlock()
	if algorithm != "plain" {
		return dbus.MakeVariant(""), "/", dbus.MakeFailedError(errors.New("test vault only supports plain sessions"))
	}
	v.opened++
	return dbus.MakeVariant(""), testSession, nil
}
func (s *testSessionObject) Close() *dbus.Error {
	s.vault.mu.Lock()
	defer s.vault.mu.Unlock()
	s.vault.closed++
	return nil
}
func (v *testVault) ReadAlias(name string) (dbus.ObjectPath, *dbus.Error) {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.alias && name == "default" {
		return testCollection, nil
	}
	return "/", nil
}
func (v *testVault) CreateCollection(properties map[string]dbus.Variant, alias string) (dbus.ObjectPath, dbus.ObjectPath, *dbus.Error) {
	v.mu.Lock()
	defer v.mu.Unlock()
	if alias != "default" || properties[collectionInterface+".Label"].Value() == nil {
		return "/", "/", dbus.MakeFailedError(errors.New("missing default collection properties"))
	}
	if v.createPrompt {
		v.promptResult = dbus.MakeVariant(testCollection)
		v.promptOperation = func() { v.alias = true }
		return "/", testPrompt, nil
	}
	v.alias = true
	return testCollection, "/", nil
}
func (v *testVault) matching(attributes map[string]string) (unlocked, locked []dbus.ObjectPath) {
	for path, item := range v.items {
		match := true
		for key, value := range attributes {
			if item.attributes[key] != value {
				match = false
				break
			}
		}
		if match {
			if item.locked || v.locked {
				locked = append(locked, path)
			} else {
				unlocked = append(unlocked, path)
			}
		}
	}
	sort.Slice(unlocked, func(i, j int) bool { return unlocked[i] < unlocked[j] })
	sort.Slice(locked, func(i, j int) bool { return locked[i] < locked[j] })
	return
}
func (v *testVault) SearchItems(attributes map[string]string) ([]dbus.ObjectPath, []dbus.ObjectPath, *dbus.Error) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.lookups = append(v.lookups, "service")
	unlocked, locked := v.matching(attributes)
	return unlocked, locked, nil
}
func (c *testCollectionObject) SearchItems(attributes map[string]string) ([]dbus.ObjectPath, *dbus.Error) {
	c.vault.mu.Lock()
	defer c.vault.mu.Unlock()
	c.vault.lookups = append(c.vault.lookups, "collection")
	a, b := c.vault.matching(attributes)
	return append(a, b...), nil
}
func (v *testVault) Unlock(objects []dbus.ObjectPath) ([]dbus.ObjectPath, dbus.ObjectPath, *dbus.Error) {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.denyUnlock {
		return nil, "/", nil
	}
	locked := v.locked
	for _, path := range objects {
		if item := v.items[path]; item != nil {
			locked = locked || item.locked
		}
	}
	if !locked {
		return objects, "/", nil
	}
	v.promptResult = dbus.MakeVariant(objects)
	v.promptOperation = func() {
		v.locked = false
		for _, path := range objects {
			if item := v.items[path]; item != nil {
				item.locked = false
			}
		}
	}
	return nil, testPrompt, nil
}
func (p *testPromptObject) Prompt(window string) *dbus.Error {
	v := p.vault
	v.mu.Lock()
	defer v.mu.Unlock()
	if !v.dismiss && v.promptOperation != nil {
		v.promptOperation()
	}
	if err := v.conn.Emit(testPrompt, promptInterface+".Completed", v.dismiss, v.promptResult); err != nil {
		return dbus.MakeFailedError(err)
	}
	return nil
}
func (p *testProperties) Get(iface, name string) (dbus.Variant, *dbus.Error) {
	v := p.vault
	v.mu.Lock()
	defer v.mu.Unlock()
	switch name {
	case "Locked":
		locked := v.locked
		if item := v.items[p.path]; item != nil {
			locked = locked || item.locked
		}
		return dbus.MakeVariant(locked), nil
	case "Collections":
		return dbus.MakeVariant([]dbus.ObjectPath{testCollection}), nil
	case "Items":
		a, b := v.matching(nil)
		return dbus.MakeVariant(append(a, b...)), nil
	}
	return dbus.MakeVariant(""), dbus.MakeFailedError(errors.New("unknown property"))
}
func (c *testCollectionObject) CreateItem(properties map[string]dbus.Variant, secret wireSecret, replace bool) (dbus.ObjectPath, dbus.ObjectPath, *dbus.Error) {
	v := c.vault
	v.mu.Lock()
	defer v.mu.Unlock()
	attributes, ok := properties[itemInterface+".Attributes"].Value().(map[string]string)
	if !ok || secret.Session != testSession {
		return "/", "/", dbus.MakeFailedError(errors.New("invalid item properties or session"))
	}
	if replace {
		for path, item := range v.items {
			if reflect.DeepEqual(item.attributes, attributes) {
				item.value = append([]byte{}, secret.Value...)
				return path, "/", nil
			}
		}
	}
	v.next++
	path := dbus.ObjectPath("/org/freedesktop/secrets/collection/login/item" + strconv.Itoa(v.next))
	item := &testItem{v, path, attributes, append([]byte{}, secret.Value...), false}
	if err := v.conn.Export(item, path, itemInterface); err != nil {
		return "/", "/", dbus.MakeFailedError(err)
	}
	if err := v.conn.Export(&testProperties{v, path}, path, "org.freedesktop.DBus.Properties"); err != nil {
		return "/", "/", dbus.MakeFailedError(err)
	}
	if v.itemPrompt {
		v.promptResult = dbus.MakeVariant(path)
		v.promptOperation = func() { v.items[path] = item }
		return "/", testPrompt, nil
	}
	v.items[path] = item
	return path, "/", nil
}
func (i *testItem) GetSecret(session dbus.ObjectPath) (wireSecret, *dbus.Error) {
	v := i.vault
	v.mu.Lock()
	defer v.mu.Unlock()
	if i.locked || v.locked {
		return wireSecret{}, dbus.MakeFailedError(errors.New("item locked"))
	}
	secret := wireSecret{session, []byte{}, append([]byte{}, i.value...), "text/plain"}
	if v.badSecret {
		secret.Parameters = []byte{1}
	}
	return secret, nil
}

func (i *testItem) Delete() (dbus.ObjectPath, *dbus.Error) {
	v := i.vault
	v.mu.Lock()
	defer v.mu.Unlock()
	if i.locked || v.locked {
		return "/", dbus.MakeFailedError(errors.New("item locked"))
	}
	delete(v.items, i.path)
	return "/", nil
}

func TestSecretServiceRoundTrip(t *testing.T) {
	v := newTestVault(t)
	if _, err := Get("wails-test", "missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing Get: %v", err)
	}
	if err := Delete("wails-test", "missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing Delete: %v", err)
	}
	for _, value := range []string{"", "secret\n ", "日本語🔒\x00password"} {
		if err := Set("wails-test", "signing", value); err != nil {
			t.Fatal(err)
		}
		if got, err := Get("wails-test", "signing"); err != nil || got != value {
			t.Fatalf("Get mismatch: %q, %v", got, err)
		}
	}
	if _, err := Get("wails-test", "different-user"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("unrelated user matched: %v", err)
	}
	if err := Delete("wails-test", "signing"); err != nil {
		t.Fatal(err)
	}
	if _, err := Get("wails-test", "signing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("deleted Get: %v", err)
	}
	v.mu.Lock()
	defer v.mu.Unlock()
	if len(v.items) != 0 || v.opened != v.closed {
		t.Fatalf("leaked items or sessions: %d items, %d/%d sessions", len(v.items), v.opened, v.closed)
	}
}

func TestSecretServiceIsolation(t *testing.T) {
	v := newTestVault(t)
	if err := Set("wails-test", "signing", "secret"); err != nil {
		t.Fatal(err)
	}
	if err := Set("other-service", "signing", "other"); err != nil {
		t.Fatal(err)
	}
	if err := Delete("wails-test", "signing"); err != nil {
		t.Fatal(err)
	}
	if got, err := Get("other-service", "signing"); err != nil || got != "other" {
		t.Fatalf("unrelated service changed: %q, %v", got, err)
	}
	v.mu.Lock()
	defer v.mu.Unlock()
	for _, scope := range v.lookups {
		if scope != "collection" {
			t.Fatalf("lookup escaped the default collection: %s", scope)
		}
	}
	if len(v.items) != 1 || v.opened != v.closed {
		t.Fatalf("leaked items or sessions: %d items, %d/%d sessions", len(v.items), v.opened, v.closed)
	}
}

func TestSecretServicePromptsAndInvalidSecrets(t *testing.T) {
	v := newTestVault(t)
	v.mu.Lock()
	v.alias = false
	v.createPrompt = true
	v.mu.Unlock()
	if err := Set("wails-test", "signing", "secret"); err != nil {
		t.Fatal(err)
	}
	v.mu.Lock()
	v.locked = true
	v.mu.Unlock()
	if got, err := Get("wails-test", "signing"); err != nil || got != "secret" {
		t.Fatalf("unlock prompt: %q, %v", got, err)
	}
	v.mu.Lock()
	v.locked = true
	v.dismiss = true
	v.mu.Unlock()
	if _, err := Get("wails-test", "signing"); !errors.Is(err, ErrDismissed) {
		t.Fatalf("dismissed prompt: %v", err)
	}
	v.mu.Lock()
	v.dismiss = false
	v.denyUnlock = true
	v.mu.Unlock()
	if _, err := Get("wails-test", "signing"); err == nil {
		t.Fatal("accepted failed unlock")
	}
	v.mu.Lock()
	v.locked = false
	v.denyUnlock = false
	v.badSecret = true
	v.mu.Unlock()
	if _, err := Get("wails-test", "signing"); err == nil {
		t.Fatal("accepted invalid plain-session secret")
	}
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.opened != v.closed {
		t.Fatalf("failed operations leaked sessions: %d/%d", v.opened, v.closed)
	}
}

func TestSecretServiceLockedDeleteAndDismissedCreation(t *testing.T) {
	v := newTestVault(t)
	if err := Set("wails-test", "signing", "secret"); err != nil {
		t.Fatal(err)
	}
	v.mu.Lock()
	v.locked = true
	v.mu.Unlock()
	if err := Delete("wails-test", "signing"); err != nil {
		t.Fatal(err)
	}
	v.mu.Lock()
	v.alias = false
	v.createPrompt = true
	v.dismiss = true
	v.mu.Unlock()
	if err := Set("wails-test", "signing", "not stored"); !errors.Is(err, ErrDismissed) {
		t.Fatalf("dismissed creation: %v", err)
	}
	v.mu.Lock()
	defer v.mu.Unlock()
	if len(v.items) != 0 || v.opened != v.closed {
		t.Fatalf("failed operation leaked data or sessions: items=%d sessions=%d/%d", len(v.items), v.opened, v.closed)
	}
}

func BenchmarkSecretServiceGet(b *testing.B) {
	newTestVault(b)
	if err := Set("wails-test", "signing", "sample password"); err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if got, err := Get("wails-test", "signing"); err != nil || got != "sample password" {
			b.Fatalf("Get = %q, %v", got, err)
		}
	}
}

func TestSecretServiceMissingCollection(t *testing.T) {
	v := newTestVault(t)
	v.mu.Lock()
	v.alias = false
	v.mu.Unlock()
	if _, err := Get("wails-test", "signing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing collection Get: %v", err)
	}
	if err := Delete("wails-test", "signing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing collection Delete: %v", err)
	}
	v.mu.Lock()
	created := v.alias
	v.mu.Unlock()
	if created {
		t.Fatal("read-only operations created a collection")
	}
	if err := Set("wails-test", "signing", "secret"); err != nil {
		t.Fatal(err)
	}
	if got, err := Get("wails-test", "signing"); err != nil || got != "secret" {
		t.Fatalf("new collection Get: %v", err)
	}
}

func TestSecretServiceItemPrompts(t *testing.T) {
	v := newTestVault(t)
	v.mu.Lock()
	v.itemPrompt = true
	v.mu.Unlock()
	if err := Set("wails-test", "signing", "secret"); err != nil {
		t.Fatal(err)
	}
	if got, err := Get("wails-test", "signing"); err != nil || got != "secret" {
		t.Fatalf("created item Get: %v", err)
	}
	if err := Delete("wails-test", "signing"); err != nil {
		t.Fatal(err)
	}
	v.mu.Lock()
	v.dismiss = true
	v.mu.Unlock()
	if err := Set("wails-test", "signing", "not stored"); !errors.Is(err, ErrDismissed) {
		t.Fatalf("dismissed item creation: %v", err)
	}
	if _, err := Get("wails-test", "signing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("dismissed creation stored an item: %v", err)
	}
}

func TestSecretServiceUnavailable(t *testing.T) {
	t.Setenv("DBUS_SESSION_BUS_ADDRESS", "unix:path="+t.TempDir()+"/missing-bus")
	if _, err := Get("wails-test", "signing"); err == nil || errors.Is(err, ErrNotFound) {
		t.Fatalf("unavailable vault Get: %v", err)
	}
	if err := Set("wails-test", "signing", "value"); err == nil {
		t.Fatal("unavailable vault accepted a write")
	}
	if err := Delete("wails-test", "signing"); err == nil {
		t.Fatal("unavailable vault accepted a delete")
	}
}

func TestSecretServiceIndividuallyLockedItem(t *testing.T) {
	deleteItem := func(service, user string) (string, error) { return "", Delete(service, user) }
	for _, tc := range []struct {
		name                 string
		operation            func(string, string) (string, error)
		dismiss, denyUnlock  bool
		wantValue, wantError string
		wantItems            int
	}{
		{"get", Get, false, false, "secret", "", 1},
		{"delete", deleteItem, false, false, "", "", 0},
		{"dismiss get", Get, true, false, "", ErrDismissed.Error(), 1},
		{"dismiss delete", deleteItem, true, false, "", ErrDismissed.Error(), 1},
		{"deny get", Get, false, true, "", "credential service did not unlock the requested item", 1},
		{"deny delete", deleteItem, false, true, "", "credential service did not unlock the requested item", 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			v := newTestVault(t)
			if err := Set("wails-test", "signing", "secret"); err != nil {
				t.Fatal(err)
			}
			v.mu.Lock()
			v.dismiss, v.denyUnlock = tc.dismiss, tc.denyUnlock
			for _, item := range v.items {
				item.locked = true
			}
			v.mu.Unlock()
			got, err := tc.operation("wails-test", "signing")
			if got != tc.wantValue {
				t.Fatalf("operation returned %q, want %q", got, tc.wantValue)
			}
			gotError := ""
			if err != nil {
				gotError = err.Error()
			}
			if gotError != tc.wantError {
				t.Fatalf("operation error = %q, want %q", gotError, tc.wantError)
			}
			v.mu.Lock()
			defer v.mu.Unlock()
			if len(v.items) != tc.wantItems || v.opened != v.closed {
				t.Fatalf("items=%d, want %d; sessions=%d/%d", len(v.items), tc.wantItems, v.opened, v.closed)
			}
		})
	}
}
