import * as $ from 'jquery'

$(function() {
    // Create HTML table showing all lobby IDs with clickable rows
    updateTable()

    $("#create-lobby").click(function(event) {
        event.preventDefault();
        let lobbyUrl = "/lobbies/create"
        $.ajax({
          url: lobbyUrl,
          type: "POST"
        })
        .done(function() {
          updateTable()
        })
    })
})

var lobbyClick = function() {
    $("#join-btn").removeAttr("disabled")
    $("#lobbies-table tr.lobby").removeClass("selected")
    $(this).addClass("selected")
    // Get lobby id from row
    let id = $(this).find("td")[0].innerHTML
    console.log("Clicked on: " + id)
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
      // TODO: When lobbies have names include them here
      let table = "<table id='players'>\
                     <thead>\
                       <tr>\
                          <th class='left'>" + id + "</th>\
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

var updateTable = function() {
  // Get lobby IDs from server
  $.ajax({
    url: "/peekLobbies",
    dataType: 'json'
  })
  .done(function(data) {
      let lobbies = data

      // Generate table of lobby IDs
      let table = "<table id='lobbies-table'>\
                    <thead>\
                      <tr>\
                        <th class='left'>Lobbies</th>\
                      </tr>\
                    </thead>"
      table += "<tbody>"
      for (var i = 0; i < lobbies.length; i++) {
        table += "<tr class='lobby'>\
                    <td class='left'>" + lobbies[i] + "</td>"
            + "<td class='right'>&#8594;</td></tr>"
      }
      if (lobbies.length == 0) {
        table += "<tr class='no-lobbies'><td>No available lobbies</td></tr>"
      }
      table += "</tbody>"

      document.getElementById("lobbies-table").innerHTML = table

      $("#lobbies-table tr.lobby").click(lobbyClick)
  })
}
