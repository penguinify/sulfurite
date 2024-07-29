
urlParams = new URLSearchParams(window.location.search);
title = document.getElementById('filename');

if (urlParams.has('file')) {
    var file = urlParams.get('file');

    title.innerHTML = file;

    fetch(window.location.origin + '/api/file?file=' + file).then(response => {
        if (response.ok) {
            response.text().then(text => {
                var editor = new Editor(document.getElementById('editor'), text.split('\n').filter(line => line.length > 0));
            });
        } else {
            alert('Error: ' + response.status);
            console.error(response);
        }
    })
}

