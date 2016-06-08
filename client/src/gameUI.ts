import * as $ from 'jquery'

export class Timer {
  private timeTotal: number
  private bar: string

  constructor(time: number, bar: string) {
    this.timeTotal = time
    this.bar = bar
  }

  public start() {
    $(this.bar).animate({width: "100%"}, this.timeTotal, "linear")
  }

  public reset() {
    $(this.bar).stop().css("width", "0")
  }
}

export class UserPanel {
    username: string;
    resources: number;

    constructor(usernameEl: string) {
      this.username = this.getUserID()

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
