function showMovie(){
	document.querySelectorAll(".correctanswer").forEach(function(el) {
		el.style.color =  "#559955";
	});
}

function refreshGif(movieName){
	img = document.getElementById("js-guess-gif")
	img.src = "/static/images/loading.gif"
	var xhr = new XMLHttpRequest();

	xhr.onload = function () {
		if (xhr.status >= 200 && xhr.status < 300) {
			newGif = xhr.response
			img.src = newGif
		} else {
			console.log('The request failed!');
		}
	};
	xhr.open('GET', '/refresh' + '?title=' + movieName);
	xhr.send();
}
