import * as $ from 'jquery'

$(function() {
    updateTable()
})

var updateTable = function() {
    $.ajax({
        url: "/peekLeaderboard",
        dataType: 'json'
    })
    .done(function(data) {
        let usernames = data.Usernames
        console.log(usernames)
        let scores = data.Scores
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
         + "<td>" + usernames[i] + "</td>"
         + "<td>" + scores[i] + "</td></tr>"
        }
        table +=  "</tbody>"
        document.getElementById("leaderboardTable").innerHTML = table
    })
}
