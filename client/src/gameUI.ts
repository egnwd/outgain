import * as $ from 'jquery'
import * as sweetalert from 'sweetalert'

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

export module ModalPopUp {
    export function mainModal() {
      sweetalert({
          title: "<h1 id=\"title\"></h1>",
          text: "\
          <a href=\"/login\" class=\"btn btn--action\">Login with Github</a>\
          ",
          html: true,
          confirmButtonText: "How to Play",
          closeOnConfirm: false
      }, function() {
          howToPlay()
      })
    }

    export function howToPlay() {
      sweetalert({
          title: "How to Play",
          text: "<p>This is how you play the game</p>",
          html: true,
          confirmButtonText: "Done",
          closeOnConfirm: false
      }, function() {
          mainModal()
      })
    }
}
