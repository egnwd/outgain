import { UserPanel } from './gameUI'
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

var userPanel = new UserPanel("#user-id", "#user-resources")

$(function() {
  if (!userPanel.isUserAuthenticated()) {
    ModalPopUp.mainModal()
  } else {
    userPanel.setUserID()
  }
})

$(function() {
  // Create HTML table showing all lobby IDs with clickable rows

  // Get lobby IDs from server
  $.ajax({
    url: "/peekLobbies",
    dataType: 'json'
  })
  .done(function(data) {
    let lobbies = data

    // Generate table of lobby IDs
    let table = "<table id='lobbies'><thead><tr><th class='right'>"
            + "Lobbies<a class='btn'>+</a></th></tr></thead>"
    for (var i = 0; i < lobbies.length; i++) {
      table += "<tr><td class='right'>" + lobbies[i] + "</td></tr>"
    }
    if (lobbies.length == 0) {
      table += "<tr><td class='right'>No available lobbies</td></tr>"
    }
    table += "</table>"
    document.getElementById("lobbies-table").innerHTML = table

    // Add clickable functions to table rows
    let rows = document.getElementById("lobbies").getElementsByTagName("tr")
    for (i = 1; i < rows.length; i++) {
      let curr = rows[i]
      let createClickHandler = function(row) {
        return function() { 
          let id = row.getElementsByTagName("td")[0].innerHTML
          // TODO: replace with proper onclick function
          // Keep in mind filler row when lobbies.length == 0
          alert("id:" + id)
        }
      }
      curr.onclick = createClickHandler(curr);
    }
  })

  // TODO: Row click -> create html table of players in that lobby with join button
  // TODO: Join click -> add user to selected lobby, redirect to game view
})

$(function() {
  $(".lobby").click(function() {
    console.log("row clicked")
  })
})
