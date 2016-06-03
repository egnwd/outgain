import * as $ from 'jquery'

export class UserPanel {
    username: string;
    resources: number;
    private usernameEl: string;
    private resourcesEl: string;

    constructor(usernameEl: string, resourcesEl: string) {
      this.usernameEl = usernameEl
      this.resourcesEl = resourcesEl
      this.username = this.getUserID()
    }

    public setUserID() {
      $(this.usernameEl).html(this.username)
    }

    public getUserID() {
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
