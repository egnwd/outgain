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
      table += "<tr><td class='left'>No available lobbies</td></tr>"
    }
    table += "</tbody>"

    document.getElementById("lobbies-table").innerHTML = table

    // Add onclick function to lobbies table rows
    $("#lobbies-table tr.lobby").click(function() {
      $("#join-btn").removeAttr("disabled")

      // Get lobby id from row
      let id = $(this).find("td")[0].innerHTML

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
          table += "<tr><td class='left'>No players in lobby</td></tr>"
        }
        table += "</tbody>"

        document.getElementById("players-table").innerHTML = table
        document.getElementById("id-field").setAttribute("value", id)
      })
    })
    // $("#lobbies-table tr.lobby")[0].click()
  })
})
