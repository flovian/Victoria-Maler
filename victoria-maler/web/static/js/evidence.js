// evidence.js: evidence upload handling
document.addEventListener('DOMContentLoaded', function () {
  var form = document.getElementById('evidence-form');
  if (!form) return;

  var msg = document.getElementById('evidence-msg');
  form.addEventListener('submit', function (e) {
    e.preventDefault();
    if (msg) msg.textContent = 'Uploading and anchoring evidence...';

    var fd = new FormData(form);

    fetch('/api/evidence?campaign_id=' + fd.get('campaign_id') + '&report_id=0', {
      method: 'POST',
      credentials: 'same-origin',
      body: fd
    })
      .then(function (res) { return res.json().then(function (data) { return { ok: res.ok, data: data }; }); })
      .then(function (r) {
        if (r.ok) {
          if (msg) {
            msg.style.color = '#0a7d32';
            msg.textContent = 'Evidence uploaded and anchored successfully!';
          }
          setTimeout(function () { location.reload(); }, 900);
        } else {
          if (msg) {
            msg.style.color = '#c62828';
            msg.textContent = r.data.message || 'Upload failed.';
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
