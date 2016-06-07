import * as $ from 'jquery'

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
