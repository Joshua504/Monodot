const dropZone = document.getElementById("drop-zone");
const fileInput = document.getElementById("file-input");
const fileName = document.getElementById("file-name");

// Open file picker
dropZone.addEventListener("click", () => {
    fileInput.click();
});

// File selected manually
fileInput.addEventListener("change", () => {
    updateFileName();
});

// Prevent browser from opening the file
["dragenter", "dragover", "dragleave", "drop"].forEach(eventName => {
    dropZone.addEventListener(eventName, preventDefaults, false);
});

function preventDefaults(e) {
    e.preventDefault();
    e.stopPropagation();
}

// Highlight while dragging
["dragenter", "dragover"].forEach(eventName => {
    dropZone.addEventListener(eventName, () => {
        dropZone.classList.add("drag-over");
    });
});

// Remove highlight
["dragleave", "drop"].forEach(eventName => {
    dropZone.addEventListener(eventName, () => {
        dropZone.classList.remove("drag-over");
    });
});

// Handle dropped file
dropZone.addEventListener("drop", e => {
    const files = e.dataTransfer.files;

    if (files.length > 0) {
        fileInput.files = files;
        updateFileName();
    }
});

function updateFileName() {
    if (fileInput.files.length > 0) {
        fileName.textContent = fileInput.files[0].name;
        showPreview(fileInput.files[0]);
    } else {
        fileName.textContent = "No file selected";
    }
}

function showPreview(file) {

    const previewContainer = document.getElementById("preview-container");
    const previewImage = document.getElementById("preview-image");

    const imageSize = document.getElementById("image-size");
    const imageDimensions = document.getElementById("image-dimensions");

    const reader = new FileReader();

    reader.onload = function (e) {

        previewImage.src = e.target.result;

        const img = new Image();

        img.onload = function () {

            imageDimensions.textContent =
                `${img.width} × ${img.height}`;

            imageSize.textContent =
                `${(file.size / 1024).toFixed(1)} KB`;

            previewContainer.hidden = false;

        };

        img.src = e.target.result;

    };

    reader.readAsDataURL(file);

}

const uploadForm = document.getElementById("upload-form");
const generateButton = document.getElementById("generate-button");
const loading = document.getElementById("loading");

uploadForm.addEventListener("submit", () => {

    generateButton.disabled = true;

    generateButton.textContent = "Generating...";

    loading.hidden = false;

});