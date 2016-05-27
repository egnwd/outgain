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

  sweetalert.setDefaults({
    closeOnConfirm: false,
    customClass: "modal",
    html: true, });


    export function mainModal() {
      sweetalert({
          title: "<h1 id=\"title\"></h1>",
          text: "<a href=\"/login\" class=\"btn btn--action\">Login with Github</a>",
          confirmButtonText: "How to Play",
      }, function() {
          howToPlay()
      })
    }

    export function howToPlay() {
      sweetalert({
          title: "<h2 class=\"modal-title\">How to Play</h2>",
          text: "<p class=\"howtoplay left\">Lorem ipsum dolor sit amet, consectetur adipiscing elit. Maecenas sollicitudin felis sed interdum aliquam. Sed volutpat quis turpis ac tincidunt. Donec interdum enim id enim congue, sed condimentum nibh pellentesque. Phasellus vulputate magna efficitur ante tristique commodo.<br><br> Sed consequat neque quis rutrum ultrices. Vestibulum hendrerit tincidunt ligula, vel rhoncus nibh dapibus euismod. In eget viverra nulla. Vestibulum porta commodo lectus, vulputate laoreet lectus. Mauris sed sollicitudin sem. Aliquam erat volutpat. Nulla ullamcorper ante a lacinia accumsan. Proin placerat pharetra lectus, et rhoncus nibh mattis sed.<br><br></p>",
          confirmButtonText: "Back",
      }, function() {
          mainModal()
      })
    }
}
