package binding

import (
	"fmt"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/logger"
)

type Bindings struct {
	db     *DB
	logger logger.CustomLogger
}

// NewBindings returns a new Bindings object
func NewBindings(logger *logger.Logger) *Bindings {
	return &Bindings{
		db:     newDB(),
		logger: logger.CustomLogger("Bindings"),
	}
}

// Add the given struct methods to the Bindings
func (b *Bindings) Add(structPtr interface{}) error {

	methods, err := getMethods(structPtr)
	if err != nil {
		return fmt.Errorf("cannout bind value to app: %s", err.Error())
	}

	for _, method := range methods {
		splitName := strings.Split(method.Name, ".")
		packageName := splitName[0]
		structName := splitName[1]
		methodName := splitName[2]

		// Is this WailsInit?
		if method.IsWailsInit() {
			err := b.db.AddWailsInit(method)
			if err != nil {
				return err
			}
			b.logger.Trace("Registered WailsInit method: %s", method.Name)
			continue
		}

		// Is this WailsShutdown?
		if method.IsWailsShutdown() {
			err := b.db.AddWailsShutdown(method)
			if err != nil {
				return err
			}
			b.logger.Trace("Registered WailsShutdown method: %s", method.Name)
			continue
		}

		// Add it as a regular method
		b.db.AddMethod(packageName, structName, methodName, method)

	}
	return nil
}

func (b *Bindings) DB() *DB {
	return b.db
}

func (b *Bindings) ToJSON() (string, error) {
	return b.db.ToJSON()
}
