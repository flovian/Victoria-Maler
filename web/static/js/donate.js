// donate.js: donation form handling
document.addEventListener('DOMContentLoaded', function () {
  var form = document.getElementById('donate-form');
  if (!form) return;

  var msg = document.getElementById('donate-msg');
  form.addEventListener('submit', function (e) {
    e.preventDefault();
    if (msg) msg.textContent = 'Processing donation...';

    var fd = new FormData(form);
    var payload = {
      campaign_id: parseInt(fd.get('campaign_id')),
      donor_name: fd.get('donor_name'),
      amount: parseFloat(fd.get('amount')),
      message: fd.get('message')
    };

    fetch('/api/donations', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'same-origin',
      body: JSON.stringify(payload)
    })
      .then(function (res) { return res.json().then(function (data) { return { ok: res.ok, data: data }; }); })
      .then(function (r) {
        if (r.ok) {
          if (msg) {
            msg.style.color = '#0a7d32';
            msg.textContent = 'Thank you! Donation recorded.';
          }
          setTimeout(function () { location.reload(); }, 900);
        } else {
          if (msg) {
            msg.style.color = '#c62828';
            msg.textContent = r.data.message || 'Could not record donation.';
          }
        }
      })
      .catch(function () {
        if (msg) {
          msg.style.color = '#c62828';
          msg.textContent = 'Network error. Please try again.';
        }
      });
  });
});
