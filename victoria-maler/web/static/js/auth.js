// auth.js: login and register form handling
document.addEventListener('DOMContentLoaded', function () {
  var msg = document.getElementById('auth-msg');

  function showMsg(text, isError) {
    if (!msg) return;
    msg.style.color = isError ? '#c62828' : '#0a7d32';
    msg.textContent = text;
  }

  var loginForm = document.getElementById('login-form');
  if (loginForm) {
    loginForm.addEventListener('submit', function (e) {
      e.preventDefault();
      showMsg('Logging in...', false);

      var fd = new FormData(loginForm);
      fetch('/api/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'same-origin',
        body: JSON.stringify({ email: fd.get('email'), password: fd.get('password') })
      })
        .then(function (res) { return res.json().then(function (data) { return { ok: res.ok, data: data }; }); })
        .then(function (r) {
          if (r.ok) {
            showMsg('Login successful, redirecting...', false);
            setTimeout(function () { window.location.href = '/dashboard.html'; }, 600);
          } else {
            showMsg(r.data.message || 'Login failed.', true);
          }
        })
        .catch(function () { showMsg('Network error. Please try again.', true); });
    });
  }

  var registerForm = document.getElementById('register-form');
  if (registerForm) {
    registerForm.addEventListener('submit', function (e) {
      e.preventDefault();
      showMsg('Creating account...', false);

      var fd = new FormData(registerForm);
      fetch('/api/auth/register', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'same-origin',
        body: JSON.stringify({ name: fd.get('name'), email: fd.get('email'), password: fd.get('password') })
      })
        .then(function (res) { return res.json().then(function (data) { return { ok: res.ok, data: data }; }); })
        .then(function (r) {
          if (r.ok) {
            showMsg('Account created, logging you in...', false);
            setTimeout(function () { window.location.href = '/dashboard.html'; }, 600);
          } else {
            showMsg(r.data.message || 'Registration failed.', true);
          }
        })
        .catch(function () { showMsg('Network error. Please try again.', true); });
    });
  }
});
