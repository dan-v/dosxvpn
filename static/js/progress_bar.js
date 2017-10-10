$(document).ready(function() {
    var currentPct = 0;
    function updateProgressBar(pct, extraMax, extraDuration) {
        if (currentPct != pct) {
            currentPct = pct;
            $('#current-progress').stop();
            $('#current-progress').animate({width: pct+'%'}, 1000, 'swing', function() {
            if (extraMax && extraDuration) {
                $('#current-progress').animate({width: extraMax+'%'}, extraDuration);
            }
            });
        }
    }

    function updateStatus() {
        $.get("/status/", function(resp) {
            console.log("Status: ", resp);
            if (resp.status == 'pending auth') {
                $('#status-line').text('Provisioning droplet...');
                updateProgressBar(5);
            } else if (resp.status == 'waiting for ssh') {
                $('#status-line').text('Setting up droplet...');
                if (resp.initial_ip) {
                    $('#initial-ip').text("Current Public IP: " + resp.initial_ip);
                    $('#initial-ip').css('display', 'block');
                }
                updateProgressBar(10, 65, 60000);
            } else if (resp.status == 'configuring vpn') {
                $('#status-line').text('Configuring VPN...');
                updateProgressBar(65, 90, 30000);
            } else if (resp.status == 'adding vpn to osx') {
                $('#status-line').text('Adding VPN to OSX...');
                updateProgressBar(93);
            } else if (resp.status == 'waiting for ip address change') {
                $('#status-line').text('Waiting for active VPN connection...');
                updateProgressBar(97);
            } else if (resp.status == 'done') {
                updateProgressBar(100);
                window.location.replace("/complete?final_ip=" + resp.final_ip);
            } else {
                updateProgressBar(0);
                $('#status-line').text('Install failed: ' + resp.status);
                $('#initial-ip').text('');
                $('#initial-ip').wrap('<a href="/"/>Retry Installation</a>');
                return;
            }
            if (resp.status != 'done') {
                setTimeout(updateStatus, 1000);
            }
        });
    }

    updateStatus();
});
