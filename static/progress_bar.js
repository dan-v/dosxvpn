$(document).ready(function() {

  var currentPct = 0;

  function updateProgressBar(pct, extraMax, extraDuration) {
      if (currentPct != pct) {
          currentPct = pct;
          $('#current-progress').stop();
          $('#current-progress').animate({width: pct+'%'}, 1000, 'swing', function() {

            // If there's an extra easing, then animate that too.
            if (extraMax && extraDuration) {
              $('#current-progress').animate({width: extraMax+'%'}, extraDuration);
            }
          });
      }
  }

    $('#exit').click(function() {
        $.get("/exit", function(resp) {
            return
        });
    });

  function updateUI() {
      $.get("/status/"+window.installID, function(resp) {
          console.log(resp);
          if (resp.status == 'pending auth') {
              $('#status-line').text('Provisioning droplet…');
              updateProgressBar(5);
          } else if (resp.status == 'waiting for ssh') {
              $('#status-line').text('Setting up droplet…');
              updateProgressBar(10, 65, 60000);
          } else if (resp.status == 'configuring vpn') {
              $('#status-line').text('Configuring VPN…');
              updateProgressBar(70, 95, 30000);
          } else if (resp.status == 'adding vpn to osx') {
              $('#status-line').text('Adding VPN to OSX');
              updateProgressBar(98);
          } else if (resp.status == 'done') {
              $('#status-line').text('Install complete');
              $('#mobileconfig').css('display', 'block');
              $('#exit').css('display', 'block')
              updateProgressBar(100);;
          } else {
              $('status-line').text('Install failed: ' + resp.status);
              updateProgressBar(0);
          }

          if (resp.status != 'done') {
              setTimeout(updateUI, 1000);
          }
      });
  }

    updateUI();
});
