//go:build linux || freebsd || openbsd || netbsd || dragonfly

package credentialstore

import (
	"errors"
	"fmt"

	"github.com/godbus/dbus/v5"
)

const secretService = "org.freedesktop.secrets"
const serviceInterface = "org.freedesktop.Secret.Service"
const collectionInterface = "org.freedesktop.Secret.Collection"
const itemInterface = "org.freedesktop.Secret.Item"
const promptInterface = "org.freedesktop.Secret.Prompt"
const servicePath dbus.ObjectPath = "/org/freedesktop/secrets"

type wireSecret struct {
	Session     dbus.ObjectPath
	Parameters  []byte
	Value       []byte
	ContentType string
}

type secretClient struct {
	conn    *dbus.Conn
	session dbus.ObjectPath
}

func openSecrets() (*secretClient, error) {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		return nil, err
	}
	client := &secretClient{conn: conn}
	var output dbus.Variant
	err = client.object(servicePath).Call(serviceInterface+".OpenSession", 0,
		"plain", dbus.MakeVariant("")).Store(&output, &client.session)
	if err == nil && client.session == "/" {
		err = errors.New("credential service did not open a session")
	}
	if err != nil {
		conn.Close()
		return nil, err
	}
	return client, nil
}

func (c *secretClient) object(path dbus.ObjectPath) dbus.BusObject {
	return c.conn.Object(secretService, path)
}

func (c *secretClient) close() {
	c.object(c.session).Call("org.freedesktop.Secret.Session.Close", 0)
	c.conn.Close()
}

func attributes(service, user string) map[string]string {
	return map[string]string{"service": service, "username": user}
}

func (c *secretClient) prompt(path dbus.ObjectPath) (dbus.Variant, error) {
	if path == "/" {
		return dbus.MakeVariant(""), nil
	}
	match := []dbus.MatchOption{
		dbus.WithMatchSender(secretService), dbus.WithMatchObjectPath(path),
		dbus.WithMatchInterface(promptInterface), dbus.WithMatchMember("Completed"),
	}
	if err := c.conn.AddMatchSignal(match...); err != nil {
		return dbus.Variant{}, err
	}
	defer c.conn.RemoveMatchSignal(match...)
	signals := make(chan *dbus.Signal, 1)
	c.conn.Signal(signals)
	defer c.conn.RemoveSignal(signals)
	// Register before Prompt: a vault may emit Completed before replying.
	if err := c.object(path).Call(promptInterface+".Prompt", 0, "").Err; err != nil {
		return dbus.Variant{}, err
	}
	for {
		select {
		case signal, ok := <-signals:
			if !ok {
				return dbus.Variant{}, errors.New("credential service connection closed")
			}
			if signal.Path != path || signal.Name != promptInterface+".Completed" {
				continue
			}
			var dismissed bool
			var result dbus.Variant
			if err := dbus.Store(signal.Body, &dismissed, &result); err != nil {
				return dbus.Variant{}, err
			}
			if dismissed {
				return dbus.Variant{}, ErrDismissed
			}
			return result, nil
		case <-c.conn.Context().Done():
			return dbus.Variant{}, c.conn.Context().Err()
		}
	}
}

func (c *secretClient) locked(path dbus.ObjectPath, iface string) (bool, error) {
	property, err := c.object(path).GetProperty(iface + ".Locked")
	if err != nil {
		return false, err
	}
	locked, ok := property.Value().(bool)
	if !ok {
		return false, errors.New("credential service returned an invalid Locked property")
	}
	return locked, nil
}

func (c *secretClient) unlock(path dbus.ObjectPath, iface string) error {
	locked, err := c.locked(path, iface)
	if err != nil || !locked {
		return err
	}
	var unlocked []dbus.ObjectPath
	var prompt dbus.ObjectPath
	if err := c.object(servicePath).Call(serviceInterface+".Unlock", 0,
		[]dbus.ObjectPath{path}).Store(&unlocked, &prompt); err != nil {
		return err
	}
	if _, err := c.prompt(prompt); err != nil {
		return err
	}
	locked, err = c.locked(path, iface)
	if err != nil {
		return err
	}
	if locked {
		return errors.New("credential service did not unlock the requested item")
	}
	return nil
}

func (c *secretClient) collection(create bool) (dbus.ObjectPath, error) {
	var collection dbus.ObjectPath
	if err := c.object(servicePath).Call(serviceInterface+".ReadAlias", 0, "default").Store(&collection); err != nil {
		return "", err
	}
	if collection == "/" && !create {
		return "", ErrNotFound
	}
	if collection == "/" {
		var prompt dbus.ObjectPath
		properties := map[string]dbus.Variant{collectionInterface + ".Label": dbus.MakeVariant("Login")}
		if err := c.object(servicePath).Call(serviceInterface+".CreateCollection", 0,
			properties, "default").Store(&collection, &prompt); err != nil {
			return "", err
		}
		if prompt != "/" {
			result, err := c.prompt(prompt)
			if err != nil {
				return "", err
			}
			if err := dbus.Store([]any{result.Value()}, &collection); err != nil {
				return "", err
			}
		}
	}
	if collection == "/" {
		return "", errors.New("credential service has no default collection")
	}
	return collection, c.unlock(collection, collectionInterface)
}

func Set(service, user, value string) error {
	client, err := openSecrets()
	if err != nil {
		return err
	}
	defer client.close()
	collection, err := client.collection(true)
	if err != nil {
		return err
	}
	properties := map[string]dbus.Variant{
		itemInterface + ".Attributes": dbus.MakeVariant(attributes(service, user)),
		itemInterface + ".Label":      dbus.MakeVariant(fmt.Sprintf("Password for '%s' on '%s'", user, service)),
	}
	secret := wireSecret{client.session, []byte{}, []byte(value), "text/plain; charset=utf-8"}
	var item, prompt dbus.ObjectPath
	if err := client.object(collection).Call(collectionInterface+".CreateItem", 0,
		properties, secret, true).Store(&item, &prompt); err != nil {
		return err
	}
	if prompt != "/" {
		result, err := client.prompt(prompt)
		if err != nil {
			return err
		}
		if err := dbus.Store([]any{result.Value()}, &item); err != nil {
			return err
		}
	}
	if item == "/" {
		return errors.New("credential service did not store the item")
	}
	return nil
}

func (c *secretClient) search(service, user string) ([]dbus.ObjectPath, error) {
	collection, err := c.collection(false)
	if err != nil {
		return nil, err
	}
	var items []dbus.ObjectPath
	err = c.object(collection).Call(collectionInterface+".SearchItems", 0,
		attributes(service, user)).Store(&items)
	if err == nil && len(items) == 0 {
		err = ErrNotFound
	}
	return items, err
}

func Get(service, user string) (string, error) {
	client, err := openSecrets()
	if err != nil {
		return "", err
	}
	defer client.close()
	items, err := client.search(service, user)
	if err != nil {
		return "", err
	}
	item := items[0]
	if err := client.unlock(item, itemInterface); err != nil {
		return "", err
	}
	var secret wireSecret
	if err := client.object(item).Call(itemInterface+".GetSecret", 0, client.session).Store(&secret); err != nil {
		return "", err
	}
	if secret.Session != client.session || len(secret.Parameters) != 0 {
		return "", errors.New("credential service returned an invalid plain-session secret")
	}
	return string(secret.Value), nil
}

func Delete(service, user string) error {
	client, err := openSecrets()
	if err != nil {
		return err
	}
	defer client.close()
	items, err := client.search(service, user)
	if err != nil {
		return err
	}
	for _, item := range items {
		if err := client.unlock(item, itemInterface); err != nil {
			return err
		}
		var prompt dbus.ObjectPath
		if err := client.object(item).Call(itemInterface+".Delete", 0).Store(&prompt); err != nil {
			return err
		}
		if _, err := client.prompt(prompt); err != nil {
			return err
		}
	}
	return nil
}
