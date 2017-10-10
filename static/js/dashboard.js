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

    $("#advanced").submit(function(e){
        return false;
    });

    $('#region-dropdown option[value="sfo2"]').attr("selected", true);

    $('.go-btn').click(function() {
        var selected = $('#region-dropdown option:selected');
        window.location.replace("/install?region=" + selected.val());
    });
});

function removeVPN(name) {
    $('.loader').show();
    console.log("Called: ", name);
    $.get("/delete?droplet=" + name, function(resp) {
        console.log("Delete response: ", resp);
        window.location.replace("/dashboard?delete=true");
        $('.loader').hide();
    });
}