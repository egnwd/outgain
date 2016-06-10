import * as $ from 'jquery'

interface LeaderboardEntry {
  Name: string
  Score: number
}

function render(entry: LeaderboardEntry, i: number) {
    let className = getClassName(i)
    let body = document.createElement("tr")
    body.setAttribute("class", className)

    let position = document.createElement("td")
    let name = document.createElement("td")
    let score = document.createElement("td")

    position.innerHTML = i+1 + ". "
    position.setAttribute("class", "position")

    name.innerHTML = entry.Name
    name.setAttribute("class", "name")

    score.innerHTML = String(entry.Score)
    score.setAttribute("class", "score")

    body.appendChild(position)
    body.appendChild(name)
    body.appendChild(score)

    return body
}

function getClassName(i: number) {
  switch (i+1) {
    case 1:
      return "first"
    case 2:
      return "second"
    case 3:
      return "third"
    default:
      return ""
  }
}

export default class GameLeaderboard {
  entries: LeaderboardEntry[] = []
  element: string

  constructor(element: string) {
    this.element = element
    this.refresh()
  }

  refresh() {
    $.when(this.update()).then(() => this.render())
  }

  private update() {
    let path = window.location.pathname + "/leaderboard"
    return $.ajax({
      url: path,
    }).done((entries) => {
      this.entries = []
      entries.forEach((entry) => {
        this.entries.push(entry as LeaderboardEntry)
      })
    })
  }

  private render() {
    let elements = this.entries.map((entry, i) => { return render(entry, i) })

    let topHalf = $(this.element).find(".top-half")
    let bottomHalf = $(this.element).find(".bottom-half")

    topHalf.text("")
    bottomHalf.text("")

    elements.forEach((element, i) => {
      if (i < elements.length / 2) {
          topHalf.append(element as any)
      } else {
          bottomHalf.append(element as any)
      }
    })
  }
}
