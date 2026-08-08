// campaigns.js: create campaign form handling
document.addEventListener('DOMContentLoaded', function () {
  var form = document.getElementById('create-campaign-form');
  if (!form) return;

  var msg = document.getElementById('campaign-msg');
  form.addEventListener('submit', function (e) {
    e.preventDefault();
    msg.textContent = 'Creating campaign...';

    var fd = new FormData(form);
    var payload = {
      title: fd.get('title'),
      location: fd.get('location'),
      description: fd.get('description'),
      target_amount: parseFloat(fd.get('target_amount'))
    };

    fetch('/api/campaigns', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'same-origin',
      body: JSON.stringify(payload)
    })
      .then(function (res) { return res.json().then(function (data) { return { ok: res.ok, data: data }; }); })
      .then(function (r) {
        if (r.ok) {
          msg.style.color = '#0a7d32';
          msg.textContent = 'Campaign created successfully!';
          setTimeout(function () { location.reload(); }, 800);
        } else {
          msg.style.color = '#c62828';
          msg.textContent = r.data.message || 'Could not create campaign.';
        }
      })
      .catch(function () {
        msg.style.color = '#c62828';
        msg.textContent = 'Network error. Please try again.';
      });
  });
});
