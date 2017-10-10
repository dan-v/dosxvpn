function getParameterByName(name, url) {
    if (!url) {
        url = window.location.hash.replace("#", "?");
    }
    name = name.replace(/[\[\]]/g, "\\$&");
    var regex = new RegExp("[?&]" + name + "(=([^&#]*)|&|#|$)"),
        results = regex.exec(url);
    if (!results) return null;
    if (!results[2]) return '';
    return decodeURIComponent(results[2].replace(/\+/g, " "));
}

$(document).ready(function() {
    var access_token = getParameterByName('access_token');
    var state = getParameterByName('state');
    window.location.replace("/dashboard?access_token=" + access_token + "&state=" + state);
});
