package engine

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func dummyEntity(id uint64, x, y, radius float64) Entity {
	return &Resource{EntityBase: EntityBase{
		ID:     id,
		X:      x,
		Y:      y,
		Radius: radius,
	}}
}

func TestEntityListSort(t *testing.T) {
	list := EntityList{
		dummyEntity(1, 1, 0, 0.5),
		dummyEntity(2, 1, 0, 0.7),
		dummyEntity(3, 2, 0, 0.5),
	}

	// Lets not sort an already sorted list
	assert.False(t, sort.IsSorted(list))

	list.Sort()

	assert.True(t, sort.IsSorted(list), "insertion should keep list sorted")
}

func TestEntityListInsert(t *testing.T) {
	list := EntityList{
		dummyEntity(1, 1, 0, 0.7),
		dummyEntity(2, 1, 0, 0.5),
		dummyEntity(3, 2, 0, 0.5),
	}

	// Ensure precondition for list.Insert
	assert.True(t, sort.IsSorted(list))

	list = list.Insert(dummyEntity(4, 1, 0, 0.3))

	assert.True(t, sort.IsSorted(list), "insertion should keep list sorted")
}

func assertSingleCollision(t *testing.T, expectedA, expectedB uint64, list EntityList) {
	hadCollision := false
	for collision := range list.Collisions() {
		gotA, gotB := collision.a.Base().ID, collision.b.Base().ID

		ok := (gotA == expectedA && gotB == expectedB) || (gotA == expectedB && gotB == expectedA)
		assert.True(t, ok, "Wrong collision reported, expected (%d, %d), got (%d, %d)",
			expectedA, expectedB, gotA, gotB)

		if ok {
			assert.False(t, hadCollision, "Collision (%d, %d) reported twice", expectedA, expectedB)
		}

		hadCollision = true
	}

	assert.True(t, hadCollision, "Collision (%d, %d) not reported", expectedA, expectedB)
}

func TestCollisions(t *testing.T) {
	assertSingleCollision(t, 0, 3, EntityList{
		dummyEntity(0, 3.76, 5.79, 2.12),
		dummyEntity(1, 3.49, 9.01, 0.50),
		dummyEntity(2, 4.65, 8.32, 0.50),
		dummyEntity(3, 4.25, 7.03, 0.10),
	})

	assertSingleCollision(t, 1, 4, EntityList{
		dummyEntity(0, 0.4, 6.2, 0.10),
		dummyEntity(1, 5.8, 5.6, 3.64),
		dummyEntity(2, 6.1, 0.0, 0.10),
		dummyEntity(3, 6.4, 9.3, 0.10),
		dummyEntity(4, 6.7, 5.8, 0.10),
	})
}
