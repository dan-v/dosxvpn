$(document).ready(function() {
    var getUrlParameter = function getUrlParameter(sParam) {
        var sPageURL = decodeURIComponent(window.location.search.substring(1)),
            sURLVariables = sPageURL.split('&'),
            sParameterName,
            i;

        for (i = 0; i < sURLVariables.length; i++) {
            sParameterName = sURLVariables[i].split('=');

            if (sParameterName[0] === sParam) {
                return sParameterName[1] === undefined ? true : sParameterName[1];
            }
        }
    };
    var token = getUrlParameter('access_token');
    var state = getUrlParameter('state');

    $("#regions").submit(function(e){
        return false;
    });

    $('#region-dropdown option[value="sfo2"]').attr("selected", true);

    $('.advanced').click(function() {
        $('#remove').css('display', 'block');
    });

    $('.go-btn').click(function() {
        var selected = $('#region-dropdown option:selected');
        window.location.replace("http://localhost:8999/install/" + state + "?region=" + selected.val());
    });

    $('.rem-btn').click(function() {
        var selected = $('#region-dropdown option:selected');
        window.location.replace("http://localhost:8999/uninstall?access_token=" + token);
    });

});