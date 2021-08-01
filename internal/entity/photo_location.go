package entity

import (
	"time"

	"github.com/photoprism/photoprism/internal/classify"
	"github.com/photoprism/photoprism/internal/maps"
	"gopkg.in/photoprism/go-tz.v2/tz"
)

// GetTimeZone uses PhotoLat and PhotoLng to guess the time zone of the photo.
func (m *Photo) GetTimeZone() string {
	result := "UTC"

	if m.HasLatLng() {
		zones, err := tz.GetZone(tz.Point{
			Lat: float64(m.PhotoLat),
			Lon: float64(m.PhotoLng),
		})

		if err == nil && len(zones) > 0 {
			result = zones[0]
		}
	}

	return result
}

// CountryName returns the photo country name.
func (m *Photo) CountryName() string {
	if name, ok := maps.CountryNames[m.CountryCode()]; ok {
		return name
	}

	return UnknownCountry.CountryName
}

// CountryCode returns the photo country code.
func (m *Photo) CountryCode() string {
	if len(m.PhotoCountry) != 2 {
		m.PhotoCountry = UnknownCountry.ID
	}

	return m.PhotoCountry
}

// GetTakenAt returns UTC time for TakenAtLocal.
func (m *Photo) GetTakenAt() time.Time {
	location, err := time.LoadLocation(m.TimeZone)

	if err != nil {
		return m.TakenAt
	}

	if takenAt, err := time.ParseInLocation("2006-01-02T15:04:05", m.TakenAtLocal.Format("2006-01-02T15:04:05"), location); err != nil {
		return m.TakenAt
	} else {
		return takenAt.UTC()
	}
}

// GetTakenAtLocal returns local time for TakenAt.
func (m *Photo) GetTakenAtLocal() time.Time {
	location, err := time.LoadLocation(m.TimeZone)

	if err != nil {
		return m.TakenAtLocal
	}

	if takenAtLocal, err := time.ParseInLocation("2006-01-02T15:04:05", m.TakenAt.In(location).Format("2006-01-02T15:04:05"), time.UTC); err != nil {
		return m.TakenAtLocal
	} else {
		return takenAtLocal.UTC()
	}
}

// UpdateLocation updates location and labels based on latitude and longitude.
func (m *Photo) UpdateLocation() (keywords []string, labels classify.Labels) {
	if m.HasLatLng() {
		var location = NewCell(m.PhotoLat, m.PhotoLng)

		err := location.Find(GeoApi)

		if location.Place == nil {
			log.Warnf("photo: failed fetching geo data (uid %s, cell %s)", m.PhotoUID, location.ID)
		} else if err == nil && location.ID != UnknownLocation.ID {
			m.Cell = location
			m.CellID = location.ID
			m.Place = location.Place
			m.PlaceID = location.PlaceID
			m.PhotoCountry = location.CountryCode()

			if m.TakenSrc != SrcManual {
				m.TimeZone = m.GetTimeZone()
				m.TakenAt = m.GetTakenAt()
			}

			FirstOrCreateCountry(NewCountry(location.CountryCode(), location.CountryName()))

			locCategory := location.Category()
			keywords = append(keywords, location.Keywords()...)

			// Append category from reverse location lookup
			if locCategory != "" {
				labels = append(labels, classify.LocationLabel(locCategory, 0))
			}

			return keywords, labels
		}
	}

	keywords = []string{}
	labels = classify.Labels{}

	if m.UnknownLocation() {
		m.Cell = &UnknownLocation
		m.CellID = UnknownLocation.ID

		// Remove place estimate if better data is available.
		if SrcPriority[m.PlaceSrc] > SrcPriority[SrcEstimate] {
			m.Place = &UnknownPlace
			m.PlaceID = UnknownPlace.ID
		}
	} else if err := m.LoadLocation(); err == nil {
		m.Place = m.Cell.Place
		m.PlaceID = m.Cell.PlaceID
	} else {
		log.Warnf("photo: location %s not found in %s", m.CellID, m.PhotoName)
	}

	if m.UnknownPlace() {
		m.Place = &UnknownPlace
		m.PlaceID = UnknownPlace.ID
	} else if err := m.LoadPlace(); err == nil {
		m.PhotoCountry = m.Place.CountryCode()
	} else {
		log.Warnf("photo: place %s not found in %s", m.PlaceID, m.PhotoName)
	}

	if m.UnknownCountry() {
		m.EstimateCountry()
	}

	if m.HasCountry() {
		FirstOrCreateCountry(NewCountry(m.CountryCode(), m.CountryName()))
	}

	return keywords, labels
}
