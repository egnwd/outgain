import * as $ from 'jquery'
import { ILogEvent } from "./protocol";

export class Timer {
  private timeTotal: number
  private bar: string

  constructor(time: number, bar: string) {
    this.timeTotal = time
    this.bar = bar
    this.reset()
  }

  public start(progress?: number) {
    progress = progress || 0
    console.log(progress)
    $(this.bar).stop().css("width", progress+"%")
    $(this.bar).animate({width: "100%"}, this.timeTotal, "linear")
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
      // The colours should probbaly be done with CSS,
      // leaving it here so that someone who cares more
      // than me can play with them easier
      switch (logEvent.logType) {
          case 0:
              this.log.innerHTML = "A new game has started, good luck!\n"
              break
          case 1:
              this.log.innerHTML = this.log.innerHTML
              + "<span style='color:#9FC155'>"
              + "Yum, "+ logEvent.protagName
              + " ate a resource\n" + "</span>"
              break
          case 2:
              this.log.innerHTML = this.log.innerHTML
              + "<span style='color:#AAE2E8'>"
              + logEvent.protagName + " ate " + logEvent.antagName
              + "\n" + "</span>"
              break
          case 3:
              this.log.innerHTML = this.log.innerHTML
              + "<span style='color:#F6A27F'>Oh no, "
              + logEvent.protagName + " hit a spike!\n"
              +  "</span>"
      }

      if (scrollUpdate) {
          this.log.scrollTop = this.log.scrollHeight - this.log.clientHeight
      }

    }
}
