import * as $ from 'jquery'

$(function() {
    updateTable()
})

var updateTable = function() {
    $.ajax({
        url: window.location.pathname.replace("summary", "leaderboard"),
    })
    .done(function(data) {
        let table = "<table id='leaderboard-table'>\
                      <thead>\
                        <tr>\
                          <th></th>\
                          <th>Name</th>\
                          <th>Score</th>\
                        </tr>\
                      </thead>"
        table += "<tbody>"
        for (var i = 0; i < 10; i++) {
            table += "<tr>\
            <td>" + (i + 1).toString() + "</td>"
         + "<td>" + data[i].Name + "</td>"
         + "<td>" + data[i].Score + "</td></tr>"
        }
        table +=  "</tbody>"
        document.getElementById("leaderboardTable").innerHTML = table
    })
}

function getLobbyId() {
  let url = window.location.href.toString()
  let re = /([0-9]+)\/summary$/g
  return re.exec(url)[1]
}

$(function() {
  let idField = document.getElementById("id-field")
  let lobbyId = getLobbyId()
  idField.setAttribute("value", lobbyId)
})
