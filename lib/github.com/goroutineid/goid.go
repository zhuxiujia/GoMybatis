package goroutineid

import "runtime"

var offsetDict = map[string]int64{
	"go1.4":    128,
	"go1.4.1":  128,
	"go1.4.2":  128,
	"go1.4.3":  128,
	"go1.5":    184,
	"go1.5.1":  184,
	"go1.5.2":  184,
	"go1.5.3":  184,
	"go1.5.4":  184,
	"go1.6":    192,
	"go1.6.1":  192,
	"go1.6.2":  192,
	"go1.6.3":  192,
	"go1.6.4":  192,
	"go1.7":    192,
	"go1.7.1":  192,
	"go1.7.2":  192,
	"go1.7.3":  192,
	"go1.7.4":  192,
	"go1.7.5":  192,
	"go1.7.6":  192,
	"go1.8":    192,
	"go1.8.1":  192,
	"go1.8.2":  192,
	"go1.8.3":  192,
	"go1.8.4":  192,
	"go1.8.5":  192,
	"go1.9":    152,
	"go1.9.1":  152,
	"go1.9.2":  152,
	"go1.10":   152,
	"go1.10.1": 152,
}

var offset = offsetDict[runtime.Version()]

// GetGoID returns the goroutine id
func GetGoID() int64
