import * as $ from 'jquery'

$(function() {
    updateList()
})

var updateList = function() {
    $.ajax({
        url: "/peekAchievements"
    })
    .done(function(data) {
        document.getElementById("achievements").innerHTML = data
    })
}
