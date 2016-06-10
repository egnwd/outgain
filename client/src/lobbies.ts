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
        animation: "slide-from-top",
        inputPlaceholder: "Foo bar..."
      },
      function(name){
        if (name === false) return false;
        name = (name as string).trim()

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
          updateTable(<string>name)
        })
      });
}


var lobbyClick = function() {
    $("#join-btn").removeAttr("disabled")
    $("#lobbies-table tr.lobby").removeClass("selected")
    $(this).addClass("selected")
    // Get lobby name from row
    let name =  $(this).attr("data-name")
    var id = $(this).attr("data-id")
    // Get users in lobby from server
    let lobbyUrl = "/lobbies/" + id + "/users"
    $.ajax({
      url: lobbyUrl,
      dataType: 'json'
    })
    .done(function(data) {
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

var updateTable = function(select?: string) {
  // Get lobby IDs from server
  $.ajax({
    url: "/peekLobbies",
    dataType: 'json'
  })
  .done(redrawTable(select))
}

var redrawTable = function(select?: string) {
  return function(data) {
      let lobbies = <Array<Lobby>> data

      let table = ""
      if (lobbies == null || lobbies.length == 0) {
        table += "<tr class='no-lobbies'><td>No available lobbies</td></tr>"
      } else {
        for (var i = 0; i < lobbies.length; i++) {
          table += "<tr class='lobby' data-name='" + lobbies[i].Name + "' data-id='" + lobbies[i].ID + "'>\
          <td class='left'>"
          + lobbies[i].Name + "</td>\
          <td class='right'>&#8594;</td></tr>"
        }
      }

      $("#lobbies-table").find("tbody")[0].innerHTML = table
      $("#lobbies-table tr.lobby").click(lobbyClick)
      // Select the lobby we just created
      if (select) {
        $("#lobbies-table tr.lobby[data-name=\"" + select + "\"]").click()
      }
  }
}
