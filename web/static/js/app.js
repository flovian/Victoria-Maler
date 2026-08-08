// app.js: small UI helpers for header behavior
document.addEventListener('DOMContentLoaded', function(){
	var btn = document.querySelector('.nav-toggle');
	var nav = document.querySelector('.nav');
	if(btn && nav){
		btn.addEventListener('click', function(){
			if(nav.style.display === 'block') nav.style.display = '';
			else nav.style.display = 'block';
		});
	}
});

