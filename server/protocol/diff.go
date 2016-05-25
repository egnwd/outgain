package protocol

func DiffWorld(previous WorldState, next WorldState) WorldDiff {
	diff := WorldDiff{
		PreviousTime: previous.Time,
		Time:         next.Time,
		Modified:     make([]Entity, 0),
		Removed:      make([]uint64, 0),
	}

	previousEntities := make(map[uint64]Entity)

	for _, entity := range previous.Entities {
		previousEntities[entity.Id] = entity
	}

	for _, entity := range next.Entities {
		previousEntity, ok := previousEntities[entity.Id]
		if ok {
			delete(previousEntities, entity.Id)
			if previousEntity != entity {
				diff.Modified = append(diff.Modified, entity)
			}
		} else {
			diff.Modified = append(diff.Modified, entity)
		}
	}

	for id, _ := range previousEntities {
		diff.Removed = append(diff.Removed, id)
	}

	return diff
}
