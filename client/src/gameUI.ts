import * as $ from 'jquery'
import { ILogEvent, IWorldState } from "./protocol"
import { lerp } from "./util"

export class Timer {
  private bar: string

  private previousTime: number
  private currentTime: number
  private dt: number
  private lastUpdate: number

  private currentProgress: number
  private previousProgress: number

  constructor(bar: string) {
    this.bar = bar
    this.reset()
  }

  public render() {
    let interpolation = this.interpolation()
    let progress = lerp(this.previousProgress, this.currentProgress, interpolation)
    $(this.bar).css('width', progress+"%")
  }

  interpolation() : number {
      if (this.previousTime != null) {
          let elapsed = Date.now() - this.lastUpdate
          return elapsed / this.dt
      } else {
          return 1
      }
  }

  public pushState(state: IWorldState) {
    let interpolation = this.interpolation()

    this.previousTime = this.currentTime
    this.currentTime = state.time

    this.previousProgress = this.currentProgress
    this.currentProgress = state.progress

    this.lastUpdate = Date.now()
    if (this.previousTime != null) {
        this.dt = this.currentTime - this.previousTime
    }
  }

  public reset() {
    $(this.bar).stop().css("width", "0")
  }
}

export class UserPanel {
    username: string;
    resources: number;
    resourcesEl: string

    constructor(usernameEl: string, resourcesEl: string) {
      this.username = this.getUserID()
      this.resourcesEl = resourcesEl
      $(usernameEl).html(this.username)
    }

    public updateScore(gains: number) {
      $(this.resourcesEl).html(UserPanel.pad(gains, 5))
    }

    private static pad(num: number, size: number) {
       return ('000000000' + num).substr(-size);
    }

    private getUserID() {
      var request = new XMLHttpRequest();
      request.open('GET', '/currentUser', false);
      request.send(null);

      let username = ""

      if (request.status == 200) {
        console.log("User: " + request.responseText);
        username = request.responseText
      }
      if (request.status == 401) {
        console.log("Not logged in");
        username = ""
      }

      return username
    }

    public isUserAuthenticated() {
      return this.username != ""
    }
}

export class GameLog {

    private log: HTMLElement

    constructor(logEl: string) {
      this.log = document.getElementById(logEl)
    }

    public update(logEvent: ILogEvent) {
      let scrollUpdate = this.log.scrollHeight - this.log.clientHeight <= this.log.scrollTop + 1

      switch (logEvent.logType) {
          case 0:
              this.log.innerHTML = "A new game has started, good luck!\n"
              break
          case 1:
              this.log.innerHTML = this.log.innerHTML
              + "<span class='collectResource'>"
              + "Yum, "+ logEvent.protagName
              + " ate a resource\n" + "</span>"
              break
          case 2:
              this.log.innerHTML = this.log.innerHTML
              + "<span class='eatCreature'>"
              + logEvent.protagName + " ate " + logEvent.antagName
              + "\n" + "</span>"
              break
          case 3:
              this.log.innerHTML = this.log.innerHTML
              + "<span class='hitSpike'>Oh no, "
              + logEvent.protagName + " hit a spike!\n"
              +  "</span>"
      }

      if (scrollUpdate) {
          this.log.scrollTop = this.log.scrollHeight - this.log.clientHeight
      }

    }
}
