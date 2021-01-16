package html

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPropertyFromLink(t *testing.T) {
	t.Run("get property", func(t *testing.T) {
		client := http.Client{Timeout: 10 * time.Second}
		c := NewCrawler(client)
		prop, err := c.propertyFromLink("fakeId", "http://www.presnc.org/properties/godette-hotel/")
		assert.Nil(t, err)
		assert.Equal(t, "fakeId", prop.ID)
		assert.Equal(t, "Godette Hotel", prop.Name)
		assert.Equal(t, "400 Pollock Street", prop.Location.Address)
		assert.Equal(t, "Beaufort", prop.Location.City)
		assert.Equal(t, "NC", prop.Location.State)
		assert.Equal(t, "28516", prop.Location.Zip)
		assert.Equal(t, "Carteret County", prop.Location.County)
		assert.Equal(t, float64(250000), prop.Price)
	})
}
