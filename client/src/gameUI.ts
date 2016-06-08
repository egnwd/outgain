import * as $ from 'jquery'

export class Timer {
  private percent = 0
  private timeTotal: number
  private timeInterval = 1
  private bar: string

  constructor(time: number, bar: string) {
    this.timeTotal = time
    this.bar = bar
  }

  private animateUpdate() {
    console.log(this.TotalTime)
    if (this.percent < this.timeTotal) {
      this.percent++;
      this.updateProgress(this.percent);
      setTimeout(this.animateUpdate, this.timeInterval);
    }
  }

  private updateProgress(percentage: number) {
    let x = (percentage/this.timeTotal) * 100
    console.log(this.bar)
    console.log("X: " + x)
    $(this.bar).css("width", x + "%");
  }

  public start() {
    console.log("start timer")
    this.animateUpdate()
  }

  public reset() {
    this.percent = 0
  }
}

export class UserPanel {
    username: string;
    resources: number;

    constructor(usernameEl: string) {
      this.username = this.getUserID()

      $(usernameEl).html(this.username)
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
