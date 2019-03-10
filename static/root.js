function updateNickname() {
    document.getElementById("nicknameStatus").innerHTML = "Updating...";
    var xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function() {
        if (xhr.readyState == 4) {
            document.getElementById("nicknameStatus").innerHTML = xhr.statusText;
        }
    }
    xhr.open("PUT", "/api/nickname", true);
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.send(JSON.stringify(
        document.getElementById("nicknameInput").value.trim()));
}

function clickTeam(season, bowlId, thisTeamId, thatTeamId) {
    var thisElement = document.getElementById(bowlId + "-" + thisTeamId);
    var thatElement = document.getElementById(bowlId + "-" + thatTeamId);
    // Save original settings.
    var thisOrigClass = thisElement.className;
    var thatOrigClass = thatElement.className;
    // Update UI.
    if (thisElement.className == "") {
        thisElement.className = "selected";
        thatElement.className = "";
    } else {
        thisElement.className = "";
    }
    // Update datastore.
    var path = "/api/seasons/" + escape(season) + "/bowls/" + escape(bowlId) + "/pick";
    var xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function() {
        if (xhr.readyState == 4) {
            if (xhr.status != 200) {
                // Revert UI on failure.
                thisElement.className = thisOrigClass;
                thatElement.className = thatOrigClass;
            }
        }
    }
    if (thisOrigClass == "") {
        xhr.open("PUT", path, true);
        xhr.setRequestHeader('Content-Type', 'application/json');
        xhr.send(JSON.stringify(thisTeamId))
    } else {
        xhr.open("DELETE", path, true);
        xhr.send();
    }
}
