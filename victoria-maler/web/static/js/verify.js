// verify.js: bitcoin hash verification
document.addEventListener('DOMContentLoaded', function () {
  var form = document.getElementById('verify-form');
  if (!form) return;

  var result = document.getElementById('verify-result');
  form.addEventListener('submit', function (e) {
    e.preventDefault();
    if (result) {
      result.textContent = 'Verifying...';
      result.className = '';
    }

    var fd = new FormData(form);
    var payload = {
      hash: fd.get('hash'),
      txid: fd.get('txid')
    };

    fetch('/api/verify', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    })
      .then(function (res) { return res.json().then(function (data) { return { ok: res.ok, data: data }; }); })
      .then(function (r) {
        if (!result) return;
        var d = r.data && r.data.data ? r.data.data : {};
        if (r.ok && d.verified) {
          result.className = 'verify-result ok';
          result.textContent = 'VERIFIED: hash matches the anchored Bitcoin record.';
        } else if (r.ok) {
          result.className = 'verify-result warn';
          result.textContent = (d.note || 'Hash does not match any anchored record.');
        } else {
          result.className = 'verify-result warn';
          result.textContent = (r.data && r.data.message) || 'Verification failed.';
        }
      })
      .catch(function () {
        if (result) {
          result.className = 'verify-result warn';
          result.textContent = 'Network error. Please try again.';
        }
      });
  });
});
