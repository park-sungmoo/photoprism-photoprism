/*
Package config contains filesystem related utility functions.

Additional information can be found in our Developer Guide:

https://github.com/photoprism/photoprism/wiki
*/
package util

import "github.com/sirupsen/logrus"

var log *logrus.Logger

func init() {
	log = logrus.StandardLogger()
}
