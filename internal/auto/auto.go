/*

Package auto contains index & import background workers.

Copyright (c) 2018 - 2021 Michael Mayer <hello@photoprism.org>

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU Affero General Public License as published
    by the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU Affero General Public License for more details.

    You should have received a copy of the GNU Affero General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.

    PhotoPrism® is a registered trademark of Michael Mayer.  You may use it as required
    to describe our software, run your own server, for educational purposes, but not for
    offering commercial goods, products, or services without prior written permission.
    In other words, please ask.

Feel free to send an e-mail to hello@photoprism.org if you have questions,
want to support our work, or just want to say hello.

Additional information can be found in our Developer Guide:
https://docs.photoprism.org/developer-guide/

*/
package auto

import (
	"time"

	"github.com/photoprism/photoprism/internal/config"

	"github.com/photoprism/photoprism/internal/event"
)

var log = event.Log

var stop = make(chan bool, 1)

// Wait starts waiting for indexing & importing opportunities.
func Start(conf *config.Config) {
	// Don't start ticker if both are disabled.
	if conf.AutoIndex().Seconds() <= 0 && conf.AutoImport().Seconds() <= 0 {
		return
	}

	ticker := time.NewTicker(time.Minute)

	go func() {
		for {
			select {
			case <-stop:
				ticker.Stop()
				return
			case <-ticker.C:
				if mustIndex(conf.AutoIndex()) {
					log.Debugf("auto-index: starting")
					ResetIndex()
					if err := Index(); err != nil {
						log.Errorf("auto-index: %s", err)
					}
				} else if mustImport(conf.AutoImport()) {
					log.Debugf("auto-import: starting")
					ResetImport()
					if err := Import(); err != nil {
						log.Errorf("auto-import: %s", err)
					}
				}
			}
		}
	}()
}

// Stop stops waiting for indexing & importing opportunities.
func Stop() {
	stop <- true
}
