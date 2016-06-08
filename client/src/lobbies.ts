import * as $ from 'jquery'
import * as sweetalert from 'sweetalert'

$(function() {
    // Create HTML table showing all lobby IDs with clickable rows
    updateTable()

    $("#create-lobby").click(createLobby)
})


var createLobby = function(event) {
  event.preventDefault();

  sweetalert({title: "Create a Lobby",
        text: "Name your lobby:",
        type: "input",
        showCancelButton: true,
        closeOnConfirm: false,
        showLoaderOnConfirm: true,
        animation: "slide-from-top",
        inputPlaceholder: "Foo bar..."
      },
      function(name){
        if (name === false) return false;
        if (name === "") {
          sweetalert.showInputError("You need to give your lobby a name!");
          return false
        }

        let lobbyUrl = "/lobbies/create"
        $.ajax({
          url: lobbyUrl,
          data: {name: name},
          type: "POST"
        })
        .done(function() {
          updateTable()
          swal({title: "Created Lobby '" + name + "'!",
                timer: 1500,
                showConfirmButton: false,
                type: "success"
              });
        })
      });
}


var lobbyClick = function() {
    $("#join-btn").removeAttr("disabled")
    $("#lobbies-table tr.lobby").removeClass("selected")
    $(this).addClass("selected")
    // Get lobby name from row
    let row = $(this).find("td")[0]
    let name = row.innerHTML
    var id = row.getAttribute("data-id")
    // Get users in lobby from server
    let lobbyUrl = "/lobbies/" + id + "/users"
    $.ajax({
      url: lobbyUrl,
      dataType: 'json'
    })
    .done(function(data) {
      console.log(data)
      let users = data
      // Generate table of users in specified lobby
      let table = "<table id='players'>\
                     <thead>\
                       <tr>\
                          <th class='left'>" + name + "</th>\
                       </tr>\
                     </thead>"

      table += "<tbody>"
      for (var i = 0; i < users.length; i++) {
        table += "<tr><td class='left'>" + users[i] + "</td>"
      }
      if (users.length == 0) {
        table += "<tr class='no-lobbies'><td>No players in lobby</td></tr>"
      }
      table += "</tbody>"

      document.getElementById("players-table").innerHTML = table
      document.getElementById("id-field").setAttribute("value", id)
    })
}

interface Lobby {
  ID: number
  Name: string
}

var updateTable = function() {
  // Get lobby IDs from server
  $.ajax({
    url: "/peekLobbies",
    dataType: 'json'
  })
  .done(function(data) {
      console.log(data)
      let lobbies = <Array<Lobby>> data

      // Generate table of lobby IDs

      let table = ""
      if (lobbies == null || lobbies.length == 0) {
        table += "<tr class='no-lobbies'><td>No available lobbies</td></tr>"
      } else {
        for (var i = 0; i < lobbies.length; i++) {
          table += "<tr class='lobby'>\
          <td class='left' data-id='" + lobbies[i].ID + "'>"
          + lobbies[i].Name + "</td>\
          <td class='right'>&#8594;</td></tr>"
        }
      }

      $("#lobbies-table").find("tbody")[0].innerHTML = table
      $("#lobbies-table tr.lobby").click(lobbyClick)
  })
}
