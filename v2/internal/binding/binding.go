package binding

import (
	"fmt"
	"strings"

	"github.com/wailsapp/wails/v2/internal/logger"
)

type Bindings struct {
	db     *DB
	logger logger.CustomLogger
}

// NewBindings returns a new Bindings object
func NewBindings(logger *logger.Logger, structPointersToBind []interface{}) *Bindings {
	result := &Bindings{
		db:     newDB(),
		logger: logger.CustomLogger("Bindings"),
	}

	// Add the structs to bind
	for _, ptr := range structPointersToBind {
		err := result.Add(ptr)
		if err != nil {
			logger.Fatal("Error during binding: " + err.Error())
		}
	}

	return result
}

// Add the given struct methods to the Bindings
func (b *Bindings) Add(structPtr interface{}) error {

	methods, err := getMethods(structPtr)
	if err != nil {
		return fmt.Errorf("cannot bind value to app: %s", err.Error())
	}

	for _, method := range methods {
		splitName := strings.Split(method.Name, ".")
		packageName := splitName[0]
		structName := splitName[1]
		methodName := splitName[2]

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
