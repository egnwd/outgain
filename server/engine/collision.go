package engine

type EntityList []Entity

func (list EntityList) Len() int {
	return len(list)
}

func (list EntityList) Less(i, j int) bool {
	return list[i].Base().Left() < list[j].Base().Left()
}

func (list EntityList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list EntityList) Insert(entity Entity) EntityList {
	result := append(list, entity)

	for i := result.Len() - 1; i > 0 && result.Less(i, i-1); i-- {
		result.Swap(i-1, i)
	}

	return result
}

func (list EntityList) Sort() {
	for i := 1; i < list.Len(); i++ {
		for j := i; j > 0 && list.Less(j, j-1); j-- {
			list.Swap(j-1, j)
		}
	}
}

func XOverlap(a, b Entity) bool {
	return a.Base().Right() > b.Base().Left() &&
		b.Base().Right() > a.Base().Left()
}

func YOverlap(a, b Entity) bool {
	return a.Base().Bottom() > b.Base().Top() &&
		b.Base().Bottom() > a.Base().Top()
}

func Collide(a, b Entity) bool {
	radii := a.Base().Radius + b.Base().Radius
	dx := a.Base().X - b.Base().X
	dy := a.Base().Y - b.Base().Y

	distSq := dx*dx + dy*dy
	radiiSq := radii * radii

	return distSq < radiiSq
}

func (list EntityList) Collisions(onCollision func(Entity, Entity)) {
	leftmost := 0

	for i := 1; i < list.Len(); i++ {
		for j := leftmost; j < i; j++ {
			if XOverlap(list[i], list[j]) {
				if YOverlap(list[i], list[j]) && Collide(list[i], list[j]) {
					onCollision(list[i], list[j])
				}
			} else {
				leftmost++
			}
		}
	}
}
