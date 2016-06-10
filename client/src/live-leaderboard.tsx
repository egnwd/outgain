import * as $ from 'jquery'
import * as moment from 'moment';
import * as React from './dom';

class Entry {
  position: number
  name: string
  score: number

  constructor(props: any) {
    this.position = props.position
    this.name = props.name
    this.score = props.score
  }

  render() {
     let cls = getClassName(this.position)

     let el =
       <tr class={cls}>
         <td class="position">{this.position+1 + ". "}</td>
         <td class="name">{this.name}</td>
         <td class="score">{String(this.score)}</td>
       </tr>;

     return el;
 }
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
  entries: any[] = []
  element: string

  constructor(element: string) {
    this.element = element
    this.refresh()
  }

  refresh() {
    $.when(this.update()).then(() => this.draw())
  }

  private update() {
    let path = window.location.pathname + "/leaderboard"
    return $.ajax({
      url: path,
    }).done((entries) => {
      this.entries = []
      entries.forEach((entry) => {
        this.entries.push(entry)
      })
    })
  }

  private draw() {
    let elements = this.entries.map((entry, i) => {
      let props = {position: i, name: entry.Name, score: entry.Score}
      return new Entry(props).render()
    })

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
