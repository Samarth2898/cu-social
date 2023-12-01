// Define the API endpoint for the profile image
const apiUrl = 'https://picsum.photos/200/300'; // Replace with your actual API endpoint

// Function to fetch the profile image
function updateProfileImage() {
    // Find the <a> element by its class name
    const profilePhotoLink = document.querySelector('.profile-photo');

    // Check if the <a> element exists
    if (profilePhotoLink) {
        // Fetch the profile image from the API
        fetch(apiUrl)
            .then(response => {
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                return response.blob();
            })
            .then(imageBlob => {
                const imageUrl = URL.createObjectURL(imageBlob);
                
                // Update the href attribute of the <a> element with the image URL
                profilePhotoLink.setAttribute('href', imageUrl);
            })
            .catch(error => {
                console.error('Error:', error);
            });
    }
}

// Call the function to update the profile image
updateProfileImage();