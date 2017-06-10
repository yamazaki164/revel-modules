package gormctx

import (
	"github.com/jinzhu/gorm"

	"github.com/revel/revel"
)

// Database connection variables
var (
	Db     *gorm.DB
	Driver string
	Spec   string
)

// Initialize function
func Init() {
	// Read configuration.
	var found bool
	if Driver, found = revel.Config.String("db.driver"); !found {
		revel.ERROR.Fatal("db.driver not configured")
	}
	if Spec, found = revel.Config.String("db.spec"); !found {
		revel.ERROR.Fatal("db.spec not configured")
	}

	// Open a connection.
	var err error
	Db, err = gorm.Open(Driver, Spec)
	if err != nil {
		revel.ERROR.Fatal(err)
	}
}

// Transactional definition for database transaction
type Transactional struct {
	*revel.Controller
	Txn *gorm.DB
}

// Begin a transaction
func (c *Transactional) Begin() revel.Result {
	txn := Db.Begin()
	if txn.Error != nil {
		panic(txn.Error)
	}
	c.Txn = txn
	return nil
}

// Rollback if it's still going (must have panicked).
func (c *Transactional) Rollback() revel.Result {
	if c.Txn != nil {
		txn := c.Txn.Rollback()
		if txn.Error != nil {
			panic(txn.Error)
		}
		c.Txn = nil
	}
	return nil
}

// Commit the transaction.
func (c *Transactional) Commit() revel.Result {
	if c.Txn != nil {
		txn := c.Txn.Commit()
		if txn.Error != nil {
			panic(txn.Error)
		}
		c.Txn = nil
	}
	return nil
}

func init() {
	revel.InterceptMethod((*Transactional).Begin, revel.BEFORE)
	revel.InterceptMethod((*Transactional).Commit, revel.AFTER)
	revel.InterceptMethod((*Transactional).Rollback, revel.FINALLY)
}
