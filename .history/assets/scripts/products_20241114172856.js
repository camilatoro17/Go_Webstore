$(document).ready(function() {
    // Handle Add Car button
    $("#add-car").on("click", function() {
        addCar();
    });

    // Handle Update Car button
    $("#update-car").on("click", function() {
        updateCar();
    });

    // Handle Delete Car button
    $("#delete-car").on("click", function() {
        deleteCar();
    });
});

// Function to add a car
function addCar() {
    // Code to validate and gather form data here
    const data = $("#product-form").serialize();

    $.ajax({
        type: "POST",
        url: "/add_car", // Update with the correct endpoint for adding
        data: data,
        success: function(response) {
            alert("Car added successfully");
            loadProductList(); // Refresh product list
            $("#product-form")[0].reset(); // Clear form
        },
        error: function(xhr) {
            alert("Error adding car: " + xhr.responseText);
        }
    });
}

// Function to update a car
function updateCar() {
    // Code to validate and gather form data here
    const data = $("#product-form").serialize();

    $.ajax({
        type: "POST",
        url: "/update_car", // Update with the correct endpoint for updating
        data: data,
        success: function(response) {
            alert("Car updated successfully");
            loadProductList(); // Refresh product list
            $("#product-form")[0].reset(); // Clear form
        },
        error: function(xhr) {
            alert("Error updating car: " + xhr.responseText);
        }
    });
}

// Function to delete a car
function deleteCar() {
    const carId = $("#product-form").find("input[name='id']").val(); // Get car ID from form if available

    if (!carId) {
        alert("Please select a car to delete.");
        return;
    }

    if (confirm("Are you sure you want to delete this car?")) {
        $.ajax({
            type: "POST",
            url: "/delete_car",
            data: { id: carId },
            success: function(response) {
                alert("Car deleted successfully");
                loadProductList();
                $("#product-form")[0].reset(); // Clear form
            },
            error: function(xhr) {
                alert("Error deleting car: " + xhr.responseText);
            }
        });
    }
}

// Function to load product list (refreshes list without page reload)
function loadProductList() {
    $.ajax({
        type: "GET",
        url: "/products", // Update URL if needed
        success: function(products) {
            const productTableBody = $("#product-table tbody"); // Make sure this ID matches the table in your HTML
            productTableBody.empty(); // Clear the existing table rows

            // Populate table with new data
            products.forEach(function(product) {
                const inactiveStatus = product.inactive === 1 ? "Yes" : "No";
                const row = `
                    <tr>
                        <td>${product.name}</td>
                        <td>${product.image}</td>
                        <td>${product.quantity}</td>
                        <td>${product.price.toFixed(2)}</td>
                        <td>${inactiveStatus}</td>
                    </tr>
                `;
                productTableBody.append(row);
            });
        },
        error: function(xhr) {
            alert("Error loading product list: " + xhr.responseText);
        }
    });
}
