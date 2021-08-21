package entity

import (
	"fmt"
	"reflect"
	"strings"
)

// Save updates an entity in the database, or inserts if it doesn't exist.
func Save(m interface{}, primaryKeys ...string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("save: %s (panic)", r)
			log.Error(err)
		}
	}()

	if err := Update(m, primaryKeys...); err == nil {
		return nil
	} else if err := UnscopedDb().Save(m).Error; err == nil {
		return nil
	} else if !strings.Contains(strings.ToLower(err.Error()), "lock") {
		return err
	} else if err := UnscopedDb().Save(m).Error; err != nil {
		return err
	}

	return nil
}

// Update updates an existing entity in the database.
func Update(m interface{}, primaryKeys ...string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("update: %s (panic)", r)
			log.Error(err)
		}
	}()

	v := reflect.ValueOf(m).Elem()

	// Abort if a primary key is zero.
	for _, k := range primaryKeys {
		if field := v.FieldByName(k); field.IsZero() {
			return fmt.Errorf("key '%s' not found", k)
		}
	}

	// Update all values except primary keys.
	if res := UnscopedDb().Model(m).Updates(GetValues(m, primaryKeys...)); res.Error != nil {
		return res.Error
	} else if res.RowsAffected > 1 {
		log.Warnf("update: more than one row affected")
	} else if res.RowsAffected == 0 {
		// MariaDB may report zero rows in case no data was actually changed, even though the row exists.
		log.Tracef("update: no rows affected")
	}

	return nil
}
