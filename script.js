document.getElementById('uploadForm').addEventListener('submit', function(event) {
    event.preventDefault();

    const fileInput = document.getElementById('fileInput');
    const formData = new FormData();
    formData.append('file', fileInput.files[0]);

    fetch('/upload', {
        method: 'POST',
        body: formData
    })
    .then(response => response.text())
    .then(data => {
        document.getElementById('message').innerText = data;
        listFiles();
    })
    .catch(error => console.error('Error:', error));
});

function listFiles() {
    fetch('/uploads')
    .then(response => response.json())
    .then(files => {
        const fileList = document.getElementById('fileList');
        fileList.innerHTML = '';

        files.forEach(file => {
            const link = document.createElement('a');
            link.href = `/download/${file}`;
            link.innerText = file;
            link.target = '_blank';

            const listItem = document.createElement('div');
            listItem.appendChild(link);

            fileList.appendChild(listItem);
        });
    })
    .catch(error => console.error('Error:', error));
}

// Initial listing of files
listFiles();
