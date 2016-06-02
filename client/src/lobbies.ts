import * as $ from 'jquery'

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
    let table = "<table id='lobbies'><thead><tr><th class='left'>Lobbies</th>"
          + "<th class='right'><a class='btn'>+</a></th></thead>"
    for (var i = 0; i < lobbies.length; i++) {
      table += "<tr><td class='left'>" + lobbies[i] + "</td>"
          + "<td class='right'>></td></tr>"
    }
    // Find way to ignore this from onclick functions before uncommenting
    /*
    if (lobbies.length == 0) {
      table += "<tr><td class='left'>No available lobbies</td></tr>"
    }
    */
    table += "</table>"
    document.getElementById("lobbies-table").innerHTML = table

    // Add onclick function to lobbies table rows
    let rows = document.getElementById("lobbies").getElementsByTagName("tr")
    for (i = 1; i < rows.length; i++) {
      let curr = rows[i]
      curr.onclick = createClickHandler(curr)
    }
  })
})

function createClickHandler(row) {
  return function() {
    // Get lobby id from row
    let id = row.getElementsByTagName("td")[0].innerHTML

    // Get users in lobby from server
    let lobbyUrl = "/lobbies/" + id + "/users"
    $.ajax({
      url: lobbyUrl,
      dataType: 'json'
    })
    .done(function(data) {
      let users = data

      // Generate table of users in specified lobby
      // TODO: When lobbies have names include them here
      let table = "<table id='players'><thead><tr><th class='left'>"
              + "Lobby name here" + "</th></thead>"
      for (var i = 0; i < users.length; i++) {
        table += "<tr><td class='left'>" + users[i] + "</td>"
      }
      if (users.length == 0) {
        table += "<tr><td class='left'>No players in lobby</td></tr>"
      }
        table += "</table>"
      document.getElementById("players-table").innerHTML = table

      document.getElementById("id-field").setAttribute("value", id)
    })
  }
}
