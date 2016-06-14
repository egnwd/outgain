package achievements

import (
	"github.com/egnwd/outgain/server/database"
)

type Achievements []Achievement

type ConditionFunc func(*database.AchievementData, *LiveData) bool

type Achievement struct {
	ID                uint64
	Title             string
	Condition         ConditionFunc
	Description       string
	DescriptionLocked string
	Tweet             string
	Slug              string
	Unlocked          bool
}

// Must be less than 32 without having to change from uint32 bitmap
const numAchievements = 5

type LiveData struct {
	Spikes    int
	Resources int
	Creatures int
}

type LiveDataMap map[string]LiveData

var achievements = Achievements{
	Achievement{
		ID:    1,
		Title: "Top of the Leaderboard",
		// FIXME: don't just return false..
		Condition: func(data *database.AchievementData, liveData *LiveData) bool {
			return false
		},
		Description:       "Congratulations on topping the leaderboard!",
		DescriptionLocked: "Unlock this achievement by topping the global leaderboard!",
		Tweet:             "I reached the top of the leaderboard on Outgain, an AI based game!",
		Slug:              "top",
		Unlocked:          true,
	},
	Achievement{
		ID:    2,
		Title: "Total 1000 points",
		Condition: func(data *database.AchievementData, liveData *LiveData) bool {
			return data.TotalScore > 1000
		},
		Description:       "Well Done! You've gained over 1000 points.",
		DescriptionLocked: "Unlock this achievement by accumulating a total of 1000 points!",
		Tweet:             "I have totalled over 1000 points on Outgain, an AI based game!",
		Slug:              "thousand",
		Unlocked:          true,
	},
	Achievement{
		ID:    3,
		Title: "Score 50 points",
		Condition: func(data *database.AchievementData, liveData *LiveData) bool {
			return data.HighScore > 50
		},
		Description:       "You scored 50 points in a single game!",
		DescriptionLocked: "Unlock this achievement by scoring 50 points in a single game!",
		Tweet:             "I scored 50 points on Outgain, an AI based game!",
		Slug:              "fifty",
		Unlocked:          true,
	},
	Achievement{
		ID:    4,
		Title: "Create an AI",
		// Unlocked another way, if checked by this function should always be false
		Condition: func(data *database.AchievementData, liveData *LiveData) bool {
			return false
		},
		Description:       "Nice! You've made your own AI.",
		DescriptionLocked: "Unlock this achievement by saving your own AI to Github!",
		Tweet:             "I created an AI on Outgain, an AI based game!",
		Slug:              "ai",
		Unlocked:          true,
	},
	Achievement{
		ID:    5,
		Title: "Avoid All the Spikes",
		Condition: func(data *database.AchievementData, liveData *LiveData) bool {
			return liveData.Spikes == 0
		},
		Description:       "Amazing! You didn't hit a single spike.",
		DescriptionLocked: "Unlock this achievement by completing a game without hitting a single spike!",
		Tweet:             "I avoided all the spikes on Outgain, an AI based game!",
		Slug:              "no-spikes",
		Unlocked:          true,
	},
}

func GetUserAchievements(username string) Achievements {
	a := make(Achievements, len(achievements))
	copy(a, achievements)
	data := database.GetAchievements(username)
	bitmap := data.Bitmap
	var i uint8
	for i = 0; i < numAchievements; i++ {
		var mask uint32 = 1 << i
		a[i].Unlocked = (bitmap & mask) == mask
	}
	return a
}

// Update changes row values to new ones
func Update(data *database.AchievementData, liveData *LiveData, gains int) {
	if data.HighScore < gains {
		data.HighScore = gains
	}
	data.TotalScore += gains
	data.Spikes += liveData.Spikes
	data.Resources += liveData.Resources
	data.Creatures += liveData.Creatures
	updateBitmap(data, liveData)
}

// UpdateBitmap changes bitmap values if achievements are now unlocked
func updateBitmap(data *database.AchievementData, liveData *LiveData) {
	// Iterate through achievements bitmap
	bitmap := data.Bitmap
	var i uint8
	for i = 0; i < numAchievements; i++ {
		// Only look at locked achievements
		var mask uint32 = 1 << i
		if (bitmap & mask) != mask {
			// Check if each locked achievement is now unlocked
			if unlocked := checkUnlock(i, data, liveData); unlocked {
				// Update bitmap value
				bitmap |= mask
			}
		}
	}
	data.Bitmap = bitmap
}

func checkUnlock(i uint8, data *database.AchievementData, liveData *LiveData) bool {
	return achievements[i].Condition(data, liveData)
}

func CreatedAI(username string) {
	unlockAchievement(username, 3)
}

func unlockAchievement(username string, achievementIndex uint8) {
	data := database.GetAchievements(username)
	var mask uint32 = 1 << achievementIndex
	data.Bitmap |= mask
	database.UpdateAchievements(data)
}
