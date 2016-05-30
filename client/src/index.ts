import * as $ from 'jquery'
import * as sweetalert from 'sweetalert'

module ModalPopUp {

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

    function howToPlay() {
      sweetalert({
          title: "<h2 class=\"modal-title\">How to Play</h2>",
          text: "<p class=\"howtoplay left\">Lorem ipsum dolor sit amet, consectetur adipiscing elit. Maecenas sollicitudin felis sed interdum aliquam. Sed volutpat quis turpis ac tincidunt. Donec interdum enim id enim congue, sed condimentum nibh pellentesque. Phasellus vulputate magna efficitur ante tristique commodo.<br><br> Sed consequat neque quis rutrum ultrices. Vestibulum hendrerit tincidunt ligula, vel rhoncus nibh dapibus euismod. In eget viverra nulla. Vestibulum porta commodo lectus, vulputate laoreet lectus. Mauris sed sollicitudin sem. Aliquam erat volutpat. Nulla ullamcorper ante a lacinia accumsan. Proin placerat pharetra lectus, et rhoncus nibh mattis sed.<br><br></p>",
          confirmButtonText: "Back",
      }, function() {
          mainModal()
      })
    }
}

$(function() {
  ModalPopUp.mainModal()
})
