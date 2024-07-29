
urlParams = new URLSearchParams(window.location.search);
title = document.getElementById('filename');
version = document.getElementById('version');

if (urlParams.has('file')) {
    var file = urlParams.get('file');

    title.value = file.split('/')[file.split('/').length - 1].split('.')[0];

    fetch(window.location.origin + '/api/file?file=' + file).then(response => {
        if (response.ok) {
            response.text().then(text => {

                version.innerText = 'v' + text.split('\n')[0];
                var editor = new Editor(document.getElementById('editor'), text.split('\n').filter(line => line.length > 0));

            });
        } else {
            alert('Error: ' + response.status);
            console.error(response);
        }
    })
}

