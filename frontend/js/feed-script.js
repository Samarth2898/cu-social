const searchInput = document.getElementById('searchInput');
const searchResults = document.getElementById('searchResults');

// Add an event listener to the search input
searchInput.addEventListener('input', function(event) {
    const searchText = event.target.value.toLowerCase().trim();
    // Perform a search based on the searchText
    // You can use AJAX/fetch requests to get data from the server or use client-side data

    // For example, if you have an array of items to search from
    const itemsToSearch = ['Item 1', 'Item 2', 'Item 3']; // Replace this with your data source
    const filteredItems = itemsToSearch.filter(item => item.toLowerCase().includes(searchText));

    // Display the search results
    displaySearchResults(filteredItems);
});

// Function to display search results
function displaySearchResults(results) {
    searchResults.innerHTML = ''; // Clear previous results
    results.forEach(result => {
        const li = document.createElement('li');
        li.textContent = result;
        searchResults.appendChild(li);
    });
}


// Call the function to update the profile image
updateProfileImage();