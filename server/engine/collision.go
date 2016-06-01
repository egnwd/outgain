package engine

func XOverlap(a, b Entity) bool {
	return a.Base().Right() > b.Base().Left() &&
		b.Base().Right() > a.Base().Left()
}

func YOverlap(a, b Entity) bool {
	return a.Base().Bottom() > b.Base().Top() &&
		b.Base().Bottom() > a.Base().Top()
}

// (Potential) Collision between two Entities.
// This may be either a candidate, or a confirmed collision, depending on how
// this was obtained. collisionsBroadPhase yields candidates, where as
// collisionsNarrowPhase only yields actual collisions.
//
// The order between a and b is irrelevant
type Collision struct {
	A, B Entity
}

// Check whether two entities collide
// Entities collide if the distance between their centers is less than the sum
// of their radii
// Order of arguments is irrelevant
func Collide(a, b Entity) bool {
	radii := a.Base().Radius + b.Base().Radius
	dx := a.Base().X - b.Base().X
	dy := a.Base().Y - b.Base().Y
	distSq := dx*dx + dy*dy
	radiiSq := radii * radii
	return distSq < radiiSq
}

// Perform the broad phase of the collision detection on the list of entities
// The returned channel will yield all pairs of Entities which overlap on the
// X axis.
// The input list must be sorted by left edge
//
// The algorithm avoids O(n^2) in most cases by taking advantage that the list
// is sorted by left edge.
//
// While iterating the list, we keep track of the leftmost entity whose right
// edge is greater than the current item's left edge.
//
// In the worst case scenario, where all entities do overlap on the X axis,
// then it ends up being O(n^2) anyway. However, this shouldn't usually happen.
func collisionsBroadPhase(list EntityList) <-chan Collision {
	out := make(chan Collision)
	leftmost := 0

	go func() {
		defer close(out)

		for i, a := range list {
			for j := leftmost; j < i; j++ {
				b := list[j]
				if XOverlap(a, b) {
					out <- Collision{a, b}
				} else if j == leftmost {
					leftmost++
				}
			}
		}
	}()

	return out
}

// Perform the narrow phase of the collision detection.
// The input channel must satisfy the postconditions of  collisionsBroadPhase,
// ie pairs of Entities which overlap on the X axis.
//
// The narrow phase checks whether these actually collide, by checking if they
// also overlap on the Y axis, and if so if the distance between the centers
// is less than the sum of the radii
func collisionsNarrowPhase(in <-chan Collision) <-chan Collision {
	out := make(chan Collision)

	go func() {
		defer close(out)

		for candidate := range in {
			if YOverlap(candidate.A, candidate.B) &&
				Collide(candidate.A, candidate.B) {

				out <- candidate
			}
		}
	}()

	return out
}

// Get all pairs of Entities which collide.
// WARNING: You must not modify the radius nor the coordinates of any entity in the
// list until the output channel is complete
func (list EntityList) Collisions() <-chan Collision {
	return collisionsNarrowPhase(collisionsBroadPhase(list))
}

func (list EntityList) SlowCollisions() <-chan Collision {
	out := make(chan Collision)

	go func() {
		defer close(out)

		for i, a := range list {
			for _, b := range list[:i] {
				if Collide(a, b) {
					out <- Collision{a, b}
				}
			}
		}
	}()

	return out
}
