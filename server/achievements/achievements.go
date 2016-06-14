package achievements

type Achievements []Achievement

type Achievement struct {
	ID                uint64
	Title             string
	Description       string
	DescriptionLocked string
	Tweet             string
	Slug              string
	Unlocked          bool
}

var dummyData = Achievements{
	Achievement{
		ID:                1,
		Title:             "Top of the Leaderboard",
		Description:       "Congratulations on topping the leaderboard!",
		DescriptionLocked: "Unlock this achievement by topping the global leaderboard!",
		Tweet:             "I reached the top of the leaderboard on Outgain, an AI based game!",
		Slug:              "top",
		Unlocked:          true,
	},
	Achievement{
		ID:                2,
		Title:             "Total 1000 points",
		Description:       "Well Done! You've gained over 1000 points.",
		DescriptionLocked: "Unlock this achievement by accumulating a total of 1000 points!",
		Tweet:             "I have totalled over 1000 points on Outgain, an AI based game!",
		Slug:              "thousand",
		Unlocked:          true,
	},
	Achievement{
		ID:                3,
		Title:             "Score 50 points",
		Description:       "You scored 50 points in a single game!",
		DescriptionLocked: "Unlock this achievement by scoring 50 points in a single game!",
		Tweet:             "I scored 50 points on Outgain, an AI based game!",
		Slug:              "fifty",
		Unlocked:          true,
	},
	Achievement{
		ID:                4,
		Title:             "Create an AI",
		Description:       "Nice! You've made your own AI.",
		DescriptionLocked: "Unlock this achievement by saving your own AI to Github!",
		Tweet:             "I created an AI on Outgain, an AI based game!",
		Slug:              "ai",
		Unlocked:          true,
	},
	Achievement{
		ID:                5,
		Title:             "Avoid All the Spikes",
		Description:       "Amazing! You didn't hit a single spike.",
		DescriptionLocked: "Unlock this achievement by completing a game without hitting a single spike!",
		Tweet:             "I avoided all the spikes on Outgain, an AI based game!",
		Slug:              "no-spikes",
		Unlocked:          true,
	},
}

func GetUserAchievements(_ string) Achievements {
	return dummyData
}
